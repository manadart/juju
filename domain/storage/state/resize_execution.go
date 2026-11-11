// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"encoding/json"
	"time"

	"github.com/canonical/sqlair"

	coreerrors "github.com/juju/juju/core/errors"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/domain/storage/internal"
	"github.com/juju/juju/internal/errors"
)

// GetStorageResizeRequestsForMachine returns pending requests attached to this
// machine, plus requests already assigned to it whose attachments have changed.
func (s *State) GetStorageResizeRequestsForMachine(ctx context.Context, machineUUID string) ([]string, error) {
	db, err := s.DB(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}
	key := resizeExecutionKey{MachineUUID: machineUUID}
	stmt, err := s.Prepare(`
SELECT DISTINCT r.uuid AS &entityUUID.uuid
FROM storage_resize_request AS r
LEFT JOIN storage_instance_volume AS iv ON iv.storage_instance_uuid = r.storage_instance_uuid
LEFT JOIN storage_volume_attachment AS va ON va.storage_volume_uuid = iv.storage_volume_uuid
LEFT JOIN storage_instance_filesystem AS ifs ON ifs.storage_instance_uuid = r.storage_instance_uuid
LEFT JOIN storage_filesystem_attachment AS fa ON fa.storage_filesystem_uuid = ifs.storage_filesystem_uuid
LEFT JOIN machine AS m ON m.net_node_uuid = COALESCE(va.net_node_uuid, fa.net_node_uuid)
WHERE r.status = 'pending'
AND (r.executor_machine_uuid = $resizeExecutionKey.machine_uuid
     OR (r.executor_machine_uuid IS NULL AND m.uuid = $resizeExecutionKey.machine_uuid))
ORDER BY r.uuid
`, key, entityUUID{})
	if err != nil {
		return nil, errors.Capture(err)
	}
	var rows []entityUUID
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt, key).GetAll(&rows)
	})
	if err != nil {
		return nil, errors.Capture(err)
	}
	result := make([]string, len(rows))
	for i, row := range rows {
		result[i] = row.UUID
	}
	return result, nil
}

// AssignStorageResize binds a request to its attached machine and identities.
// Retries revalidate that binding. Incomplete provisioning is left unassigned.
func (s *State) AssignStorageResize(ctx context.Context, uuid, machineUUID string) (internal.StorageResizeWork, error) {
	db, err := s.DB(ctx)
	if err != nil {
		return internal.StorageResizeWork{}, errors.Capture(err)
	}
	key := resizeExecutionKey{UUID: uuid, MachineUUID: machineUUID}
	stmt, err := s.Prepare(`
UPDATE storage_resize_request
SET executor_machine_uuid = $resizeExecutionKey.machine_uuid,
    execution_identity = $resizeExecutionKey.identity
WHERE uuid = $resizeExecutionKey.uuid
`, key)
	if err != nil {
		return internal.StorageResizeWork{}, errors.Capture(err)
	}
	var work resizeExecutionRow
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		record, err := s.getResizeExecutionRecord(ctx, tx, key)
		if err != nil {
			return err
		}
		if record.Executor != "" && record.Executor != machineUUID {
			return storageerrors.StorageResizeNotAssigned
		}
		if record.Status != "pending" {
			return storageerrors.StorageResizeResultAlreadyRecorded
		}
		work, err = s.getResizeExecutionWork(ctx, tx, key)
		if err != nil {
			return err
		}
		if work.MachineUUID != machineUUID {
			if record.Executor == machineUUID {
				return storageerrors.StorageResizeTargetChanged
			}
			return storageerrors.StorageResizeNotAssigned
		}
		key.Identity = work.identity()
		if record.Executor != "" {
			if record.Identity != key.Identity {
				return storageerrors.StorageResizeTargetChanged
			}
			return nil
		}
		if work.Alive && work.ProviderType == "loop" {
			if work.VolumeID != "" && (work.VolumeProviderID == "" || work.DeviceName == "") {
				return storageerrors.StorageResizeNotReady
			}
			if work.FilesystemID != "" && work.MountPoint == "" {
				return storageerrors.StorageResizeNotReady
			}
		}
		return tx.Query(ctx, stmt, key).Run()
	})
	if err != nil {
		return internal.StorageResizeWork{}, errors.Capture(err)
	}
	return internal.StorageResizeWork{
		UUID: uuid, RequestedSizeMiB: work.RequestedSizeMiB,
		ProviderType: work.ProviderType,
		VolumeID:     work.VolumeID, VolumeProviderID: work.VolumeProviderID,
		DeviceName: work.DeviceName, FilesystemID: work.FilesystemID,
		MountPoint: work.MountPoint, VolumeSizeMiB: work.VolumeSizeMiB,
		Alive: work.Alive, ReadOnly: work.ReadOnly, Shared: work.Shared,
		MachineScoped: work.MachineScoped,
	}, nil
}

