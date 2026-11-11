// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storageprovisioner

import (
	"context"
	"time"

	"github.com/juju/names/v6"

	coreerrors "github.com/juju/juju/core/errors"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/rpc/params"
)

// ResizeAccessor exposes only the executor's eligible stages. Reports are
// idempotent and may be retried after an ambiguous response.
type ResizeAccessor interface {
	Requests(context.Context) ([]string, error)
	PrepareStage(context.Context, string) (params.StorageResizeStageWork, error)
	ReportStage(context.Context, params.StorageResizeStageReport) error
}

// ResizeExecutor performs one bounded, repeatable reconciliation. It must call
// revalidate before mutation and return observed capacities, never desired sizes.
type ResizeExecutor interface {
	Execute(context.Context, params.StorageResizeStageWork, func(context.Context) error) (params.StorageResizeStageReport, error)
}

type resizeOutcome struct {
	report *params.StorageResizeStageReport
	next   time.Time
	delay  time.Duration
}

func processStorageResizes(ctx context.Context, deps *dependencies, results map[string]resizeOutcome) error {
	ids, err := deps.config.Resizes.Requests(ctx)
	if err != nil {
		return errors.Errorf("fetching resize requests: %w", err)
	}
	pending := make(map[string]bool, len(ids))
	for _, id := range ids {
		pending[id] = true
	}
	for id := range results {
		if !pending[id] {
			delete(results, id)
		}
	}
	for _, id := range ids {
		if err := ctx.Err(); err != nil {
			return err
		}
		outcome := results[id]
		if deps.config.Clock.Now().Before(outcome.next) {
			continue
		}
		if outcome.report == nil {
			work, err := deps.config.Resizes.PrepareStage(ctx, id)
			if err != nil {
				switch resizeErrorCode(err) {
				case params.CodeTryAgain, params.CodeAlreadyExists, params.CodeNotFound, params.CodeUnauthorized:
					continue
				}
				deps.config.Logger.Warningf(ctx, "preparing resize %q: %v", id, err)
				results[id] = backoffResize(deps, outcome)
				continue
			}
			timeout := 5 * time.Minute
			if _, ok := deps.config.Scope.(names.ModelTag); ok {
				timeout = 30 * time.Second
			}
			stepCtx, cancel := context.WithTimeout(ctx, timeout)
			var validationErr error
			revalidate := func(ctx context.Context) error {
				current, err := deps.config.Resizes.PrepareStage(ctx, id)
				if err == nil && (current.Identity != work.Identity || current.Stage != work.Stage || current.Phase != work.Phase) {
					err = errors.New("resize stage changed before mutation")
				}
				validationErr = err
				return err
			}
			report, err := deps.config.ResizeExecutor.Execute(stepCtx, work, revalidate)
			if stepCtx.Err() != nil {
				err = stepCtx.Err()
			}
			cancel()
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if validationErr != nil {
				results[id] = backoffResize(deps, outcome)
				continue
			}
			report.UUID = id
			report.Identity = work.Identity
			report.Stage = work.Stage
			report.Phase = work.Phase
			if err != nil {
				report.Message = err.Error()
				if errors.Is(err, coreerrors.NotSupported) || errors.Is(err, coreerrors.NotValid) {
					report.Status = "failed"
				} else if work.Phase == "check" {
					report.Status = "blocked"
				} else {
					report.Status = "pending"
				}
			}
			outcome.report = &report
		}
		report := *outcome.report
		err := deps.config.Resizes.ReportStage(ctx, report)
		if err != nil {
			switch resizeErrorCode(err) {
			case params.CodeAlreadyExists, params.CodeNotFound, params.CodeUnauthorized, params.CodeTryAgain:
				delete(results, id)
				continue
			case params.CodeNotValid:
				// A rejected completion becomes a stage failure, preserving
				// earlier growth without applying rejected observations.
				if report.Status != "failed" {
					report.Status = "failed"
					report.Message = "Resize result rejected: " + err.Error()
					report.ObservedSizeMiB = 0
					report.FilesystemSizeMiB = 0
					outcome.report = &report
				}
			}
			deps.config.Logger.Warningf(ctx, "reporting resize %q: %v", id, err)
			results[id] = backoffResize(deps, outcome)
			continue
		}
		if report.Status == "blocked" || report.Status == "pending" {
			outcome.report = nil
			results[id] = backoffResize(deps, outcome)
		} else {
			delete(results, id)
			deps.config.Logger.Infof(ctx, "storage resize %q %s %s: %s", id, report.Stage, report.Status, report.Message)
		}
	}
	return nil
}

func backoffResize(deps *dependencies, outcome resizeOutcome) resizeOutcome {
	outcome.delay = min(max(outcome.delay*2, 5*time.Second), time.Minute)
	outcome.next = deps.config.Clock.Now().Add(outcome.delay)
	return outcome
}

// Model service calls and agent RPCs use the same reconciliation decisions.
func resizeErrorCode(err error) string {
	switch {
	case errors.Is(err, storageerrors.StorageResizeNotReady):
		return params.CodeTryAgain
	case errors.Is(err, storageerrors.StorageResizeResultAlreadyRecorded):
		return params.CodeAlreadyExists
	case errors.Is(err, storageerrors.StorageResizeIntentNotFound):
		return params.CodeNotFound
	case errors.Is(err, storageerrors.StorageResizeNotAssigned):
		return params.CodeUnauthorized
	case errors.Is(err, storageerrors.StorageResizeTargetChanged), errors.Is(err, coreerrors.NotValid):
		return params.CodeNotValid
	}
	return params.ErrCode(err)
}
