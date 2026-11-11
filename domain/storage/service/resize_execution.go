// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"math"
	"strings"

	coreerrors "github.com/juju/juju/core/errors"
	coremachine "github.com/juju/juju/core/machine"
	"github.com/juju/juju/core/trace"
	domainstorage "github.com/juju/juju/domain/storage"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/domain/storage/internal"
	"github.com/juju/juju/internal/errors"
)

// GetStorageResizeRequestsForMachine returns pending work for a machine.
// Unattached requests remain pending until an attachment supplies an executor.
func (s *Service) GetStorageResizeRequestsForMachine(ctx context.Context, machineUUID coremachine.UUID) ([]domainstorage.StorageResizeUUID, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()
	if err := machineUUID.Validate(); err != nil {
		return nil, errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	ids, err := s.st.GetStorageResizeRequestsForMachine(ctx, machineUUID.String())
	if err != nil {
		return nil, errors.Errorf("getting resize requests for machine %q: %w", machineUUID, err)
	}
	result := make([]domainstorage.StorageResizeUUID, len(ids))
	for i, id := range ids {
		result[i] = domainstorage.StorageResizeUUID(id)
	}
	return result, nil
}

// PrepareMachineStorageResize binds execution and checks the persisted target.
// Unsupported or invalid requests are bound so the executor can report failure.
// StorageResizeNotReady means provisioning should be allowed to finish first.
func (s *Service) PrepareMachineStorageResize(ctx context.Context, uuid domainstorage.StorageResizeUUID, machineUUID coremachine.UUID) (domainstorage.StorageResizeWork, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()
	if err := uuid.Validate(); err != nil {
		return domainstorage.StorageResizeWork{}, errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	if err := machineUUID.Validate(); err != nil {
		return domainstorage.StorageResizeWork{}, errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	work, err := s.st.AssignStorageResize(ctx, uuid.String(), machineUUID.String())
	if err != nil {
		return domainstorage.StorageResizeWork{}, errors.Errorf("preparing resize request %q: %w", uuid, err)
	}
	if !work.Alive {
		return domainstorage.StorageResizeWork{}, storageerrors.StorageInstanceNotAlive
	}
	if work.ProviderType != "loop" || work.VolumeID == "" || !work.MachineScoped {
		return domainstorage.StorageResizeWork{}, errors.New("this prototype supports machine-local loop volumes only").Add(coreerrors.NotSupported)
	}
	if work.ReadOnly || work.Shared {
		return domainstorage.StorageResizeWork{}, errors.New("resize requires a single writable attachment").Add(coreerrors.NotSupported)
	}
	if work.RequestedSizeMiB < work.VolumeSizeMiB {
		return domainstorage.StorageResizeWork{}, errors.New("storage shrinking is not supported").Add(coreerrors.NotSupported)
	}
	if work.RequestedSizeMiB > math.MaxInt64/(1024*1024) {
		return domainstorage.StorageResizeWork{}, errors.New("requested size exceeds the supported local file size").Add(coreerrors.NotValid)
	}
	return domainstorage.StorageResizeWork{
		UUID: uuid, RequestedSizeMiB: work.RequestedSizeMiB,
		ProviderType: work.ProviderType, VolumeID: work.VolumeID,
		VolumeProviderID: work.VolumeProviderID, DeviceName: work.DeviceName,
		FilesystemID: work.FilesystemID, MountPoint: work.MountPoint,
	}, nil
}

// RecordMachineStorageResizeResult records the assigned executor's outcome.
// Successful results require observed capacity at each required layer. Failed
// results may also report partial growth; growth is never rolled back here.
func (s *Service) RecordMachineStorageResizeResult(ctx context.Context, uuid domainstorage.StorageResizeUUID, machineUUID coremachine.UUID, report domainstorage.StorageResizeReport) error {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()
	if err := uuid.Validate(); err != nil {
		return errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	if err := machineUUID.Validate(); err != nil {
		return errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	switch report.Status {
	case domainstorage.StorageResizeCompleted:
		if report.VolumeSizeMiB == 0 {
			return errors.New("completed resize requires observed volume capacity").Add(coreerrors.NotValid)
		}
	case domainstorage.StorageResizeFailed:
		if strings.TrimSpace(report.Message) == "" {
			return errors.New("failed resize requires a message").Add(coreerrors.NotValid)
		}
	default:
		return errors.New("resize result must be completed or failed").Add(coreerrors.NotValid)
	}
	if report.VolumeSizeMiB > math.MaxInt64/(1024*1024) || report.FilesystemSizeMiB > math.MaxInt64/(1024*1024) {
		return errors.New("reported capacity is out of range").Add(coreerrors.NotValid)
	}
	err := s.st.SetMachineStorageResizeResult(ctx, uuid.String(), machineUUID.String(), internal.StorageResizeReport{
		StorageResizeResult: internal.StorageResizeResult{
			Status: string(report.Status), Message: report.Message,
			CompletedAt: s.StorageService.clock.Now().UTC(),
		},
		VolumeSizeMiB: report.VolumeSizeMiB, FilesystemSizeMiB: report.FilesystemSizeMiB,
	})
	if err != nil {
		return errors.Errorf("recording machine result for resize request %q: %w", uuid, err)
	}
	return nil
}
