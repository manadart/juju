// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package asynccharmdownloader

import (
	"context"

	corecharm "github.com/juju/juju/core/charm"
	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/domain/deployment/charm/repository"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/internal/charmhub"
	"github.com/juju/juju/internal/errors"
)

// ModelConfigService provides the model configuration needed to locate
// Charmhub.
type ModelConfigService interface {
	ModelConfig(context.Context) (*config.Config, error)
}

// CharmResolver resolves a Charmhub charm for an application origin.
type CharmResolver interface {
	ResolveForDeploy(context.Context, corecharm.CharmID) (corecharm.ResolvedDataForDeploy, error)
}

// NewCharmResolverFunc constructs a Charmhub resolver.
type NewCharmResolverFunc func(
	context.Context,
	charmhub.HTTPClient,
	ModelConfigService,
	logger.Logger,
) (CharmResolver, error)

// NewCharmResolver constructs a resolver using the model's Charmhub URL.
func NewCharmResolver(
	ctx context.Context,
	httpClient charmhub.HTTPClient,
	modelConfigService ModelConfigService,
	logger logger.Logger,
) (CharmResolver, error) {
	modelConfig, err := modelConfigService.ModelConfig(ctx)
	if err != nil {
		return nil, errors.Errorf("getting model config: %w", err)
	}
	charmhubURL, _ := modelConfig.CharmHubURL()

	return repository.NewCharmHubRepository(repository.CharmHubRepositoryConfig{
		CharmhubHTTPClient: httpClient,
		CharmhubURL:        charmhubURL,
		Logger:             logger,
	})
}
