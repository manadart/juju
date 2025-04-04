// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package instancepoller

import (
	"context"
	"reflect"

	"github.com/juju/errors"

	"github.com/juju/juju/apiserver/facade"
)

// Register is called to expose a package of facades onto a given registry.
func Register(registry facade.FacadeRegistry) {
	registry.MustRegister("InstancePoller", 4, func(stdCtx context.Context, ctx facade.ModelContext) (facade.Facade, error) {
		return newFacade(ctx)
	}, reflect.TypeOf((*InstancePollerAPI)(nil)))
}

// newFacade wraps NewInstancePollerAPI for facade registration.
func newFacade(ctx facade.ModelContext) (*InstancePollerAPI, error) {
	st := ctx.State()
	m, err := st.Model()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return NewInstancePollerAPI(
		st,
		ctx.DomainServices().Network(),
		ctx.DomainServices().Machine(),
		m,
		ctx.Resources(),
		ctx.Auth(),
		ctx.DomainServices().ControllerConfig(),
		ctx.Clock(),
		ctx.Logger().Child("instancepoller"))
}
