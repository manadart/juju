// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"slices"

	"github.com/canonical/sqlair"

	applicationerrors "github.com/juju/juju/domain/application/errors"
	"github.com/juju/juju/domain/life"
	sequencestate "github.com/juju/juju/domain/sequence/state"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/domain/storage/internal"
	"github.com/juju/juju/internal/errors"
)

// GetUnitStorageResizeTargets reads a live unit's attachments and current
// application targets. Missing directives remain visible as zero targets.
func (s *State) GetUnitStorageResizeTargets(ctx context.Context, name string) ([]internal.UnitStorageResizeTarget, error) {
	db, err := s.DB(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}
	var result []internal.UnitStorageResizeTarget
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		var err error
		result, err = s.getUnitStorageResizeTargets(ctx, tx, name)
		return err
	})
	return result, errors.Capture(err)
}

// CreateUnitStorageResizeIntents records a batch only if all inputs remain
// current. A conflict rolls back the entire batch, including earlier inserts.
func (s *State) CreateUnitStorageResizeIntents(ctx context.Context, name string, expected []internal.UnitStorageResizeTarget, intents []internal.StorageResizeIntent) (map[string]uint64, error) {
	db, err := s.DB(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}
	stmt, err := s.Prepare(`INSERT INTO storage_resize_request (*) VALUES ($storageResizeIntent.*)`, storageResizeIntent{})
	if err != nil {
		return nil, errors.Capture(err)
	}
	var ids map[string]uint64
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		ids = make(map[string]uint64, len(intents))
		current, err := s.getUnitStorageResizeTargets(ctx, tx, name)
		if err != nil {
			return err
		}
		if !slices.Equal(current, expected) {
			return storageerrors.StorageResizeTargetChanged
		}
		byUUID := make(map[string]internal.UnitStorageResizeTarget, len(current))
		for _, target := range current {
			byUUID[target.StorageInstanceUUID] = target
		}
		for _, intent := range intents {
			target, ok := byUUID[intent.StorageInstanceUUID]
			if !ok || !target.Alive || target.SizeMiB != intent.RequestedSizeMiB || target.PendingUUID != "" {
				return storageerrors.StorageResizeTargetChanged
			}
			row := storageResizeIntent{UUID: intent.UUID, StorageInstanceUUID: intent.StorageInstanceUUID, StorageID: target.StorageID, RequestedSizeMiB: intent.RequestedSizeMiB, CreatedAt: intent.CreatedAt}
			row.ID, err = sequencestate.NextValue(ctx, s, tx, storageResizeSequenceNamespace)
			if err != nil {
				return errors.Capture(err)
			}
			if err := tx.Query(ctx, stmt, row).Run(); err != nil {
				return errors.Capture(err)
			}
			ids[intent.UUID] = row.ID
		}
		return nil
	})
	if err != nil {
		return nil, errors.Capture(err)
	}
	return ids, nil
}

