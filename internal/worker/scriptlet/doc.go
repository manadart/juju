// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package scriptlet implements a controller-side worker that evaluates
// Starform scriptlets for agentless charms. Instead of running a unit agent,
// the worker watches model events (config changes, relation changes, lifecycle
// events) and dispatches them to charm-provided scriptlets. Scriptlets declare
// intent (status, state, relation data) via an intent collector; intents are
// validated and applied through domain services only after successful script
// execution.
package scriptlet
