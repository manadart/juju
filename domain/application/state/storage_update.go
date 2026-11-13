// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"slices"

	"github.com/canonical/sqlair"

	coreapplication "github.com/juju/juju/core/application"
	"github.com/juju/juju/domain/application/charm"
	applicationerrors "github.com/juju/juju/domain/application/errors"
	"github.com/juju/juju/domain/application/internal"
	domainstorage "github.com/juju/juju/domain/storage"
	"github.com/juju/juju/internal/errors"
)

// GetApplicationStorageSnapshot returns current directives and their constraints.
func (st *State) GetApplicationStorageSnapshot(ctx context.Context, uuid string) (internal.ApplicationStorageSnapshot, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return internal.ApplicationStorageSnapshot{}, errors.Capture(err)
	}
	var result internal.ApplicationStorageSnapshot
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		var err error
		result, err = st.getApplicationStorageSnapshot(ctx, tx, uuid)
		return err
	})
	return result, errors.Capture(err)
}

// SetApplicationStorageDirectives applies validated edits if their snapshot is
// still current. Existing unit directives and storage instances are untouched.
func (st *State) SetApplicationStorageDirectives(ctx context.Context, uuid string, expected internal.ApplicationStorageSnapshot, updates []domainstorage.DirectiveArg) error {
	db, err := st.DB(ctx)
	if err != nil {
		return errors.Capture(err)
	}
	return errors.Capture(db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		current, err := st.getApplicationStorageSnapshot(ctx, tx, uuid)
		if err != nil {
			return err
		}
		if current.CharmUUID != expected.CharmUUID || !slices.Equal(current.Directives, expected.Directives) {
			return applicationerrors.StorageDirectivesChanged
		}
		return st.updateApplicationStorageDirectives(ctx, tx, coreapplication.UUID(uuid), current.CharmUUID, updates)
	}))
}

// checkApplicationStorageDirectives prevents stale unit allocation and refresh.
func (st *State) checkApplicationStorageDirectives(ctx context.Context, tx *sqlair.TX, uuid string, expected []domainstorage.DirectiveArg) error {
	current, err := st.getApplicationStorageSnapshot(ctx, tx, uuid)
	if err != nil {
		return err
	}
	if len(expected) != len(current.Directives) {
		return applicationerrors.StorageDirectivesChanged
	}
	byName := make(map[domainstorage.Name]domainstorage.DirectiveArg, len(expected))
	for _, directive := range expected {
		byName[directive.Name] = directive
	}
	for _, directive := range current.Directives {
		if byName[directive.Name] != directive.DirectiveArg {
			return applicationerrors.StorageDirectivesChanged
		}
	}
	return nil
}

func (st *State) getApplicationStorageSnapshot(ctx context.Context, tx *sqlair.TX, uuid string) (internal.ApplicationStorageSnapshot, error) {
	var result internal.ApplicationStorageSnapshot
	if err := st.checkApplicationAlive(ctx, tx, uuid); err != nil {
		return result, err
	}
	charmUUID, err := st.getCharmIDByApplicationUUID(ctx, tx, uuid)
	if err != nil {
		return result, err
	}
	result.CharmUUID = charmUUID
	key := entityUUID{UUID: uuid}
	stmt, err := st.Prepare(`
SELECT d.storage_name AS &storageDirectiveSnapshotRow.storage_name,
       d.storage_pool_uuid AS &storageDirectiveSnapshotRow.storage_pool_uuid,
       d.size_mib AS &storageDirectiveSnapshotRow.size_mib,
       d.count AS &storageDirectiveSnapshotRow.count,
       COALESCE(sp.type, '') AS &storageDirectiveSnapshotRow.provider_type,
       k.kind AS &storageDirectiveSnapshotRow.kind,
       cs.count_min AS &storageDirectiveSnapshotRow.count_min,
       cs.count_max AS &storageDirectiveSnapshotRow.count_max,
       COALESCE(cs.minimum_size_mib, 0) AS &storageDirectiveSnapshotRow.minimum_size_mib,
       cs.shared AS &storageDirectiveSnapshotRow.shared
FROM application AS a
JOIN application_storage_directive AS d ON d.application_uuid = a.uuid AND d.charm_uuid = a.charm_uuid
JOIN charm_storage AS cs ON cs.charm_uuid = d.charm_uuid AND cs.name = d.storage_name
JOIN charm_storage_kind AS k ON k.id = cs.storage_kind_id
LEFT JOIN storage_pool AS sp ON sp.uuid = d.storage_pool_uuid
WHERE a.uuid = $entityUUID.uuid
ORDER BY d.storage_name
`, key, storageDirectiveSnapshotRow{})
	if err != nil {
		return result, err
	}
	var rows []storageDirectiveSnapshotRow
	if err := tx.Query(ctx, stmt, key).GetAll(&rows); err != nil && !errors.Is(err, sqlair.ErrNoRows) {
		return result, err
	}
	for _, row := range rows {
		result.Directives = append(result.Directives, internal.ApplicationStorageDirective{
			DirectiveArg: domainstorage.DirectiveArg{Name: domainstorage.Name(row.Name), PoolUUID: domainstorage.StoragePoolUUID(row.PoolUUID), Size: row.Size, Count: row.Count},
			Definition:   internal.CharmStorageDefinitionForValidation{Name: row.Name, Type: charm.StorageType(row.Kind), CountMin: row.CountMin, CountMax: row.CountMax, MinimumSize: row.MinimumSize, Shared: row.Shared},
			ProviderType: row.ProviderType,
		})
	}
	return result, nil
}

type storageDirectiveSnapshotRow struct {
	ProviderType string `db:"provider_type"`
	Name         string `db:"storage_name"`
	PoolUUID     string `db:"storage_pool_uuid"`
	Size         uint64 `db:"size_mib"`
	Count        uint32 `db:"count"`
	Kind         string `db:"kind"`
	CountMin     int    `db:"count_min"`
	CountMax     int    `db:"count_max"`
	MinimumSize  uint64 `db:"minimum_size_mib"`
	Shared       bool   `db:"shared"`
}
