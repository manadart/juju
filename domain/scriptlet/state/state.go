// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"sort"
	"time"

	"github.com/canonical/sqlair"

	coreapplication "github.com/juju/juju/core/application"
	corecharm "github.com/juju/juju/core/charm"
	"github.com/juju/juju/core/database"
	coreerrors "github.com/juju/juju/core/errors"
	corenetwork "github.com/juju/juju/core/network"
	"github.com/juju/juju/core/watcher/eventsource"
	"github.com/juju/juju/domain"
	scriptletservice "github.com/juju/juju/domain/scriptlet/service"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/internal/uuid"
)

// State defines the access mechanism for interacting with scriptlet
// state in the context of the model database.
type State struct {
	*domain.StateBase
}

// NewState constructs a new state for interacting with the underlying
// scriptlet state of a model.
func NewState(factory database.TxnRunnerFactory) *State {
	return &State{
		StateBase: domain.NewStateBase(factory),
	}
}

// GetScriptletApplicationNames returns the reference_name for all scriptlet
// charms in the model.
func (st *State) GetScriptletApplicationNames(ctx context.Context) ([]string, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}

	stmt, err := st.Prepare(`
SELECT &scriptletCharmName.reference_name
FROM charm
JOIN scriptlet_charm ON scriptlet_charm.charm_uuid = charm.uuid
`, scriptletCharmName{})
	if err != nil {
		return nil, errors.Errorf("preparing scriptlet application names query: %w", err)
	}

	var rows []scriptletCharmName
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		if err := tx.Query(ctx, stmt).GetAll(&rows); err != nil && !errors.Is(err, sqlair.ErrNoRows) {
			return errors.Errorf("getting scriptlet application names: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Capture(err)
	}

	names := make([]string, len(rows))
	for i, r := range rows {
		names[i] = r.ReferenceName
	}
	return names, nil
}

// NamespaceForWatchScriptletApplications returns the namespace and
// initial query for watching scriptlet application changes. The initial
// query returns the UUIDs of all applications whose charm has a row in
// the scriptlet_charm table.
func (st *State) InitialWatchStatementScriptletApplications() (string, eventsource.NamespaceQuery) {
	queryFunc := func(ctx context.Context, runner database.TxnRunner) ([]string, error) {
		stmt, err := st.Prepare(`
SELECT a.uuid AS &applicationUUID.uuid
FROM   application AS a
JOIN   charm AS c ON c.uuid = a.charm_uuid
JOIN   scriptlet_charm AS sc ON sc.charm_uuid = c.uuid
`, applicationUUID{})
		if err != nil {
			return nil, errors.Capture(err)
		}

		var result []applicationUUID
		err = runner.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
			if err := tx.Query(ctx, stmt).GetAll(&result); err != nil && !errors.Is(err, sqlair.ErrNoRows) {
				return errors.Errorf("querying for scriptlet applications: %w", err)
			}
			return nil
		})
		if err != nil {
			return nil, errors.Capture(err)
		}
		uuids := make([]string, len(result))
		for i, r := range result {
			uuids[i] = r.UUID
		}
		return uuids, nil
	}
	return "application", queryFunc
}

// IsScriptletApplication returns true if the application with the given
// UUID has a charm that is a scriptlet charm (has a row in scriptlet_charm).
func (st *State) IsScriptletApplication(ctx context.Context, appUUID string) (bool, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return false, errors.Capture(err)
	}

	entity := applicationUUID{UUID: appUUID}
	stmt, err := st.Prepare(`
SELECT a.uuid AS &applicationUUID.uuid
FROM   application AS a
JOIN   charm AS c ON c.uuid = a.charm_uuid
JOIN   scriptlet_charm AS sc ON sc.charm_uuid = c.uuid
WHERE  a.uuid = $applicationUUID.uuid
`, entity)
	if err != nil {
		return false, errors.Capture(err)
	}

	var result applicationUUID
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt, entity).Get(&result)
	})
	if errors.Is(err, sqlair.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, errors.Capture(err)
	}
	return true, nil
}

// ApplicationExists returns true if an application with the given name
// already exists.
func (st *State) ApplicationExists(ctx context.Context, name string) (bool, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return false, errors.Capture(err)
	}

	ref := deployRef{Name: name}
	stmt, err := st.Prepare(`
SELECT a.uuid AS &applicationUUID.uuid
FROM   application AS a
WHERE  a.name = $deployRef.name
`, ref, applicationUUID{})
	if err != nil {
		return false, errors.Capture(err)
	}

	var result applicationUUID
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, ref).Get(&result)
		if errors.Is(err, sqlair.ErrNoRows) {
			return nil
		}
		return err
	})
	if err != nil {
		return false, errors.Capture(err)
	}
	return result.UUID != "", nil
}

