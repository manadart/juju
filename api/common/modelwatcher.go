// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common

import (
	"context"
	"time"

	"github.com/juju/errors"

	"github.com/juju/juju/api/base"
	apiwatcher "github.com/juju/juju/api/watcher"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/rpc/params"
)

// ModelConfigWatcher provides common client-side API functions
// to call into apiserver.common.ModelConfigWatcher.
type ModelConfigWatcher struct {
	facade base.FacadeCaller
}

// NewModelConfigWatcher creates a ModelConfigWatcher on the specified facade,
// and uses this name when calling through the caller.
func NewModelConfigWatcher(facade base.FacadeCaller) *ModelConfigWatcher {
	return &ModelConfigWatcher{facade}
}

// WatchForModelConfigChanges return a NotifyWatcher waiting for the
// model configuration to change.
func (e *ModelConfigWatcher) WatchForModelConfigChanges(ctx context.Context) (watcher.NotifyWatcher, error) {
	var result params.NotifyWatchResult
	err := e.facade.FacadeCall(ctx, "WatchForModelConfigChanges", nil, &result)
	if err != nil {
		return nil, err
	}
	return apiwatcher.NewNotifyWatcher(e.facade.RawAPICaller(), result), nil
}

// ModelConfig returns the current model configuration.
func (e *ModelConfigWatcher) ModelConfig(ctx context.Context) (*config.Config, error) {
	var result params.ModelConfigResult
	err := e.facade.FacadeCall(ctx, "ModelConfig", nil, &result)
	if err != nil {
		return nil, errors.Trace(err)
	}
	conf, err := config.New(config.NoDefaults, result.Config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return conf, nil
}

// UpdateStatusHookInterval returns the current update status hook interval.
func (e *ModelConfigWatcher) UpdateStatusHookInterval(ctx context.Context) (time.Duration, error) {
	// TODO(wallyworld) - lp:1602237 - this needs to have it's own backend implementation.
	// For now, we'll piggyback off the ModelConfig API.
	modelConfig, err := e.ModelConfig(ctx)
	if err != nil {
		return 0, err
	}
	return modelConfig.UpdateStatusHookInterval(), nil
}

// WatchUpdateStatusHookInterval returns a NotifyWatcher that fires when the
// update status hook interval changes.
func (e *ModelConfigWatcher) WatchUpdateStatusHookInterval(ctx context.Context) (watcher.NotifyWatcher, error) {
	// TODO(wallyworld) - lp:1602237 - this needs to have it's own backend implementation.
	// For now, we'll piggyback off the ModelConfig API.
	return e.WatchForModelConfigChanges(ctx)
}
