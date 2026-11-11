// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package params

import "time"

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

// StorageResizeRequestArgs identifies a recorded resize request.
type StorageResizeRequestArgs struct {
	UUID string `json:"uuid"`
}

const (
	// StorageResizePending indicates that no terminal result has been reported.
	StorageResizePending = "pending"
	// StorageResizeCompleted indicates a successful resize result.
	StorageResizeCompleted = "completed"
	// StorageResizeFailed indicates an unsuccessful resize result.
	StorageResizeFailed = "failed"
)

// StorageResizeDetails describes a request and its pending or terminal result.
type StorageResizeDetails struct {
	UUID             string     `json:"uuid"`
	StorageTag       string     `json:"storage-tag"`
	RequestedSizeMiB uint64     `json:"requested-size-mib"`
	Status           string     `json:"status"`
	Message          string     `json:"message,omitempty"`
	CreatedAt        time.Time  `json:"created-at"`
	CompletedAt      *time.Time `json:"completed-at,omitempty"`
}
