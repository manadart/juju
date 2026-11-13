// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"

	coreapplication "github.com/juju/juju/core/application"
	coreerrors "github.com/juju/juju/core/errors"
	"github.com/juju/juju/domain/storageprovisioning"
	"github.com/juju/juju/internal/errors"
)

// GetK8sFilesystemAllocations returns existing allocation sizes and bindings.
// New ordinals without unit rows will use the application defaults instead.
func (s *Service) GetK8sFilesystemAllocations(ctx context.Context, appUUID coreapplication.UUID) ([]storageprovisioning.K8sFilesystemAllocation, error) {
	if err := appUUID.Validate(); err != nil {
		return nil, errors.Capture(err)
	}
	allocations, err := s.st.GetK8sFilesystemAllocations(ctx, appUUID.String())
	if err != nil {
		return nil, errors.Errorf("getting filesystem allocations for application %q: %w", appUUID, err)
	}
	result := make([]storageprovisioning.K8sFilesystemAllocation, 0, len(allocations))
	seen := make(map[[2]string]bool)
	for _, allocation := range allocations {
		key := [2]string{allocation.StorageName, allocation.UnitName}
		if seen[key] {
			return nil, errors.Errorf("multiple allocations of storage %q on unit %q are not supported", allocation.StorageName, allocation.UnitName).Add(coreerrors.NotSupported)
		}
		seen[key] = true
		result = append(result, storageprovisioning.K8sFilesystemAllocation{
			StorageName: allocation.StorageName, UnitName: allocation.UnitName,
			SizeMiB: allocation.SizeMiB, VolumeID: allocation.VolumeID,
			ClaimName: allocation.ClaimName,
		})
	}
	return result, nil
}