// DeployScriptlet registers a scriptlet charm and creates the application
// entity in a single atomic transaction.
func (st *State) DeployScriptlet(ctx context.Context, args scriptletservice.DeployScriptletArgs) error {
	db, err := st.DB(ctx)
	if err != nil {
		return errors.Capture(err)
	}

	charmID, err := corecharm.NewID()
	if err != nil {
		return errors.Errorf("generating charm uuid: %w", err)
	}
	appID, err := coreapplication.NewUUID()
	if err != nil {
		return errors.Errorf("generating application uuid: %w", err)
	}

	relations, err := encodeRelations(charmID.String(), args.Relations)
	if err != nil {
		return errors.Errorf("encoding charm relations: %w", err)
	}

	// ── Insert statements ──────────────────────────────────────────────────

	charmRow := insertCharm{
		UUID:            charmID.String(),
		ReferenceName:   args.ApplicationName,
		SourceID:        0, // local
		Revision:        -1,
		ArchitectureID:  0, // amd64
		Available:       true,
		ArchivePath:     sql.NullString{},
		ObjectStoreUUID: sql.NullString{},
		Version:         sql.NullString{},
	}
	insCharmStmt, err := st.Prepare(`
INSERT INTO charm (uuid, reference_name, source_id, revision, architecture_id, available, archive_path, object_store_uuid, version)
VALUES ($insertCharm.uuid, $insertCharm.reference_name, $insertCharm.source_id, $insertCharm.revision, $insertCharm.architecture_id, $insertCharm.available, $insertCharm.archive_path, $insertCharm.object_store_uuid, $insertCharm.version)
`, charmRow)
	if err != nil {
		return errors.Errorf("preparing insert charm: %w", err)
	}

	scriptletRow := insertScriptletCharm{
		CharmUUID: charmID.String(),
		Scriptlet: args.Scriptlet,
	}
	insScriptletStmt, err := st.Prepare(`
INSERT INTO scriptlet_charm (*) VALUES ($insertScriptletCharm.*)
`, scriptletRow)
	if err != nil {
		return errors.Errorf("preparing insert scriptlet_charm: %w", err)
	}

	insRelStmt, err := st.Prepare(`
INSERT INTO charm_relation (*) VALUES ($insertCharmRelation.*)
`, insertCharmRelation{})
	if err != nil {
		return errors.Errorf("preparing insert charm_relation: %w", err)
	}

	configRows := encodeCharmConfig(charmID.String(), args.Config)
	insConfigStmt, err := st.Prepare(`
INSERT INTO charm_config (*) VALUES ($insertCharmConfig.*)
`, insertCharmConfig{})
	if err != nil {
		return errors.Errorf("preparing insert charm_config: %w", err)
	}

	charmMetaRow := insertCharmMetadata{
		CharmUUID:   charmID.String(),
		Name:        args.ApplicationName,
		Subordinate: false,
		RunAsID:     0,
	}
	insCharmMetaStmt, err := st.Prepare(`
INSERT INTO charm_metadata (charm_uuid, name, subordinate, run_as_id)
VALUES ($insertCharmMetadata.charm_uuid, $insertCharmMetadata.name, $insertCharmMetadata.subordinate, $insertCharmMetadata.run_as_id)
`, charmMetaRow)
	if err != nil {
		return errors.Errorf("preparing insert charm_metadata: %w", err)
	}

	appRow := insertApplication{
		UUID:      appID.String(),
		Name:      args.ApplicationName,
		LifeID:    0, // alive
		CharmUUID: charmID.String(),
		SpaceUUID: corenetwork.AlphaSpaceId.String(),
	}
	insAppStmt, err := st.Prepare(`
INSERT INTO application (uuid, name, life_id, charm_uuid, charm_modified_version, space_uuid)
VALUES ($insertApplication.uuid, $insertApplication.name, $insertApplication.life_id, $insertApplication.charm_uuid, 0, $insertApplication.space_uuid)
`, appRow)
	if err != nil {
		return errors.Errorf("preparing insert application: %w", err)
	}

	platformRow := insertApplicationPlatform{
		ApplicationUUID: appID.String(),
		OSID:            0, // ubuntu
		Channel:         sql.NullString{},
		ArchitectureID:  0, // amd64
	}
	insPlatformStmt, err := st.Prepare(`
INSERT INTO application_platform (application_uuid, os_id, channel, architecture_id)
VALUES ($insertApplicationPlatform.application_uuid, $insertApplicationPlatform.os_id, $insertApplicationPlatform.channel, $insertApplicationPlatform.architecture_id)
`, platformRow)
	if err != nil {
		return errors.Errorf("preparing insert application_platform: %w", err)
	}

	statusRow := insertApplicationStatus{
		ApplicationUUID: appID.String(),
		StatusID:        1, // unknown
		Message:         sql.NullString{},
		Data:            sql.NullString{},
		UpdatedAt:       time.Now().UTC(),
	}
	insStatusStmt, err := st.Prepare(`
INSERT INTO application_status (application_uuid, status_id, message, data, updated_at)
VALUES ($insertApplicationStatus.application_uuid, $insertApplicationStatus.status_id, $insertApplicationStatus.message, $insertApplicationStatus.data, $insertApplicationStatus.updated_at)
`, statusRow)
	if err != nil {
		return errors.Errorf("preparing insert application_status: %w", err)
	}

	insEndpointStmt, err := st.Prepare(`
INSERT INTO application_endpoint (uuid, application_uuid, space_uuid, charm_relation_uuid)
VALUES ($insertApplicationEndpoint.uuid, $insertApplicationEndpoint.application_uuid, $insertApplicationEndpoint.space_uuid, $insertApplicationEndpoint.charm_relation_uuid)
`, insertApplicationEndpoint{})
	if err != nil {
		return errors.Errorf("preparing insert application_endpoint: %w", err)
	}

	channelRow := insertApplicationChannel{
		ApplicationUUID: appID.String(),
		Track:           "",
		Risk:            "stable",
		Branch:          "",
	}
	insChannelStmt, err := st.Prepare(`
INSERT INTO application_channel (application_uuid, track, risk, branch)
VALUES ($insertApplicationChannel.application_uuid, $insertApplicationChannel.track, $insertApplicationChannel.risk, $insertApplicationChannel.branch)
`, channelRow)
	if err != nil {
		return errors.Errorf("preparing insert application_channel: %w", err)
	}

	configHashRow := insertApplicationConfigHash{
		ApplicationUUID: appID.String(),
		SHA256:          hashCharmConfigDefaults(configRows),
	}
	insConfigHashStmt, err := st.Prepare(`
INSERT INTO application_config_hash (application_uuid, sha256)
VALUES ($insertApplicationConfigHash.application_uuid, $insertApplicationConfigHash.sha256)
`, configHashRow)
	if err != nil {
		return errors.Errorf("preparing insert application_config_hash: %w", err)
	}

	return db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		// Insert charm rows.
		if err := tx.Query(ctx, insCharmStmt, charmRow).Run(); err != nil {
			return errors.Errorf("inserting charm: %w", err)
		}
		if err := tx.Query(ctx, insScriptletStmt, scriptletRow).Run(); err != nil {
			return errors.Errorf("inserting scriptlet_charm: %w", err)
		}
		if len(relations) > 0 {
			if err := tx.Query(ctx, insRelStmt, relations).Run(); err != nil {
				return errors.Errorf("inserting charm_relation: %w", err)
			}
		}
		if len(configRows) > 0 {
			if err := tx.Query(ctx, insConfigStmt, configRows).Run(); err != nil {
				return errors.Errorf("inserting charm_config: %w", err)
			}
		}
		if err := tx.Query(ctx, insCharmMetaStmt, charmMetaRow).Run(); err != nil {
			return errors.Errorf("inserting charm_metadata: %w", err)
		}

		// Insert application rows.
		if err := tx.Query(ctx, insAppStmt, appRow).Run(); err != nil {
			return errors.Errorf("inserting application: %w", err)
		}
		if err := tx.Query(ctx, insPlatformStmt, platformRow).Run(); err != nil {
			return errors.Errorf("inserting application_platform: %w", err)
		}
		if err := tx.Query(ctx, insStatusStmt, statusRow).Run(); err != nil {
			return errors.Errorf("inserting application_status: %w", err)
		}
		if err := tx.Query(ctx, insChannelStmt, channelRow).Run(); err != nil {
			return errors.Errorf("inserting application_channel: %w", err)
		}
		if err := tx.Query(ctx, insConfigHashStmt, configHashRow).Run(); err != nil {
			return errors.Errorf("inserting application_config_hash: %w", err)
		}
		for _, rel := range relations {
			epID, err := uuid.NewUUID()
			if err != nil {
				return errors.Errorf("generating endpoint uuid: %w", err)
			}
			epRow := insertApplicationEndpoint{
				UUID:              epID.String(),
				ApplicationUUID:   appID.String(),
				SpaceUUID:         sql.NullString{},
				CharmRelationUUID: rel.UUID,
			}
			if err := tx.Query(ctx, insEndpointStmt, epRow).Run(); err != nil {
				return errors.Errorf("inserting application_endpoint: %w", err)
			}
		}

		return nil
	})
}

