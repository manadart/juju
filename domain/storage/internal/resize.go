// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package internal

import "time"

// StorageResizeIntent records a pending request to grow a storage instance.
// RequestedSizeMiB is a target capacity, not an increment or observed size.
type StorageResizeIntent struct {
	UUID                string
	StorageInstanceUUID string
	RequestedSizeMiB    uint64
	CreatedAt           time.Time
}

// StorageResizeRequest includes the retained identity and result of a request.
type StorageResizeRequest struct {
	StorageResizeIntent
	StorageID   string
	Status      string
	Message     string
	CompletedAt *time.Time
}

// StorageResizeResult is a terminal result reported by the resize executor.
type StorageResizeResult struct {
	Status      string
	Message     string
	CompletedAt time.Time
}

// StorageResizeWork is the current machine-local view of a resize target.
type StorageResizeWork struct {
	UUID             string
	RequestedSizeMiB uint64
	ProviderType     string
	VolumeID         string
	VolumeProviderID string
	DeviceName       string
	FilesystemID     string
	MountPoint       string
	VolumeSizeMiB    uint64
	Alive            bool
	ReadOnly         bool
	Shared           bool
	MachineScoped    bool
}

// StorageResizeReport records an execution outcome and measured capacities.
// Zero capacities mean no observation was made at that layer.
type StorageResizeReport struct {
	StorageResizeResult
	VolumeSizeMiB     uint64
	FilesystemSizeMiB uint64
}

// UnitStorageResizeTarget is a snapshot used to stage a unit's requests atomically.
// Attachment and charm identities prevent applying stale target selections.
type UnitStorageResizeTarget struct {
	UnitUUID             string
	ApplicationCharmUUID string
	UnitCharmUUID        string
	StorageInstanceUUID  string
	StorageID            string
	AttachmentUUID       string
	SizeMiB              uint64
	ObservedSizeMiB      uint64
	PendingUUID          string
	PendingSizeMiB       uint64
	Alive                bool
}
