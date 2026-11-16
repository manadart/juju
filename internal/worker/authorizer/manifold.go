// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package authorizer

import (
	"context"

	"github.com/juju/clock"
	"github.com/juju/worker/v5"
	"github.com/juju/worker/v5/dependency"
	"github.com/openfga/openfga/pkg/storage"

	"github.com/juju/juju/core/database"
	"github.com/juju/juju/core/permission"
	"github.com/juju/juju/domain/access/service"
	"github.com/juju/juju/domain/access/state"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/internal/openfga"
	"github.com/juju/juju/internal/worker/common"
)

// ManifoldConfig supplies the database dependency and construction functions.
// The backing service is built directly, avoiding a domainservices dependency.
type ManifoldConfig struct {
	DBAccessorName    string
	Clock             clock.Clock
	NewWorker         func(context.Context, Config) (worker.Worker, error)
	NewBackingService func(database.DBGetter, clock.Clock) storage.OpenFGADatastore
	NewAuthorizer     func(context.Context, storage.OpenFGADatastore) (Authorizer, error)
}

// Validate checks the dependencies required to start the worker.
func (c ManifoldConfig) Validate() error {
	if c.DBAccessorName == "" {
		return errors.New("empty DBAccessorName")
	}
	if c.Clock == nil {
		return errors.New("nil Clock")
	}
	if c.NewWorker == nil {
		return errors.New("nil NewWorker")
	}
	if c.NewBackingService == nil {
		return errors.New("nil NewBackingService")
	}
	if c.NewAuthorizer == nil {
		return errors.New("nil NewAuthorizer")
	}
	return nil
}

// Manifold exposes one shared authorizer per controller process. Its caller
// must gate startup on completion of controller database upgrades.
func Manifold(config ManifoldConfig) dependency.Manifold {
	return dependency.Manifold{
		Inputs: []string{config.DBAccessorName},
		Start:  config.start,
		Output: output,
	}
}

func (c ManifoldConfig) start(ctx context.Context, getter dependency.Getter) (worker.Worker, error) {
	if err := c.Validate(); err != nil {
		return nil, errors.Capture(err)
	}
	var dbGetter database.DBGetter
	if err := getter.Get(c.DBAccessorName, &dbGetter); err != nil {
		return nil, errors.Capture(err)
	}
	backing := c.NewBackingService(dbGetter, c.Clock)
	if backing == nil {
		return nil, errors.New("nil authorisation backing service")
	}
	return c.NewWorker(ctx, Config{
		NewAuthorizer: func(ctx context.Context) (Authorizer, error) {
			return c.NewAuthorizer(ctx, backing)
		},
	})
}

func output(in worker.Worker, out any) error {
	if w, ok := in.(*common.CleanupWorker); ok {
		in = w.Worker
	}
	w, ok := in.(*authorizerWorker)
	if !ok {
		return errors.Errorf("expected authorizer worker, got %T", in)
	}
	switch out := out.(type) {
	case *permission.Authorizer:
		*out = w
		return nil
	default:
		return errors.Errorf("unsupported authorizer output type %T", out)
	}
}

// NewBackingService supplies the access service and SQLair state directly from
// the accessor, independently of the domainservices worker.
func NewBackingService(dbGetter database.DBGetter, clock clock.Clock) storage.OpenFGADatastore {
	factory := database.NewTxnRunnerFactoryForNamespace(dbGetter.GetDB, database.ControllerNS)
	return service.NewAuthorisationService(state.NewAuthorisationState(factory), clock)
}

// NewAuthorizer constructs the embedded engine against its backing service.
func NewAuthorizer(ctx context.Context, backing storage.OpenFGADatastore) (Authorizer, error) {
	return openfga.NewAuthorizer(ctx, backing)
}