func encodeRelations(charmUUID string, rels []scriptletservice.ScriptletRelation) ([]insertCharmRelation, error) {
	result := make([]insertCharmRelation, 0, len(rels))
	for _, r := range rels {
		relUUID, err := uuid.NewUUID()
		if err != nil {
			return nil, errors.Errorf("generating relation uuid: %w", err)
		}
		roleID, err := encodeRole(r.Role)
		if err != nil {
			return nil, err
		}
		scopeID := 0 // global
		if r.Scope == "container" {
			scopeID = 1
		}
		result = append(result, insertCharmRelation{
			UUID:      relUUID.String(),
			CharmUUID: charmUUID,
			Name:      r.Name,
			RoleID:    roleID,
			ScopeID:   scopeID,
			Interface: r.Interface,
			Optional:  r.Optional,
			Capacity:  r.Limit,
		})
	}
	return result, nil
}

func encodeRole(role string) (int, error) {
	switch role {
	case "provider":
		return 0, nil
	case "requirer":
		return 1, nil
	case "peer":
		return 2, nil
	default:
		return -1, errors.Errorf("unknown relation role %q", role)
	}
}

func encodeCharmConfig(charmUUID string, opts []scriptletservice.ScriptletConfigOption) []insertCharmConfig {
	rows := make([]insertCharmConfig, 0, len(opts))
	for _, opt := range opts {
		var defaultValue *string
		if opt.DefaultValue != "" {
			v := opt.DefaultValue
			defaultValue = &v
		}
		rows = append(rows, insertCharmConfig{
			CharmUUID:    charmUUID,
			Key:          opt.Key,
			TypeID:       encodeConfigType(opt.Type),
			DefaultValue: defaultValue,
			Description:  opt.Description,
		})
	}
	return rows
}

