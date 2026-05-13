// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package scriptletcharm

import (
	"context"
	"reflect"

	"github.com/juju/names/v6"

	"github.com/juju/juju/apiserver/common"
	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/apiserver/facade"
)

// Register is called to expose the facade onto a registry.
func Register(registry facade.FacadeRegistry) {
	registry.MustRegister("ScriptletCharm", 1, func(stdCtx context.Context, ctx facade.ModelContext) (facade.Facade, error) {
		return NewAPI(ctx)
	}, reflect.TypeFor[*API]())
}

// NewAPI returns a new scriptlet charm API facade.
func NewAPI(ctx facade.ModelContext) (*API, error) {
	authorizer := ctx.Auth()
	if !authorizer.AuthClient() {
		return nil, apiservererrors.ErrPerm
	}

	domainServices := ctx.DomainServices()
	return &API{
		modelTag:   names.NewModelTag(ctx.ModelUUID().String()),
		authorizer: authorizer,
		check:      common.NewBlockChecker(domainServices.BlockCommand()),
		service:    domainServices.Scriptlet(),
	}, nil
}
