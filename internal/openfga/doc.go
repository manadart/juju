// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package openfga embeds OpenFGA v1.22.1 behind Juju's database transaction
// runner. Its datastore owns neither connections nor database lifecycle. The
// domain-services worker owns the server and closes its background resources.
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
