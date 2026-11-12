// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"context"

	"github.com/juju/names/v6"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	coremodel "github.com/juju/juju/core/model"
	coreunit "github.com/juju/juju/core/unit"
	"github.com/juju/juju/rpc/params"
)

// ResizeUnitStorage stages requests using the current application directives.
func (a *StorageAPI) ResizeUnitStorage(ctx context.Context, args params.UnitStorageResizeArgs) (params.UnitStorageResizeResults, error) {
	if err := a.checkCanWrite(ctx); err != nil {
		return params.UnitStorageResizeResults{}, err
	}
	if err := a.blockChecker.ChangeAllowed(ctx); err != nil {
		return params.UnitStorageResizeResults{}, err
	}
	if a.modelType != coremodel.IAAS {
		return params.UnitStorageResizeResults{}, apiservererrors.ParamsErrorf(params.CodeNotSupported, "storage resize requires a machine model")
	}
	tag, err := names.ParseUnitTag(args.UnitTag)
	if err != nil {
		return params.UnitStorageResizeResults{}, apiservererrors.ParamsErrorf(params.CodeNotValid, "%s", err)
	}
	results, err := a.storageService.RequestUnitStorageResize(ctx, coreunit.Name(tag.Id()))
	if err != nil {
		return params.UnitStorageResizeResults{}, apiservererrors.ServerError(err)
	}
	out := params.UnitStorageResizeResults{Results: make([]params.UnitStorageResizeResult, 0, len(results))}
	for _, result := range results {
		out.Results = append(out.Results, params.UnitStorageResizeResult{StorageTag: names.NewStorageTag(result.StorageID.String()).String(), UUID: result.UUID.String(), RequestedSizeMiB: result.RequestedSizeMiB, AlreadySatisfied: result.AlreadySatisfied})
	}
	return out, nil
}

// ResizeUnitStorage is unavailable before version 10.
func (*StorageAPIv6) ResizeUnitStorage(_, _, _ struct{}) {}
func (*StorageAPIv7) ResizeUnitStorage(_, _, _ struct{}) {}
func (*StorageAPIv8) ResizeUnitStorage(_, _, _ struct{}) {}
func (*StorageAPIv9) ResizeUnitStorage(_, _, _ struct{}) {}
