// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storageprovisioner

import (
	"context"

	domainstorage "github.com/juju/juju/domain/storage"
	"github.com/juju/juju/rpc/params"
)

// ModelResizeService coordinates provider stages without an agent connection.
type ModelResizeService interface {
	GetModelStorageResizeRequests(context.Context) ([]string, error)
	PrepareModelStorageResize(context.Context, domainstorage.StorageResizeUUID) (domainstorage.ResizeStageWork, error)
	RecordModelStorageResizeStage(context.Context, domainstorage.ResizeStageReport) error
}

// modelResizeAdapter converts domain values to the shared provisioner work
// representation; provider configuration never travels through the agent API.
type modelResizeAdapter struct{ service ModelResizeService }

func (a *modelResizeAdapter) Requests(ctx context.Context) ([]string, error) {
	return a.service.GetModelStorageResizeRequests(ctx)
}
func (a *modelResizeAdapter) PrepareStage(ctx context.Context, uuid string) (params.StorageResizeStageWork, error) {
	work, err := a.service.PrepareModelStorageResize(ctx, domainstorage.StorageResizeUUID(uuid))
	if err != nil {
		return params.StorageResizeStageWork{}, err
	}
	return params.StorageResizeStageWork{UUID: work.UUID, Identity: work.Identity, Stage: work.Stage, Phase: work.Phase, Target: params.StorageResizeStageTarget(work.Target)}, nil
}
func (a *modelResizeAdapter) ReportStage(ctx context.Context, report params.StorageResizeStageReport) error {
	return a.service.RecordModelStorageResizeStage(ctx, domainstorage.ResizeStageReport(report))
}
