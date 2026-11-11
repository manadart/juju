// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storageresizer

import (
	"context"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/rpc/params"
)

// Client provides access to resize work for the authenticated machine.
type Client struct {
	facade base.FacadeCaller
}

// NewClient returns a machine resize client. Call Supported before using it
// with a controller that may predate the StorageResizer facade.
func NewClient(caller base.APICaller) *Client {
	return &Client{facade: base.NewFacadeCaller(caller, "StorageResizer")}
}

// Supported reports whether the controller exposes this facade.
func (c *Client) Supported() bool { return c.facade.BestAPIVersion() >= 1 }

// Requests returns pending request IDs for this machine.
func (c *Client) Requests(ctx context.Context) ([]string, error) {
	var result params.StringsResult
	err := c.facade.FacadeCall(ctx, "Requests", nil, &result)
	return result.Result, err
}

// Prepare binds execution and returns a freshly validated target.
func (c *Client) Prepare(ctx context.Context, uuid string) (params.StorageResizeWork, error) {
	var result params.StorageResizeWork
	err := c.facade.FacadeCall(ctx, "Prepare", params.StorageResizeRequestArgs{UUID: uuid}, &result)
	return result, err
}

// Report persists the outcome and any observed capacities.
func (c *Client) Report(ctx context.Context, result params.StorageResizeReport) error {
	return c.facade.FacadeCall(ctx, "Report", result, nil)
}
