// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"context"
	"reflect"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/internal/errors"
)

// Register is called to expose a package of facades onto a given registry.
func Register(registry facade.FacadeRegistry) {
	registry.MustRegister("Storage", 6, func(stdCtx context.Context, ctx facade.ModelContext) (facade.Facade, error) {
		return newStorageAPIV6(ctx) // modify Remove to support force and maxWait; add DetachStorage to support force and maxWait.
	}, reflect.TypeOf((*StorageAPIv6)(nil)))

	registry.MustRegister("Storage", 7, func(stdCtx context.Context, ctx facade.ModelContext) (facade.Facade, error) {
		return newStorageAPI(ctx) // support force option on import-fileystem.
	}, reflect.TypeOf((*StorageAPI)(nil)))
}

func newStorageAPIV6(ctx facade.ModelContext) (*StorageAPIv6, error) {
	storageAPI, err := newStorageAPI(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}
	return &StorageAPIv6{
		storageAPI,
	}, nil
}

// newStorageAPI returns a new storage API facade.
func newStorageAPI(ctx facade.ModelContext) (*StorageAPI, error) {
	domainServices := ctx.DomainServices()

	authorizer := ctx.Auth()
	if !authorizer.AuthClient() {
		return nil, apiservererrors.ErrPerm
	}

	storageService := domainServices.Storage()
	return NewStorageAPI(
		ctx.ControllerUUID(),
		ctx.ModelUUID(),
		domainServices.BlockDevice(),
		storageService,
		domainServices.Application(),
		domainServices.Removal(),
		authorizer,
		domainServices.BlockCommand()), nil
}
