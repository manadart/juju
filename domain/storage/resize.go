// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"time"

	corestorage "github.com/juju/juju/core/storage"
)

// StorageResizeUUID uniquely identifies a storage resize intent.
type StorageResizeUUID baseUUID

// StorageResizeIntent describes the original resize request. Recording intent
// does not approve the request or change the storage's observed capacity.
type StorageResizeIntent struct {
	UUID                StorageResizeUUID
	StorageInstanceUUID StorageInstanceUUID
	RequestedSizeMiB    uint64
	CreatedAt           time.Time
}

// StorageResizeStatus describes whether a resize result has been reported.
type StorageResizeStatus string

const (
	// StorageResizePending indicates that no terminal result has been reported.
	StorageResizePending StorageResizeStatus = "pending"
	// StorageResizeCompleted indicates a successful resize result.
	StorageResizeCompleted StorageResizeStatus = "completed"
	// StorageResizeFailed indicates an unsuccessful resize result.
	StorageResizeFailed StorageResizeStatus = "failed"
)

// StorageResizeRequest includes a request's historical identity and result.
// Pending requests have no completion time. Results survive storage removal.
type StorageResizeRequest struct {
	StorageResizeIntent
	StorageID   corestorage.ID
	Status      StorageResizeStatus
	Message     string
	CompletedAt *time.Time
}

// StorageResizeResult is a terminal outcome reported by the resize executor.
type StorageResizeResult struct {
	Status  StorageResizeStatus
	Message string
}

// NewStorageResizeUUID creates a new storage resize intent identifier.
func NewStorageResizeUUID() (StorageResizeUUID, error) {
	u, err := newUUID()
	return StorageResizeUUID(u), err
}

// String returns the string representation of this UUID.
func (u StorageResizeUUID) String() string {
	return baseUUID(u).String()
}

// Validate returns an error if this UUID is invalid.
func (u StorageResizeUUID) Validate() error {
	return baseUUID(u).validate()
}
