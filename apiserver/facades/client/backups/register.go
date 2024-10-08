// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package backups

import (
	"context"
	"reflect"

	"github.com/juju/juju/apiserver/facade"
)

// Register is called to expose a package of facades onto a given registry.
func Register(registry facade.FacadeRegistry) {
	registry.MustRegister("Backups", 3, func(stdCtx context.Context, ctx facade.ModelContext) (facade.Facade, error) {
		return newFacade(ctx)
	}, reflect.TypeOf((*API)(nil)))
}

// newFacade provides the required signature for facade registration.
func newFacade(ctx facade.ModelContext) (*API, error) {
	return NewAPI(
		ctx.DomainServices().ControllerConfig(),
		ctx.Auth(),
		ctx.MachineTag(),
		ctx.DataDir(),
		ctx.LogDir(),
	)
}