// SetMachineStorageResizeResult atomically records observations and the result
// for the assigned machine. A successful result requires unchanged identities,
// live resources and enough observed capacity at every required layer.
func (s *State) SetMachineStorageResizeResult(ctx context.Context, uuid, machineUUID string, report internal.StorageResizeReport) error {
	db, err := s.DB(ctx)
	if err != nil {
		return errors.Capture(err)
	}
	key := resizeExecutionKey{UUID: uuid, MachineUUID: machineUUID}
	result := resizeExecutionResult{
		UUID: uuid, Status: report.Status, Message: report.Message,
		CompletedAt: report.CompletedAt, VolumeSizeMiB: report.VolumeSizeMiB,
		FilesystemSizeMiB: report.FilesystemSizeMiB,
	}
	setStmt, err := s.Prepare(`
UPDATE storage_resize_request
SET status = $resizeExecutionResult.status,
    result_message = $resizeExecutionResult.result_message,
    completed_at = $resizeExecutionResult.completed_at
WHERE uuid = $resizeExecutionResult.uuid
`, result)
	if err != nil {
		return errors.Capture(err)
	}
	return errors.Capture(db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		record, err := s.getResizeExecutionRecord(ctx, tx, key)
		if err != nil {
			return err
		}
		if record.Executor != machineUUID {
			return storageerrors.StorageResizeNotAssigned
		}
		if record.Status != "pending" {
			if record.Status == report.Status && record.Message == report.Message {
				return nil
			}
			return storageerrors.StorageResizeResultAlreadyRecorded
		}
		work, workErr := s.getResizeExecutionWork(ctx, tx, key)
		if workErr != nil && !errors.Is(workErr, storageerrors.StorageResizeTargetChanged) {
			return workErr
		}
		unchanged := workErr == nil && work.MachineUUID == machineUUID && record.Identity == work.identity() && work.Alive && !work.ReadOnly
		if report.Status == "completed" {
			if !unchanged {
				return storageerrors.StorageResizeTargetChanged
			}
			if report.VolumeSizeMiB < work.RequestedSizeMiB || (work.FilesystemID != "" && report.FilesystemSizeMiB < work.RequestedSizeMiB) {
				return errors.New("reported resize capacity is below the requested target").Add(coreerrors.NotValid)
			}
		}
		if unchanged {
			if err := s.updateResizeObservations(ctx, tx, work, result); err != nil {
				return err
			}
		}
		return tx.Query(ctx, setStmt, result).Run()
	}))
}

func (s *State) updateResizeObservations(ctx context.Context, tx *sqlair.TX, work resizeExecutionRow, result resizeExecutionResult) error {
	volumeStmt, err := s.Prepare(`
UPDATE storage_volume SET size_mib = MAX(COALESCE(size_mib, 0), $resizeExecutionResult.volume_size_mib)
WHERE uuid = $resizeExecutionRow.volume_uuid
`, result, work)
	if err != nil {
		return errors.Capture(err)
	}
	fsStmt, err := s.Prepare(`
UPDATE storage_filesystem SET size_mib = MAX(COALESCE(size_mib, 0), $resizeExecutionResult.filesystem_size_mib)
WHERE uuid = $resizeExecutionRow.filesystem_uuid
`, result, work)
	if err != nil {
		return errors.Capture(err)
	}
	if result.VolumeSizeMiB != 0 {
		if err := tx.Query(ctx, volumeStmt, result, work).Run(); err != nil {
			return err
		}
	}
	if result.FilesystemSizeMiB != 0 {
		if err := tx.Query(ctx, fsStmt, result, work).Run(); err != nil {
			return err
		}
	}
	if result.Status != "completed" {
		return nil
	}
	stmt, err := s.Prepare(`
UPDATE storage_instance SET requested_size_mib = MAX(requested_size_mib, $resizeExecutionRow.requested_size_mib)
WHERE uuid = $resizeExecutionRow.instance_uuid
`, work)
	if err != nil {
		return errors.Capture(err)
	}
	return tx.Query(ctx, stmt, work).Run()
}

