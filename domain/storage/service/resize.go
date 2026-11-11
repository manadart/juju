// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"math"
	"strings"

	coreerrors "github.com/juju/juju/core/errors"
	corestorage "github.com/juju/juju/core/storage"
	"github.com/juju/juju/core/trace"
	domainstorage "github.com/juju/juju/domain/storage"
	"github.com/juju/juju/domain/storage/internal"
	"github.com/juju/juju/internal/errors"
)

// ResizeState describes persistence of storage resize requests and results.
type ResizeState interface {
	// CreateStorageResizeIntent records intent for a live storage instance.
	// It returns StorageInstanceNotFound, StorageInstanceNotAlive, or
	// StorageResizeIntentAlreadyExists if the request cannot be recorded.
	CreateStorageResizeIntent(context.Context, internal.StorageResizeIntent) error

	// GetStorageResizeIntent returns intent by UUID, or
	// StorageResizeIntentNotFound if it does not exist.
	GetStorageResizeIntent(context.Context, string) (internal.StorageResizeIntent, error)

	// GetStorageResizeRequest returns the original request and its result, or
	// StorageResizeIntentNotFound when the request does not exist.
	GetStorageResizeRequest(context.Context, string) (internal.StorageResizeRequest, error)

	// SetStorageResizeResult records the first terminal result. Identical
	// reports are idempotent; a conflicting result returns
	// StorageResizeResultAlreadyRecorded. Unknown requests return
	// StorageResizeIntentNotFound.
	SetStorageResizeResult(context.Context, string, internal.StorageResizeResult) error

	// GetStorageResizeRequestsForMachine returns attached or assigned requests.
	GetStorageResizeRequestsForMachine(context.Context, string) ([]string, error)

	// AssignStorageResize binds execution and revalidates resource identities.
	AssignStorageResize(context.Context, string, string) (internal.StorageResizeWork, error)

	// SetMachineStorageResizeResult validates ownership and records observed
	// capacities and the result atomically.
	SetMachineStorageResizeResult(context.Context, string, string, internal.StorageResizeReport) error
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

// GetStorageResizeIntent returns the original request. It returns NotValid for an
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

// GetStorageResizeRequest returns a request and its pending or terminal result.
// It returns NotValid for an invalid UUID, or StorageResizeIntentNotFound if
// the request does not exist. Removing storage does not remove its requests.
func (s *Service) GetStorageResizeRequest(
	ctx context.Context, uuid domainstorage.StorageResizeUUID,
) (domainstorage.StorageResizeRequest, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	if err := uuid.Validate(); err != nil {
		return domainstorage.StorageResizeRequest{}, errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	request, err := s.st.GetStorageResizeRequest(ctx, uuid.String())
	if err != nil {
		return domainstorage.StorageResizeRequest{}, errors.Errorf("getting storage resize request %q: %w", uuid, err)
	}
	return domainstorage.StorageResizeRequest{
		StorageResizeIntent: domainstorage.StorageResizeIntent{
			UUID:                domainstorage.StorageResizeUUID(request.UUID),
			StorageInstanceUUID: domainstorage.StorageInstanceUUID(request.StorageInstanceUUID),
			RequestedSizeMiB:    request.RequestedSizeMiB,
			CreatedAt:           request.CreatedAt,
		},
		StorageID:   corestorage.ID(request.StorageID),
		Status:      domainstorage.StorageResizeStatus(request.Status),
		Message:     request.Message,
		CompletedAt: request.CompletedAt,
	}, nil
}

// RecordStorageResizeResult records a completed or failed result and permits
// another resize request for the storage. A failed result requires a message.
// Retrying the same report preserves the original result and completion time.
//
// It returns NotValid for invalid input, StorageResizeIntentNotFound for an
// unknown request, or StorageResizeResultAlreadyRecorded for a conflicting
// terminal report. Recording a result does not update observed capacity.
func (s *Service) RecordStorageResizeResult(
	ctx context.Context, uuid domainstorage.StorageResizeUUID,
	result domainstorage.StorageResizeResult,
) error {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	if err := uuid.Validate(); err != nil {
		return errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	switch result.Status {
	case domainstorage.StorageResizeCompleted:
	case domainstorage.StorageResizeFailed:
		if strings.TrimSpace(result.Message) == "" {
			return errors.New("failed resize result requires a message").Add(coreerrors.NotValid)
		}
	default:
		return errors.Errorf("resize result status %q is not terminal", result.Status).Add(coreerrors.NotValid)
	}
	err := s.st.SetStorageResizeResult(ctx, uuid.String(), internal.StorageResizeResult{
		Status:      string(result.Status),
		Message:     result.Message,
		CompletedAt: s.StorageService.clock.Now().UTC(),
	})
	if err != nil {
		return errors.Errorf("recording result for storage resize request %q: %w", uuid, err)
	}
	return nil
}
