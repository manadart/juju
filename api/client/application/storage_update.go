// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"context"

	"github.com/juju/errors"
	"github.com/juju/names/v6"

	"github.com/juju/juju/core/storage"
	"github.com/juju/juju/rpc/params"
)

// SetApplicationStorage updates allocation sizes without resizing instances.
// Version 23 identifies this behaviour and preserves the pool and count.
func (c *Client) SetApplicationStorage(ctx context.Context, name string, updates map[string]storage.DirectiveUpdate) error {
	if c.BestAPIVersion() < 23 {
		return errors.NotSupportedf("allocation-only storage directive updates on this version of Juju")
	}
	if !names.IsValidApplication(name) {
		return errors.NotValidf("application name %q", name)
	}
	directives := make(map[string]params.StorageDirectives, len(updates))
	for name, update := range updates {
		directives[name] = params.StorageDirectives{SizeMiB: update.SizeMiB}
	}
	args := params.ApplicationStorageUpdateRequest{ApplicationStorageUpdates: []params.ApplicationStorageUpdate{{ApplicationTag: names.NewApplicationTag(name).String(), StorageDirectives: directives}}}
	var result params.ErrorResults
	if err := c.facade.FacadeCall(ctx, "UpdateApplicationStorage", args, &result); err != nil {
		return errors.Trace(err)
	}
	return result.OneError()
}
