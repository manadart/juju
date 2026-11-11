// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storageresizer

import (
	"context"
	"reflect"

	"github.com/juju/names/v6"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/apiserver/facade"
	coremachine "github.com/juju/juju/core/machine"
	domainstorage "github.com/juju/juju/domain/storage"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/rpc/params"
)

// StorageService provides machine-specific resize operations.
type StorageService interface {
	GetMachineStorageResizeStageRequests(context.Context, coremachine.UUID) ([]string, error)
	PrepareMachineStorageResizeStage(context.Context, domainstorage.StorageResizeUUID, coremachine.UUID) (domainstorage.ResizeStageWork, error)
	RecordMachineStorageResizeStage(context.Context, coremachine.UUID, domainstorage.ResizeStageReport) error

	GetStorageResizeRequestsForMachine(context.Context, coremachine.UUID) ([]domainstorage.StorageResizeUUID, error)
	PrepareMachineStorageResize(context.Context, domainstorage.StorageResizeUUID, coremachine.UUID) (domainstorage.StorageResizeWork, error)
	RecordMachineStorageResizeResult(context.Context, domainstorage.StorageResizeUUID, coremachine.UUID, domainstorage.StorageResizeReport) error
}

// API exposes only the authenticated machine's resize work.
type API struct {
	service     StorageService
	machineUUID coremachine.UUID
}

// Register exposes the machine-only StorageResizer facade.
func Register(registry facade.FacadeRegistry) {
	registry.MustRegister("StorageResizer", 1, func(ctx context.Context, f facade.ModelContext) (facade.Facade, error) {
		if !f.Auth().AuthMachineAgent() {
			return nil, apiservererrors.ErrPerm
		}
		tag, ok := f.Auth().GetAuthTag().(names.MachineTag)
		if !ok {
			return nil, apiservererrors.ErrPerm
		}
		machineUUID, err := f.DomainServices().Machine().GetMachineUUID(ctx, coremachine.Name(tag.Id()))
		if err != nil {
			return nil, apiservererrors.ServerError(err)
		}
		return &API{service: f.DomainServices().Storage(), machineUUID: machineUUID}, nil
	}, reflect.TypeFor[*API]())
	registry.MustRegister("StorageResizer", 2, func(ctx context.Context, f facade.ModelContext) (facade.Facade, error) {
		if !f.Auth().AuthMachineAgent() {
			return nil, apiservererrors.ErrPerm
		}
		tag, ok := f.Auth().GetAuthTag().(names.MachineTag)
		if !ok {
			return nil, apiservererrors.ErrPerm
		}
		machineUUID, err := f.DomainServices().Machine().GetMachineUUID(ctx, coremachine.Name(tag.Id()))
		if err != nil {
			return nil, apiservererrors.ServerError(err)
		}
		return &APIv2{API: &API{service: f.DomainServices().Storage(), machineUUID: machineUUID}}, nil
	}, reflect.TypeFor[*APIv2]())
}

// Requests returns pending work for this machine, including interrupted work.
func (a *API) Requests(ctx context.Context) (params.StringsResult, error) {
	ids, err := a.service.GetStorageResizeRequestsForMachine(ctx, a.machineUUID)
	if err != nil {
		return params.StringsResult{}, resizeError(err)
	}
	result := params.StringsResult{Result: make([]string, len(ids))}
	for i, id := range ids {
		result.Result[i] = id.String()
	}
	return result, nil
}

// Prepare binds the request and revalidates its target before execution.
func (a *API) Prepare(ctx context.Context, args params.StorageResizeRequestArgs) (params.StorageResizeWork, error) {
	work, err := a.service.PrepareMachineStorageResize(ctx, domainstorage.StorageResizeUUID(args.UUID), a.machineUUID)
	if err != nil {
		return params.StorageResizeWork{}, resizeError(err)
	}
	return params.StorageResizeWork{
		UUID: work.UUID.String(), RequestedSizeMiB: work.RequestedSizeMiB,
		ProviderType: work.ProviderType, VolumeID: work.VolumeID,
		VolumeProviderID: work.VolumeProviderID, DeviceName: work.DeviceName,
		FilesystemID: work.FilesystemID, MountPoint: work.MountPoint,
	}, nil
}

// Report records an outcome only for a request assigned to this machine.
func (a *API) Report(ctx context.Context, args params.StorageResizeReport) error {
	err := a.service.RecordMachineStorageResizeResult(ctx, domainstorage.StorageResizeUUID(args.UUID), a.machineUUID, domainstorage.StorageResizeReport{
		StorageResizeResult: domainstorage.StorageResizeResult{
			Status: domainstorage.StorageResizeStatus(args.Status), Message: args.Message,
		},
		VolumeSizeMiB: args.VolumeSizeMiB, FilesystemSizeMiB: args.FilesystemSizeMiB,
	})
	return resizeError(err)
}

func resizeError(err error) error {
	if err == nil {
		return nil
	}
	var code string
	switch {
	case errors.Is(err, storageerrors.StorageResizeNotAssigned):
		code = params.CodeUnauthorized
	case errors.Is(err, storageerrors.StorageResizeIntentNotFound):
		code = params.CodeNotFound
	case errors.Is(err, storageerrors.StorageResizeNotReady):
		code = params.CodeTryAgain
	case errors.Is(err, storageerrors.StorageResizeResultAlreadyRecorded):
		code = params.CodeAlreadyExists
	case errors.Is(err, storageerrors.StorageResizeTargetChanged), errors.Is(err, storageerrors.StorageInstanceNotAlive):
		code = params.CodeNotValid
	default:
		return apiservererrors.ServerError(err)
	}
	return apiservererrors.ParamsErrorf(code, "%s", err)
}

// APIv2 adds stage reports while retaining the guarded legacy methods.
type APIv2 struct{ *API }

// StageRequests returns only work assigned to the authenticated machine.
func (a *APIv2) StageRequests(ctx context.Context) (params.StringsResult, error) {
	ids, err := a.service.GetMachineStorageResizeStageRequests(ctx, a.machineUUID)
	return params.StringsResult{Result: ids}, resizeError(err)
}

// PrepareStage revalidates and returns the next eligible machine step.
func (a *APIv2) PrepareStage(ctx context.Context, args params.StorageResizeRequestArgs) (params.StorageResizeStageWork, error) {
	work, err := a.service.PrepareMachineStorageResizeStage(ctx, domainstorage.StorageResizeUUID(args.UUID), a.machineUUID)
	if err != nil {
		return params.StorageResizeStageWork{}, resizeError(err)
	}
	// Provider pool configuration is needed only by the model executor.
	work.Target.PoolAttributes = nil
	return params.StorageResizeStageWork{UUID: work.UUID, Identity: work.Identity, Stage: work.Stage, Phase: work.Phase, Target: params.StorageResizeStageTarget(work.Target)}, nil
}

// ReportStage records a check or observed execution step for this machine.
func (a *APIv2) ReportStage(ctx context.Context, args params.StorageResizeStageReport) error {
	return resizeError(a.service.RecordMachineStorageResizeStage(ctx, a.machineUUID, domainstorage.ResizeStageReport(args)))
}
