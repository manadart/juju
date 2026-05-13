// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package scriptlet

import (
	"context"
	"time"

	"github.com/canonical/starform/starform"
	"github.com/canonical/starlark/starlark"
	"github.com/juju/errors"
	"github.com/juju/worker/v5"
	"github.com/juju/worker/v5/catacomb"

	"github.com/juju/juju/core/application"
	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/domain/deployment/charm"
	internalerrors "github.com/juju/juju/internal/errors"
	internalworker "github.com/juju/juju/internal/worker"
)

const (
	// States which report the state of the worker.
	stateStarted = "started"
)

const setStatusSafety = starlark.NotSafe

type StatusUpdate struct {
	Status  string
	Message string
}

func (u *StatusUpdate) statusUpdate() *StatusUpdate {
	return u
}

type StateUpdate struct {
	Name  string
	Value starlark.Value
}

type StateUpdates struct {
	Updates []StateUpdate
}

func (u *StateUpdates) stateUpdates() *StateUpdates {
	return u
}

type OnConfigChangedState struct {
	StatusUpdate
	StateUpdates
}

type OnRelationCreatedState struct {
	StatusUpdate
	StateUpdates
}

type OnRelationJoinedState struct {
	StatusUpdate
	StateUpdates
}

type OnRelationChangedState struct {
	StatusUpdate
	StateUpdates
}

type OnRelationDepartedState struct {
	StatusUpdate
	StateUpdates
}

type OnRelationBrokenState struct {
	StatusUpdate
	StateUpdates
}

// Worker is a controller-side worker that watches for scriptlet
// applications and dispatches model events to Starform scriptlets.
// It uses a worker.Runner to manage per-application child workers.
type Worker struct {
	internalStates chan string
	catacomb       catacomb.Catacomb
	config         Config

	runner *worker.Runner
}

// NewWorker returns a new scriptlet worker.
func NewWorker(config Config) (*Worker, error) {
	return newWorker(config, nil)
}

func newWorker(config Config, internalStates chan string) (*Worker, error) {
	if err := config.Validate(); err != nil {
		return nil, internalerrors.Capture(err)
	}

	runner, err := worker.NewRunner(worker.RunnerParams{
		Name:  "scriptlet",
		Clock: config.Clock,
		IsFatal: func(err error) bool {
			return false
		},
		ShouldRestart: internalworker.ShouldRunnerRestart,
		RestartDelay:  time.Second * 10,
		Logger:        internalworker.WrapLogger(config.Logger),
	})
	if err != nil {
		return nil, internalerrors.Capture(err)
	}

	w := &Worker{
		internalStates: internalStates,
		config:         config,
		runner:         runner,
	}

	if err := catacomb.Invoke(catacomb.Plan{
		Name: "scriptlet",
		Site: &w.catacomb,
		Work: w.loop,
		Init: []worker.Worker{
			w.runner,
		},
	}); err != nil {
		return nil, internalerrors.Capture(err)
	}

	return w, nil
}

// Kill is part of the worker.Worker interface.
func (w *Worker) Kill() {
	w.catacomb.Kill(nil)
}

// Wait is part of the worker.Worker interface.
func (w *Worker) Wait() error {
	return w.catacomb.Wait()
}

func (w *Worker) loop() error {
	ctx, cancel := w.scopedContext()
	defer cancel()

	logger := w.config.Logger

	logger.Infof(ctx, "starting scriptlet worker")

	// Watch for scriptlet application changes.
	appWatcher, err := w.config.ScriptletService.WatchScriptletApplications(ctx)
	if err != nil {
		return internalerrors.Errorf("watching scriptlet applications: %w", err)
	}
	if err := w.catacomb.Add(appWatcher); err != nil {
		return internalerrors.Capture(err)
	}

	w.reportInternalState(stateStarted)

	for {
		select {
		case <-w.catacomb.Dying():
			return w.catacomb.ErrDying()

		case appUUIDs, ok := <-appWatcher.Changes():
			if !ok {
				return internalerrors.New("application watcher closed")
			}
			if err := w.handleApplicationChanges(ctx, appUUIDs); err != nil {
				return internalerrors.Errorf("handling application changes: %w", err)
			}
		}
	}
}

