// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storageprovisioner

import (
	"context"
	"time"

	"github.com/juju/names/v6"

	"github.com/juju/juju/internal/errors"
	storageresize "github.com/juju/juju/internal/storage/resize"
	"github.com/juju/juju/rpc/params"
)

// ResizeAccessor provides authoritative work for the authenticated machine.
// Prepare binds and revalidates the target. Report is idempotent for an
// identical terminal outcome, allowing retries after an ambiguous RPC failure.
type ResizeAccessor interface {
	Requests(context.Context) ([]string, error)
	Prepare(context.Context, string) (params.StorageResizeWork, error)
	Report(context.Context, params.StorageResizeReport) error
}

// ResizeExecutor grows an existing local resource and measures its capacity.
// Implementations must honour cancellation, never shrink or replace storage,
// and call revalidate before each mutation. Stages must be safe to repeat after
// interruption. Partial observations may accompany a failure.
type ResizeExecutor interface {
	Resize(context.Context, storageresize.Request, func(context.Context) error) (storageresize.Observation, error)
}

type resizeOutcome struct {
	work   params.StorageResizeWork
	report params.StorageResizeReport
}

func processStorageResizes(ctx context.Context, deps *dependencies, results map[string]resizeOutcome) error {
	ids, err := deps.config.Resizes.Requests(ctx)
	if err != nil {
		return errors.Errorf("fetching resize requests: %w", err)
	}
	// A result may already have reached the controller despite an RPC error.
	pending := make(map[string]bool, len(ids))
	for _, id := range ids {
		pending[id] = true
	}
	for id := range results {
		if !pending[id] {
			delete(results, id)
		}
	}
	for _, id := range ids {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, ok := results[id]; !ok {
			outcome, err := executeStorageResize(ctx, deps, id)
			if err != nil {
				return errors.Errorf("executing storage resize %q: %w", id, err)
			}
			if outcome.report.UUID == "" {
				continue
			}
			results[id] = outcome
		}
		outcome := results[id]
		err := deps.config.Resizes.Report(ctx, outcome.report)
		switch params.ErrCode(err) {
		case params.CodeNotFound, params.CodeAlreadyExists, params.CodeUnauthorized:
			delete(results, id)
			continue
		case params.CodeNotValid:
			// The target may have changed after local execution. Do not retry
			// local mutations or apply old observations to a new attachment.
			if outcome.report.Status == "completed" {
				outcome.report = params.StorageResizeReport{
					UUID: id, Status: "failed",
					Message: "Storage grew locally but completion was rejected: " + err.Error(),
				}
				results[id] = outcome
				continue
			}
		}
		if err != nil {
			return errors.Errorf("reporting storage resize %q: %w", id, err)
		}
		delete(results, id)
		if outcome.report.Status == "completed" {
			updateResizeCache(deps, outcome)
		}
		deps.config.Logger.Infof(ctx, "storage resize %q %s: %s", id, outcome.report.Status, outcome.report.Message)
	}
	return nil
}

func executeStorageResize(ctx context.Context, deps *dependencies, id string) (resizeOutcome, error) {
	work, err := deps.config.Resizes.Prepare(ctx, id)
	if err != nil {
		switch params.ErrCode(err) {
		case params.CodeTryAgain, params.CodeAlreadyExists, params.CodeNotFound, params.CodeUnauthorized:
			return resizeOutcome{}, nil
		case params.CodeNotValid, params.CodeNotSupported:
			return resizeOutcome{report: params.StorageResizeReport{
				UUID: id, Status: "failed", Message: err.Error(),
			}}, nil
		default:
			return resizeOutcome{}, err
		}
	}
	operationCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	var retry bool
	revalidate := func(ctx context.Context) error {
		current, err := deps.config.Resizes.Prepare(ctx, id)
		if err != nil {
			switch params.ErrCode(err) {
			case params.CodeNotValid, params.CodeNotSupported, params.CodeAlreadyExists, params.CodeNotFound, params.CodeUnauthorized:
			default:
				retry = true
			}
			return err
		}
		if current != work {
			return errors.New("resize target changed during execution")
		}
		return nil
	}
	observed, err := deps.config.ResizeExecutor.Resize(operationCtx, storageresize.Request{
		ProviderType: work.ProviderType, VolumeID: work.VolumeID,
		VolumeProviderID: work.VolumeProviderID, DeviceName: work.DeviceName,
		FilesystemID: work.FilesystemID, MountPoint: work.MountPoint,
		SizeMiB: work.RequestedSizeMiB,
	}, revalidate)
	if ctx.Err() != nil {
		return resizeOutcome{}, ctx.Err()
	}
	if retry {
		// API unavailability is not a terminal storage failure. Repeating
		// completed growth stages is safe when the connection recovers.
		return resizeOutcome{}, err
	}
	report := params.StorageResizeReport{
		UUID: id, Status: "completed", Message: "Storage resize completed",
		VolumeSizeMiB:     observed.VolumeSizeMiB,
		FilesystemSizeMiB: observed.FilesystemSizeMiB,
	}
	if err != nil {
		report.Status = "failed"
		report.Message = "Storage resize failed; storage may have partially grown: " + err.Error()
	}
	return resizeOutcome{work: work, report: report}, nil
}

// updateResizeCache keeps the provisioner's local view in step with accepted
// completion. Watcher refreshes also read the persisted observed capacities.
func updateResizeCache(deps *dependencies, outcome resizeOutcome) {
	volumeTag := names.NewVolumeTag(outcome.work.VolumeID)
	if volume, ok := deps.volumes[volumeTag]; ok {
		volume.Size = max(volume.Size, outcome.report.VolumeSizeMiB)
		deps.volumes[volumeTag] = volume
	}
	if device, ok := deps.volumeBlockDevices[volumeTag]; ok {
		device.SizeMiB = max(device.SizeMiB, outcome.report.VolumeSizeMiB)
		deps.volumeBlockDevices[volumeTag] = device
	}
	if outcome.work.FilesystemID != "" {
		tag := names.NewFilesystemTag(outcome.work.FilesystemID)
		if filesystem, ok := deps.filesystems[tag]; ok {
			filesystem.Size = max(filesystem.Size, outcome.report.FilesystemSizeMiB)
			deps.filesystems[tag] = filesystem
		}
	}
}
