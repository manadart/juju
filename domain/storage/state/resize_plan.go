// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/canonical/sqlair"

	domainstorage "github.com/juju/juju/domain/storage"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/internal/errors"
)

// GetStorageResizeStageRequests returns model planning/backend work or the
// stages assigned to a machine. Machines may also plan their local loop work.
func (s *State) GetStorageResizeStageRequests(ctx context.Context, machineUUID string) ([]string, error) {
	db, err := s.DB(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}
	key := resizeExecutionKey{MachineUUID: machineUUID}
	stmt, err := s.Prepare(`
SELECT r.uuid AS &entityUUID.uuid
FROM storage_resize_request AS r
LEFT JOIN storage_resize_plan AS p ON p.request_uuid = r.uuid
WHERE r.status = 'pending'
AND (($resizeExecutionKey.machine_uuid = '' AND
      (p.request_uuid IS NULL OR p.backend_machine_uuid = ''))
 OR ($resizeExecutionKey.machine_uuid <> '' AND
      (p.backend_machine_uuid = $resizeExecutionKey.machine_uuid
       OR p.node_machine_uuid = $resizeExecutionKey.machine_uuid)))
ORDER BY r.created_at, r.uuid`, key, entityUUID{})
	if err != nil {
		return nil, errors.Capture(err)
	}
	var rows []entityUUID
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, key).GetAll(&rows)
		if errors.Is(err, sqlair.ErrNoRows) {
			return nil
		}
		return err
	})
	if err != nil {
		return nil, errors.Capture(err)
	}
	ids := make([]string, len(rows))
	for i, row := range rows {
		ids[i] = row.UUID
	}
	if machineUUID != "" {
		local, err := s.GetStorageResizeRequestsForMachine(ctx, machineUUID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, local...)
	}
	return ids, nil
}

// GetStorageResizePlan returns durable progress together with current inputs.
// Identity is empty until a plan is committed; Target.Identity is always fresh.
func (s *State) GetStorageResizePlan(ctx context.Context, uuid string) (domainstorage.ResizePlan, error) {
	db, err := s.DB(ctx)
	if err != nil {
		return domainstorage.ResizePlan{}, errors.Capture(err)
	}
	var plan domainstorage.ResizePlan
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		var err error
		plan, _, err = s.getResizePlan(ctx, tx, uuid)
		return err
	})
	return plan, errors.Capture(err)
}

// CreateStorageResizePlan commits checked resource bindings and initial stages.
// It revalidates the snapshot and preserves any original legacy binding.
func (s *State) CreateStorageResizePlan(ctx context.Context, plan domainstorage.ResizePlan) error {
	db, err := s.DB(ctx)
	if err != nil {
		return errors.Capture(err)
	}
	return errors.Capture(db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		current, work, err := s.getResizePlan(ctx, tx, plan.UUID)
		if err != nil {
			return err
		}
		if current.Identity != "" {
			return nil
		}
		if current.Target.Identity != plan.Identity {
			return storageerrors.StorageResizeTargetChanged
		}
		legacy, err := s.getResizeExecutionRecord(ctx, tx, resizeExecutionKey{UUID: plan.UUID})
		if err != nil {
			return err
		}
		if legacy.Executor != "" && legacy.Identity != work.identity() {
			return storageerrors.StorageResizeTargetChanged
		}
		row := resizePlanRow{UUID: plan.UUID, Identity: plan.Identity, BackendMachineUUID: plan.BackendMachineUUID, NodeMachineUUID: plan.NodeMachineUUID}
		stmt, err := s.Prepare(`INSERT INTO storage_resize_plan (*) VALUES ($resizePlanRow.*)`, row)
		if err != nil {
			return err
		}
		if err := tx.Query(ctx, stmt, row).Run(); err != nil {
			return err
		}
		for _, stage := range []string{"backend", "attachment"} {
			state := "checking"
			if stage == "attachment" && plan.NodeMachineUUID == "" {
				state = "skipped"
			}
			row := resizeStageRow{UUID: plan.UUID, Stage: stage, State: state, Prerequisites: "[]", FollowUp: "[]"}
			stmt, err := s.Prepare(`INSERT INTO storage_resize_stage (*) VALUES ($resizeStageRow.*)`, row)
			if err != nil {
				return err
			}
			if err := tx.Query(ctx, stmt, row).Run(); err != nil {
				return err
			}
		}
		return nil
	}))
}

