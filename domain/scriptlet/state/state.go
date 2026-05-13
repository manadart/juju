// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"database/sql"

	"github.com/canonical/sqlair"

	corecharm "github.com/juju/juju/core/charm"
	"github.com/juju/juju/core/database"
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
WHERE is_scriptlet = TRUE
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
// initial query for watching scriptlet application changes.
func (st *State) NamespaceForWatchScriptletApplications() (string, eventsource.NamespaceQuery) {
	return "charm", func(ctx context.Context, runner database.TxnRunner) ([]string, error) {
		return st.GetScriptletApplicationNames(ctx)
	}
}

// RegisterScriptlet inserts a scriptlet charm row, the scriptlet body, and all
// its relations into the model database. Re-registering an existing application
// name replaces the prior charm and its relations.
func (st *State) RegisterScriptlet(ctx context.Context, args scriptletservice.RegisterScriptletArgs) error {
	db, err := st.DB(ctx)
	if err != nil {
		return errors.Capture(err)
	}

	charmID, err := corecharm.NewID()
	if err != nil {
		return errors.Errorf("generating charm uuid: %w", err)
	}

	relations, err := encodeRelations(charmID.String(), args.Relations)
	if err != nil {
		return errors.Errorf("encoding charm relations: %w", err)
	}

	// Delete-then-insert for idempotency: remove the old scriptlet_charm and
	// charm_relation rows first (they FK to charm), then the charm row.
	delCharmRelStmt, err := st.Prepare(`
DELETE FROM charm_relation
WHERE charm_uuid IN (
    SELECT uuid FROM charm
    WHERE reference_name = $insertCharm.reference_name AND is_scriptlet = TRUE
)`, insertCharm{})
	if err != nil {
		return errors.Errorf("preparing delete charm_relation: %w", err)
	}

	delScriptletStmt, err := st.Prepare(`
DELETE FROM scriptlet_charm
WHERE charm_uuid IN (
    SELECT uuid FROM charm
    WHERE reference_name = $insertCharm.reference_name AND is_scriptlet = TRUE
)`, insertCharm{})
	if err != nil {
		return errors.Errorf("preparing delete scriptlet_charm: %w", err)
	}

	delCharmStmt, err := st.Prepare(`
DELETE FROM charm
WHERE reference_name = $insertCharm.reference_name AND is_scriptlet = TRUE
`, insertCharm{})
	if err != nil {
		return errors.Errorf("preparing delete charm: %w", err)
	}

	charmRow := insertCharm{
		UUID:          charmID.String(),
		ReferenceName: args.ApplicationName,
		SourceID:      0, // local
		Revision:      -1,
		ArchitectureID: 0, // amd64 — required by chk_charm_architecture for local source
		Available:     true,
		IsScriptlet:   true,
		ArchivePath:   sql.NullString{},
		ObjectStoreUUID: sql.NullString{},
		Version:       sql.NullString{},
	}
	insCharmStmt, err := st.Prepare(`
INSERT INTO charm (uuid, reference_name, source_id, revision, architecture_id, available, is_scriptlet, archive_path, object_store_uuid, version)
VALUES ($insertCharm.uuid, $insertCharm.reference_name, $insertCharm.source_id, $insertCharm.revision, $insertCharm.architecture_id, $insertCharm.available, $insertCharm.is_scriptlet, $insertCharm.archive_path, $insertCharm.object_store_uuid, $insertCharm.version)
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

	return db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		ref := insertCharm{ReferenceName: args.ApplicationName}
		if err := tx.Query(ctx, delCharmRelStmt, ref).Run(); err != nil {
			return errors.Errorf("deleting old charm_relation: %w", err)
		}
		if err := tx.Query(ctx, delScriptletStmt, ref).Run(); err != nil {
			return errors.Errorf("deleting old scriptlet_charm: %w", err)
		}
		if err := tx.Query(ctx, delCharmStmt, ref).Run(); err != nil {
			return errors.Errorf("deleting old charm: %w", err)
		}
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
