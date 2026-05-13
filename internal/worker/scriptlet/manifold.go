// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package scriptlet

import (
	"context"

	"github.com/juju/clock"
	jujuerrors "github.com/juju/errors"
	"github.com/juju/worker/v5"
	"github.com/juju/worker/v5/dependency"

	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/internal/services"
)

// ManifoldConfig describes how to create a scriptlet worker.
type ManifoldConfig struct {
	DomainServicesName    string
	ClockName             string
	NewWorker             func(Config) (*Worker, error)
	GetScriptletService   func(dependency.Getter, string) (ScriptletService, error)
	GetApplicationService func(dependency.Getter, string) (ApplicationService, error)
	Logger                logger.Logger
}

// Validate checks the manifold configuration for obvious errors.
func (cfg ManifoldConfig) Validate() error {
	if cfg.DomainServicesName == "" {
		return jujuerrors.NotValidf("empty DomainServicesName")
	}
	if cfg.ClockName == "" {
		return jujuerrors.NotValidf("empty ClockName")
	}
	if cfg.NewWorker == nil {
		return jujuerrors.NotValidf("nil NewWorker")
	}
	if cfg.GetScriptletService == nil {
		return jujuerrors.NotValidf("nil GetScriptletService")
	}
	if cfg.GetApplicationService == nil {
		return jujuerrors.NotValidf("nil GetApplicationService")
	}
	if cfg.Logger == nil {
		return jujuerrors.NotValidf("nil Logger")
	}
	return nil
}

// Manifold returns a dependency.Manifold that runs the scriptlet worker.
func Manifold(cfg ManifoldConfig) dependency.Manifold {
	return dependency.Manifold{
		Inputs: []string{
			cfg.DomainServicesName,
			cfg.ClockName,
		},
		Start: func(ctx context.Context, getter dependency.Getter) (worker.Worker, error) {
			if err := cfg.Validate(); err != nil {
				return nil, errors.Capture(err)
			}

			scriptletService, err := cfg.GetScriptletService(getter, cfg.DomainServicesName)
			if err != nil {
				return nil, errors.Capture(err)
			}

			applicationService, err := cfg.GetApplicationService(getter, cfg.DomainServicesName)
			if err != nil {
				return nil, errors.Capture(err)
			}

			var clk clock.Clock
			if err := getter.Get(cfg.ClockName, &clk); err != nil {
				return nil, errors.Capture(err)
			}

			w, err := cfg.NewWorker(Config{
				ScriptletService:   scriptletService,
				ApplicationService: applicationService,
				Clock:              clk,
				Logger:             cfg.Logger,
			})
			if err != nil {
				return nil, errors.Errorf("creating scriptlet worker: %w", err)
			}
			return w, nil
		},
	}
}

// GetScriptletService extracts the ScriptletService from the
// dependency engine via the DomainServices.
func GetScriptletService(getter dependency.Getter, name string) (ScriptletService, error) {
	var domainServices services.DomainServices
	if err := getter.Get(name, &domainServices); err != nil {
		return nil, errors.Capture(err)
	}
	return domainServices.Scriptlet(), nil
}

// GetApplicationService extracts the ApplicationService from the
// dependency engine via the DomainServices.
func GetApplicationService(getter dependency.Getter, name string) (ApplicationService, error) {
	var domainServices services.DomainServices
	if err := getter.Get(name, &domainServices); err != nil {
		return nil, errors.Capture(err)
	}
	return domainServices.Application(), nil
}

// Config holds the dependencies needed by the scriptlet worker.
type Config struct {
	ScriptletService   ScriptletService
	ApplicationService ApplicationService
	Clock              clock.Clock
	Logger             logger.Logger
}

// Validate checks the worker configuration.
func (c Config) Validate() error {
	if c.ScriptletService == nil {
		return jujuerrors.NotValidf("nil ScriptletService")
	}
	if c.ApplicationService == nil {
		return jujuerrors.NotValidf("nil ApplicationService")
	}
	if c.Clock == nil {
		return jujuerrors.NotValidf("nil Clock")
	}
	if c.Logger == nil {
		return jujuerrors.NotValidf("nil Logger")
	}
	return nil
}