// SetStorageResizeStage atomically persists a stage and advances approval or
// completion. Reports cannot change identity, skip stages or regress progress.
func (s *State) SetStorageResizeStage(ctx context.Context, machineUUID string, report domainstorage.ResizeStageReport, now time.Time) error {
	db, err := s.DB(ctx)
	if err != nil {
		return errors.Capture(err)
	}
	return errors.Capture(db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		plan, work, err := s.getResizePlan(ctx, tx, report.UUID)
		if err != nil {
			return err
		}
		if plan.Identity == "" {
			return storageerrors.StorageResizeNotReady
		}
		owner := plan.BackendMachineUUID
		current := plan.Backend
		if report.Stage == "attachment" {
			owner = plan.NodeMachineUUID
			current = plan.Attachment
			if owner == "" {
				return storageerrors.StorageResizeNotAssigned
			}
		}
		if owner != machineUUID {
			return storageerrors.StorageResizeNotAssigned
		}
		if report.Identity != plan.Identity {
			return storageerrors.StorageResizeTargetChanged
		}
		unchanged := plan.Identity == plan.Target.Identity && plan.Target.Alive
		if !unchanged && report.Status != "failed" {
			return storageerrors.StorageResizeTargetChanged
		}
		// A lost response may be retried after another stage advanced the request.
		if current.State == "completed" || (report.Phase == "check" && current.State != "checking") {
			return nil
		}
		stage, phase := plan.Next(machineUUID)
		if stage != report.Stage || phase != report.Phase {
			return storageerrors.StorageResizeNotReady
		}
		state := current.State
		switch report.Status {
		case "blocked":
			state = "checking"
		case "checked":
			state = "checked"
		case "pending":
			state = "running"
		case "completed":
			state = "completed"
		case "failed":
			state = "failed"
		}
		if report.Phase == "ensure" {
			report.Prerequisites = current.Prerequisites
			report.FollowUp = current.FollowUp
		}
		prereq, _ := json.Marshal(report.Prerequisites)
		follow, _ := json.Marshal(report.FollowUp)
		row := resizeStageRow{UUID: report.UUID, Stage: report.Stage, State: state, Message: report.Message, Prerequisites: string(prereq), FollowUp: string(follow), ObservedSizeMiB: max(current.ObservedSizeMiB, report.ObservedSizeMiB), FilesystemSizeMiB: max(current.FilesystemSizeMiB, report.FilesystemSizeMiB)}
		if report.Status == "completed" && (report.ObservedSizeMiB < plan.Target.SizeMiB || (report.Stage == "attachment" && plan.Target.FilesystemID != "" && report.FilesystemSizeMiB < plan.Target.SizeMiB)) {
			return errors.New("observed resize capacity is below target")
		}
		stmt, err := s.Prepare(`UPDATE storage_resize_stage SET
state=$resizeStageRow.state, message=$resizeStageRow.message,
prerequisites=$resizeStageRow.prerequisites, follow_up=$resizeStageRow.follow_up,
observed_size_mib=$resizeStageRow.observed_size_mib,
filesystem_size_mib=$resizeStageRow.filesystem_size_mib
WHERE request_uuid=$resizeStageRow.request_uuid AND stage=$resizeStageRow.stage`, row)
		if err != nil {
			return err
		}
		if err := tx.Query(ctx, stmt, row).Run(); err != nil {
			return err
		}
		if unchanged {
			if err := s.setResizeStageObservations(ctx, tx, work, report); err != nil {
				return err
			}
		}
		if report.Status == "failed" {
			return s.setResizePlanResult(ctx, tx, report.UUID, "failed", report.Message, now)
		}
		updated := domainstorage.ResizeStage{State: state, ObservedSizeMiB: row.ObservedSizeMiB, FilesystemSizeMiB: row.FilesystemSizeMiB}
		if report.Stage == "backend" {
			plan.Backend = updated
		} else {
			plan.Attachment = updated
		}
		if plan.ApprovedSizeMiB == 0 && plan.Backend.State == "checked" && (plan.Attachment.State == "checked" || plan.Attachment.State == "skipped") {
			row := resizePlanRow{UUID: report.UUID, ApprovedSizeMiB: plan.Target.SizeMiB}
			stmt, err := s.Prepare(`UPDATE storage_resize_plan SET approved_size_mib=$resizePlanRow.approved_size_mib WHERE request_uuid=$resizePlanRow.request_uuid`, row)
			if err != nil {
				return err
			}
			if err := tx.Query(ctx, stmt, row).Run(); err != nil {
				return err
			}
		}
		if plan.Backend.State == "completed" && (plan.Attachment.State == "completed" || plan.Attachment.State == "skipped") {
			if err := s.setResizePlanResult(ctx, tx, report.UUID, "completed", "Storage resize completed", now); err != nil {
				return err
			}
			stmt, err := s.Prepare(`UPDATE storage_instance SET requested_size_mib=MAX(requested_size_mib,$resizeExecutionRow.requested_size_mib) WHERE uuid=$resizeExecutionRow.instance_uuid`, work)
			if err != nil {
				return err
			}
			return tx.Query(ctx, stmt, work).Run()
		}
		return nil
	}))
}