// handleApplicationChanges processes changes to scriptlet applications,
// starting per-application workers via the Runner. appUUIDs are the
// UUIDs of applications whose charm is a scriptlet.
func (w *Worker) handleApplicationChanges(ctx context.Context, appUUIDs []string) error {
	logger := w.config.Logger

	for _, appUUID := range appUUIDs {
		logger.Infof(ctx, "ensuring scriptlet runner for application %q", appUUID)

		id := application.UUID(appUUID)
		err := w.runner.StartWorker(ctx, appUUID, func(ctx context.Context) (worker.Worker, error) {
			return newApplicationRunner(applicationRunnerConfig{
				AppUUID:            id,
				ScriptletService:   w.config.ScriptletService,
				ApplicationService: w.config.ApplicationService,
				RelationService:    w.config.RelationService,
				Logger:             w.config.Logger,
			})
		})
		if errors.Is(err, errors.AlreadyExists) {
			continue
		}
		if err != nil {
			return internalerrors.Errorf("starting runner for %q: %w", appUUID, err)
		}
	}

	return nil
}

func (w *Worker) scopedContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(w.catacomb.Context(context.Background()))
}

func (w *Worker) reportInternalState(state string) {
	if w.internalStates != nil {
		select {
		case w.internalStates <- state:
		default:
		}
	}
}

// applicationRunner handles events for a single scriptlet application.
// It is managed by the parent Worker's Runner.
type applicationRunner struct {
	catacomb catacomb.Catacomb
	config   applicationRunnerConfig

	scriptSet *starform.ScriptSet
}

type applicationRunnerConfig struct {
	AppUUID            application.UUID
	ScriptletService   ScriptletService
	ApplicationService ApplicationService
	RelationService    RelationService
	Logger             logger.Logger
}

func newApplicationRunner(config applicationRunnerConfig) (*applicationRunner, error) {
	r := &applicationRunner{
		config: config,
	}

	if err := catacomb.Invoke(catacomb.Plan{
		Name: "scriptlet-" + config.AppUUID.String(),
		Site: &r.catacomb,
		Work: r.loop,
	}); err != nil {
		return nil, internalerrors.Capture(err)
	}

	return r, nil
}

// Kill is part of the worker.Worker interface.
func (r *applicationRunner) Kill() {
	r.catacomb.Kill(nil)
}

// Wait is part of the worker.Worker interface.
func (r *applicationRunner) Wait() error {
	return r.catacomb.Wait()
}

func (r *applicationRunner) loop() error {
	ctx, cancel := r.scopedContext()
	defer cancel()

	logger := r.config.Logger
	appUUID := r.config.AppUUID

	logger.Infof(ctx, "scriptlet runner started for %q", appUUID)

	// Load the scriptlet source from the database.
	scriptletSrc, err := r.config.ScriptletService.GetApplicationScriptlet(ctx, appUUID)
	if err != nil {
		return internalerrors.Errorf("loading scriptlet for %q: %w", appUUID, err)
	}

	scriptSet, err := starform.NewScriptSet(&starform.ScriptSetOptions{
		App: &starform.AppObject{
			Name: "juju",
		},
	})
	if err != nil {
		return internalerrors.Errorf("creating script set: %w", err)
	}
	if err := scriptSet.LoadSources(ctx, []starform.ScriptSource{
		&dbScriptSource{content: scriptletSrc},
	}); err != nil {
		return internalerrors.Errorf("loading script set sources: %w", err)
	}
	r.scriptSet = scriptSet

	// Watch for config changes on this application.
	configWatcher, err := r.config.ApplicationService.WatchApplicationConfigChangeByUUID(ctx, appUUID)
	if err != nil {
		return internalerrors.Errorf("watching config for %q: %w", appUUID, err)
	}
	if err := r.catacomb.Add(configWatcher); err != nil {
		return internalerrors.Capture(err)
	}

	// Watch for relation lifecycle changes on this application.
	relationWatcher, err := r.config.RelationService.WatchRelationsLifeSuspendedStatusForApplication(ctx, appUUID)
	if err != nil {
		return internalerrors.Errorf("watching relations for %q: %w", appUUID, err)
	}
	if err := r.catacomb.Add(relationWatcher); err != nil {
		return internalerrors.Capture(err)
	}

	// Track known relations for detecting created/broken.
	knownRelations := make(map[string]relationState)

	for {
		select {
		case <-r.catacomb.Dying():
			return r.catacomb.ErrDying()

		case _, ok := <-configWatcher.Changes():
			if !ok {
				return internalerrors.New("config watcher closed")
			}
			if err := r.handleConfigChanged(ctx, appUUID); err != nil {
				logger.Errorf(ctx, "handling config-changed for %q: %v", appUUID, err)
			}

		case changes, ok := <-relationWatcher.Changes():
			if !ok {
				return internalerrors.New("relation watcher closed")
			}
			if err := r.handleRelationChanges(ctx, appUUID, changes, knownRelations); err != nil {
				logger.Errorf(ctx, "handling relation changes for %q: %v", appUUID, err)
			}
		}
	}
}

