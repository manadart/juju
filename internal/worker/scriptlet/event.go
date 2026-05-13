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
