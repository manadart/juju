// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"

	"github.com/canonical/sqlair"

	coreapplication "github.com/juju/juju/core/application"
	applicationerrors "github.com/juju/juju/domain/application/errors"
	"github.com/juju/juju/domain/storageprovisioning/internal"
	"github.com/juju/juju/internal/errors"
)

// GetK8sFilesystemAllocations returns live allocations, including filesystems
// whose PV and PVC identities have not yet been reported by Kubernetes.
func (st *State) GetK8sFilesystemAllocations(ctx context.Context, appUUID string) ([]internal.K8sFilesystemAllocation, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}
	key := entityUUID{UUID: appUUID}
	stmt, err := st.Prepare(`
SELECT si.storage_name AS &K8sFilesystemAllocation.storage_name,
       u.name AS &K8sFilesystemAllocation.unit_name,
       si.requested_size_mib AS &K8sFilesystemAllocation.size_mib,
       COALESCE(f.provider_id, '') AS &K8sFilesystemAllocation.volume_id,
       COALESCE(fa.provider_id, '') AS &K8sFilesystemAllocation.claim_name
FROM unit AS u
JOIN storage_attachment AS sa ON sa.unit_uuid = u.uuid AND sa.life_id = 0
JOIN storage_instance AS si ON si.uuid = sa.storage_instance_uuid AND si.life_id = 0
JOIN storage_instance_filesystem AS sif ON sif.storage_instance_uuid = si.uuid
JOIN storage_filesystem AS f ON f.uuid = sif.storage_filesystem_uuid AND f.life_id = 0
LEFT JOIN storage_filesystem_attachment AS fa
    ON fa.storage_filesystem_uuid = f.uuid AND fa.net_node_uuid = u.net_node_uuid AND fa.life_id = 0
WHERE u.application_uuid = $entityUUID.uuid AND u.life_id = 0
ORDER BY si.storage_name, u.name, si.storage_id
`, key, internal.K8sFilesystemAllocation{})
	if err != nil {
		return nil, errors.Capture(err)
	}
	var result []internal.K8sFilesystemAllocation
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		exists, err := st.checkApplicationExists(ctx, tx, coreapplication.UUID(appUUID))
		if err != nil {
			return err
		} else if !exists {
			return applicationerrors.ApplicationNotFound
		}
		if err := tx.Query(ctx, stmt, key).GetAll(&result); err != nil && !errors.Is(err, sqlair.ErrNoRows) {
			return err
		}
		return nil
	})
	return result, errors.Capture(err)
}
