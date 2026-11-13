// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"database/sql"
	"time"

	"github.com/canonical/sqlair"

	"github.com/juju/juju/domain/life"
	"github.com/juju/juju/domain/sequence"
	sequencestate "github.com/juju/juju/domain/sequence/state"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/domain/storage/internal"
	"github.com/juju/juju/internal/errors"
)

const storageResizeSequenceNamespace = sequence.StaticNamespace("storage_resize")

// CreateStorageResizeIntent records one pending intent for a live storage
// instance. It returns StorageInstanceNotFound, StorageInstanceNotAlive, or
// StorageResizeIntentAlreadyExists when the request cannot be recorded.
// Completed requests remain available and do not prevent another request.
func (s *State) CreateStorageResizeIntent(
	ctx context.Context, intent internal.StorageResizeIntent,
) error {
	db, err := s.DB(ctx)
	if err != nil {
		return errors.Capture(err)
	}
	row := storageResizeIntent{
		UUID:                intent.UUID,
		StorageInstanceUUID: intent.StorageInstanceUUID,
		RequestedSizeMiB:    intent.RequestedSizeMiB,
		CreatedAt:           intent.CreatedAt,
	}
	type storageLife struct {
		Life      life.Life `db:"life_id"`
		StorageID string    `db:"storage_id"`
	}
	lifeStmt, err := s.Prepare(`
SELECT si.life_id AS &storageLife.life_id,
       si.storage_id AS &storageLife.storage_id
FROM storage_instance AS si
WHERE si.uuid = $storageResizeIntent.storage_instance_uuid
`, row, storageLife{})
	if err != nil {
		return errors.Capture(err)
	}
	existingStmt, err := s.Prepare(`
SELECT sri.uuid AS &entityUUID.uuid
FROM storage_resize_request AS sri
WHERE sri.storage_instance_uuid = $storageResizeIntent.storage_instance_uuid
AND sri.status = 'pending'
`, row, entityUUID{})
	if err != nil {
		return errors.Capture(err)
	}
	insertStmt, err := s.Prepare(`
INSERT INTO storage_resize_request (*) VALUES ($storageResizeIntent.*)
`, row)
	if err != nil {
		return errors.Capture(err)
	}
	return errors.Capture(db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		var current storageLife
		err := tx.Query(ctx, lifeStmt, row).Get(&current)
		if errors.Is(err, sqlair.ErrNoRows) {
			return storageerrors.StorageInstanceNotFound
		} else if err != nil {
			return errors.Capture(err)
		}
		if current.Life != life.Alive {
			return storageerrors.StorageInstanceNotAlive
		}
		row.StorageID = current.StorageID
		var existing entityUUID
		err = tx.Query(ctx, existingStmt, row).Get(&existing)
		if err == nil {
			return storageerrors.StorageResizeIntentAlreadyExists
		} else if !errors.Is(err, sqlair.ErrNoRows) {
			return errors.Capture(err)
		}
		row.ID, err = sequencestate.NextValue(ctx, s, tx, storageResizeSequenceNamespace)
		if err != nil {
			return errors.Capture(err)
		}
		return errors.Capture(tx.Query(ctx, insertStmt, row).Run())
	}))
}

// GetStorageResizeIntent returns the original intent by UUID. It returns
// StorageResizeIntentNotFound when the intent does not exist.
func (s *State) GetStorageResizeIntent(
	ctx context.Context, uuid string,
) (internal.StorageResizeIntent, error) {
	db, err := s.DB(ctx)
	if err != nil {
		return internal.StorageResizeIntent{}, errors.Capture(err)
	}
	row := storageResizeIntent{UUID: uuid}
	stmt, err := s.Prepare(`
SELECT sri.uuid AS &storageResizeIntent.uuid,
       sri.storage_instance_uuid AS &storageResizeIntent.storage_instance_uuid,
       sri.requested_size_mib AS &storageResizeIntent.requested_size_mib,
       sri.created_at AS &storageResizeIntent.created_at
FROM storage_resize_request AS sri
WHERE sri.uuid = $storageResizeIntent.uuid
`, row)
	if err != nil {
		return internal.StorageResizeIntent{}, errors.Capture(err)
	}
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, row).Get(&row)
		if errors.Is(err, sqlair.ErrNoRows) {
			return storageerrors.StorageResizeIntentNotFound
		}
		return errors.Capture(err)
	})
	if err != nil {
		return internal.StorageResizeIntent{}, errors.Capture(err)
	}
	return internal.StorageResizeIntent{
		UUID:                row.UUID,
		StorageInstanceUUID: row.StorageInstanceUUID,
		RequestedSizeMiB:    row.RequestedSizeMiB,
		CreatedAt:           row.CreatedAt,
	}, nil
}

type storageResizeIntent struct {
	UUID                string    `db:"uuid"`
	ID                  uint64    `db:"request_id"`
	StorageInstanceUUID string    `db:"storage_instance_uuid"`
	StorageID           string    `db:"storage_id"`
	RequestedSizeMiB    uint64    `db:"requested_size_mib"`
	CreatedAt           time.Time `db:"created_at"`
}

