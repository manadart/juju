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
