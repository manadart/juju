// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"math"

	coreapplication "github.com/juju/juju/core/application"
	coreerrors "github.com/juju/juju/core/errors"
	"github.com/juju/juju/core/model"
	corestorage "github.com/juju/juju/core/storage"
	applicationerrors "github.com/juju/juju/domain/application/errors"
	"github.com/juju/juju/domain/application/internal"
	domainstorage "github.com/juju/juju/domain/storage"
	"github.com/juju/juju/internal/errors"
)

// ApplicationStorageState persists validated updates against a stable snapshot.
type ApplicationStorageState interface {
	GetApplicationStorageSnapshot(context.Context, string) (internal.ApplicationStorageSnapshot, error)
	SetApplicationStorageDirectives(context.Context, string, internal.ApplicationStorageSnapshot, []domainstorage.DirectiveArg) error
}

// UpdateApplicationStorageDirectives updates allocation defaults without
// resizing existing instances or changing existing unit directive copies.
// Omitted fields retain their values. Concurrent edits are merged on retry.
func (s *ProviderService) UpdateApplicationStorageDirectives(ctx context.Context, uuid coreapplication.UUID, overrides map[string]StorageDirectiveOverrides) error {
	if err := uuid.Validate(); err != nil {
		return errors.Errorf("%w", err).Add(coreerrors.NotValid)
	}
	modelType, err := s.st.GetModelType(ctx)
	if err != nil {
		return errors.Capture(err)
	}
	if modelType != model.IAAS {
		return errors.New("application storage updates require a machine model in this prototype").Add(coreerrors.NotSupported)
	}
	if len(overrides) == 0 {
		return errors.New("no storage directive updates supplied").Add(coreerrors.NotValid)
	}
	for name, override := range overrides {
		if err := corestorage.Name(name).Validate(); err != nil {
			return errors.Errorf("%w", err).Add(coreerrors.NotValid)
		}
		if override.Size != nil && (*override.Size == 0 || *override.Size > math.MaxInt64/(1024*1024)) {
			return errors.Errorf("storage %q size is outside the supported range", name).Add(coreerrors.NotValid)
		}
		if override.PoolUUID != nil {
			if err := override.PoolUUID.Validate(); err != nil {
				return errors.Errorf("%w", err).Add(coreerrors.NotValid)
			}
		}
		if override.Size == nil && override.Count == nil && override.PoolUUID == nil {
			return errors.Errorf("storage %q has no fields to update", name).Add(coreerrors.NotValid)
		}
	}
	for range 3 {
		snapshot, err := s.st.GetApplicationStorageSnapshot(ctx, uuid.String())
		if err != nil {
			return errors.Capture(err)
		}
		definitions := make(map[string]internal.CharmStorageDefinitionForValidation, len(snapshot.Directives))
		for _, directive := range snapshot.Directives {
			definitions[directive.Name.String()] = directive.Definition
		}
		for name := range overrides {
			if _, ok := definitions[name]; !ok {
				return errors.Errorf("storage %q: %w", name, applicationerrors.StorageNameNotSupported)
			}
		}
		if err := s.storageService.ValidateApplicationStorageDirectiveOverrides(ctx, definitions, overrides); err != nil {
			return errors.Errorf("validating storage directives: %w", err).Add(coreerrors.NotValid)
		}
		updates := make([]domainstorage.DirectiveArg, 0, len(overrides))
		for _, directive := range snapshot.Directives {
			override, ok := overrides[directive.Name.String()]
			if !ok {
				continue
			}
			arg := directive.DirectiveArg
			if override.Size != nil {
				arg.Size = *override.Size
			}
			if override.Count != nil {
				arg.Count = *override.Count
			}
			if override.PoolUUID != nil {
				arg.PoolUUID = *override.PoolUUID
			}
			updates = append(updates, arg)
		}
		err = s.st.SetApplicationStorageDirectives(ctx, uuid.String(), snapshot, updates)
		if errors.Is(err, applicationerrors.StorageDirectivesChanged) {
			continue
		}
		return errors.Capture(err)
	}
	return applicationerrors.StorageDirectivesChanged
}
