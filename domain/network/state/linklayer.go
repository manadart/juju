// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"

	"github.com/canonical/sqlair"

	"github.com/juju/juju/domain/network"
	"github.com/juju/juju/internal/errors"
)

func (st *State) SetMachineNetConfig(
	ctx context.Context, nodeUUID string, nics []network.NetInterface, addrs []network.NetAddr,
) error {
	db, err := st.DB()
	if err != nil {
		return errors.Capture(err)
	}

	nUUID := entityUUID{UUID: nodeUUID}

	// TODO (manadart 2025-04-29): This is temporary and serves to get us
	// operational with addresses on Dqlite.
	// We will set devices and addresses for any given machine *one time*
	// with subsequent updates being a no-op until we add the full
	// reconciliation logic.
	var devCount countResult
	devCountSql := "SELECT COUNT(*) AS &countResult.* FROM link_layer_device WHERE net_node_uuid = $entityUUID.uuid"
	devCountStmt, err := st.Prepare(devCountSql, nUUID, devCount)
	if err != nil {
		return errors.Errorf("preparing device count statement: %w", err)
	}

	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		if err := tx.Query(ctx, devCountStmt, nUUID).Get(&devCount); err != nil {
			return errors.Errorf("running device count statement: %w", err)
		}

		// If we've done it for this machine before, bug out.
		if devCount.Count > 0 {
			return nil
		}

		// Otherwise, insert the data.

		return nil
	})

	return errors.Capture(err)
}

func (st *State) reconcileNetConfig(
	nodeUUID string, nics []network.NetInterface,
) ([]linkLayerDeviceDML, error) {
	// TODO (manadart 2025-04-30): This will have to return more types for DNS
	// nameservers/addresses, provider ID entries etc.

	// The idea here will be to set the ones that we know from querying existing
	// devices, then generate new ones for devices we don't have yet.
	nameToUUID := make(map[string]network.InterfaceUUID, len(nics))
	for _, n := range nics {
		nicUUID, err := network.NewInterfaceUUID()
		if err != nil {
			return nil, errors.Capture(err)
		}
		nameToUUID[n.Name] = nicUUID
	}

	nicsDML := make([]linkLayerDeviceDML, len(nics))
	for i, n := range nics {
		var err error
		if nicsDML[i], err = netInterfaceToDML(n, nodeUUID, nameToUUID); err != nil {
			return nil, errors.Capture(err)
		}
	}

	return nicsDML, nil
}

func (st *State) InsertLinkLayerDevices(ctx context.Context, tx sqlair.TX, devs []linkLayerDeviceDML) error {
	stmt, err := st.Prepare(
		"INSERT INTO link_layer_device (*) VALUES ($linkLayerDeviceDML.*)", devs[0])
	if err != nil {
		return errors.Errorf("preparing device insert statement: %w", err)
	}

	err = tx.Query(ctx, stmt, devs).Run()
	if err != nil {
		return errors.Errorf("running device insert statement: %w", err)
	}

	return nil
}