func (r *applicationRunner) handleConfigChanged(ctx context.Context, appUUID application.UUID) error {
	logger := r.config.Logger
	logger.Debugf(ctx, "dispatching config_changed for %q", appUUID)

	// Query the current application config.
	config, err := r.config.ApplicationService.GetApplicationConfigWithDefaults(ctx, appUUID)
	if err != nil {
		return internalerrors.Errorf("getting config for %q: %w", appUUID, err)
	}

	// TODO(hackathon): Serialise config and pass to starform scriptlet.
	_ = config

	event := starform.EventObject{
		Name:  "config_changed",
		// TODO(kcza): complete this
		State: &OnConfigChangedState{},
	}
	return r.scriptSet.Handle(ctx, &event)
}

func (r *applicationRunner) handleRelationChanges(
	ctx context.Context,
	appUUID application.UUID,
	relationUUIDs []string,
	knownRelations map[string]relationState,
) error {
	logger := r.config.Logger
	logger.Debugf(ctx, "relation changes for %q: %v", appUUID, relationUUIDs)

	// For each relation UUID in the change set, determine if it's new
	// (relation-created) or removed (relation-broken).
	seen := make(map[string]bool, len(relationUUIDs))
	for _, relUUID := range relationUUIDs {
		seen[relUUID] = true

		if _, known := knownRelations[relUUID]; !known {
			// New relation — dispatch relation-created.
			knownRelations[relUUID] = relationState{
				RelationUUID: relUUID,
				Members:      make(map[string]bool),
			}

			event := starform.EventObject{
				Name: "relation_created",
			}
			if err := r.scriptSet.Handle(ctx, &event); err != nil {
				logger.Errorf(ctx, "dispatching relation_created for %q rel %s: %v",
					appUUID, relUUID, err)
			}

			// Also fire relation-joined for the application itself
			// (the scriptlet app "joins" the relation).
			joinEvent := starform.EventObject{
				Name: "relation_joined",
			}
			if err := r.scriptSet.Handle(ctx, &joinEvent); err != nil {
				logger.Errorf(ctx, "dispatching relation_joined for %q rel %s: %v",
					appUUID, relUUID, err)
			}
		} else {
			// Existing relation changed — dispatch relation-changed.
			event := starform.EventObject{
				Name: "relation_changed",
			}
			if err := r.scriptSet.Handle(ctx, &event); err != nil {
				logger.Errorf(ctx, "dispatching relation_changed for %q rel %s: %v",
					appUUID, relUUID, err)
			}
		}
	}

	// Check for relations that have been removed (broken).
	for relUUID := range knownRelations {
		if !seen[relUUID] {
			// Relation is gone — dispatch relation-departed then
			// relation-broken.
			departEvent := starform.EventObject{
				Name: "relation_departed",
			}
			if err := r.scriptSet.Handle(ctx, &departEvent); err != nil {
				logger.Errorf(ctx, "dispatching relation_departed for %q rel %s: %v",
					appUUID, relUUID, err)
			}

			brokenEvent := starform.EventObject{
				Name: "relation_broken",
			}
			if err := r.scriptSet.Handle(ctx, &brokenEvent); err != nil {
				logger.Errorf(ctx, "dispatching relation_broken for %q rel %s: %v",
					appUUID, relUUID, err)
			}

			delete(knownRelations, relUUID)
		}
	}

	event := starform.EventObject{
		Name:  "relation_created",
		// TODO(kcza): complete this
		State: &OnRelationCreatedState{},
	}
	return r.scriptSet.Handle(ctx, &event)
}