func encodeConfigType(t string) int {
	switch t {
	case "string":
		return 0
	case "int":
		return 1
	case "float":
		return 2
	case "boolean":
		return 3
	case "secret":
		return 4
	default:
		return 0 // default to string
	}
}

// hashCharmConfigDefaults produces a SHA-256 hash of the config defaults
// keyed by option name. This seeds the application_config_hash row so the
// change stream watcher can detect subsequent config mutations.
func hashCharmConfigDefaults(rows []insertCharmConfig) string {
	h := sha256.New()
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Key < rows[j].Key
	})
	for _, r := range rows {
		_, _ = h.Write([]byte(r.Key))
		if r.DefaultValue != nil {
			_, _ = h.Write([]byte(*r.DefaultValue))
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

// coreerrors is used above for NotFound — keep the import used.
var _ = coreerrors.NotFound

// GetApplicationScriptlet returns the scriptlet source for the application
// identified by its UUID. It follows the path: application → charm →
// scriptlet_charm.
func (st *State) GetApplicationScriptlet(ctx context.Context, appUUID string) (string, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return "", errors.Capture(err)
	}

	entity := applicationUUID{UUID: appUUID}
	stmt, err := st.Prepare(`
SELECT sc.scriptlet AS &scriptletContent.scriptlet
FROM   application AS a
JOIN   charm AS c ON c.uuid = a.charm_uuid
JOIN   scriptlet_charm AS sc ON sc.charm_uuid = c.uuid
WHERE  a.uuid = $applicationUUID.uuid
`, entity, scriptletContent{})
	if err != nil {
		return "", errors.Errorf("preparing get application scriptlet: %w", err)
	}

	var result scriptletContent
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt, entity).Get(&result)
	})
	if errors.Is(err, sqlair.ErrNoRows) {
		return "", errors.Errorf("scriptlet not found for application %q", appUUID).Add(coreerrors.NotFound)
	}
	if err != nil {
		return "", errors.Capture(err)
	}
	return result.Scriptlet, nil
}
