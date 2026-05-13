// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"strings"

	"github.com/juju/juju/core/changestream"
	coreerrors "github.com/juju/juju/core/errors"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/core/watcher/eventsource"
	"github.com/juju/juju/domain/application"
	"github.com/juju/juju/internal/errors"
)

// ScriptletRelation describes one relation endpoint from a scriptlet charm's
// metadata.
type ScriptletRelation struct {
	Name      string
	Role      string // "provider" | "requirer" | "peer"
	Interface string
	Scope     string // "global" | "container"; defaults to "global" if empty
	Optional  bool
	Limit     int
}

// RegisterScriptletArgs contains everything needed to record a scriptlet charm.
type RegisterScriptletArgs struct {
	ApplicationName string
	Scriptlet       string
	Relations       []ScriptletRelation
}

// State describes retrieval and persistence methods for scriptlet
// applications.
type State interface {
	// GetScriptletApplicationNames returns the names of all applications
	// that use scriptlet charms.
	GetScriptletApplicationNames(ctx context.Context) ([]string, error)

	// NamespaceForWatchScriptletApplications returns the namespace and
	// initial query for watching scriptlet application changes.
	NamespaceForWatchScriptletApplications() (string, eventsource.NamespaceQuery)

	// RegisterScriptlet records the scriptlet charm into the charm table
	// and its relations into charm_relation.
	RegisterScriptlet(ctx context.Context, args RegisterScriptletArgs) error
}

// ApplicationService provides access to the application domain service
// for configuration operations that the scriptlet service delegates.
type ApplicationService interface {
	// WatchApplicationConfigHash watches for changes to an
	// application's configuration hash.
	WatchApplicationConfigHash(ctx context.Context, appName string) (watcher.StringsWatcher, error)

	// GetApplicationConfigWithDefaults returns the application config
	// with defaults applied for the given application name.
	GetApplicationConfigWithDefaults(ctx context.Context, appName string) (map[string]interface{}, error)
}

// WatcherFactory instances return watchers for a given namespace.
type WatcherFactory interface {
	// NewNamespaceWatcher returns a new watcher that watches for changes
	// in the given namespace.
	NewNamespaceWatcher(
		ctx context.Context,
		initialQuery eventsource.NamespaceQuery,
		summary string,
		filterOption eventsource.FilterOption,
		filterOptions ...eventsource.FilterOption,
	) (watcher.StringsWatcher, error)
}

// Service provides the API for managing scriptlet applications.
type Service struct {
	st State
}

// NewService returns a new service reference wrapping the input state.
func NewService(st State) *Service {
	return &Service{st: st}
}

// GetScriptletApplicationNames returns the names of all applications
// that use scriptlet charms.
func (s *Service) GetScriptletApplicationNames(ctx context.Context) ([]string, error) {
	names, err := s.st.GetScriptletApplicationNames(ctx)
	if err != nil {
		return nil, errors.Errorf("getting scriptlet application names: %w", err)
	}
	return names, nil
}

// RegisterScriptlet records the scriptlet charm in the model database.
func (s *Service) RegisterScriptlet(ctx context.Context, args RegisterScriptletArgs) error {
	if !application.IsValidApplicationName(args.ApplicationName) {
		return errors.Errorf("application name %q is not valid", args.ApplicationName).
			Add(coreerrors.NotValid)
	}
	if strings.TrimSpace(args.Scriptlet) == "" {
		return errors.Errorf("scriptlet is empty").Add(coreerrors.NotValid)
	}
	for _, r := range args.Relations {
		if r.Name == "" {
			return errors.Errorf("relation name is empty").Add(coreerrors.NotValid)
		}
		switch r.Role {
		case "provider", "requirer", "peer":
		default:
			return errors.Errorf("unknown relation role %q for %q", r.Role, r.Name).Add(coreerrors.NotValid)
		}
		if r.Scope != "" && r.Scope != "global" && r.Scope != "container" {
			return errors.Errorf("unknown relation scope %q for %q", r.Scope, r.Name).Add(coreerrors.NotValid)
		}
	}
	return s.st.RegisterScriptlet(ctx, args)
}

// WatchableService provides the API for managing scriptlet applications
// and the ability to create watchers.
type WatchableService struct {
	*Service
	watcherFactory WatcherFactory
}

// NewWatchableService returns a new watchable service reference wrapping
// the input state.
func NewWatchableService(st State, watcherFactory WatcherFactory) *WatchableService {
	return &WatchableService{
		Service:        NewService(st),
		watcherFactory: watcherFactory,
	}
}

// WatchScriptletApplications returns a watcher that emits application
// names when scriptlet applications are added or removed.
func (s *WatchableService) WatchScriptletApplications(ctx context.Context) (watcher.StringsWatcher, error) {
	namespace, query := s.st.NamespaceForWatchScriptletApplications()
	return s.watcherFactory.NewNamespaceWatcher(
		ctx,
		query,
		"scriptlet applications watcher",
		eventsource.NamespaceFilter(namespace, changestream.All),
	)
}
