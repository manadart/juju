// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package params

// StorageResizeStageTarget carries a staged resize reconciliation.
type StorageResizeStageTarget struct {
	Identity       string            `json:"identity"`
	SizeMiB        uint64            `json:"size-mib"`
	ProviderType   string            `json:"provider-type"`
	PoolName       string            `json:"pool-name"`
	PoolAttributes map[string]string `json:"pool-attributes"`
	Kind           string            `json:"kind"`
	MachineScoped  bool              `json:"machine-scoped"`
	MachineUUID    string            `json:"machine-uuid"`
	InstanceID     string            `json:"instance-id"`
	AttachmentID   string            `json:"attachment-id"`
	VolumeID       string            `json:"volume-id"`
	ProviderID     string            `json:"provider-id"`
	DeviceName     string            `json:"device-name"`
	DeviceLink     string            `json:"device-link"`
	FilesystemID   string            `json:"filesystem-id"`
	FilesystemUUID string            `json:"filesystem-uuid"`
	MountPoint     string            `json:"mount-point"`
	Alive          bool              `json:"alive"`
	ReadOnly       bool              `json:"read-only"`
	Shared         bool              `json:"shared"`
}

// StorageResizeStageWork carries a staged resize reconciliation.
type StorageResizeStageWork struct {
	UUID     string                   `json:"uuid"`
	Identity string                   `json:"identity"`
	Stage    string                   `json:"stage"`
	Phase    string                   `json:"phase"`
	Target   StorageResizeStageTarget `json:"target"`
}

// StorageResizeStageReport carries a staged resize reconciliation.
type StorageResizeStageReport struct {
	UUID              string   `json:"uuid"`
	Identity          string   `json:"identity"`
	Stage             string   `json:"stage"`
	Phase             string   `json:"phase"`
	Status            string   `json:"status"`
	Message           string   `json:"message"`
	Prerequisites     []string `json:"prerequisites"`
	FollowUp          []string `json:"follow-up"`
	ObservedSizeMiB   uint64   `json:"observed-size-mib"`
	FilesystemSizeMiB uint64   `json:"filesystem-size-mib"`
}
