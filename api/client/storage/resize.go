// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"context"

	"github.com/juju/errors"
	"github.com/juju/names/v6"

	"github.com/juju/juju/rpc/params"
)

// ResizeStorage records a resize intent and returns its UUID. The size is the
// requested target capacity in MiB. Acceptance does not mean storage has grown.
func (c *Client) ResizeStorage(
	ctx context.Context, storageID string, sizeMiB uint64,
) (string, error) {
	if c.BestAPIVersion() < 8 {
		return "", errors.NotSupportedf("storage resize on this version of Juju")
	}
	if !names.IsValidStorage(storageID) {
		return "", errors.NotValidf("storage ID %q", storageID)
	}
	args := params.StorageResizeArgs{
		StorageTag: names.NewStorageTag(storageID).String(),
		SizeMiB:    sizeMiB,
	}
	var result params.StorageResizeResult
	if err := c.facade.FacadeCall(ctx, "ResizeStorage", args, &result); err != nil {
		return "", errors.Trace(err)
	}
	return result.UUID, nil
}
