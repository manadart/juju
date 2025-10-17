// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"

	"github.com/canonical/sqlair"

	"github.com/juju/juju/core/storage"
	domainstorage "github.com/juju/juju/domain/storage"
	storageprovisioningerrors "github.com/juju/juju/domain/storageprovisioning/errors"
	"github.com/juju/juju/internal/errors"
)

// GetStorageAttachmentUUIDByInstanceID returns the unit storage attachment
// UUID for the storage instance with the input ID.
// At the time of writing a storage instance has at most one unit attachment.
func (st *State) GetStorageAttachmentUUIDByInstanceID(ctx context.Context, id string) (string, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return "", errors.Capture(err)
	}

	var attUUID entityUUID
	stID := storageID{ID: id}

	q := `
SELECT sa.uuid as &entityUUID.uuid
FROM storage_attachment sa
JOIN storage_instance si ON sa.storage_instance_uuid = si.uuid
WHERE si.storage_id = $storageID.storage_id`
	stmt, err := st.Prepare(q, attUUID, stID)
	if err != nil {
		return "", errors.Errorf("preparing storage attachment query: %w", err)
	}

	if err := db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		if err := tx.Query(ctx, stmt, stID).Get(&attUUID); err != nil {
			if errors.Is(err, sqlair.ErrNoRows) {
				return storageprovisioningerrors.StorageAttachmentNotFound
			}
			return errors.Errorf("running storage attachment query: %w", err)
		}
		return nil
	}); err != nil {
		return "", errors.Capture(err)
	}

	return attUUID.UUID, nil
}

func (st *State) GetModelDetails() (domainstorage.ModelDetails, error) {
	//TODO implement me
	return domainstorage.ModelDetails{}, errors.New("not implemented")
}

func (st *State) ImportFilesystem(
	ctx context.Context, name storage.Name, filesystem domainstorage.FilesystemInfo,
) (storage.ID, error) {
	//TODO implement me
	return "", errors.New("not implemented")
}
