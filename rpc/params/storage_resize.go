// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package params

// StorageResizeArgs identifies storage and its requested target capacity.
type StorageResizeArgs struct {
	StorageTag string `json:"storage-tag"`
	SizeMiB    uint64 `json:"size-mib"`
}

// StorageResizeResult identifies a recorded resize intent, not a completed
// resize or an approved target capacity.
type StorageResizeResult struct {
	UUID string `json:"uuid"`
}
