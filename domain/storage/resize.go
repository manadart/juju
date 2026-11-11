// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import "time"

// StorageResizeUUID uniquely identifies a storage resize intent.
type StorageResizeUUID baseUUID

// StorageResizeIntent describes a pending resize request. Recording intent
// does not approve the request or change the storage's observed capacity.
type StorageResizeIntent struct {
	UUID                StorageResizeUUID
	StorageInstanceUUID StorageInstanceUUID
	RequestedSizeMiB    uint64
	CreatedAt           time.Time
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
