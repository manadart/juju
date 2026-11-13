// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package internal

import "github.com/juju/juju/domain/storage"

// ApplicationStorageSnapshot binds directive validation to a charm revision.
type ApplicationStorageSnapshot struct {
	CharmUUID  string
	Directives []ApplicationStorageDirective
}

// ApplicationStorageDirective includes the constraints used to validate edits.
type ApplicationStorageDirective struct {
	storage.DirectiveArg
	Definition   CharmStorageDefinitionForValidation
	ProviderType string
}