// FailStorageResizePlan retains a terminal failure when identities or lifecycle
// change. The model coordinator or an assigned machine may report this failure.
func (s *State) FailStorageResizePlan(ctx context.Context, uuid, machineUUID, message string, now time.Time) error {
	db, err := s.DB(ctx)
	if err != nil {
		return errors.Capture(err)
	}
	return errors.Capture(db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		plan, _, err := s.getResizePlan(ctx, tx, uuid)
		if err != nil {
			return err
		}
		if machineUUID != "" {
			if plan.Identity == "" {
				legacy, err := s.getResizeExecutionRecord(ctx, tx, resizeExecutionKey{UUID: uuid})
				if err != nil {
					return err
				}
				if legacy.Executor != machineUUID && (!plan.Target.MachineScoped || plan.Target.MachineUUID != machineUUID) {
					return storageerrors.StorageResizeNotAssigned
				}
			} else if plan.BackendMachineUUID != machineUUID && plan.NodeMachineUUID != machineUUID {
				return storageerrors.StorageResizeNotAssigned
			}
		}
		return s.setResizePlanResult(ctx, tx, uuid, "failed", message, now)
	}))
}

func (s *State) getResizePlan(ctx context.Context, tx *sqlair.TX, uuid string) (domainstorage.ResizePlan, resizeExecutionRow, error) {
	plan := domainstorage.ResizePlan{UUID: uuid}
	key := resizeExecutionKey{UUID: uuid}
	record, err := s.getResizeExecutionRecord(ctx, tx, key)
	if err != nil {
		return plan, resizeExecutionRow{}, err
	}
	if record.Status != "pending" {
		return plan, resizeExecutionRow{}, storageerrors.StorageResizeResultAlreadyRecorded
	}
	work, err := s.getResizeExecutionWork(ctx, tx, key)
	if err != nil {
		return plan, work, err
	}
	attrs, err := s.resizePoolAttributes(ctx, tx, work.PoolUUID)
	if err != nil {
		return plan, work, err
	}
	target := domainstorage.ResizeTarget{
		SizeMiB: work.RequestedSizeMiB, ProviderType: work.ProviderType, PoolName: work.PoolName, PoolAttributes: attrs,
		Kind: "volume", MachineScoped: work.MachineScoped, MachineUUID: work.MachineUUID, InstanceID: work.InstanceID,
		AttachmentID: work.VolumeAttachmentUUID, VolumeID: work.VolumeID, ProviderID: work.VolumeProviderID,
		DeviceName: work.DeviceName, DeviceLink: work.DeviceLink, FilesystemID: work.FilesystemID,
		FilesystemUUID: work.HostFilesystemUUID, MountPoint: work.MountPoint, Alive: work.Alive, ReadOnly: work.ReadOnly, Shared: work.Shared,
	}
	if work.VolumeUUID == "" {
		target.Kind = "filesystem"
		target.ProviderID = work.FilesystemProviderID
		target.AttachmentID = work.FilesystemAttachmentUUID
	}
	// Include database resource identities as well as the worker-visible inputs.
	data, _ := json.Marshal([]any{work.identity(), target})
	sum := sha256.Sum256(data)
	target.Identity = hex.EncodeToString(sum[:])
	plan.Target = target
	row := resizePlanRow{UUID: uuid}
	stmt, err := s.Prepare(`SELECT (p.*) AS (&resizePlanRow.*) FROM storage_resize_plan AS p WHERE p.request_uuid=$resizePlanRow.request_uuid`, row)
	if err != nil {
		return plan, work, err
	}
	err = tx.Query(ctx, stmt, row).Get(&row)
	if errors.Is(err, sqlair.ErrNoRows) {
		return plan, work, nil
	}
	if err != nil {
		return plan, work, err
	}
	plan.Identity = row.Identity
	plan.BackendMachineUUID = row.BackendMachineUUID
	plan.NodeMachineUUID = row.NodeMachineUUID
	plan.ApprovedSizeMiB = row.ApprovedSizeMiB
	stageKey := resizeStageRow{UUID: uuid}
	stmt, err = s.Prepare(`SELECT (s.*) AS (&resizeStageRow.*) FROM storage_resize_stage AS s WHERE s.request_uuid=$resizeStageRow.request_uuid`, stageKey)
	if err != nil {
		return plan, work, err
	}
	var rows []resizeStageRow
	if err := tx.Query(ctx, stmt, stageKey).GetAll(&rows); err != nil {
		return plan, work, err
	}
	for _, row := range rows {
		stage := domainstorage.ResizeStage{State: row.State, Message: row.Message, ObservedSizeMiB: row.ObservedSizeMiB, FilesystemSizeMiB: row.FilesystemSizeMiB}
		if err := json.Unmarshal([]byte(row.Prerequisites), &stage.Prerequisites); err != nil {
			return plan, work, err
		}
		if err := json.Unmarshal([]byte(row.FollowUp), &stage.FollowUp); err != nil {
			return plan, work, err
		}
		if row.Stage == "backend" {
			plan.Backend = stage
		} else {
			plan.Attachment = stage
		}
	}
	return plan, work, nil
}

