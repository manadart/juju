// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// authorisation generates the persisted permission model from the versioned
// OpenFGA DSL in domain/access/authorisation. It uses the pinned language
// transformer and the embedded engine's model validator. Output is canonical
// JSON inside the controller schema's marked seed block; all other SQL is
// preserved. The authorizer reads the resulting model from the controller DB.
//
// Run from the repository root with go run ./generate/authorisation, or use
// make rebuild-authorisation. The -root flag supports go generate from the
// model package. The -check flag reports stale output without writing files.
// A different definition cannot overwrite the current immutable model ID.
package main
