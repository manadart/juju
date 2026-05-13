// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package scriptlet

// Event represents a model event dispatched to a scriptlet.
type Event struct {
	// Kind identifies the type of event.
	Kind EventKind

	// Attrs holds read-only event data exposed to the scriptlet
	// as event attributes.
	Attrs map[string]interface{}
}

// EventKind identifies the type of scriptlet event.
type EventKind string

const (
	// EventInstall is dispatched when a scriptlet application is
	// first deployed.
	EventInstall EventKind = "install"

	// EventStart is dispatched after install or upgrade.
	EventStart EventKind = "start"

	// EventStop is dispatched when the application is being removed.
	EventStop EventKind = "stop"

	// EventConfigChanged is dispatched when the application config
	// changes.
	EventConfigChanged EventKind = "config_changed"

	// EventRelationCreated is dispatched when a relation is first
	// established for the application.
	EventRelationCreated EventKind = "relation_created"

	// EventRelationJoined is dispatched when a new unit joins a
	// relation.
	EventRelationJoined EventKind = "relation_joined"

	// EventRelationChanged is dispatched when relation data changes.
	EventRelationChanged EventKind = "relation_changed"

	// EventRelationDeparted is dispatched when a unit leaves a
	// relation.
	EventRelationDeparted EventKind = "relation_departed"

	// EventRelationBroken is dispatched when a relation is removed.
	EventRelationBroken EventKind = "relation_broken"

	// EventUpdateStatus is dispatched periodically for status
	// updates.
	EventUpdateStatus EventKind = "update_status"
)

// StarformEventName returns the Starlark-safe event name used in
// Starform scriptlets.
func (k EventKind) StarformEventName() string {
	return string(k)
}

// RelationInfo describes a relation that the scriptlet application
// participates in. Used for tracking and dispatching relation events.
type RelationInfo struct {
	// RelationUUID is the unique identifier of the relation.
	RelationUUID string

	// RelationID is the integer relation ID used in hook contexts.
	RelationID int

	// Endpoint is the local endpoint name.
	Endpoint string

	// RemoteAppName is the name of the remote application.
	RemoteAppName string

	// Life is the lifecycle status of the relation ("alive", "dying", "dead").
	Life string
}

// relationState tracks the local view of a relation for the scriptlet worker.
// This is a simplified version of the uniter's relation.State, adapted for
// unit-less operation.
type relationState struct {
	// RelationUUID is the unique identifier of the relation.
	RelationUUID string

	// RelationID is the integer relation ID.
	RelationID int

	// Endpoint is the local endpoint name.
	Endpoint string

	// RemoteAppName is the remote application name.
	RemoteAppName string

	// Created tracks whether relation_created has been fired.
	Created bool

	// Members tracks remote units currently in scope, keyed by unit name.
	Members map[string]bool
}