func (s *State) resizePoolAttributes(ctx context.Context, tx *sqlair.TX, uuid string) (map[string]string, error) {
	key := entityUUID{UUID: uuid}
	type attribute struct {
		Key   string `db:"key"`
		Value string `db:"value"`
	}
	stmt, err := s.Prepare(`SELECT (a.key,a.value) AS (&attribute.*) FROM storage_pool_attribute AS a WHERE a.storage_pool_uuid=$entityUUID.uuid`, key, attribute{})
	if err != nil {
		return nil, err
	}
	var rows []attribute
	if err := tx.Query(ctx, stmt, key).GetAll(&rows); err != nil && !errors.Is(err, sqlair.ErrNoRows) {
		return nil, err
	}
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		result[row.Key] = row.Value
	}
	return result, nil
}

func (s *State) setResizeStageObservations(ctx context.Context, tx *sqlair.TX, work resizeExecutionRow, report domainstorage.ResizeStageReport) error {
	type observation struct {
		UUID string `db:"uuid"`
		Size uint64 `db:"size_mib"`
	}
	updates := []struct {
		table string
		value observation
	}{}
	if report.Stage == "backend" {
		if work.VolumeUUID != "" {
			updates = append(updates, struct {
				table string
				value observation
			}{"storage_volume", observation{work.VolumeUUID, report.ObservedSizeMiB}})
		} else {
			updates = append(updates, struct {
				table string
				value observation
			}{"storage_filesystem", observation{work.FilesystemUUID, report.ObservedSizeMiB}})
		}
	} else {
		updates = append(updates, struct {
			table string
			value observation
		}{"block_device", observation{work.BlockDeviceUUID, report.ObservedSizeMiB}})
		updates = append(updates, struct {
			table string
			value observation
		}{"storage_filesystem", observation{work.FilesystemUUID, report.FilesystemSizeMiB}})
	}
	for _, update := range updates {
		if update.value.Size == 0 || update.value.UUID == "" {
			continue
		}
		stmt, err := s.Prepare(`UPDATE `+update.table+` SET size_mib=MAX(COALESCE(size_mib,0),$observation.size_mib) WHERE uuid=$observation.uuid`, update.value)
		if err != nil {
			return err
		}
		if err := tx.Query(ctx, stmt, update.value).Run(); err != nil {
			return err
		}
	}
	return nil
}

func (s *State) setResizePlanResult(ctx context.Context, tx *sqlair.TX, uuid, status, message string, now time.Time) error {
	row := resizeExecutionResult{UUID: uuid, Status: status, Message: message, CompletedAt: now}
	stmt, err := s.Prepare(`UPDATE storage_resize_request SET status=$resizeExecutionResult.status,result_message=$resizeExecutionResult.result_message,completed_at=$resizeExecutionResult.completed_at WHERE uuid=$resizeExecutionResult.uuid AND status='pending'`, row)
	if err != nil {
		return err
	}
	return tx.Query(ctx, stmt, row).Run()
}

func (s *State) rejectStagedResize(ctx context.Context, tx *sqlair.TX, uuid string) error {
	row := resizePlanRow{UUID: uuid}
	stmt, err := s.Prepare(`SELECT p.request_uuid AS &resizePlanRow.request_uuid FROM storage_resize_plan AS p WHERE p.request_uuid=$resizePlanRow.request_uuid`, row)
	if err != nil {
		return err
	}
	err = tx.Query(ctx, stmt, row).Get(&row)
	if errors.Is(err, sqlair.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	return storageerrors.StorageResizeNotReady
}

type resizePlanRow struct {
	UUID               string `db:"request_uuid"`
	Identity           string `db:"identity"`
	BackendMachineUUID string `db:"backend_machine_uuid"`
	NodeMachineUUID    string `db:"node_machine_uuid"`
	ApprovedSizeMiB    uint64 `db:"approved_size_mib"`
}
type resizeStageRow struct {
	UUID              string `db:"request_uuid"`
	Stage             string `db:"stage"`
	State             string `db:"state"`
	Message           string `db:"message"`
	Prerequisites     string `db:"prerequisites"`
	FollowUp          string `db:"follow_up"`
	ObservedSizeMiB   uint64 `db:"observed_size_mib"`
	FilesystemSizeMiB uint64 `db:"filesystem_size_mib"`
}
