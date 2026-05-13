// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"

	"github.com/canonical/sqlair"

	"github.com/juju/juju/core/database"
	"github.com/juju/juju/core/watcher/eventsource"
	"github.com/juju/juju/domain"
	"github.com/juju/juju/internal/errors"
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

// GetScriptletApplicationNames returns the names of all applications
// that use scriptlet charms.
func (st *State) GetScriptletApplicationNames(ctx context.Context) ([]string, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}

	stmt, err := st.Prepare(`
SELECT &scriptletCharm.application_name
FROM scriptlet_charm
`, scriptletCharm{})
	if err != nil {
		return nil, errors.Errorf("preparing scriptlet application names query: %w", err)
	}

	var charms []scriptletCharm
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		if err := tx.Query(ctx, stmt).GetAll(&charms); err != nil && !errors.Is(err, sqlair.ErrNoRows) {
			return errors.Errorf("getting scriptlet application names: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Capture(err)
	}

	names := make([]string, len(charms))
	for i, charm := range charms {
		names[i] = charm.ApplicationName
	}
	return names, nil
}

// NamespaceForWatchScriptletApplications returns the namespace and
// initial query for watching scriptlet application changes.
func (st *State) NamespaceForWatchScriptletApplications() (string, eventsource.NamespaceQuery) {
	return "scriptlet_charm", func(ctx context.Context, runner database.TxnRunner) ([]string, error) {
		return st.GetScriptletApplicationNames(ctx)
	}
}

// RegisterScriptlet records the raw scriptlet text for an application name.
func (st *State) RegisterScriptlet(ctx context.Context, applicationName, scriptlet string) error {
	db, err := st.DB(ctx)
	if err != nil {
		return errors.Capture(err)
	}

	charm := scriptletCharm{
		ApplicationName: applicationName,
		Scriptlet:       scriptlet,
	}
	stmt, err := st.Prepare(`
INSERT INTO scriptlet_charm (*) VALUES ($scriptletCharm.*)
ON CONFLICT (application_name) DO UPDATE SET scriptlet = excluded.scriptlet;
`, charm)
	if err != nil {
		return errors.Errorf("preparing register scriptlet statement: %w", err)
	}

	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		if err := tx.Query(ctx, stmt, charm).Run(); err != nil {
			return errors.Errorf("registering scriptlet: %w", err)
		}
		return nil
	})
	return errors.Capture(err)
}
