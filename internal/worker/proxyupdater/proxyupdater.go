// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package proxyupdater

import (
	"context"
	"io"
	stdos "os"
	stdexec "os/exec"
	"strings"

	"github.com/juju/errors"
	"github.com/juju/proxy"
	"github.com/juju/worker/v4"

	"github.com/juju/juju/api/agent/proxyupdater"
	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/core/snap"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/internal/packaging/commands"
	"github.com/juju/juju/internal/packaging/config"
)

type Config struct {
	SupportLegacyValues bool
	EnvFiles            []string
	SystemdFiles        []string
	API                 API
	ExternalUpdate      func(proxy.Settings) error
	InProcessUpdate     func(proxy.Settings) error
	RunFunc             func(string, string, ...string) (string, error)
	Logger              logger.Logger
}

// Validate ensures that all the required fields have values.
func (c *Config) Validate() error {
	if c.API == nil {
		return errors.NotValidf("missing API")
	}
	if c.InProcessUpdate == nil {
		return errors.NotValidf("missing InProcessUpdate")
	}
	if c.Logger == nil {
		return errors.NotValidf("missing Logger")
	}
	return nil
}

// API is an interface that is provided to New
// which can be used to fetch the API host ports
type API interface {
	ProxyConfig(context.Context) (proxyupdater.ProxyConfiguration, error)
	WatchForProxyConfigAndAPIHostPortChanges(context.Context) (watcher.NotifyWatcher, error)
}

// proxyWorker is responsible for monitoring the juju environment
// configuration and making changes on the physical (or virtual) machine as
// necessary to match the environment changes.  Examples of these types of
// changes are apt proxy configuration and the juju proxies stored in the juju
// proxy file.
type proxyWorker struct {
	aptProxy  proxy.Settings
	aptMirror string
	proxy     proxy.Settings

	snapProxy           proxy.Settings
	snapStoreProxy      string
	snapStoreAssertions string
	snapStoreProxyURL   string

	// The whole point of the first value is to make sure that the files
	// are written out the first time through, even if they are the same as
	// "last" time, as the initial value for last time is the zeroed struct.
	// There is the possibility that the files exist on disk with old
	// settings, and the environment has been updated to now not have them. We
	// need to make sure that the disk reflects the environment, so the first
	// time through, even if the proxies are empty, we write the files to
	// disk.
	first  bool
	config Config
}

