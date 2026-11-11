// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"time"

	"github.com/canonical/sqlair"

	"github.com/juju/juju/domain/life"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/domain/storage/internal"
	"github.com/juju/juju/internal/errors"
)

// CreateStorageResizeIntent records one pending intent for a live storage
// instance. It returns StorageInstanceNotFound, StorageInstanceNotAlive, or
// StorageResizeIntentAlreadyExists when the request cannot be recorded.
// Removing the storage instance also removes its pending intent.
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
		Life life.Life `db:"life_id"`
	}
	lifeStmt, err := s.Prepare(`
SELECT si.life_id AS &storageLife.life_id
FROM storage_instance AS si
WHERE si.uuid = $storageResizeIntent.storage_instance_uuid
`, row, storageLife{})
	if err != nil {
		return errors.Capture(err)
	}
	existingStmt, err := s.Prepare(`
SELECT sri.uuid AS &entityUUID.uuid
FROM storage_resize_intent AS sri
WHERE sri.storage_instance_uuid = $storageResizeIntent.storage_instance_uuid
`, row, entityUUID{})
	if err != nil {
		return errors.Capture(err)
	}
	insertStmt, err := s.Prepare(`
INSERT INTO storage_resize_intent (*) VALUES ($storageResizeIntent.*)
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
		var existing entityUUID
		err = tx.Query(ctx, existingStmt, row).Get(&existing)
		if err == nil {
			return storageerrors.StorageResizeIntentAlreadyExists
		} else if !errors.Is(err, sqlair.ErrNoRows) {
			return errors.Capture(err)
		}
		return errors.Capture(tx.Query(ctx, insertStmt, row).Run())
	}))
}

// GetStorageResizeIntent returns a pending intent by UUID. It returns
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
FROM storage_resize_intent AS sri
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
	StorageInstanceUUID string    `db:"storage_instance_uuid"`
	RequestedSizeMiB    uint64    `db:"requested_size_mib"`
	CreatedAt           time.Time `db:"created_at"`
}
