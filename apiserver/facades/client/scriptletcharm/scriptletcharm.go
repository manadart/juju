// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package scriptletcharm

import (
	"context"

	"github.com/juju/names/v6"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/core/permission"
	scriptletservice "github.com/juju/juju/domain/scriptlet/service"
	"github.com/juju/juju/rpc/params"
)

// ScriptletCharmService defines the domain service methods required by this facade.
type ScriptletCharmService interface {
	RegisterScriptlet(context.Context, scriptletservice.RegisterScriptletArgs) error
}

// BlockChecker defines the block-checking functionality required by the facade.
type BlockChecker interface {
	ChangeAllowed(context.Context) error
}

// Authorizer defines the authorization methods required by the facade.
type Authorizer interface {
	HasPermission(ctx context.Context, operation permission.Access, target names.Tag) error
}

// API implements the ScriptletCharm facade.
type API struct {
	modelTag   names.ModelTag
	authorizer Authorizer
	check      BlockChecker
	service    ScriptletCharmService
}

func (api *API) checkCanWrite(ctx context.Context) error {
	return api.authorizer.HasPermission(ctx, permission.WriteAccess, api.modelTag)
}

// Register records a scriptlet charm's raw scriptlet text.
func (api *API) Register(ctx context.Context, args params.RegisterScriptletCharmArgs) params.ErrorResult {
	if err := api.checkCanWrite(ctx); err != nil {
		return params.ErrorResult{Error: apiservererrors.ServerError(err)}
	}
	if err := api.check.ChangeAllowed(ctx); err != nil {
		return params.ErrorResult{Error: apiservererrors.ServerError(err)}
	}

	err := api.service.RegisterScriptlet(ctx, scriptletservice.RegisterScriptletArgs{
		ApplicationName: args.ApplicationName,
		Scriptlet:       args.Scriptlet,
	})
	return params.ErrorResult{Error: apiservererrors.ServerError(err)}
}
