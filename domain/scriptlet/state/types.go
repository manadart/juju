// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import "database/sql"

type scriptletCharmName struct {
	ReferenceName string `db:"reference_name"`
}

type applicationUUID struct {
	UUID string `db:"uuid"`
}

type charmUUID struct {
	UUID string `db:"uuid"`
}

type insertCharm struct {
	UUID           string         `db:"uuid"`
	ReferenceName  string         `db:"reference_name"`
	SourceID       int            `db:"source_id"`
	Revision       int            `db:"revision"`
	ArchitectureID int            `db:"architecture_id"`
	Available      bool           `db:"available"`
	IsScriptlet    bool           `db:"is_scriptlet"`
	ArchivePath    sql.NullString `db:"archive_path"`
	ObjectStoreUUID sql.NullString `db:"object_store_uuid"`
	Version        sql.NullString `db:"version"`
}

type insertScriptletCharm struct {
	CharmUUID string `db:"charm_uuid"`
	Scriptlet string `db:"scriptlet"`
}

type insertCharmRelation struct {
	UUID      string `db:"uuid"`
	CharmUUID string `db:"charm_uuid"`
	Name      string `db:"name"`
	RoleID    int    `db:"role_id"`
	ScopeID   int    `db:"scope_id"`
	Interface string `db:"interface"`
	Optional  bool   `db:"optional"`
	Capacity  int    `db:"capacity"`
}