func (s *State) getResizeExecutionRecord(ctx context.Context, tx *sqlair.TX, key resizeExecutionKey) (resizeExecutionRecord, error) {
	stmt, err := s.Prepare(`
SELECT COALESCE(r.executor_machine_uuid, '') AS &resizeExecutionRecord.executor,
       COALESCE(r.execution_identity, '') AS &resizeExecutionRecord.identity,
       r.status AS &resizeExecutionRecord.status,
       r.result_message AS &resizeExecutionRecord.result_message
FROM storage_resize_request AS r WHERE r.uuid = $resizeExecutionKey.uuid
`, key, resizeExecutionRecord{})
	if err != nil {
		return resizeExecutionRecord{}, errors.Capture(err)
	}
	var result resizeExecutionRecord
	err = tx.Query(ctx, stmt, key).Get(&result)
	if errors.Is(err, sqlair.ErrNoRows) {
		return result, storageerrors.StorageResizeIntentNotFound
	}
	return result, errors.Capture(err)
}

func (s *State) getResizeExecutionWork(ctx context.Context, tx *sqlair.TX, key resizeExecutionKey) (resizeExecutionRow, error) {
	stmt, err := s.Prepare(`
WITH fs_counts AS (
    SELECT storage_filesystem_uuid AS filesystem_uuid, COUNT(*) AS count
    FROM storage_filesystem_attachment AS fa GROUP BY storage_filesystem_uuid
)
SELECT COALESCE(m.uuid, '') AS &resizeExecutionRow.machine_uuid,
       r.storage_instance_uuid AS &resizeExecutionRow.instance_uuid,
       r.requested_size_mib AS &resizeExecutionRow.requested_size_mib,
       COALESCE(p.uuid, '') AS &resizeExecutionRow.pool_uuid,
       COALESCE(p.type, '') AS &resizeExecutionRow.provider_type,
       COALESCE(v.uuid, '') AS &resizeExecutionRow.volume_uuid,
       COALESCE(v.volume_id, '') AS &resizeExecutionRow.volume_id,
       COALESCE(v.provider_id, '') AS &resizeExecutionRow.volume_provider_id,
       COALESCE(v.size_mib, 0) AS &resizeExecutionRow.volume_size_mib,
       COALESCE(va.uuid, '') AS &resizeExecutionRow.volume_attachment_uuid,
       COALESCE(b.uuid, '') AS &resizeExecutionRow.block_device_uuid,
       COALESCE(b.name, '') AS &resizeExecutionRow.device_name,
       COALESCE(f.uuid, '') AS &resizeExecutionRow.filesystem_uuid,
       COALESCE(f.filesystem_id, '') AS &resizeExecutionRow.filesystem_id,
       COALESCE(fa.uuid, '') AS &resizeExecutionRow.filesystem_attachment_uuid,
       COALESCE(fa.mount_point, '') AS &resizeExecutionRow.mount_point,
       COALESCE(va.read_only, 0) OR COALESCE(fa.read_only, 0) AS &resizeExecutionRow.read_only,
       COALESCE(fc.count, 0) > 1 AS &resizeExecutionRow.shared,
       COALESCE(v.provision_scope_id = 1, 0) AS &resizeExecutionRow.machine_scoped,
       COALESCE(si.life_id = 0 AND m.life_id = 0
           AND (v.uuid IS NULL OR v.life_id = 0)
           AND (va.uuid IS NULL OR va.life_id = 0)
           AND (f.uuid IS NULL OR f.life_id = 0)
           AND (fa.uuid IS NULL OR fa.life_id = 0), 0) AS &resizeExecutionRow.alive
FROM storage_resize_request AS r
LEFT JOIN storage_instance AS si ON si.uuid = r.storage_instance_uuid
LEFT JOIN storage_pool AS p ON p.uuid = si.storage_pool_uuid
LEFT JOIN storage_instance_volume AS iv ON iv.storage_instance_uuid = si.uuid
LEFT JOIN storage_volume AS v ON v.uuid = iv.storage_volume_uuid
LEFT JOIN storage_volume_attachment AS va ON va.storage_volume_uuid = v.uuid
LEFT JOIN block_device AS b ON b.uuid = va.block_device_uuid
LEFT JOIN storage_instance_filesystem AS ifs ON ifs.storage_instance_uuid = si.uuid
LEFT JOIN storage_filesystem AS f ON f.uuid = ifs.storage_filesystem_uuid
LEFT JOIN storage_filesystem_attachment AS fa ON fa.storage_filesystem_uuid = f.uuid
    AND (va.net_node_uuid IS NULL OR fa.net_node_uuid = va.net_node_uuid)
LEFT JOIN fs_counts AS fc ON fc.filesystem_uuid = f.uuid
LEFT JOIN machine AS m ON m.net_node_uuid = COALESCE(va.net_node_uuid, fa.net_node_uuid)
WHERE r.uuid = $resizeExecutionKey.uuid
AND (m.uuid = $resizeExecutionKey.machine_uuid OR m.uuid IS NULL)
`, key, resizeExecutionRow{})
	if err != nil {
		return resizeExecutionRow{}, errors.Capture(err)
	}
	var row resizeExecutionRow
	err = tx.Query(ctx, stmt, key).Get(&row)
	if errors.Is(err, sqlair.ErrNoRows) {
		return row, storageerrors.StorageResizeTargetChanged
	}
	return row, errors.Capture(err)
}

