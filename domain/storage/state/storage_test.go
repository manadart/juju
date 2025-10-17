// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	stdtesting "testing"

	"github.com/juju/tc"

	"github.com/juju/juju/core/network"
	"github.com/juju/juju/domain/schema/testing"
	storageprovisioningerrors "github.com/juju/juju/domain/storageprovisioning/errors"
)

type storageSuite struct {
	testing.ModelSuite
}

func TestStorageSuite(t *stdtesting.T) {
	tc.Run(t, &storageSuite{})
}

func (s *storageSuite) TestGetStorageAttachmentUUIDByInstanceIDNotFound(c *tc.C) {
	st := NewState(s.TxnRunnerFactory())

	_, err := st.GetStorageAttachmentUUIDByInstanceID(c.Context(), "pgdata/0")
	c.Check(err, tc.ErrorIs, storageprovisioningerrors.StorageAttachmentNotFound)
}

func (s *storageSuite) TestGetStorageAttachmentUUIDByInstanceID(c *tc.C) {
	st := NewState(s.TxnRunnerFactory())

	_, siID, saUUID := s.addAppUnitStorage(c)

	got, err := st.GetStorageAttachmentUUIDByInstanceID(c.Context(), siID)
	c.Assert(err, tc.ErrorIsNil)

	c.Check(got, tc.Equals, saUUID)
}

// addAppUnitStorage sets up a unit with a storage attachment.
// The storage instance and attachmentIDs are returned with the.
func (s *storageSuite) addAppUnitStorage(c *tc.C) (string, string, string) {
	ctx := c.Context()

	charm := "charm-uuid"
	_, err := s.DB().ExecContext(ctx, "INSERT INTO charm (uuid, reference_name, architecture_id) VALUES (?, ?, ?)",
		charm, charm, 0)
	c.Assert(err, tc.ErrorIsNil)

	app := "app-uuid"
	_, err = s.DB().ExecContext(
		ctx, "INSERT INTO application (uuid, name, life_id, charm_uuid, space_uuid) VALUES (?, ?, ?, ?, ?)",
		app, app, 0, charm, network.AlphaSpaceId,
	)
	c.Assert(err, tc.ErrorIsNil)

	node := "net-node-uuid"
	_, err = s.DB().ExecContext(ctx, "INSERT INTO net_node (uuid) VALUES (?)", node)
	c.Assert(err, tc.ErrorIsNil)

	unit := "unit-uuid"
	_, err = s.DB().ExecContext(
		ctx,
		"INSERT INTO unit (uuid, name, life_id, application_uuid, charm_uuid, net_node_uuid) VALUES (?, ?, ?, ?, ?, ?)",
		unit, unit, 0, app, charm, node)
	c.Assert(err, tc.ErrorIsNil)

	storagePool := "storage-pool-uuid"
	_, err = s.DB().ExecContext(ctx, "INSERT INTO storage_pool (uuid, name, type) VALUES (?, ?, ?)",
		storagePool, "loop", "loop")
	c.Assert(err, tc.ErrorIsNil)

	storageInstanceUUID := "storage-instance-uuid"
	storageInstanceID := charm + "/0"
	_, err = s.DB().Exec(`
INSERT INTO storage_instance (
    uuid, storage_id, storage_kind_id, storage_pool_uuid, requested_size_mib,
    charm_name, storage_name, life_id
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		storageInstanceUUID, storageInstanceID, 1, storagePool, 100, charm, "storage", 0)
	c.Assert(err, tc.ErrorIsNil)

	storageAttachment := "storage-attachment-uuid"
	_, err = s.DB().ExecContext(ctx,
		"INSERT INTO storage_attachment (uuid, storage_instance_uuid, unit_uuid, life_id) VALUES (?, ?, ?, ?)",
		storageAttachment, storageInstanceUUID, unit, 0)
	c.Assert(err, tc.ErrorIsNil)

	return storageInstanceUUID, storageInstanceID, storageAttachment
}
