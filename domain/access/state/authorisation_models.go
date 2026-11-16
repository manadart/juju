// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"

	"github.com/canonical/sqlair"

	"github.com/juju/juju/internal/errors"
)

// CreateAuthorisationModel writes an immutable model and its replay marker.
func (s *AuthorisationState) CreateAuthorisationModel(ctx context.Context, operation, store, id, model string) error {
	in := authModel{Store: store, ID: id, Model: model}
	stmt, err := s.Prepare(`
INSERT INTO openfga_authorization_model (store_id, model_id, model)
VALUES ($authModel.store_id, $authModel.model_id, $authModel.model)
`, in)
	if err != nil {
		return errors.Capture(err)
	}
	return s.write(ctx, operation, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt, in).Run()
	})
}

// GetAuthorisationModel retrieves a model by its immutable ID.
func (s *AuthorisationState) GetAuthorisationModel(ctx context.Context, store, id string) (string, error) {
	db, err := s.factory(ctx)
	if err != nil {
		return "", errors.Capture(err)
	}
	in := authModel{Store: store, ID: id}
	stmt, err := s.Prepare(`
SELECT m.model AS &authModel.model
FROM openfga_authorization_model AS m
WHERE m.store_id = $authModel.store_id AND m.model_id = $authModel.model_id
`, in)
	if err != nil {
		return "", errors.Capture(err)
	}
	var out authModel
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt, in).Get(&out)
	})
	return out.Model, authorisationError(err)
}

// GetAuthorisationModels retrieves a page ordered from newest to oldest.
func (s *AuthorisationState) GetAuthorisationModels(ctx context.Context, store, before string, limit int) ([]string, error) {
	db, err := s.factory(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}
	in := authModelPage{Store: store, Before: before, Limit: limit}
	stmt, err := s.Prepare(`
SELECT m.model AS &authModel.model
FROM openfga_authorization_model AS m
WHERE m.store_id = $authModelPage.store_id
AND ($authModelPage.before = '' OR m.model_id < $authModelPage.before)
ORDER BY m.model_id DESC LIMIT $authModelPage.limit
`, in, authModel{})
	if err != nil {
		return nil, errors.Capture(err)
	}
	var rows []authModel
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, in).GetAll(&rows)
		if errors.Is(err, sqlair.ErrNoRows) {
			return nil
		}
		return err
	})
	if err != nil {
		return nil, errors.Capture(err)
	}
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.Model)
	}
	return result, nil
}

// SetAuthorisationAssertions replaces the assertions for an immutable model.
func (s *AuthorisationState) SetAuthorisationAssertions(ctx context.Context, operation, store, model, assertions string) error {
	in := authAssertions{Store: store, Model: model, Assertions: assertions}
	stmt, err := s.Prepare(`
INSERT INTO openfga_assertion (store_id, model_id, assertions)
VALUES ($authAssertions.store_id, $authAssertions.model_id, $authAssertions.assertions)
ON CONFLICT (store_id, model_id) DO UPDATE SET assertions = excluded.assertions
`, in)
	if err != nil {
		return errors.Capture(err)
	}
	return s.write(ctx, operation, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt, in).Run()
	})
}

// GetAuthorisationAssertions retrieves the serialized assertions for a model.
func (s *AuthorisationState) GetAuthorisationAssertions(ctx context.Context, store, model string) (string, error) {
	db, err := s.factory(ctx)
	if err != nil {
		return "", errors.Capture(err)
	}
	in := authAssertions{Store: store, Model: model}
	stmt, err := s.Prepare(`
SELECT a.assertions AS &authAssertions.assertions
FROM openfga_assertion AS a
WHERE a.store_id = $authAssertions.store_id AND a.model_id = $authAssertions.model_id
`, in)
	if err != nil {
		return "", errors.Capture(err)
	}
	var out authAssertions
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt, in).Get(&out)
	})
	return out.Assertions, authorisationError(err)
}

type authModel struct {
	Store string `db:"store_id"`
	ID    string `db:"model_id"`
	Model string `db:"model"`
}
type authModelPage struct {
	Store  string `db:"store_id"`
	Before string `db:"before"`
	Limit  int    `db:"limit"`
}
type authAssertions struct {
	Store      string `db:"store_id"`
	Model      string `db:"model_id"`
	Assertions string `db:"assertions"`
}
