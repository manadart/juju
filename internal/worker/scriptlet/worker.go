// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package scriptlet

import (
	"context"
	"time"

	"github.com/canonical/starform/starform"
	"github.com/juju/errors"
	"github.com/juju/worker/v5"
	"github.com/juju/worker/v5/catacomb"

	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/core/watcher"
	internalerrors "github.com/juju/juju/internal/errors"
	internalworker "github.com/juju/juju/internal/worker"
)

const (
	// States which report the state of the worker.
	stateStarted = "started"
)

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

		case appNames, ok := <-appWatcher.Changes():
			if !ok {
				return internalerrors.New("application watcher closed")
			}
			if err := w.handleApplicationChanges(ctx, appNames); err != nil {
				return internalerrors.Errorf("handling application changes: %w", err)
			}
		}
	}
}

// handleApplicationChanges processes changes to scriptlet applications,
// starting per-application workers via the Runner.
func (w *Worker) handleApplicationChanges(ctx context.Context, appNames []string) error {
	logger := w.config.Logger

	for _, appName := range appNames {
		logger.Infof(ctx, "ensuring scriptlet runner for application %q", appName)

		err := w.runner.StartWorker(ctx, appName, func(ctx context.Context) (worker.Worker, error) {
			return newApplicationRunner(applicationRunnerConfig{
				AppName: appName,
				Logger:  w.config.Logger,
			})
		})
		if errors.Is(err, errors.AlreadyExists) {
			continue
		}
		if err != nil {
			return internalerrors.Errorf("starting runner for %q: %w", appName, err)
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
	AppName string
	Logger  logger.Logger
}

func newApplicationRunner(config applicationRunnerConfig) (*applicationRunner, error) {
	scriptSet, err := starform.NewScriptSet(&starform.ScriptSetOptions{
		App: &starform.AppObject{
			Name: "juju",
		},
	})
	if err != nil {
		return nil, internalerrors.Errorf("creating script set: %w", err)
	}
	if err := scriptSet.LoadSources(context.Background(), []starform.ScriptSource{
		hardcodedScriptSource{},
	}); err != nil {
		return nil, internalerrors.Errorf("loading script set sources: %w", err)
	}

	r := &applicationRunner{
		config:    config,
		scriptSet: scriptSet,
	}

	if err := catacomb.Invoke(catacomb.Plan{
		Name: "scriptlet-" + config.AppName,
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
	logger.Infof(ctx, "scriptlet runner started for %q", r.config.AppName)

	// TODO(hackathon): Watch for config changes, relation changes,
	// and lifecycle events for this application. Dispatch events to
	// the Starform scriptlet when they occur.
	if err := r.handleConfigChanged(ctx); err != nil {
		return internalerrors.Errorf("handling scriptlet event: %w", err)
	}

	select {
	case <-r.catacomb.Dying():
		return r.catacomb.ErrDying()
	}
}

func (r *applicationRunner) handleConfigChanged(ctx context.Context) error {
	event := starform.EventObject{
		Name: "config_changed",
	}
	return r.scriptSet.Handle(ctx, &event)
}

func (r *applicationRunner) scopedContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(r.catacomb.Context(context.Background()))
}

// ScriptletService defines the scriptlet service methods needed by
// the scriptlet worker.
type ScriptletService interface {
	// WatchScriptletApplications returns a watcher that emits
	// application names when scriptlet applications are added,
	// removed, or changed.
	WatchScriptletApplications(ctx context.Context) (watcher.StringsWatcher, error)
}

type hardcodedScriptSource struct{}

var _ starform.ScriptSource = &hardcodedScriptSource{}

func (hardcodedScriptSource) Path() string {
	return "hooks.star"
}

func (hardcodedScriptSource) Content(context.Context) ([]byte, error) {
	ret := []byte(`
def init():
    juju.observe("config_changed", on_config_changed)

def on_config_changed(event):
	print('hello, world')
`)
	return ret, nil
}
