// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package openfga configures the embedded OpenFGA v1.22.1 engine. The access
// domain supplies its datastore through an authorisation backing service and
// SQLair state. A dedicated authorizer worker owns startup, cancellation and
// closure. Database connections remain owned by the accessor.
//
// The versioned DSL under domain/access/authorisation generates the persisted
// model seed. The engine reads that model from the controller database.
//
// The controller schema owns the persisted store, immutable model, tuples,
// assertions and changelog. Domain mutations write tuples inside their own
// SQLair transactions; triggers record every actual tuple change atomically.
// Read-only permission views preserve Juju's existing list and migration
// representations, without storing a second copy of grants. Access checks use
// OpenFGA, with cross-request tuple and decision caches disabled.
//
// The initial model preserves controller, cloud, model and offer access levels.
// Public grants use external-user wildcards, and disabled/removed users are
// rejected by the access service before reaching the engine. Credential
// ownership and externally delegated JWT permissions retain their contracts.
//
// This prototype requires a fresh controller. It changes the beta base schema
// and native backup payload; upgrading an existing database or restoring a
// backup from the previous permission schema is not supported. Recovery
// restores stores, models, assertions and tuples, rebuilding the changelog.
// Existing permission tests have not been ported or run for this prototype.
package openfga
