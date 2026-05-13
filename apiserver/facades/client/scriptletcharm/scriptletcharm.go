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
	DeployScriptlet(context.Context, scriptletservice.DeployScriptletArgs) error
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

// Deploy registers the scriptlet charm and creates the application in one shot.
func (api *API) Deploy(ctx context.Context, args params.DeployScriptletCharmArgs) params.ErrorResult {
	if err := api.checkCanWrite(ctx); err != nil {
		return params.ErrorResult{Error: apiservererrors.ServerError(err)}
	}
	if err := api.check.ChangeAllowed(ctx); err != nil {
		return params.ErrorResult{Error: apiservererrors.ServerError(err)}
	}

	relations := make([]scriptletservice.ScriptletRelation, len(args.Relations))
	for i, r := range args.Relations {
		relations[i] = scriptletservice.ScriptletRelation{
			Name:      r.Name,
			Role:      r.Role,
			Interface: r.Interface,
			Scope:     r.Scope,
			Optional:  r.Optional,
			Limit:     r.Limit,
		}
	}

	err := api.service.DeployScriptlet(ctx, scriptletservice.DeployScriptletArgs{
		ApplicationName: args.ApplicationName,
		Scriptlet:       args.Scriptlet,
		Relations:       relations,
	})
	return params.ErrorResult{Error: apiservererrors.ServerError(err)}
}
