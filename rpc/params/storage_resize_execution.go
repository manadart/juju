// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package params

// StorageResizeWork describes a machine-local resize target.
type StorageResizeWork struct {
	UUID             string `json:"uuid"`
	RequestedSizeMiB uint64 `json:"requested-size-mib"`
	ProviderType     string `json:"provider-type"`
	VolumeID         string `json:"volume-id"`
	VolumeProviderID string `json:"volume-provider-id"`
	DeviceName       string `json:"device-name"`
	FilesystemID     string `json:"filesystem-id,omitempty"`
	MountPoint       string `json:"mount-point,omitempty"`
}

// StorageResizeReport includes the execution outcome and measured capacities.
type StorageResizeReport struct {
	UUID              string `json:"uuid"`
	Status            string `json:"status"`
	Message           string `json:"message"`
	VolumeSizeMiB     uint64 `json:"volume-size-mib"`
	FilesystemSizeMiB uint64 `json:"filesystem-size-mib"`
}
