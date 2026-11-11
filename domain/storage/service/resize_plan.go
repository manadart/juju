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
	"github.com/juju/juju/internal/errors"
)

// GetModelStorageResizeRequests returns requests needing planning or backend
// reconciliation by the model storage provisioner.
func (s *Service) GetModelStorageResizeRequests(ctx context.Context) ([]string, error) {
	return s.st.GetStorageResizeStageRequests(ctx, "")
}

// GetMachineStorageResizeStageRequests returns this machine's assigned stages.
func (s *Service) GetMachineStorageResizeStageRequests(ctx context.Context, machineUUID coremachine.UUID) ([]string, error) {
	if err := machineUUID.Validate(); err != nil {
		return nil, errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	return s.st.GetStorageResizeStageRequests(ctx, machineUUID.String())
}

// PrepareModelStorageResize creates a plan or returns the next provider step.
func (s *Service) PrepareModelStorageResize(ctx context.Context, uuid domainstorage.StorageResizeUUID) (domainstorage.ResizeStageWork, error) {
	return s.prepareStorageResizeStage(ctx, uuid, "")
}

// PrepareMachineStorageResizeStage returns only this machine's eligible work.
func (s *Service) PrepareMachineStorageResizeStage(ctx context.Context, uuid domainstorage.StorageResizeUUID, machineUUID coremachine.UUID) (domainstorage.ResizeStageWork, error) {
	if err := machineUUID.Validate(); err != nil {
		return domainstorage.ResizeStageWork{}, errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	return s.prepareStorageResizeStage(ctx, uuid, machineUUID.String())
}

func (s *Service) prepareStorageResizeStage(ctx context.Context, uuid domainstorage.StorageResizeUUID, machineUUID string) (domainstorage.ResizeStageWork, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()
	if err := uuid.Validate(); err != nil {
		return domainstorage.ResizeStageWork{}, errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	plan, err := s.st.GetStorageResizePlan(ctx, uuid.String())
	if err != nil {
		return domainstorage.ResizeStageWork{}, errors.Errorf("getting resize plan %q: %w", uuid, err)
	}
	if machineUUID != "" {
		if plan.Identity == "" {
			if !plan.Target.MachineScoped || plan.Target.MachineUUID != machineUUID {
				return domainstorage.ResizeStageWork{}, storageerrors.StorageResizeNotAssigned
			}
		} else if plan.BackendMachineUUID != machineUUID && plan.NodeMachineUUID != machineUUID {
			return domainstorage.ResizeStageWork{}, storageerrors.StorageResizeNotAssigned
		}
	}
	fail := func(message string) (domainstorage.ResizeStageWork, error) {
		if err := s.st.FailStorageResizePlan(ctx, uuid.String(), machineUUID, message, s.StorageService.clock.Now().UTC()); err != nil {
			return domainstorage.ResizeStageWork{}, err
		}
		return domainstorage.ResizeStageWork{}, storageerrors.StorageResizeResultAlreadyRecorded
	}
	if !plan.Target.Alive {
		return fail("Storage or its attachment is no longer alive")
	}
	if plan.Identity != "" && plan.Identity != plan.Target.Identity {
		return fail("Storage resize target or attachment changed")
	}
	if plan.Target.ReadOnly || plan.Target.Shared {
		return fail("Resize requires a single writable attachment")
	}
	if plan.Target.SizeMiB > math.MaxInt64/(1024*1024) {
		return fail("Resize target exceeds the supported size range")
	}
	if plan.Identity == "" {
		t := plan.Target
		if t.ProviderID == "" || (t.MachineScoped && t.MachineUUID == "") ||
			(t.Kind == "volume" && ((t.FilesystemID != "" && (t.MachineUUID == "" || t.MountPoint == "")) || (t.MachineUUID != "" && t.DeviceName == ""))) {
			return domainstorage.ResizeStageWork{}, storageerrors.StorageResizeNotReady
		}
		plan.Identity = t.Identity
		if t.MachineScoped {
			plan.BackendMachineUUID = t.MachineUUID
		}
		if t.Kind == "volume" {
			plan.NodeMachineUUID = t.MachineUUID
		}
		if err := s.st.CreateStorageResizePlan(ctx, plan); err != nil {
			if errors.Is(err, storageerrors.StorageResizeTargetChanged) {
				return fail("Storage resize target changed during planning")
			}
			return domainstorage.ResizeStageWork{}, errors.Errorf("creating resize plan %q: %w", uuid, err)
		}
		plan, err = s.st.GetStorageResizePlan(ctx, uuid.String())
		if err != nil {
			return domainstorage.ResizeStageWork{}, err
		}
	}
	stage, phase := plan.Next(machineUUID)
	if stage == "" {
		return domainstorage.ResizeStageWork{}, storageerrors.StorageResizeNotReady
	}
	return domainstorage.ResizeStageWork{UUID: uuid.String(), Identity: plan.Identity, Stage: stage, Phase: phase, Target: plan.Target}, nil
}

// RecordModelStorageResizeStage records one provider check or ensure result.
func (s *Service) RecordModelStorageResizeStage(ctx context.Context, report domainstorage.ResizeStageReport) error {
	return s.recordStorageResizeStage(ctx, "", report)
}

// RecordMachineStorageResizeStage records one step for the assigned machine.
func (s *Service) RecordMachineStorageResizeStage(ctx context.Context, machineUUID coremachine.UUID, report domainstorage.ResizeStageReport) error {
	if err := machineUUID.Validate(); err != nil {
		return errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	return s.recordStorageResizeStage(ctx, machineUUID.String(), report)
}

func (s *Service) recordStorageResizeStage(ctx context.Context, machineUUID string, report domainstorage.ResizeStageReport) error {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()
	if err := domainstorage.StorageResizeUUID(report.UUID).Validate(); err != nil {
		return errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	if report.Identity == "" || (report.Stage != "backend" && report.Stage != "attachment") || (report.Phase != "check" && report.Phase != "ensure") {
		return errors.New("invalid resize stage identity or phase").Add(coreerrors.NotValid)
	}
	switch report.Status {
	case "blocked", "checked":
		if report.Phase != "check" {
			return errors.New("check result outside check phase").Add(coreerrors.NotValid)
		}
	case "pending", "completed":
		if report.Phase != "ensure" {
			return errors.New("execution result outside ensure phase").Add(coreerrors.NotValid)
		}
	case "failed":
		if strings.TrimSpace(report.Message) == "" {
			return errors.New("failed stage requires a message").Add(coreerrors.NotValid)
		}
	default:
		return errors.New("invalid resize stage result").Add(coreerrors.NotValid)
	}
	if report.ObservedSizeMiB > math.MaxInt64/(1024*1024) || report.FilesystemSizeMiB > math.MaxInt64/(1024*1024) {
		return errors.New("observed capacity is out of range").Add(coreerrors.NotValid)
	}
	plan, err := s.st.GetStorageResizePlan(ctx, report.UUID)
	if err != nil {
		return err
	}
	// The state layer repeats ownership, phase and identity checks atomically.
	if report.Status == "checked" {
		reason := ""
		if len(report.Prerequisites) != 0 {
			reason = "Resize prerequisites are not supported by this prototype"
		}
		for _, requirement := range report.FollowUp {
			switch requirement {
			case "device-refresh":
				if plan.NodeMachineUUID == "" {
					reason = "Device refresh requires a machine attachment"
				}
			case "filesystem-grow":
				if plan.NodeMachineUUID == "" || plan.Target.FilesystemID == "" {
					reason = "Filesystem growth requires a managed filesystem attachment"
				}
			default:
				reason = "Unsupported resize follow-up requirement: " + requirement
			}
		}
		if reason != "" {
			report.Status = "failed"
			report.Message = reason
		}
	}
	if report.Status == "completed" && (report.ObservedSizeMiB < plan.Target.SizeMiB || (report.Stage == "attachment" && plan.Target.FilesystemID != "" && report.FilesystemSizeMiB < plan.Target.SizeMiB)) {
		return errors.New("completed stage requires sufficient observed capacity").Add(coreerrors.NotValid)
	}
	if report.Status == "failed" {
		report.Message = report.Stage + " " + report.Phase + ": " + report.Message
	}
	if err := s.st.SetStorageResizeStage(ctx, machineUUID, report, s.StorageService.clock.Now().UTC()); err != nil {
		return errors.Errorf("recording resize stage for %q: %w", report.UUID, err)
	}
	return nil
}
