// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

type entityUUID struct {
	UUID string `db:"uuid"`
}

type storageID struct {
	ID string `db:"storage_id"`
}