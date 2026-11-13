// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package internal

// K8sFilesystemAllocation preserves an existing unit's allocation and binding.
type K8sFilesystemAllocation struct {
	StorageName string `db:"storage_name"`
	UnitName    string `db:"unit_name"`
	SizeMiB     uint64 `db:"size_mib"`
	VolumeID    string `db:"volume_id"`
	ClaimName   string `db:"claim_name"`
}
