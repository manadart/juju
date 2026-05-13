// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"database/sql"
	"time"
)

type scriptletCharmName struct {
	ReferenceName string `db:"reference_name"`
}

type applicationUUID struct {
	UUID string `db:"uuid"`
}

type charmUUID struct {
	UUID string `db:"uuid"`
}

type scriptletContent struct {
	Scriptlet string `db:"scriptlet"`
}

type insertCharm struct {
	UUID            string         `db:"uuid"`
	ReferenceName   string         `db:"reference_name"`
	SourceID        int            `db:"source_id"`
	Revision        int            `db:"revision"`
	ArchitectureID  int            `db:"architecture_id"`
	Available       bool           `db:"available"`
	ArchivePath     sql.NullString `db:"archive_path"`
	ObjectStoreUUID sql.NullString `db:"object_store_uuid"`
	Version         sql.NullString `db:"version"`
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

// deployRef holds an application name and charm reference_name (both equal to
// ApplicationName) for use as SQL input parameters in deploy queries.
type deployRef struct {
	Name          string `db:"name"`
	ReferenceName string `db:"reference_name"`
}

type insertCharmMetadata struct {
	CharmUUID   string `db:"charm_uuid"`
	Name        string `db:"name"`
	Subordinate bool   `db:"subordinate"`
	RunAsID     int    `db:"run_as_id"`
}

type insertApplication struct {
	UUID      string `db:"uuid"`
	Name      string `db:"name"`
	LifeID    int    `db:"life_id"`
	CharmUUID string `db:"charm_uuid"`
	SpaceUUID string `db:"space_uuid"`
}

type insertApplicationPlatform struct {
	ApplicationUUID string         `db:"application_uuid"`
	OSID            int            `db:"os_id"`
	Channel         sql.NullString `db:"channel"`
	ArchitectureID  int            `db:"architecture_id"`
}

type insertApplicationStatus struct {
	ApplicationUUID string         `db:"application_uuid"`
	StatusID        int            `db:"status_id"`
	Message         sql.NullString `db:"message"`
	Data            sql.NullString `db:"data"`
	UpdatedAt       time.Time      `db:"updated_at"`
}

type insertApplicationEndpoint struct {
	UUID              string         `db:"uuid"`
	ApplicationUUID   string         `db:"application_uuid"`
	SpaceUUID         sql.NullString `db:"space_uuid"`
	CharmRelationUUID string         `db:"charm_relation_uuid"`
}

type insertApplicationChannel struct {
	ApplicationUUID string `db:"application_uuid"`
	Track           string `db:"track"`
	Risk            string `db:"risk"`
	Branch          string `db:"branch"`
}

type insertCharmConfig struct {
	CharmUUID    string  `db:"charm_uuid"`
	Key          string  `db:"key"`
	TypeID       int     `db:"type_id"`
	DefaultValue *string `db:"default_value"`
	Description  string  `db:"description"`
}

type insertApplicationConfigHash struct {
	ApplicationUUID string `db:"application_uuid"`
	SHA256          string `db:"sha256"`
}

type charmRelationUUID struct {
	UUID string `db:"uuid"`
}
