// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storageprovisioner

import (
	"context"

	coreerrors "github.com/juju/juju/core/errors"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/internal/storage"
	storageresize "github.com/juju/juju/internal/storage/resize"
	"github.com/juju/juju/rpc/params"
)

type localResizeExecutor struct{ local *storageresize.Executor }

func (e *localResizeExecutor) Execute(ctx context.Context, work params.StorageResizeStageWork, revalidate func(context.Context) error) (params.StorageResizeStageReport, error) {
	t := work.Target
	request := storageresize.Request{ProviderType: t.ProviderType, VolumeID: t.VolumeID, VolumeProviderID: t.ProviderID, DeviceName: t.DeviceName, DeviceLink: t.DeviceLink, FilesystemID: t.FilesystemID, FilesystemUUID: t.FilesystemUUID, MountPoint: t.MountPoint, SizeMiB: t.SizeMiB}
	var observed storageresize.Observation
	var err error
	if work.Phase == "check" {
		if work.Stage == "backend" {
			observed, err = e.local.CheckBackend(ctx, request)
		} else {
			observed, err = e.local.CheckAttachment(ctx, request)
		}
	} else {
		if work.Stage == "backend" {
			observed, err = e.local.EnsureBackend(ctx, request, revalidate)
		} else {
			observed, err = e.local.EnsureAttachment(ctx, request, revalidate)
		}
	}
	report := params.StorageResizeStageReport{ObservedSizeMiB: observed.VolumeSizeMiB, FilesystemSizeMiB: observed.FilesystemSizeMiB}
	if err != nil {
		report.Status = "failed"
		report.Message = err.Error()
		return report, nil
	}
	if work.Phase == "check" {
		report.Status = "checked"
	} else if observed.VolumeSizeMiB >= t.SizeMiB && (work.Stage != "attachment" || t.FilesystemID == "" || observed.FilesystemSizeMiB >= t.SizeMiB) {
		report.Status = "completed"
	} else {
		report.Status = "pending"
		report.Message = "Waiting for local capacity to reach the approved target"
	}
	return report, nil
}

type providerResizeExecutor struct{ registry storage.ProviderRegistry }

func (e *providerResizeExecutor) Execute(ctx context.Context, work params.StorageResizeStageWork, revalidate func(context.Context) error) (params.StorageResizeStageReport, error) {
	var report params.StorageResizeStageReport
	t := work.Target
	if work.Stage != "backend" {
		return report, errors.New("provider executor received attachment work").Add(coreerrors.NotValid)
	}
	provider, err := e.registry.StorageProvider(storage.ProviderType(t.ProviderType))
	if err != nil {
		return report, err
	}
	if provider.Scope() != storage.ScopeEnviron {
		return report, errors.New("provider resize scope does not match the plan").Add(coreerrors.NotSupported)
	}
	attrs := make(map[string]any, len(t.PoolAttributes))
	for k, v := range t.PoolAttributes {
		attrs[k] = v
	}
	config, err := storage.NewConfig(t.PoolName, storage.ProviderType(t.ProviderType), attrs)
	if err != nil {
		return report, err
	}
	var attachments []storage.ResizeAttachment
	if t.AttachmentID != "" {
		attachments = []storage.ResizeAttachment{{AttachmentID: t.AttachmentID, InstanceID: t.InstanceID}}
	}
	var requirements storage.ResizeRequirements
	var progress storage.ResizeProgress
	// Resolving resize sources intentionally does not consult Dynamic(), which
	// describes creation after machine provisioning, not resize capability.
	switch t.Kind {
	case "volume":
		source, sourceErr := provider.VolumeSource(config)
		if sourceErr != nil {
			return report, sourceErr
		}
		resizer, ok := source.(storage.VolumeResizer)
		if !ok {
			return report, errors.New("volume source does not implement resizing").Add(coreerrors.NotSupported)
		}
		request := storage.VolumeResizeRequest{OperationID: work.UUID, VolumeID: t.ProviderID, SizeMiB: t.SizeMiB, Attachments: attachments}
		if work.Phase == "check" {
			requirements, err = resizer.CheckVolumeResize(ctx, request)
		} else {
			if err := revalidate(ctx); err != nil {
				return report, err
			}
			progress, err = resizer.EnsureVolumeSize(ctx, request)
		}
	case "filesystem":
		source, sourceErr := provider.FilesystemSource(config)
		if sourceErr != nil {
			return report, sourceErr
		}
		resizer, ok := source.(storage.FilesystemResizer)
		if !ok {
			return report, errors.New("filesystem source does not implement resizing").Add(coreerrors.NotSupported)
		}
		request := storage.FilesystemResizeRequest{OperationID: work.UUID, FilesystemID: t.ProviderID, SizeMiB: t.SizeMiB, Attachments: attachments}
		if work.Phase == "check" {
			requirements, err = resizer.CheckFilesystemResize(ctx, request)
		} else {
			if err := revalidate(ctx); err != nil {
				return report, err
			}
			progress, err = resizer.EnsureFilesystemSize(ctx, request)
		}
	default:
		return report, errors.New("unknown resize resource kind").Add(coreerrors.NotValid)
	}
	if work.Phase == "check" {
		report.Status = "checked"
		if requirements.Blocked {
			report.Status = "blocked"
		}
		report.Message = requirements.Message
		report.ObservedSizeMiB = requirements.ObservedSizeMiB
		for _, v := range requirements.Prerequisites {
			report.Prerequisites = append(report.Prerequisites, string(v))
		}
		for _, v := range requirements.FollowUp {
			report.FollowUp = append(report.FollowUp, string(v))
		}
		if err == nil && requirements.AlreadySatisfied && requirements.ObservedSizeMiB < t.SizeMiB {
			return report, errors.New("provider reported satisfaction below the target capacity").Add(coreerrors.NotValid)
		}
	} else {
		report.Status = "pending"
		report.Message = progress.Message
		report.ObservedSizeMiB = progress.ObservedSizeMiB
		if progress.Complete {
			if progress.ObservedSizeMiB < t.SizeMiB {
				return report, errors.New("provider reported completion below the target capacity").Add(coreerrors.NotValid)
			}
			report.Status = "completed"
		}
	}
	return report, err
}
