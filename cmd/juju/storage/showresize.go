// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"context"
	"time"

	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"github.com/juju/names/v6"
	"github.com/juju/utils/v4"

	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/cmd"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/core/output"
	"github.com/juju/juju/rpc/params"
)

// NewShowResizeCommand returns a command for retrieving resize results.
func NewShowResizeCommand() cmd.Command {
	c := &showResizeCommand{}
	c.newAPIFunc = func(ctx context.Context) (storageResizeResultAPI, error) {
		return c.NewStorageAPI(ctx)
	}
	return modelcmd.Wrap(c)
}

type storageResizeResultAPI interface {
	Close() error
	StorageResizeDetails(context.Context, string) (params.StorageResizeDetails, error)
}

type showResizeCommand struct {
	StorageCommandBase
	uuid       string
	wait       bool
	out        cmd.Output
	newAPIFunc func(context.Context) (storageResizeResultAPI, error)
}

// Info implements Command.Info.
func (c *showResizeCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    "show-storage-resize",
		Args:    "<request-id>",
		Purpose: "Shows the status and result of a storage resize request.",
		Doc: `
Retrieve the request identified by the UUID returned by resize-storage.
The result includes the original storage ID, requested size in MiB, status,
creation time and, when reported, the completion time and result message.

Without --wait, pending requests are shown immediately. With --wait, the command
waits until a completed or failed result is reported. Stopping the wait does
not cancel the request. An unknown request ID is an error.

Results remain available after completion and after the storage is removed.
This command exits successfully when it retrieves a result, including a failed
resize result; inspect the status field to determine the resize outcome.
`,
		Examples: "juju show-storage-resize <request-id> --wait --format json",
		SeeAlso:  []string{"resize-storage", "show-storage"},
	})
}

// SetFlags implements Command.SetFlags.
func (c *showResizeCommand) SetFlags(f *gnuflag.FlagSet) {
	c.StorageCommandBase.SetFlags(f)
	f.BoolVar(&c.wait, "wait", false, "Wait for the resize result")
	c.out.AddFlags(f, "yaml", output.DefaultFormatters)
}

// Init implements Command.Init.
func (c *showResizeCommand) Init(args []string) error {
	if len(args) != 1 {
		return errors.New("show-storage-resize requires one request ID")
	}
	if !utils.IsValidUUIDString(args[0]) {
		return errors.NotValidf("resize request ID %q", args[0])
	}
	c.uuid = args[0]
	return nil
}

// Run implements Command.Run.
func (c *showResizeCommand) Run(ctx *cmd.Context) error {
	api, err := c.newAPIFunc(ctx)
	if err != nil {
		return err
	}
	defer api.Close()
	result, err := getStorageResizeResult(ctx, api, c.uuid, c.wait, clock.WallClock)
	if err != nil {
		return err
	}
	formatted, err := formatStorageResizeResult(result)
	if err != nil {
		return err
	}
	return c.out.Write(ctx, formatted)
}

// getStorageResizeResult polls without holding a server call open. Polling is
// read-only; cancellation never cancels or resubmits the persisted request.
func getStorageResizeResult(
	ctx context.Context, api storageResizeResultAPI, uuid string, wait bool,
	clk clock.Clock,
) (params.StorageResizeDetails, error) {
	for {
		if err := ctx.Err(); err != nil {
			return params.StorageResizeDetails{}, err
		}
		result, err := api.StorageResizeDetails(ctx, uuid)
		if err != nil {
			return params.StorageResizeDetails{}, err
		}
		switch result.Status {
		case params.StorageResizeCompleted, params.StorageResizeFailed:
			return result, nil
		case params.StorageResizePending:
			if !wait {
				return result, nil
			}
		default:
			return params.StorageResizeDetails{}, errors.Errorf("unknown resize status %q", result.Status)
		}
		timer := clk.NewTimer(2 * time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return params.StorageResizeDetails{}, ctx.Err()
		case <-timer.Chan():
			timer.Stop()
		}
	}
}

type storageResizeOutput struct {
	ID               string     `json:"id" yaml:"id"`
	Storage          string     `json:"storage" yaml:"storage"`
	RequestedSizeMiB uint64     `json:"requested-size-mib" yaml:"requested-size-mib"`
	Status           string     `json:"status" yaml:"status"`
	Message          string     `json:"message,omitempty" yaml:"message,omitempty"`
	CreatedAt        time.Time  `json:"created-at" yaml:"created-at"`
	CompletedAt      *time.Time `json:"completed-at,omitempty" yaml:"completed-at,omitempty"`
}

func formatStorageResizeResult(result params.StorageResizeDetails) (storageResizeOutput, error) {
	tag, err := names.ParseStorageTag(result.StorageTag)
	if err != nil {
		return storageResizeOutput{}, err
	}
	return storageResizeOutput{
		ID:               result.UUID,
		Storage:          tag.Id(),
		RequestedSizeMiB: result.RequestedSizeMiB,
		Status:           result.Status,
		Message:          result.Message,
		CreatedAt:        result.CreatedAt,
		CompletedAt:      result.CompletedAt,
	}, nil
}
