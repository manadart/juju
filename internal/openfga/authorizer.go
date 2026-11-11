// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package openfga

import (
	"context"
	"strings"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/server"
	"github.com/openfga/openfga/pkg/typesystem"

	"github.com/juju/juju/core/database"
	"github.com/juju/juju/core/permission"
)

// Authorizer owns the embedded server. Its owner must cancel its context and
// call Close before releasing it. No tuple or decision cache spans requests.
type Authorizer struct {
	server *server.Server
}

// NewAuthorizer validates the persisted model before starting the server.
func NewAuthorizer(ctx context.Context, factory database.TxnRunnerFactory) (*Authorizer, error) {
	ds := NewDatastore(factory)
	model, err := ds.ReadAuthorizationModel(ctx, StoreID, ModelID)
	if err != nil {
		return nil, err
	}
	if _, err := typesystem.NewAndValidate(ctx, model); err != nil {
		return nil, err
	}
	s, err := server.NewServerWithOpts(
		server.WithContext(ctx),
		server.WithDatastore(ds),
		server.WithContextPropagationToDatastore(true),
		server.WithCheckQueryCacheEnabled(false),
		server.WithCheckIteratorCacheEnabled(false),
		server.WithListObjectsIteratorCacheEnabled(false),
		server.WithCacheControllerEnabled(false),
		server.WithSharedIteratorEnabled(false),
		server.WithDispatchThrottlingCheckResolverEnabled(false),
		server.WithListObjectsDispatchThrottlingEnabled(false),
		server.WithListUsersDispatchThrottlingEnabled(false),
	)
	if err != nil {
		return nil, err
	}
	return &Authorizer{server: s}, nil
}

// Close releases the server's managed background resources.
func (a *Authorizer) Close() { a.server.Close() }

// Check evaluates the persisted Juju authorisation model with fresh tuples.
func (a *Authorizer) Check(ctx context.Context, subject string, target permission.ID, access permission.Access) (bool, error) {
	result, err := a.server.Check(ctx, &openfgav1.CheckRequest{
		StoreId:              StoreID,
		AuthorizationModelId: ModelID,
		TupleKey: &openfgav1.CheckRequestTupleKey{
			User:     subject,
			Object:   target.ObjectType.String() + ":" + target.Key,
			Relation: "can_" + strings.ReplaceAll(access.String(), "-", "_"),
		},
		Consistency: openfgav1.ConsistencyPreference_HIGHER_CONSISTENCY,
	})
	if err != nil {
		return false, err
	}
	return result.Allowed, nil
}
