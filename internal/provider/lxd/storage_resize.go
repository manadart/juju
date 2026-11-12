// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd

import (
	"context"
	"fmt"
	"maps"
	"math"
	"strings"

	lxdclient "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared"
	"github.com/canonical/lxd/shared/api"
	"github.com/canonical/lxd/shared/units"

	coreerrors "github.com/juju/juju/core/errors"
	"github.com/juju/juju/environs/tags"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/internal/storage"
)

var _ storage.FilesystemResizer = (*lxdFilesystemSource)(nil)

// CheckFilesystemResize checks online growth of a ZFS filesystem quota.
// LXD owns the filesystem, so no machine attachment work is required.
func (s *lxdFilesystemSource) CheckFilesystemResize(
	ctx context.Context, request storage.FilesystemResizeRequest,
) (storage.ResizeRequirements, error) {
	current, err := s.inspectFilesystemResize(ctx, request)
	if err != nil {
		return storage.ResizeRequirements{}, s.env.HandleCredentialError(ctx, err)
	}
	return storage.ResizeRequirements{
		AlreadySatisfied: current.sizeMiB >= request.SizeMiB,
		ObservedSizeMiB:  current.sizeMiB,
	}, nil
}

// EnsureFilesystemSize grows the existing dataset's quota in place. Reading
// the current quota makes retries safe after a lost response or worker restart;
// an already satisfied target never lowers a larger quota.
func (s *lxdFilesystemSource) EnsureFilesystemSize(
	ctx context.Context, request storage.FilesystemResizeRequest,
) (storage.ResizeProgress, error) {
	current, err := s.inspectFilesystemResize(ctx, request)
	if err != nil {
		return storage.ResizeProgress{}, s.env.HandleCredentialError(ctx, err)
	}
	progress := storage.ResizeProgress{
		Complete:        current.sizeMiB >= request.SizeMiB,
		ObservedSizeMiB: current.sizeMiB,
	}
	if progress.Complete {
		return progress, nil
	}

	update := current.volume.Writable()
	update.Config = maps.Clone(update.Config)
	update.Config["size"] = fmt.Sprintf("%dMiB", request.SizeMiB)
	// Keep all other settings and use the ETag to reject intervening changes.
	op, err := current.server.UpdateStoragePoolVolume(
		current.pool, storagePoolVolumeType, current.volume.Name,
		update, current.etag,
	)
	if err == nil {
		// Older LXD versions return a completed no-op operation after the
		// synchronous update. Newer versions return an asynchronous operation.
		// The worker bounds this call with its reconciliation context.
		err = op.WaitContext(ctx)
	}
	if err != nil {
		return progress, errors.Errorf("updating LXD filesystem quota: %w",
			s.env.HandleCredentialError(ctx, err))
	}

	observed, err := s.inspectFilesystemResize(ctx, request)
	if err != nil {
		return progress, s.env.HandleCredentialError(ctx, err)
	}
	if observed.volume.Config["volatile.uuid"] != current.volume.Config["volatile.uuid"] {
		return progress, errors.New("LXD volume identity changed during resize").Add(coreerrors.NotValid)
	}
	return storage.ResizeProgress{
		Complete:        observed.sizeMiB >= request.SizeMiB,
		ObservedSizeMiB: observed.sizeMiB,
	}, nil
}

type filesystemResizeInfo struct {
	server  lxdclient.InstanceServer
	pool    string
	volume  *api.StorageVolume
	etag    string
	sizeMiB uint64
}

// inspectFilesystemResize is read-only and is repeated immediately before
// mutation. A quota is provider capacity, not a reservation of free pool space.
func (s *lxdFilesystemSource) inspectFilesystemResize(
	ctx context.Context, request storage.FilesystemResizeRequest,
) (filesystemResizeInfo, error) {
	var result filesystemResizeInfo
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if request.OperationID == "" || request.SizeMiB == 0 || request.SizeMiB > math.MaxInt64/(1024*1024) {
		return result, errors.New("invalid LXD filesystem resize request").Add(coreerrors.NotValid)
	}
	poolName, volumeName, err := parseFilesystemId(request.FilesystemID)
	if err != nil || poolName == "" || volumeName == "" || strings.ContainsAny(poolName+volumeName, "/\\") {
		return result, errors.New("invalid LXD filesystem provider ID").Add(coreerrors.NotValid)
	}
	server := s.env.server()
	if server.IsClustered() {
		return result, errors.New("LXD filesystem resize currently requires a standalone server").Add(coreerrors.NotSupported)
	}
	result.server, err = server.WithContext(ctx)
	if err != nil {
		return result, err
	}
	result.pool = poolName
	pool, _, err := result.server.GetStoragePool(poolName)
	if err != nil {
		return result, errors.Errorf("reading LXD storage pool: %w", err)
	}
	if pool.Driver != "zfs" {
		return result, errors.New("LXD filesystem resize currently requires a ZFS pool").Add(coreerrors.NotSupported)
	}
	volume, etag, err := result.server.GetStoragePoolVolume(poolName, storagePoolVolumeType, volumeName)
	if err != nil {
		return result, errors.Errorf("reading LXD storage volume: %w", err)
	}
	if volume.Name != volumeName || volume.Type != storagePoolVolumeType || volume.ContentType != "filesystem" {
		return result, errors.New("LXD resize requires a custom filesystem volume").Add(coreerrors.NotSupported)
	}
	blockMode := volume.Config["zfs.block_mode"]
	if blockMode == "" {
		blockMode = pool.Config["volume.zfs.block_mode"]
	}
	if shared.IsTrue(blockMode) {
		return result, errors.New("LXD resize does not support ZFS block mode").Add(coreerrors.NotSupported)
	}
	if len(request.Attachments) > 1 || len(volume.UsedBy) > 1 {
		return result, errors.New("LXD resize does not support shared filesystems").Add(coreerrors.NotSupported)
	}
	if volume.Config["user."+tags.JujuModel] != s.env.uuid {
		return result, errors.New("LXD storage volume does not belong to this model").Add(coreerrors.NotValid)
	}
	if volume.Config["volatile.uuid"] == "" {
		return result, errors.New("LXD storage volume has no stable UUID").Add(coreerrors.NotSupported)
	}
	// An absent or zero quota is unlimited. Assigning a finite target would
	// restrict capacity, so it is not a supported growth operation.
	sizeBytes, err := units.ParseByteSizeString(volume.Config["size"])
	if err != nil {
		return result, errors.Errorf("reading LXD filesystem quota: %w", err)
	}
	if sizeBytes <= 0 {
		return result, errors.New("LXD resize requires an existing finite filesystem quota").Add(coreerrors.NotSupported)
	}
	result.volume = volume
	result.etag = etag
	result.sizeMiB = uint64(sizeBytes) / (1024 * 1024)
	return result, nil
}

// WithContext adapts the LXD client's context capability, which is not exposed
// by its InstanceServer interface. The clone retains the project's scope and
// shares the connection without changing the environment's client context.
func (s *providerServer) WithContext(ctx context.Context) (lxdclient.InstanceServer, error) {
	server, ok := s.InstanceServer.(interface {
		WithContext(context.Context) lxdclient.InstanceServer
	})
	if !ok {
		return nil, errors.New("LXD client does not support request contexts").Add(coreerrors.NotSupported)
	}
	return server.WithContext(ctx), nil
}