// GetStorageResizeRequest returns a request and its result, including requests
// whose storage has been removed. It returns StorageResizeIntentNotFound when
// the request does not exist.
func (s *State) GetStorageResizeRequest(
	ctx context.Context, uuid string,
) (internal.StorageResizeRequest, error) {
	db, err := s.DB(ctx)
	if err != nil {
		return internal.StorageResizeRequest{}, errors.Capture(err)
	}
	row := storageResizeRequest{UUID: uuid}
	stmt, err := s.Prepare(`
SELECT (srr.uuid, srr.request_id, srr.storage_instance_uuid, srr.storage_id,
        srr.requested_size_mib, srr.created_at, srr.status,
        srr.result_message, srr.completed_at) AS (&storageResizeRequest.*)
FROM storage_resize_request AS srr
WHERE srr.uuid = $storageResizeRequest.uuid
`, row)
	if err != nil {
		return internal.StorageResizeRequest{}, errors.Capture(err)
	}
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, row).Get(&row)
		if errors.Is(err, sqlair.ErrNoRows) {
			return storageerrors.StorageResizeIntentNotFound
		}
		return errors.Capture(err)
	})
	if err != nil {
		return internal.StorageResizeRequest{}, errors.Capture(err)
	}
	result := internal.StorageResizeRequest{
		StorageResizeIntent: internal.StorageResizeIntent{
			UUID:                row.UUID,
			StorageInstanceUUID: row.StorageInstanceUUID,
			RequestedSizeMiB:    row.RequestedSizeMiB,
			CreatedAt:           row.CreatedAt,
		},
		ID:        row.ID,
		StorageID: row.StorageID,
		Status:    row.Status,
		Message:   row.Message,
	}
	if row.CompletedAt.Valid {
		result.CompletedAt = &row.CompletedAt.Time
	}
	return result, nil
}

// GetStorageResizeUUIDByID resolves a user-facing request ID, including retained
// results whose storage has been removed.
func (s *State) GetStorageResizeUUIDByID(ctx context.Context, id uint64) (string, error) {
	db, err := s.DB(ctx)
	if err != nil {
		return "", errors.Capture(err)
	}
	row := storageResizeRequest{ID: id}
	stmt, err := s.Prepare(`
SELECT r.uuid AS &storageResizeRequest.uuid
FROM storage_resize_request AS r
WHERE r.request_id = $storageResizeRequest.request_id
`, row)
	if err != nil {
		return "", errors.Capture(err)
	}
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		if err := tx.Query(ctx, stmt, row).Get(&row); errors.Is(err, sqlair.ErrNoRows) {
			return storageerrors.StorageResizeIntentNotFound
		} else {
			return errors.Capture(err)
		}
	})
	return row.UUID, errors.Capture(err)
}

// SetStorageResizeResult records a terminal result and releases the pending
// request constraint atomically. Reporting the same status and message again
// is harmless and preserves the original completion time. A different result
// returns StorageResizeResultAlreadyRecorded. Unknown requests return
// StorageResizeIntentNotFound.
func (s *State) SetStorageResizeResult(
	ctx context.Context, uuid string, result internal.StorageResizeResult,
) error {
	db, err := s.DB(ctx)
	if err != nil {
		return errors.Capture(err)
	}
	row := storageResizeRequest{
		UUID:        uuid,
		Status:      result.Status,
		Message:     result.Message,
		CompletedAt: sql.NullTime{Time: result.CompletedAt, Valid: true},
	}
	getStmt, err := s.Prepare(`
SELECT srr.status AS &storageResizeRequest.status,
       srr.result_message AS &storageResizeRequest.result_message
FROM storage_resize_request AS srr
WHERE srr.uuid = $storageResizeRequest.uuid
`, row)
	if err != nil {
		return errors.Capture(err)
	}
	setStmt, err := s.Prepare(`
UPDATE storage_resize_request
SET status = $storageResizeRequest.status,
    result_message = $storageResizeRequest.result_message,
    completed_at = $storageResizeRequest.completed_at
WHERE uuid = $storageResizeRequest.uuid AND status = 'pending'
`, row)
	if err != nil {
		return errors.Capture(err)
	}
	return errors.Capture(db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		var current storageResizeRequest
		err := tx.Query(ctx, getStmt, row).Get(&current)
		if errors.Is(err, sqlair.ErrNoRows) {
			return storageerrors.StorageResizeIntentNotFound
		} else if err != nil {
			return errors.Capture(err)
		}
		if current.Status != "pending" {
			if current.Status == row.Status && current.Message == row.Message {
				return nil
			}
			return storageerrors.StorageResizeResultAlreadyRecorded
		}
		return errors.Capture(tx.Query(ctx, setStmt, row).Run())
	}))
}

type storageResizeRequest struct {
	UUID                string       `db:"uuid"`
	ID                  uint64       `db:"request_id"`
	StorageInstanceUUID string       `db:"storage_instance_uuid"`
	StorageID           string       `db:"storage_id"`
	RequestedSizeMiB    uint64       `db:"requested_size_mib"`
	CreatedAt           time.Time    `db:"created_at"`
	Status              string       `db:"status"`
	Message             string       `db:"result_message"`
	CompletedAt         sql.NullTime `db:"completed_at"`
}
