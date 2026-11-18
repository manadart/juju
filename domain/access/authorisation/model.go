// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package authorisation

//go:generate go run ../../../generate/authorisation -root ../../..

const (
	// StoreID identifies the controller's Juju authorisation store.
	StoreID = "01J00000000000000000000000"

	// ModelID identifies the immutable model used for permission checks.
	// A changed model requires a new ID and the appropriate schema migration.
	ModelID = "01J00000000000000000000001"

	// ModelFile is the versioned DSL source relative to this package.
	ModelFile = "model_v1.fga"
)
