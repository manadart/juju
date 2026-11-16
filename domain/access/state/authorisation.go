// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"

	"github.com/canonical/sqlair"

	"github.com/juju/juju/core/database"
	"github.com/juju/juju/domain"
	"github.com/juju/juju/domain/access/internal"
	internaldatabase "github.com/juju/juju/internal/database"
	"github.com/juju/juju/internal/errors"
)

// AuthorisationState persists the backing data for the authorisation engine.
// Connections, retries and database lifetime remain with the database accessor.
type AuthorisationState struct {
	*domain.StateBase
	// Use the accessor directly so errors can be translated after its
	// retries, before StateBase masks database errors.
	factory database.TxnRunnerFactory
}

// NewAuthorisationState supplies persistence without starting background work.
func NewAuthorisationState(factory database.TxnRunnerFactory) *AuthorisationState {
	return &AuthorisationState{StateBase: domain.NewStateBase(factory), factory: factory}
}

// GetAuthorisationReady checks the database lifetime and installed schema.
func (s *AuthorisationState) GetAuthorisationReady(ctx context.Context) (bool, error) {
	db, err := s.factory(ctx)
	if err != nil {
		return false, errors.Capture(err)
	}
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	case <-db.Dying():
		return false, nil
	default:
	}
	stmt, err := s.Prepare(`SELECT COUNT(*) AS &authCount.count FROM openfga_store AS s`, authCount{})
	if err != nil {
		return false, errors.Capture(err)
	}
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		var count authCount
		return tx.Query(ctx, stmt).Get(&count)
	})
	return err == nil, errors.Capture(err)
}

// write commits a replay marker together with the mutation. The operation ID
// is supplied by the service and is stable across all accessor retries.
func (s *AuthorisationState) write(ctx context.Context, operation string, fn func(context.Context, *sqlair.TX) error) error {
	db, err := s.factory(ctx)
	if err != nil {
		return errors.Capture(err)
	}
	op := authOperation{ID: operation}
	get, err := s.Prepare(`SELECT o.ulid AS &authOperation.ulid FROM openfga_operation AS o WHERE o.ulid = $authOperation.ulid`, op)
	if err != nil {
		return errors.Capture(err)
	}
	insert, err := s.Prepare(`INSERT INTO openfga_operation (ulid) VALUES ($authOperation.ulid)`, op)
	if err != nil {
		return errors.Capture(err)
	}
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		var found authOperation
		err := tx.Query(ctx, get, op).Get(&found)
		if err == nil {
			return nil
		}
		if !errors.Is(err, sqlair.ErrNoRows) {
			return err
		}
		if err := fn(ctx, tx); err != nil {
			return err
		}
		return tx.Query(ctx, insert, op).Run()
	})
	return authorisationError(err)
}

// Translate errors only after the accessor has finished retrying.
func authorisationError(err error) error {
	if errors.Is(err, sqlair.ErrNoRows) {
		return internal.ErrAuthorisationNotFound
	}
	if internaldatabase.IsErrConstraintUnique(err) || internaldatabase.IsErrConstraintPrimaryKey(err) {
		return errors.Errorf("%w: %w", internal.ErrAuthorisationCollision, err)
	}
	return errors.Capture(err)
}

type authOperation struct {
	ID string `db:"ulid"`
}
type authCount struct {
	Count int `db:"count"`
}
