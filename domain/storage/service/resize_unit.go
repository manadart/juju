// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"math"

	coreerrors "github.com/juju/juju/core/errors"
	corestorage "github.com/juju/juju/core/storage"
	coreunit "github.com/juju/juju/core/unit"
	domainstorage "github.com/juju/juju/domain/storage"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/domain/storage/internal"
	"github.com/juju/juju/internal/errors"
)

// RequestUnitStorageResize stages growth to the current application directives.
// Targets are fixed at submission. All requests are committed together, matching
// pending requests are reused, and sufficient observed capacity is a no-op.
func (s *Service) RequestUnitStorageResize(ctx context.Context, name coreunit.Name) ([]domainstorage.UnitStorageResizeResult, error) {
	if err := name.Validate(); err != nil {
		return nil, errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	for range 3 {
		targets, err := s.st.GetUnitStorageResizeTargets(ctx, name.String())
		if err != nil {
			return nil, errors.Errorf("getting storage for unit %q: %w", name, err)
		}
		results := make([]domainstorage.UnitStorageResizeResult, 0, len(targets))
		intents := make([]internal.StorageResizeIntent, 0, len(targets))
		for _, target := range targets {
			if !target.Alive {
				return nil, errors.Errorf("storage %q: %w", target.StorageID, storageerrors.StorageInstanceNotAlive)
			}
			if target.SizeMiB == 0 || target.SizeMiB > math.MaxInt64/(1024*1024) {
				return nil, errors.Errorf("storage %q has no valid application capacity target", target.StorageID).Add(coreerrors.NotValid)
			}
			result := domainstorage.UnitStorageResizeResult{StorageID: corestorage.ID(target.StorageID), RequestedSizeMiB: target.SizeMiB}
			switch {
			case target.PendingUUID != "":
				if target.PendingSizeMiB != target.SizeMiB {
					return nil, errors.Errorf("storage %q has a pending target of %d MiB, application target is %d MiB: %w", target.StorageID, target.PendingSizeMiB, target.SizeMiB, storageerrors.StorageResizeIntentAlreadyExists)
				}
				result.UUID = domainstorage.StorageResizeUUID(target.PendingUUID)
			case target.ObservedSizeMiB >= target.SizeMiB:
				result.AlreadySatisfied = true
			default:
				uuid, err := domainstorage.NewStorageResizeUUID()
				if err != nil {
					return nil, errors.Capture(err)
				}
				result.UUID = uuid
				intents = append(intents, internal.StorageResizeIntent{UUID: uuid.String(), StorageInstanceUUID: target.StorageInstanceUUID, RequestedSizeMiB: target.SizeMiB, CreatedAt: s.StorageService.clock.Now().UTC()})
			}
			results = append(results, result)
		}
		err = s.st.CreateUnitStorageResizeIntents(ctx, name.String(), targets, intents)
		if errors.Is(err, storageerrors.StorageResizeTargetChanged) {
			continue
		}
		if err != nil {
			return nil, errors.Errorf("requesting storage resize for unit %q: %w", name, err)
		}
		return results, nil
	}
	return nil, errors.Errorf("storage or application directives changed while staging unit %q; retry: %w", name, storageerrors.StorageResizeTargetChanged)
}
