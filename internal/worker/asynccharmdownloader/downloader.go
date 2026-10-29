// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package asynccharmdownloader

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/juju/clock"
	"github.com/juju/retry"
	"github.com/juju/worker/v5"
	"gopkg.in/tomb.v2"

	"github.com/juju/juju/core/application"
	corecharm "github.com/juju/juju/core/charm"
	"github.com/juju/juju/core/logger"
	domainapplication "github.com/juju/juju/domain/application"
	applicationerrors "github.com/juju/juju/domain/application/errors"
	deploymentcharm "github.com/juju/juju/domain/deployment/charm"
	"github.com/juju/juju/domain/deployment/charm/charmdownloader"
	"github.com/juju/juju/internal/errors"
)

const (
	retryAttempts = 3
	retryDelay    = 20 * time.Second

	// downloadTimeout is the maximum amount of time to allow for a charm
	// download before timing out. We might want to make this configurable in
	// the future
	downloadTimeout = 5 * time.Minute
)

type asyncDownloadWorker struct {
	tomb tomb.Tomb

	appID              application.UUID
	applicationService ApplicationService
	downloader         Downloader
	resolver           CharmResolver

	clock  clock.Clock
	logger logger.Logger
}

// NewAsyncDownloadWorker creates a new async worker that downloads charms for
// the specified application.
func NewAsyncDownloadWorker(
	appID application.UUID,
	applicationService ApplicationService,
	downloader Downloader,
	resolver CharmResolver,
	clock clock.Clock,
	logger logger.Logger,
) worker.Worker {
	w := &asyncDownloadWorker{
		appID:              appID,
		applicationService: applicationService,
		downloader:         downloader,
		resolver:           resolver,
		clock:              clock,
		logger:             logger,
	}
	w.tomb.Go(w.loop)
	return w
}

// Kill is part of the worker.Worker interface.
func (w *asyncDownloadWorker) Kill() {
	w.tomb.Kill(nil)
}

// Wait is part of the worker.Worker interface.
func (w *asyncDownloadWorker) Wait() error {
	return w.tomb.Wait()
}

func (w *asyncDownloadWorker) loop() error {
	ctx, cancel := w.scopedContext()
	defer cancel()

	w.logger.Infof(ctx, "downloading charm for application %q", w.appID)

	info, err := w.applicationService.GetAsyncCharmDownloadInfo(ctx, w.appID)
	if errors.Is(err, applicationerrors.CharmAlreadyAvailable) {
		// If the application is already downloading a charm, we can skip this
		// application.
		w.logger.Infof(ctx, "charm for application %q is already available", w.appID)
		return nil
	} else if err != nil {
		return errors.Capture(err)
	}

	if info.ResolutionPending {
		info, err = w.resolvePendingCharm(ctx, info)
		if err != nil {
			return errors.Capture(err)
		}
	}

	// Ensure we've got a valid URL.
	url, err := url.Parse(info.DownloadInfo.DownloadURL)
	if err != nil {
		return errors.Capture(err)
	}

	// Download the charm for the application.
	var result *charmdownloader.DownloadResult
	if err := retry.Call(retry.CallArgs{
		Func: func() error {
			ctx, cancel := context.WithTimeout(ctx, downloadTimeout)
			defer cancel()

			var err error
			result, err = w.downloader.Download(ctx, url, info.SHA256)
			if err != nil {
				return errors.Capture(err)
			}
			return nil
		},
		Attempts: retryAttempts,
		Delay:    retryDelay,
		Clock:    w.clock,
		NotifyFunc: func(err error, i int) {
			w.logger.Warningf(ctx, "failed to download charm for application %q, attempt %d: %v", w.appID, i, err)
		},
		Stop: w.tomb.Dying(),
	}); err != nil {
		return errors.Capture(err)
	}

	// Ensure the charm is removed after the worker has finished.
	defer func() {
		if err := os.Remove(result.Path); err != nil && !os.IsNotExist(err) {
			w.logger.Warningf(ctx, "failed to remove temporary file %q: %v", result.Path, err)
		}
	}()

	// The charm has been downloaded, we can now resolve the download slot.
	//
	// TODO: If this fails (which it sometimes does), the result is that the
	// charm is redownloaded on the next attempt. We should re-architecture
	// these workers to avoid this unnecessary work
	err = w.applicationService.ResolveCharmDownload(ctx, w.appID, domainapplication.ResolveCharmDownload{
		SHA256:    result.SHA256,
		SHA384:    result.SHA384,
		CharmUUID: info.CharmUUID,
		Path:      result.Path,
		Size:      result.Size,
	})
	if err != nil && !errors.Is(err, applicationerrors.CharmAlreadyResolved) {
		w.logger.Errorf(ctx, "failed to resolve charm download for application %q: %v", w.appID, err)
		return errors.Capture(err)
	}

	// Exit cleanly, so the worker doesn't get restarted.
	w.logger.Infof(ctx, "downloaded and resolved charm for application %q", w.appID)
	return nil
}

func (w *asyncDownloadWorker) resolvePendingCharm(
	ctx context.Context,
	info domainapplication.CharmDownloadInfo,
) (domainapplication.CharmDownloadInfo, error) {
	appName, err := w.applicationService.GetApplicationName(ctx, w.appID)
	if err != nil {
		return domainapplication.CharmDownloadInfo{}, errors.Errorf("getting application name: %w", err)
	}
	origin, err := w.applicationService.GetApplicationCharmOrigin(ctx, appName)
	if err != nil {
		return domainapplication.CharmDownloadInfo{}, errors.Errorf("getting charm origin: %w", err)
	}

	// The stored revision, ID and hash came from the provisional architecture.
	// Resolve the application's channel afresh using the architecture bound by
	// the relation transaction.
	origin.ID = ""
	origin.Hash = ""
	origin.Revision = nil
	charmURL, err := deploymentcharm.ParseURL(info.Name)
	if err != nil {
		return domainapplication.CharmDownloadInfo{}, errors.Errorf("parsing charm name %q: %w", info.Name, err)
	}
	resolved, err := w.resolver.ResolveForDeploy(ctx, corecharm.CharmID{
		URL:    charmURL,
		Origin: origin,
	})
	if err != nil {
		return domainapplication.CharmDownloadInfo{}, errors.Errorf("resolving charm for bound architecture: %w", err)
	}

	resolvedOrigin := resolved.EssentialMetadata.ResolvedOrigin
	downloadInfo := resolved.EssentialMetadata.DownloadInfo
	if err := w.applicationService.UpdatePendingCharmResolution(
		ctx,
		w.appID,
		domainapplication.PendingCharmResolution{
			Revision:           resolved.URL.Revision,
			SHA256:             resolvedOrigin.Hash,
			CharmhubIdentifier: downloadInfo.CharmhubIdentifier,
			DownloadURL:        downloadInfo.DownloadURL,
			DownloadSize:       downloadInfo.DownloadSize,
		},
	); err != nil {
		return domainapplication.CharmDownloadInfo{}, errors.Errorf("saving charm resolution: %w", err)
	}

	return w.applicationService.GetAsyncCharmDownloadInfo(ctx, w.appID)
}

func (w *asyncDownloadWorker) scopedContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(w.tomb.Context(context.Background()))
}
