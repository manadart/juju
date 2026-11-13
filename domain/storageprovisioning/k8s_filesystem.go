// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storageprovisioning

// K8sFilesystemAllocation describes storage already allocated to a unit.
// Empty provider identities mean Kubernetes has not yet reported the binding.
type K8sFilesystemAllocation struct {
	StorageName string
	UnitName    string
	SizeMiB     uint64
	VolumeID    string
	ClaimName   string
}
