// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package scriptletcharm

import (
	"context"

	"github.com/juju/errors"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/rpc/params"
)

// Option is a function that can be used to configure a Client.
type Option = base.Option

// WithTracer returns an Option that configures the Client to use the
// supplied tracer.
var WithTracer = base.WithTracer

// Client allows access to the scriptlet charm API endpoint.
type Client struct {
	base.ClientFacade
	facade base.FacadeCaller
}

// NewClient creates a new client for accessing the scriptlet charm API.
func NewClient(st base.APICallCloser, options ...Option) *Client {
	frontend, backend := base.NewClientFacade(st, "ScriptletCharm", options...)
	return &Client{ClientFacade: frontend, facade: backend}
}

// Deploy registers the scriptlet charm and creates the application in one shot.
func (c *Client) Deploy(ctx context.Context, args params.DeployScriptletCharmArgs) error {
	var result params.ErrorResult
	if err := c.facade.FacadeCall(ctx, "Deploy", args, &result); err != nil {
		return errors.Trace(err)
	}
	if result.Error != nil {
		return errors.Trace(result.Error)
	}
	return nil
}
