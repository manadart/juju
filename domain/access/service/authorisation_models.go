// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"encoding/json"

	"github.com/oklog/ulid/v2"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/storage"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/juju/juju/internal/errors"
)

func (s *AuthorisationService) MaxTypesPerAuthorizationModel() int {
	return storage.DefaultMaxTypesPerAuthorizationModel
}

// WriteAuthorizationModel persists an immutable model under its supplied ID.
func (s *AuthorisationService) WriteAuthorizationModel(ctx context.Context, store string, model *openfgav1.AuthorizationModel) error {
	if len(model.GetTypeDefinitions()) == 0 {
		return nil
	}
	data, err := protojson.Marshal(model)
	if err != nil {
		return err
	}
	return authorisationStorageError(s.st.CreateAuthorisationModel(ctx, ulid.Make().String(), store, model.Id, string(data)))
}

func (s *AuthorisationService) ReadAuthorizationModel(ctx context.Context, store, id string) (*openfgav1.AuthorizationModel, error) {
	data, err := s.st.GetAuthorisationModel(ctx, store, id)
	if err != nil {
		return nil, authorisationStorageError(err)
	}
	model := &openfgav1.AuthorizationModel{}
	if err := protojson.Unmarshal([]byte(data), model); err != nil {
		return nil, err
	}
	if len(model.TypeDefinitions) == 0 {
		return nil, storage.ErrNotFound
	}
	return model, nil
}

func (s *AuthorisationService) ReadAuthorizationModels(ctx context.Context, store string, options storage.ReadAuthorizationModelsOptions) ([]*openfgav1.AuthorizationModel, string, error) {
	if err := validToken(options.Pagination.From); err != nil {
		return nil, "", err
	}
	rows, err := s.st.GetAuthorisationModels(ctx, store, options.Pagination.From, pageSize(options.Pagination)+1)
	if err != nil {
		return nil, "", authorisationStorageError(err)
	}
	result := make([]*openfgav1.AuthorizationModel, 0, len(rows))
	for _, row := range rows {
		model := &openfgav1.AuthorizationModel{}
		if err := protojson.Unmarshal([]byte(row), model); err != nil {
			return nil, "", err
		}
		result = append(result, model)
	}
	var next string
	if n := pageSize(options.Pagination); len(result) > n {
		result = result[:n]
		next = result[n-1].Id
	}
	return result, next, err
}

func (s *AuthorisationService) FindLatestAuthorizationModel(ctx context.Context, store string) (*openfgav1.AuthorizationModel, error) {
	models, _, err := s.ReadAuthorizationModels(ctx, store, storage.ReadAuthorizationModelsOptions{Pagination: storage.NewPaginationOptions(1, "")})
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, storage.ErrNotFound
	}
	return models[0], nil
}

func (s *AuthorisationService) WriteAssertions(ctx context.Context, store, model string, assertions []*openfgav1.Assertion) error {
	data := make([]json.RawMessage, 0, len(assertions))
	for _, assertion := range assertions {
		b, err := protojson.Marshal(assertion)
		if err != nil {
			return err
		}
		data = append(data, b)
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return authorisationStorageError(s.st.SetAuthorisationAssertions(ctx, ulid.Make().String(), store, model, string(b)))
}

func (s *AuthorisationService) ReadAssertions(ctx context.Context, store, model string) ([]*openfgav1.Assertion, error) {
	data, err := s.st.GetAuthorisationAssertions(ctx, store, model)
	err = authorisationStorageError(err)
	if errors.Is(err, storage.ErrNotFound) {
		return []*openfgav1.Assertion{}, nil
	}
	if err != nil {
		return nil, err
	}
	var raw []json.RawMessage
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return nil, err
	}
	result := make([]*openfgav1.Assertion, 0, len(raw))
	for _, data := range raw {
		assertion := &openfgav1.Assertion{}
		if err := protojson.Unmarshal(data, assertion); err != nil {
			return nil, err
		}
		result = append(result, assertion)
	}
	return result, nil
}
