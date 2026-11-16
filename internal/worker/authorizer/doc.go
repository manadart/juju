// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package authorizer owns the embedded permission engine for one controller
// process. Its manifold constructs the access domain's authorisation service
// directly from controller state and the database accessor. It must start
// after database upgrades, before the domainservices worker consumes its
// permission.Authorizer output. It has no dependency on domainservices.
//
// Startup loads and validates the persisted model before publishing the output.
// The startup context bounds that initialisation; the worker's own context
// governs the engine thereafter. Shutdown cancels checks, closes admission,
// joins accepted calls and closes the engine. Consumers retain only the Check
// capability and cannot close the shared engine.
//
// Ordinary access services use this shared checker for permission evaluation.
// The backing service translates the OpenFGA storage contract, while access
// state owns SQLair queries, atomic tuple writes and transaction replay markers.
// Resource mutations continue to write tuples in their existing transactions.
package authorizer