func jujuSetStatus(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	if err := starlark.CheckSafety(thread, starlark.NotSafe); err != nil {
		return nil, err
	}

	var status string
	var message string
	if err := starlark.UnpackArgs(
		fn.Name(), args, kwargs,
		"status", &status,
		"message?", &message,
	); err != nil {
		return nil, err
	}

	event := starform.Event(thread)
	update := StatusUpdate{
		Status:  status,
		Message: message,
	}
	if err := thread.AddAllocs(starlark.EstimateSize(update)); err != nil {
		return nil, err
	}

	state, ok := event.State.(interface{ statusUpdate() *StatusUpdate })
	if !ok || state == nil || state.statusUpdate() == nil {
		return nil, starform.ErrUnavailable
	}
	*state.statusUpdate() = update

	return starlark.None, nil
}

func jujuSetState(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	if err := starlark.CheckSafety(thread, starlark.NotSafe); err != nil {
		return nil, err
	}

	var name string
	var value starlark.Value
	if err := starlark.UnpackArgs(
		fn.Name(), args, kwargs,
		"name", &name,
		"value", &value,
	); err != nil {
		return nil, err
	}

	event := starform.Event(thread)
	update := StateUpdate{
		Name:  name,
		Value: value,
	}

	state, ok := event.State.(interface{ stateUpdates() *StateUpdates })
	if !ok || state == nil || state.stateUpdates() == nil {
		return nil, starform.ErrUnavailable
	}

	appender := starlark.NewSafeAppender(thread, &state.stateUpdates().Updates)
	if err := appender.Append(update); err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func (r *applicationRunner) scopedContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(r.catacomb.Context(context.Background()))
}

// ScriptletService defines the scriptlet service methods needed by
// the scriptlet worker.
type ScriptletService interface {
	// WatchScriptletApplications returns a watcher that emits
	// application UUIDs when scriptlet applications are added,
	// removed, or changed.
	WatchScriptletApplications(ctx context.Context) (watcher.StringsWatcher, error)

	// GetApplicationScriptlet returns the scriptlet source for the
	// application identified by its UUID.
	GetApplicationScriptlet(ctx context.Context, appUUID application.UUID) (string, error)
}

// ApplicationService defines the application domain service methods
// needed by the scriptlet worker for per-application operations.
type ApplicationService interface {
	// WatchApplicationConfigChangeByUUID returns a watcher that emits
	// notifications when the application's config changes.
	WatchApplicationConfigChangeByUUID(ctx context.Context, appUUID application.UUID) (watcher.NotifyWatcher, error)

	// GetApplicationConfigWithDefaults returns the application config
	// with defaults applied.
	GetApplicationConfigWithDefaults(ctx context.Context, appUUID application.UUID) (charm.Config, error)
}

// RelationService defines the relation domain service methods needed
// by the scriptlet worker.
type RelationService interface {
	// WatchRelationsLifeSuspendedStatusForApplication watches for relation
	// lifecycle changes (created, broken) for the given application.
	// Returns relation UUIDs on change.
	WatchRelationsLifeSuspendedStatusForApplication(ctx context.Context, appUUID application.UUID) (watcher.StringsWatcher, error)
}

// dbScriptSource implements starform.ScriptSource by returning the
// scriptlet content loaded from the database.
type dbScriptSource struct {
	content string
}

var _ starform.ScriptSource = &dbScriptSource{}

func (s *dbScriptSource) Path() string {
	return "hooks.star"
}

func (s *dbScriptSource) Content(_ context.Context) ([]byte, error) {
	return []byte(s.content), nil
}
