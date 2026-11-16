// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"

	"github.com/canonical/sqlair"

	"github.com/juju/juju/domain/access/internal"
	"github.com/juju/juju/internal/errors"
)

// CreateAuthorisationStore persists store metadata and a replay marker.
func (s *AuthorisationState) CreateAuthorisationStore(ctx context.Context, operation string, store internal.AuthorisationStore) error {
	in := authStore{ID: store.ID, Name: store.Name, CreatedAt: store.CreatedAt, UpdatedAt: store.UpdatedAt}
	stmt, err := s.Prepare(`
INSERT INTO openfga_store (id, name, created_at, updated_at)
VALUES ($authStore.id, $authStore.name, $authStore.created_at, $authStore.updated_at)
`, in)
	if err != nil {
		return errors.Capture(err)
	}
	return s.write(ctx, operation, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt, in).Run()
	})
}

// DeleteAuthorisationStore soft-deletes store metadata.
func (s *AuthorisationState) DeleteAuthorisationStore(ctx context.Context, operation, id, when string) error {
	in := authStoreDeletion{ID: id, When: when}
	stmt, err := s.Prepare(`
UPDATE openfga_store AS s SET deleted_at = $authStoreDeletion.deleted_at
WHERE s.id = $authStoreDeletion.id
`, in)
	if err != nil {
		return errors.Capture(err)
	}
	return s.write(ctx, operation, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt, in).Run()
	})
}

// GetAuthorisationStore retrieves a store that has not been deleted.
func (s *AuthorisationState) GetAuthorisationStore(ctx context.Context, id string) (internal.AuthorisationStore, error) {
	db, err := s.factory(ctx)
	if err != nil {
		return internal.AuthorisationStore{}, errors.Capture(err)
	}
	in := authStore{ID: id}
	stmt, err := s.Prepare(`
SELECT (s.id, s.name, s.created_at, s.updated_at) AS (&authStore.*)
FROM openfga_store AS s WHERE s.id = $authStore.id AND s.deleted_at IS NULL
`, in)
	if err != nil {
		return internal.AuthorisationStore{}, errors.Capture(err)
	}
	var out authStore
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt, in).Get(&out)
	})
	return out.toStore(), authorisationError(err)
}

// GetAuthorisationStores retrieves active stores in ID order after a cursor.
func (s *AuthorisationState) GetAuthorisationStores(ctx context.Context, after string) ([]internal.AuthorisationStore, error) {
	db, err := s.factory(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}
	in := authStore{ID: after}
	stmt, err := s.Prepare(`
SELECT (s.id, s.name, s.created_at, s.updated_at) AS (&authStore.*)
FROM openfga_store AS s WHERE s.deleted_at IS NULL AND s.id > $authStore.id
ORDER BY s.id
`, in)
	if err != nil {
		return nil, errors.Capture(err)
	}
	var rows []authStore
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
	result := make([]internal.AuthorisationStore, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.toStore())
	}
	return result, nil
}

type authStore struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (s authStore) toStore() internal.AuthorisationStore {
	return internal.AuthorisationStore{ID: s.ID, Name: s.Name, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt}
}

type authStoreDeletion struct {
	ID   string `db:"id"`
	When string `db:"deleted_at"`
}
