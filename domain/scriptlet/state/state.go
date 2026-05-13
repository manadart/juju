// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"

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
//
// TODO(hackathon): This currently returns an empty list. A real
// implementation would query the database for applications whose
// charm metadata indicates a scriptlet charm type.
func (st *State) GetScriptletApplicationNames(ctx context.Context) ([]string, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}

	// Placeholder query. Once the charm format includes a scriptlet
	// indicator, this will filter by that.
	_ = db
	return nil, nil
}

// NamespaceForWatchScriptletApplications returns the namespace and
// initial query for watching scriptlet application changes.
//
// TODO(hackathon): This uses the application table namespace as a
// placeholder. A real implementation would use a dedicated table or
// view that only tracks scriptlet-type applications.
func (st *State) NamespaceForWatchScriptletApplications() (string, eventsource.NamespaceQuery) {
	return "application", func(ctx context.Context, runner database.TxnRunner) ([]string, error) {
		return st.GetScriptletApplicationNames(ctx)
	}
}
