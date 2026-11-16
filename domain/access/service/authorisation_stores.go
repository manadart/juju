// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/oklog/ulid/v2"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/storage"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/juju/juju/domain/access/internal"
)

func (s *AuthorisationService) CreateStore(ctx context.Context, store *openfgav1.Store) (*openfgav1.Store, error) {
	if store.GetId() == "" || store.GetName() == "" {
		return nil, fmt.Errorf("store ID and name are required")
	}
	now := s.clock.Now().UTC()
	err := s.st.CreateAuthorisationStore(ctx, ulid.Make().String(), internal.AuthorisationStore{
		ID: store.Id, Name: store.Name, CreatedAt: timestamp(now), UpdatedAt: timestamp(now),
	})
	if err != nil {
		return nil, authorisationStorageError(err)
	}
	return &openfgav1.Store{Id: store.Id, Name: store.Name, CreatedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}, nil
}

func (s *AuthorisationService) DeleteStore(ctx context.Context, id string) error {
	return authorisationStorageError(s.st.DeleteAuthorisationStore(ctx, ulid.Make().String(), id, timestamp(s.clock.Now())))
}

func (s *AuthorisationService) GetStore(ctx context.Context, id string) (*openfgav1.Store, error) {
	row, err := s.st.GetAuthorisationStore(ctx, id)
	if err != nil {
		return nil, authorisationStorageError(err)
	}
	return decodeStore(row.ID, row.Name, row.CreatedAt, row.UpdatedAt)
}

func (s *AuthorisationService) ListStores(ctx context.Context, options storage.ListStoresOptions) ([]*openfgav1.Store, string, error) {
	if err := validToken(options.Pagination.From); err != nil {
		return nil, "", err
	}
	rows, err := s.st.GetAuthorisationStores(ctx, options.Pagination.From)
	if err != nil {
		return nil, "", authorisationStorageError(err)
	}
	var result []*openfgav1.Store
	for _, row := range rows {
		if options.Name != "" && row.Name != options.Name {
			continue
		}
		if len(options.IDs) > 0 && !slices.Contains(options.IDs, row.ID) {
			continue
		}
		store, err := decodeStore(row.ID, row.Name, row.CreatedAt, row.UpdatedAt)
		if err != nil {
			return nil, "", err
		}
		result = append(result, store)
		if len(result) > pageSize(options.Pagination) {
			break
		}
	}
	var next string
	if n := pageSize(options.Pagination); len(result) > n {
		result = result[:n]
		next = result[n-1].Id
	}
	return result, next, err
}

func decodeStore(id, name, created, updated string) (*openfgav1.Store, error) {
	c, err := time.Parse(time.RFC3339Nano, created)
	if err != nil {
		return nil, err
	}
	u, err := time.Parse(time.RFC3339Nano, updated)
	if err != nil {
		return nil, err
	}
	return &openfgav1.Store{Id: id, Name: name, CreatedAt: timestamppb.New(c), UpdatedAt: timestamppb.New(u)}, nil
}
