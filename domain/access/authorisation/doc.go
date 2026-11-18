// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package authorisation identifies the persisted Juju permission model.
// The versioned DSL file named by ModelFile is its editable source. ModelID
// identifies that immutable definition in the controller database and is shared
// by the generator and the running authorizer.
//
// From the repository root, regenerate the model with:
//
//	make rebuild-authorisation
//
// Alternatively, run go generate ./domain/access/authorisation. The generator
// parses and validates the DSL, then updates only the marked model INSERT in
// domain/schema/controller/sql/0017-user-permission.sql. It preserves the rest
// of that file. To check freshness:
//
//	go run ./generate/authorisation -check
//
// For a changed definition, retain the old DSL file, add the next version and
// update ModelFile and ModelID together. The generator refuses different model
// content under the same ID. This beta prototype uses the base schema and needs
// fresh controllers. Supported upgrades will require an incremental schema
// migration that installs the new model before the authorizer selects it;
// released schema patches must not be regenerated in place.
//
// JIMM's model is a useful reference:
// https://github.com/canonical/jimm/blob/v3/openfga/authorisation_model.fga
// This model preserves the prototype's controller, cloud, model and offer
// grants, with separate local and external subjects. JIMM's groups, roles and
// parent-object inheritance are outside this initial model. Disabled-user
// checks and controller-wide overrides remain in their existing callers.
package authorisation
