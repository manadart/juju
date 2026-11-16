// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"time"

	"github.com/juju/clock"
	"github.com/openfga/openfga/pkg/storage"

	"github.com/juju/juju/domain/access/internal"
	"github.com/juju/juju/internal/errors"
)

// AuthorisationState is persistence for the embedded engine's backing service.
// Mutations include a stable operation ID for replay after an ambiguous commit.
type AuthorisationState interface {
	GetAuthorisationReady(context.Context) (bool, error)
	GetAuthorisationTuples(context.Context, string, internal.TupleFilter) ([]internal.Tuple, error)
	WriteAuthorisationTuples(context.Context, string, internal.TupleWrite) error
	GetAuthorisationChanges(context.Context, string, internal.ChangesFilter) ([]internal.TupleChange, error)
	CreateAuthorisationModel(context.Context, string, string, string, string) error
	GetAuthorisationModel(context.Context, string, string) (string, error)
	GetAuthorisationModels(context.Context, string, string, int) ([]string, error)
	SetAuthorisationAssertions(context.Context, string, string, string, string) error
	GetAuthorisationAssertions(context.Context, string, string) (string, error)
	CreateAuthorisationStore(context.Context, string, internal.AuthorisationStore) error
	DeleteAuthorisationStore(context.Context, string, string, string) error
	GetAuthorisationStore(context.Context, string) (internal.AuthorisationStore, error)
	GetAuthorisationStores(context.Context, string) ([]internal.AuthorisationStore, error)
}

// AuthorisationService adapts access persistence to OpenFGA's storage contract.
// It owns no connections or background work; the authorizer worker owns the
// engine which consumes it. Domain mutations still use their own transactions.
type AuthorisationService struct {
	st    AuthorisationState
	clock clock.Clock
}

var _ storage.OpenFGADatastore = (*AuthorisationService)(nil)

// NewAuthorisationService constructs the engine's backing service.
func NewAuthorisationService(st AuthorisationState, clock clock.Clock) *AuthorisationService {
	return &AuthorisationService{st: st, clock: clock}
}

// Close leaves the accessor-owned database open.
func (s *AuthorisationService) Close() {}

// IsReady checks that state can serve the persisted authorisation data.
func (s *AuthorisationService) IsReady(ctx context.Context) (storage.ReadinessStatus, error) {
	ready, err := s.st.GetAuthorisationReady(ctx)
	status := storage.ReadinessStatus{IsReady: ready}
	if !ready && err == nil {
		status.Message = "controller database is shutting down"
	}
	return status, err
}

func authorisationStorageError(err error) error {
	if errors.Is(err, internal.ErrAuthorisationNotFound) {
		return storage.ErrNotFound
	}
	if errors.Is(err, internal.ErrAuthorisationCollision) {
		return errors.Errorf("%w: %w", storage.ErrCollision, err)
	}
	return err
}

func pageSize(p storage.PaginationOptions) int {
	if p.PageSize <= 0 {
		return storage.DefaultPageSize
	}
	return p.PageSize
}

func timestamp(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000000000Z")
}