// NewWorker returns a worker.Worker that updates proxy environment variables for the
// process and for the whole machine.
var NewWorker = func(config Config) (worker.Worker, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	envWorker := &proxyWorker{
		first:  true,
		config: config,
	}
	w, err := watcher.NewNotifyWorker(watcher.NotifyConfig{
		Handler: envWorker,
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	return w, nil
}

func (w *proxyWorker) saveProxySettings(ctx context.Context) error {
	// The proxy settings are (usually) stored in three files:
	// - /etc/juju-proxy.conf - in 'env' format
	// - /etc/systemd/system.conf.d/juju-proxy.conf
	// - /etc/systemd/user.conf.d/juju-proxy.conf - both in 'systemd' format
	for _, file := range w.config.EnvFiles {
		err := stdos.WriteFile(file, []byte(w.proxy.AsScriptEnvironment()), 0644)
		if err != nil {
			w.config.Logger.Errorf(ctx, "Error updating environment file %s - %v", file, err)
		}
	}
	for _, file := range w.config.SystemdFiles {
		err := stdos.WriteFile(file, []byte(w.proxy.AsSystemdDefaultEnv()), 0644)
		if err != nil {
			w.config.Logger.Errorf(ctx, "Error updating systemd file - %v", err)
		}
	}
	return nil
}

func (w *proxyWorker) handleProxyValues(ctx context.Context, legacyProxySettings, jujuProxySettings proxy.Settings) {
	// Legacy proxy settings update the environment, and also call the
	// InProcessUpdate, which installs the proxy into the default HTTP
	// transport. The same occurs for jujuProxySettings.
	settings := jujuProxySettings
	if jujuProxySettings.HasProxySet() {
		w.config.Logger.Debugf(ctx, "applying in-process juju proxy settings %#v", jujuProxySettings)
	} else {
		settings = legacyProxySettings
		w.config.Logger.Debugf(ctx, "applying in-process legacy proxy settings %#v", legacyProxySettings)
	}

	settings.SetEnvironmentValues()
	if err := w.config.InProcessUpdate(settings); err != nil {
		w.config.Logger.Errorf(ctx, "error updating in-process proxy settings: %v", err)
	}

	// If the external update function is passed in, it is to update the LXD
	// proxies. We want to set this to the proxy specified regardless of whether
	// it was set with the legacy fields or the new juju fields.
	if externalFunc := w.config.ExternalUpdate; externalFunc != nil {
		if err := externalFunc(settings); err != nil {
			// It isn't really fatal, but we should record it.
			w.config.Logger.Errorf(ctx, "%v", err)
		}
	}

	// Here we write files to disk. This is done only for legacyProxySettings.
	if w.config.SupportLegacyValues && (legacyProxySettings != w.proxy || w.first) {
		w.config.Logger.Debugf(ctx, "saving new legacy proxy settings %#v", legacyProxySettings)
		w.proxy = legacyProxySettings
		if err := w.saveProxySettings(ctx); err != nil {
			// It isn't really fatal, but we should record it.
			w.config.Logger.Errorf(ctx, "error saving proxy settings: %v", err)
		}
	}
}

func (w *proxyWorker) handleSnapProxyValues(ctx context.Context, proxy proxy.Settings, storeID, storeAssertions, storeProxyURL string) {
	if w.config.RunFunc == nil {
		w.config.Logger.Tracef(ctx, "snap proxies not updated")
		return
	}

	var snapSettings []string
	maybeAddSettings := func(setting, value, saved string) {
		if value != saved || w.first {
			snapSettings = append(snapSettings, setting+"="+value)
		}
	}
	maybeAddSettings("proxy.http", proxy.Http, w.snapProxy.Http)
	maybeAddSettings("proxy.https", proxy.Https, w.snapProxy.Https)

	// Proxy URL changed; either a new proxy has been provided or the proxy
	// has been removed. Proxy URL changes have a higher precedence than
	// manually specifying the assertions and store ID.
	if storeProxyURL != w.snapStoreProxyURL {
		if storeProxyURL != "" {
			var err error
			if storeAssertions, storeID, err = snap.LookupAssertions(storeProxyURL); err != nil {
				w.config.Logger.Errorf(ctx, "unable to lookup snap store assertions: %v", err)
				return
			} else {
				w.config.Logger.Infof(ctx, "auto-detected snap store assertions from proxy")
				w.config.Logger.Infof(ctx, "auto-detected snap store ID as %q", storeID)
			}
		} else if storeAssertions != "" && storeID != "" {
			// The proxy URL has been removed. However, if the user
			// has manually provided assertion/store ID config
			// options we should restore them. To do this, we reset
			// the last seen values so we can force-apply the
			// previously specified manual values. Otherwise, the
			// provided storeAssertions/storeID values are empty
			// and we simply fall through to allow the code to
			// reset the proxy.store setting to an empty value.
			w.snapStoreAssertions, w.snapStoreProxy = "", ""
		}
		w.snapStoreProxyURL = storeProxyURL
	} else if storeProxyURL != "" {
		// Re-use the storeID and assertions obtained by querying the
		// proxy during the last update.
		storeAssertions, storeID = w.snapStoreAssertions, w.snapStoreProxy
	}

	maybeAddSettings("proxy.store", storeID, w.snapStoreProxy)

	// If an assertion file was provided we need to "snap ack" it before
	// configuring snap to use the store ID.
	if storeAssertions != w.snapStoreAssertions && storeAssertions != "" {
		output, err := w.config.RunFunc(storeAssertions, "snap", "ack", "/dev/stdin")
		if err != nil {
			w.config.Logger.Warningf(ctx, "unable to acknowledge assertions: %v, output: %q", err, output)
			return
		}
		w.snapStoreAssertions = storeAssertions
	}

	if len(snapSettings) > 0 {
		args := append([]string{"set", "system"}, snapSettings...)
		output, err := w.config.RunFunc(noStdIn, "snap", args...)
		if err != nil {
			w.config.Logger.Warningf(ctx, "unable to set snap core settings %v: %v, output: %q", snapSettings, err, output)
		} else {
			w.config.Logger.Debugf(ctx, "snap core settings %v updated, output: %q", snapSettings, output)
			w.snapProxy = proxy
			w.snapStoreProxy = storeID
		}
	}
}

func (w *proxyWorker) handleAptProxyValues(ctx context.Context, aptSettings proxy.Settings, aptMirror string) {
	mirrorUpdateNeeded := aptMirror != "" && aptMirror != w.aptMirror
	aptCommander := commands.NewAptPackageCommander()

	if aptSettings != w.aptProxy || w.first {
		w.config.Logger.Debugf(ctx, "new apt proxy settings %#v", aptSettings)
		w.aptProxy = aptSettings

		// Always finish with a new line.
		content := aptCommander.ProxyConfigContents(w.aptProxy) + "\n"
		err := stdos.WriteFile(config.AptProxyConfigFile, []byte(content), 0644)
		if err != nil {
			// It isn't really fatal, but we should record it.
			w.config.Logger.Errorf(ctx, "error writing apt proxy config file: %v", err)
		}
	}
	if mirrorUpdateNeeded {
		if w.config.RunFunc == nil {
			w.config.Logger.Tracef(ctx, "apt mirrors not updated")
			return
		}
		w.config.Logger.Debugf(ctx, "new apt mirror value %v", aptMirror)
		w.aptMirror = aptMirror

		cmds := aptCommander.SetMirrorCommands(aptMirror, aptMirror)
		script := []string{"#!/bin/bash", "set -e"}
		script = append(script, "(")
		script = append(script, cmds...)
		script = append(script, ")")
		w.config.Logger.Tracef(ctx, strings.Join(script, "\n"))
		if output, err := w.config.RunFunc(noStdIn, "/bin/bash", "-c", strings.Join(script, "\n")); err != nil {
			w.config.Logger.Warningf(ctx, "unable to update apt mirrors: %v, output: %q", err, output)
		}
	}
	return
}

func (w *proxyWorker) onChange(ctx context.Context) error {
	config, err := w.config.API.ProxyConfig(ctx)
	if err != nil {
		return err
	}

	w.handleProxyValues(ctx, config.LegacyProxy, config.JujuProxy)
	w.handleSnapProxyValues(ctx, config.SnapProxy, config.SnapStoreProxyId, config.SnapStoreProxyAssertions, config.SnapStoreProxyURL)
	w.handleAptProxyValues(ctx, config.APTProxy, config.AptMirror)
	return nil
}

// SetUp is defined on the worker.NotifyWatchHandler interface.
func (w *proxyWorker) SetUp(ctx context.Context) (watcher.NotifyWatcher, error) {
	// We need to set this up initially as the NotifyWorker sucks up the first
	// event.
	err := w.onChange(ctx)
	if err != nil {
		return nil, err
	}
	w.first = false
	return w.config.API.WatchForProxyConfigAndAPIHostPortChanges(ctx)
}

// Handle is defined on the worker.NotifyWatchHandler interface.
func (w *proxyWorker) Handle(ctx context.Context) error {
	return w.onChange(ctx)
}

// TearDown is defined on the worker.NotifyWatchHandler interface.
func (w *proxyWorker) TearDown() error {
	// Nothing to cleanup, only state is the watcher
	return nil
}

const noStdIn = ""

// RunWithStdIn executes the command specified with the args with optional stdin.
func RunWithStdIn(input string, command string, args ...string) (string, error) {
	cmd := stdexec.Command(command, args...)

	if input != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return "", errors.Annotate(err, "getting stdin pipe")
		}

		go func() {
			defer stdin.Close()
			_, _ = io.WriteString(stdin, input)
		}()
	}

	out, err := cmd.CombinedOutput()
	output := string(out)
	return output, err
}
