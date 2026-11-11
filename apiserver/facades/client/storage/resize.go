// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"context"

	"github.com/juju/names/v6"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	coreerrors "github.com/juju/juju/core/errors"
	coremodel "github.com/juju/juju/core/model"
	corestorage "github.com/juju/juju/core/storage"
	domainstorage "github.com/juju/juju/domain/storage"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/rpc/params"
)

// ResizeStorage records resize intent for a storage instance. A CHANGE block
// prevents submission. Recording intent does not perform or approve a resize.
func (a *StorageAPI) ResizeStorage(
	ctx context.Context, args params.StorageResizeArgs,
) (params.StorageResizeResult, error) {
	if err := a.checkCanWrite(ctx); err != nil {
		return params.StorageResizeResult{}, err
	}
	if err := a.blockChecker.ChangeAllowed(ctx); err != nil {
		return params.StorageResizeResult{}, err
	}
	if a.modelType != coremodel.IAAS {
		return params.StorageResizeResult{}, apiservererrors.ParamsErrorf(
			params.CodeNotSupported, "storage resize requires a machine model",
		)
	}
	tag, err := names.ParseStorageTag(args.StorageTag)
	if err != nil {
		return params.StorageResizeResult{}, apiservererrors.ParamsErrorf(
			params.CodeNotValid, "%s", err,
		)
	}
	uuid, err := a.storageService.RequestStorageResize(
		ctx, corestorage.ID(tag.Id()), args.SizeMiB,
	)
	if err != nil {
		var code string
		switch {
		case errors.Is(err, coreerrors.NotValid),
			errors.Is(err, storageerrors.StorageInstanceNotAlive):
			code = params.CodeNotValid
		case errors.Is(err, storageerrors.StorageInstanceNotFound):
			code = params.CodeNotFound
		case errors.Is(err, storageerrors.StorageResizeIntentAlreadyExists):
			code = params.CodeAlreadyExists
		default:
			return params.StorageResizeResult{}, apiservererrors.ServerError(err)
		}
		return params.StorageResizeResult{}, apiservererrors.ParamsErrorf(code, "%s", err)
	}
	return params.StorageResizeResult{UUID: uuid.String()}, nil
}

// StorageResizeDetails returns a recorded request and its current result.
// Results remain readable after storage removal and while changes are blocked.
func (a *StorageAPI) StorageResizeDetails(
	ctx context.Context, args params.StorageResizeRequestArgs,
) (params.StorageResizeDetails, error) {
	if err := a.checkCanRead(ctx); err != nil {
		return params.StorageResizeDetails{}, err
	}
	request, err := a.storageService.GetStorageResizeRequest(
		ctx, domainstorage.StorageResizeUUID(args.UUID),
	)
	if errors.Is(err, storageerrors.StorageResizeIntentNotFound) {
		return params.StorageResizeDetails{}, apiservererrors.ParamsErrorf(
			params.CodeNotFound, "%s", err,
		)
	} else if err != nil {
		return params.StorageResizeDetails{}, apiservererrors.ServerError(err)
	}
	return params.StorageResizeDetails{
		UUID:             request.UUID.String(),
		StorageTag:       names.NewStorageTag(request.StorageID.String()).String(),
		RequestedSizeMiB: request.RequestedSizeMiB,
		Status:           string(request.Status),
		Message:          request.Message,
		CreatedAt:        request.CreatedAt,
		CompletedAt:      request.CompletedAt,
	}, nil
}

// ResizeStorage is unavailable in version 6.
func (*StorageAPIv6) ResizeStorage(_, _, _ struct{}) {}

// ResizeStorage is unavailable in version 7.
func (*StorageAPIv7) ResizeStorage(_, _, _ struct{}) {}

// StorageResizeDetails is unavailable in version 6.
func (*StorageAPIv6) StorageResizeDetails(_, _, _ struct{}) {}

// StorageResizeDetails is unavailable in version 7.
func (*StorageAPIv7) StorageResizeDetails(_, _, _ struct{}) {}

// StorageResizeDetails is unavailable in version 8.
func (*StorageAPIv8) StorageResizeDetails(_, _, _ struct{}) {}
