// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasmodelconfigmanager

import (
	"context"

	"github.com/juju/errors"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/common"
	apiwatcher "github.com/juju/juju/api/watcher"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/rpc/params"
)

// Option is a function that can be used to configure a Client.
type Option = base.Option

// WithTracer returns an Option that configures the Client to use the
// supplied tracer.
var WithTracer = base.WithTracer

// Client allows access to the CAAS model config manager API endpoint.
type Client struct {
	facade base.FacadeCaller
	*common.ControllerConfigAPI
}

// NewClient returns a client used to access the CAAS Application Provisioner API.
func NewClient(caller base.APICaller, options ...Option) (*Client, error) {
	_, isModel := caller.ModelTag()
	if !isModel {
		return nil, errors.New("expected model specific API connection")
	}
	facadeCaller := base.NewFacadeCaller(caller, "CAASModelConfigManager", options...)
	return &Client{
		facade:              facadeCaller,
		ControllerConfigAPI: common.NewControllerConfig(facadeCaller),
	}, nil
}

// WatchControllerConfig provides a watcher for changes on controller config.
func (c *Client) WatchControllerConfig(ctx context.Context) (watcher.StringsWatcher, error) {
	var result params.StringsWatchResult
	if err := c.facade.FacadeCall(ctx, "WatchControllerConfig", nil, &result); err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return apiwatcher.NewStringsWatcher(c.facade.RawAPICaller(), result), nil
}
