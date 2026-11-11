// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"math"

	coreerrors "github.com/juju/juju/core/errors"
	corestorage "github.com/juju/juju/core/storage"
	"github.com/juju/juju/core/trace"
	domainstorage "github.com/juju/juju/domain/storage"
	"github.com/juju/juju/domain/storage/internal"
	"github.com/juju/juju/internal/errors"
)

// ResizeState describes persistence of pending storage resize intent.
type ResizeState interface {
	// CreateStorageResizeIntent records intent for a live storage instance.
	// It returns StorageInstanceNotFound, StorageInstanceNotAlive, or
	// StorageResizeIntentAlreadyExists if the request cannot be recorded.
	CreateStorageResizeIntent(context.Context, internal.StorageResizeIntent) error

	// GetStorageResizeIntent returns intent by UUID, or
	// StorageResizeIntentNotFound if it does not exist.
	GetStorageResizeIntent(context.Context, string) (internal.StorageResizeIntent, error)
}

// RequestStorageResize records a target capacity in MiB for an existing storage
// instance and returns the request UUID. It does not check provider support,
// approve the request, or perform a resize. Those checks belong to execution.
//
// It returns NotValid for invalid inputs, StorageInstanceNotFound or
// StorageInstanceNotAlive for unavailable storage, and
// StorageResizeIntentAlreadyExists when a request is already pending.
func (s *Service) RequestStorageResize(
	ctx context.Context, storageID corestorage.ID, sizeMiB uint64,
) (domainstorage.StorageResizeUUID, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	if err := storageID.Validate(); err != nil {
		return "", errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	if sizeMiB == 0 || sizeMiB > math.MaxInt64 {
		return "", errors.New("requested size must be between 1 and 9223372036854775807 MiB").Add(coreerrors.NotValid)
	}
	storageUUID, err := s.st.GetStorageInstanceUUIDByID(ctx, storageID.String())
	if err != nil {
		return "", errors.Errorf("getting storage %q: %w", storageID, err)
	}
	uuid, err := domainstorage.NewStorageResizeUUID()
	if err != nil {
		return "", errors.Capture(err)
	}
	intent := internal.StorageResizeIntent{
		UUID:                uuid.String(),
		StorageInstanceUUID: storageUUID.String(),
		RequestedSizeMiB:    sizeMiB,
		CreatedAt:           s.StorageService.clock.Now().UTC(),
	}
	if err := s.st.CreateStorageResizeIntent(ctx, intent); err != nil {
		return "", errors.Errorf("requesting resize of storage %q: %w", storageID, err)
	}
	return uuid, nil
}

// GetStorageResizeIntent returns a pending request. It returns NotValid for an
// invalid UUID, or StorageResizeIntentNotFound if no such request exists.
func (s *Service) GetStorageResizeIntent(
	ctx context.Context, uuid domainstorage.StorageResizeUUID,
) (domainstorage.StorageResizeIntent, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	if err := uuid.Validate(); err != nil {
		return domainstorage.StorageResizeIntent{}, errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	intent, err := s.st.GetStorageResizeIntent(ctx, uuid.String())
	if err != nil {
		return domainstorage.StorageResizeIntent{}, errors.Errorf("getting storage resize intent %q: %w", uuid, err)
	}
	return domainstorage.StorageResizeIntent{
		UUID:                domainstorage.StorageResizeUUID(intent.UUID),
		StorageInstanceUUID: domainstorage.StorageInstanceUUID(intent.StorageInstanceUUID),
		RequestedSizeMiB:    intent.RequestedSizeMiB,
		CreatedAt:           intent.CreatedAt,
	}, nil
}
