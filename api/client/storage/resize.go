// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"context"
	"strconv"

	"github.com/juju/errors"
	"github.com/juju/names/v6"
	"github.com/juju/utils/v4"

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

// SupportsStorageResizeResults reports whether results can be retrieved.
// Callers can check this before submitting a request they intend to wait for.
func (c *Client) SupportsStorageResizeResults() bool {
	return c.BestAPIVersion() >= 9
}

// StorageResizeDetails returns a request's current status and recorded result.
func (c *Client) StorageResizeDetails(
	ctx context.Context, uuid string,
) (params.StorageResizeDetails, error) {
	if !c.SupportsStorageResizeResults() {
		return params.StorageResizeDetails{}, errors.NotSupportedf("storage resize results on this version of Juju")
	}
	if !utils.IsValidUUIDString(uuid) {
		return params.StorageResizeDetails{}, errors.NotValidf("resize request ID %q", uuid)
	}
	var result params.StorageResizeDetails
	if err := c.facade.FacadeCall(ctx, "StorageResizeDetails", params.StorageResizeRequestArgs{UUID: uuid}, &result); err != nil {
		return params.StorageResizeDetails{}, errors.Trace(err)
	}
	return result, nil
}

// ResizeUnitStorage stages requests for a unit using application directives.
// Integer request IDs require version 11, checked before submitting any work.
func (c *Client) ResizeUnitStorage(ctx context.Context, unitName string) ([]params.UnitStorageResizeResult, error) {
	if c.BestAPIVersion() < 11 {
		return nil, errors.NotSupportedf("unit storage resize on this version of Juju")
	}
	if !names.IsValidUnit(unitName) {
		return nil, errors.NotValidf("unit name %q", unitName)
	}
	var result params.UnitStorageResizeResults
	if err := c.facade.FacadeCall(ctx, "ResizeUnitStorage", params.UnitStorageResizeArgs{UnitTag: names.NewUnitTag(unitName).String()}, &result); err != nil {
		return nil, errors.Trace(err)
	}
	return result.Results, nil
}

// StorageResizeDetailsByID returns a result using its model-scoped integer ID.
func (c *Client) StorageResizeDetailsByID(ctx context.Context, id string) (params.StorageResizeDetails, error) {
	if c.BestAPIVersion() < 11 {
		return params.StorageResizeDetails{}, errors.NotSupportedf("integer storage resize request IDs on this version of Juju")
	}
	if _, err := strconv.ParseUint(id, 10, 63); err != nil {
		return params.StorageResizeDetails{}, errors.NotValidf("resize request ID %q", id)
	}
	var result params.StorageResizeDetails
	if err := c.facade.FacadeCall(ctx, "StorageResizeDetailsByID", params.StorageResizeRequestIDArgs{ID: id}, &result); err != nil {
		return params.StorageResizeDetails{}, errors.Trace(err)
	}
	return result, nil
}
