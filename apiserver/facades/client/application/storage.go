// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"context"
	"math"

	"github.com/juju/errors"

	"github.com/juju/juju/core/storage"
	applicationservice "github.com/juju/juju/domain/application/service"
	"github.com/juju/juju/rpc/params"
)

// storageDirectives validates storage directives to override an applications
// storage directives.
func storageDirectives(
	ctx context.Context,
	storageService StorageService,
	storage map[string]storage.Directive,
) (map[string]applicationservice.StorageDirectiveOverrides, error) {
	res := map[string]applicationservice.StorageDirectiveOverrides{}
	for storageName, storageDirective := range storage {
		sdo := applicationservice.StorageDirectiveOverrides{}
		if storageDirective.Count != 0 {
			if storageDirective.Count > math.MaxUint32 {
				return nil, errors.NotValidf(
					"storage directive %s count too large", storageName,
				)
			}
			count := uint32(storageDirective.Count)
			sdo.Count = &count
		}
		if storageDirective.Size != 0 {
			sdo.Size = &storageDirective.Size
		}
		if storageDirective.Pool != "" {
			poolUUID, err := storageService.GetStoragePoolUUID(
				ctx,
				storageDirective.Pool,
			)
			if err != nil {
				return nil, errors.Annotatef(err, "storage directive %s pool",
					storageName)
			}
			sdo.PoolUUID = &poolUUID
		}
		res[storageName] = sdo
	}
	return res, nil
}

// applicationStorageDirectiveUpdates preserves explicitly supplied fields.
func applicationStorageDirectiveUpdates(ctx context.Context, storageService StorageService, directives map[string]params.StorageDirectives) (map[string]applicationservice.StorageDirectiveOverrides, error) {
	updates := make(map[string]applicationservice.StorageDirectiveOverrides, len(directives))
	for name, directive := range directives {
		update := applicationservice.StorageDirectiveOverrides{Size: directive.SizeMiB}
		if directive.Count != nil {
			if *directive.Count > math.MaxUint32 {
				return nil, errors.NotValidf("storage directive %q count too large", name)
			}
			count := uint32(*directive.Count)
			update.Count = &count
		}
		if directive.Pool != "" {
			uuid, err := storageService.GetStoragePoolUUID(ctx, directive.Pool)
			if err != nil {
				return nil, errors.Annotatef(err, "storage directive %q pool", name)
			}
			update.PoolUUID = &uuid
		}
		updates[name] = update
	}
	return updates, nil
}