func (s *State) getUnitStorageResizeTargets(ctx context.Context, tx *sqlair.TX, name string) ([]internal.UnitStorageResizeTarget, error) {
	type unitKey struct {
		Name string `db:"name"`
	}
	type unitIdentity struct {
		UUID                 string    `db:"uuid"`
		UnitLife             life.Life `db:"unit_life"`
		ApplicationLife      life.Life `db:"application_life"`
		CharmUUID            string    `db:"charm_uuid"`
		ApplicationCharmUUID string    `db:"application_charm_uuid"`
	}
	key := unitKey{Name: name}
	stmt, err := s.Prepare(`
SELECT u.uuid AS &unitIdentity.uuid, u.life_id AS &unitIdentity.unit_life,
       a.life_id AS &unitIdentity.application_life, u.charm_uuid AS &unitIdentity.charm_uuid,
       a.charm_uuid AS &unitIdentity.application_charm_uuid
FROM unit AS u JOIN application AS a ON a.uuid = u.application_uuid
WHERE u.name = $unitKey.name`, key, unitIdentity{})
	if err != nil {
		return nil, err
	}
	var unit unitIdentity
	if err := tx.Query(ctx, stmt, key).Get(&unit); errors.Is(err, sqlair.ErrNoRows) {
		return nil, applicationerrors.UnitNotFound
	} else if err != nil {
		return nil, err
	}
	if unit.UnitLife != life.Alive {
		return nil, applicationerrors.UnitNotAlive
	}
	if unit.ApplicationLife != life.Alive {
		return nil, applicationerrors.ApplicationNotAlive
	}
	// Raw block capacity must include the machine observation: backend growth
	// alone does not establish completion of attachment work.
	stmt, err = s.Prepare(`
SELECT si.uuid AS &unitResizeTargetRow.instance_uuid,
       si.storage_id AS &unitResizeTargetRow.storage_id,
       sa.uuid AS &unitResizeTargetRow.attachment_uuid,
       COALESCE(d.size_mib, 0) AS &unitResizeTargetRow.size_mib,
       COALESCE(CASE WHEN si.storage_kind_id = 1 THEN f.size_mib
                     ELSE MIN(COALESCE(v.size_mib, 0), COALESCE(b.size_mib, 0)) END, 0) AS &unitResizeTargetRow.observed_size_mib,
       COALESCE(r.uuid, '') AS &unitResizeTargetRow.pending_uuid,
       COALESCE(r.request_id, 0) AS &unitResizeTargetRow.pending_id,
       COALESCE(r.requested_size_mib, 0) AS &unitResizeTargetRow.pending_size_mib,
       COALESCE(si.life_id = 0 AND sa.life_id = 0, 0) AS &unitResizeTargetRow.alive
FROM unit AS u
JOIN application AS a ON a.uuid = u.application_uuid
JOIN storage_attachment AS sa ON sa.unit_uuid = u.uuid
JOIN storage_instance AS si ON si.uuid = sa.storage_instance_uuid
LEFT JOIN application_storage_directive AS d ON d.application_uuid = a.uuid
    AND d.charm_uuid = a.charm_uuid AND d.storage_name = si.storage_name
LEFT JOIN storage_instance_filesystem AS ifs ON ifs.storage_instance_uuid = si.uuid
LEFT JOIN storage_filesystem AS f ON f.uuid = ifs.storage_filesystem_uuid
LEFT JOIN storage_instance_volume AS iv ON iv.storage_instance_uuid = si.uuid
LEFT JOIN storage_volume AS v ON v.uuid = iv.storage_volume_uuid
LEFT JOIN storage_volume_attachment AS va ON va.storage_volume_uuid = v.uuid AND va.net_node_uuid = u.net_node_uuid
LEFT JOIN block_device AS b ON b.uuid = va.block_device_uuid
LEFT JOIN storage_resize_request AS r ON r.storage_instance_uuid = si.uuid AND r.status = 'pending'
WHERE u.name = $unitKey.name
ORDER BY si.storage_id`, key, unitResizeTargetRow{})
	if err != nil {
		return nil, err
	}
	var rows []unitResizeTargetRow
	if err := tx.Query(ctx, stmt, key).GetAll(&rows); err != nil && !errors.Is(err, sqlair.ErrNoRows) {
		return nil, err
	}
	result := make([]internal.UnitStorageResizeTarget, 0, len(rows))
	for _, row := range rows {
		result = append(result, internal.UnitStorageResizeTarget{
			UnitUUID: unit.UUID, ApplicationCharmUUID: unit.ApplicationCharmUUID, UnitCharmUUID: unit.CharmUUID,
			StorageInstanceUUID: row.InstanceUUID, StorageID: row.StorageID, AttachmentUUID: row.AttachmentUUID,
			SizeMiB: row.SizeMiB, ObservedSizeMiB: row.ObservedSizeMiB, PendingUUID: row.PendingUUID, PendingID: row.PendingID, PendingSizeMiB: row.PendingSizeMiB, Alive: row.Alive,
		})
	}
	return result, nil
}

type unitResizeTargetRow struct {
	InstanceUUID    string `db:"instance_uuid"`
	StorageID       string `db:"storage_id"`
	AttachmentUUID  string `db:"attachment_uuid"`
	SizeMiB         uint64 `db:"size_mib"`
	ObservedSizeMiB uint64 `db:"observed_size_mib"`
	PendingUUID     string `db:"pending_uuid"`
	PendingID       uint64 `db:"pending_id"`
	PendingSizeMiB  uint64 `db:"pending_size_mib"`
	Alive           bool   `db:"alive"`
}
