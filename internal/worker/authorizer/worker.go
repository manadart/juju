// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package authorizer

import (
	"context"
	"sync"

	"github.com/juju/worker/v5"
	"gopkg.in/tomb.v2"

	"github.com/juju/juju/core/permission"
	"github.com/juju/juju/internal/errors"
)

// Authorizer is the engine owned by this worker. Close is called once, after
// cancellation and completion of all checks admitted by the worker.
type Authorizer interface {
	permission.Authorizer
	Close()
}

// Config supplies the engine constructor. It must honour context cancellation
// during startup and validate the persisted model before returning the engine.
type Config struct {
	NewAuthorizer func(context.Context) (Authorizer, error)
}

// NewWorker starts an authorizer and waits until it can serve checks. The input
// context bounds startup; the worker owns the engine lifetime after startup.
func NewWorker(ctx context.Context, config Config) (worker.Worker, error) {
	if config.NewAuthorizer == nil {
		return nil, errors.New("nil NewAuthorizer")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	w := &authorizerWorker{config: config, ready: make(chan struct{})}
	w.ctx = w.tomb.Context(context.Background())
	w.tomb.Go(w.loop)
	select {
	case <-ctx.Done():
		w.Kill()
		_ = w.Wait()
		return nil, ctx.Err()
	case <-w.tomb.Dead():
		return nil, w.Wait()
	case <-w.ready:
		return w, nil
	}
}

type authorizerWorker struct {
	config Config
	tomb   tomb.Tomb
	ctx    context.Context
	ready  chan struct{}
	engine Authorizer

	mu      sync.Mutex
	stopped bool
	checks  sync.WaitGroup
}

var _ permission.Authorizer = (*authorizerWorker)(nil)

func (w *authorizerWorker) loop() error {
	engine, err := w.config.NewAuthorizer(w.ctx)
	if err != nil {
		return errors.Errorf("starting authorisation engine: %w", err)
	}
	if engine == nil {
		return errors.New("nil authorisation engine")
	}
	defer engine.Close()
	w.engine = engine
	close(w.ready)

	<-w.tomb.Dying()
	// Close admission before waiting, so no check can race Add against Wait.
	w.mu.Lock()
	w.stopped = true
	w.mu.Unlock()
	w.checks.Wait()
	return w.tomb.Err()
}

// Check binds each request to both its caller and the worker lifetime.
func (w *authorizerWorker) Check(ctx context.Context, subject string, target permission.ID, access permission.Access) (bool, error) {
	w.mu.Lock()
	if w.stopped || w.ctx.Err() != nil {
		w.mu.Unlock()
		return false, context.Canceled
	}
	w.checks.Add(1)
	w.mu.Unlock()
	defer w.checks.Done()

	ctx, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(w.ctx, cancel)
	defer func() {
		stop()
		cancel()
	}()
	return w.engine.Check(ctx, subject, target, access)
}

// Kill cancels startup, background resources and in-flight permission checks.
func (w *authorizerWorker) Kill() { w.tomb.Kill(nil) }

// Wait joins all active checks and engine shutdown.
func (w *authorizerWorker) Wait() error { return w.tomb.Wait() }