type resizeExecutionKey struct {
	UUID        string `db:"uuid"`
	MachineUUID string `db:"machine_uuid"`
	Identity    string `db:"identity"`
}

type resizeExecutionRecord struct {
	Executor string `db:"executor"`
	Identity string `db:"identity"`
	Status   string `db:"status"`
	Message  string `db:"result_message"`
}

type resizeExecutionResult struct {
	UUID              string    `db:"uuid"`
	Status            string    `db:"status"`
	Message           string    `db:"result_message"`
	CompletedAt       time.Time `db:"completed_at"`
	VolumeSizeMiB     uint64    `db:"volume_size_mib"`
	FilesystemSizeMiB uint64    `db:"filesystem_size_mib"`
}

type resizeExecutionRow struct {
	MachineUUID              string `db:"machine_uuid"`
	InstanceUUID             string `db:"instance_uuid"`
	RequestedSizeMiB         uint64 `db:"requested_size_mib"`
	PoolUUID                 string `db:"pool_uuid"`
	ProviderType             string `db:"provider_type"`
	VolumeUUID               string `db:"volume_uuid"`
	VolumeID                 string `db:"volume_id"`
	VolumeProviderID         string `db:"volume_provider_id"`
	VolumeSizeMiB            uint64 `db:"volume_size_mib"`
	VolumeAttachmentUUID     string `db:"volume_attachment_uuid"`
	BlockDeviceUUID          string `db:"block_device_uuid"`
	DeviceName               string `db:"device_name"`
	FilesystemUUID           string `db:"filesystem_uuid"`
	FilesystemID             string `db:"filesystem_id"`
	FilesystemAttachmentUUID string `db:"filesystem_attachment_uuid"`
	MountPoint               string `db:"mount_point"`
	ReadOnly                 bool   `db:"read_only"`
	Shared                   bool   `db:"shared"`
	MachineScoped            bool   `db:"machine_scoped"`
	Alive                    bool   `db:"alive"`
}

func (r resizeExecutionRow) identity() string {
	data, _ := json.Marshal([]any{r.MachineUUID, r.InstanceUUID, r.PoolUUID,
		r.ProviderType, r.VolumeUUID, r.VolumeProviderID,
		r.VolumeAttachmentUUID, r.BlockDeviceUUID, r.DeviceName,
		r.FilesystemUUID, r.FilesystemAttachmentUUID, r.MountPoint, r.ReadOnly,
		r.Shared, r.MachineScoped})
	return string(data)
}
