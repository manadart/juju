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

// applicationStorageDirectiveUpdates accepts size-only directive updates.
func applicationStorageDirectiveUpdates(directives map[string]params.StorageDirectives) (map[string]applicationservice.StorageDirectiveOverrides, error) {
	updates := make(map[string]applicationservice.StorageDirectiveOverrides, len(directives))
	for name, directive := range directives {
		if directive.Pool != "" || directive.Count != nil {
			return nil, errors.NotSupportedf("pool or count updates for storage %q; only size can be updated", name)
		}
		updates[name] = applicationservice.StorageDirectiveOverrides{Size: directive.SizeMiB}
	}
	return updates, nil
}
