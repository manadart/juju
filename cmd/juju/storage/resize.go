// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"context"
	"math"

	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"github.com/juju/utils/v4"

	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/cmd"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/core/output"
	corestorage "github.com/juju/juju/core/storage"
	"github.com/juju/juju/rpc/params"
)

// NewResizeCommand returns a command for requesting a storage resize.
func NewResizeCommand() cmd.Command {
	c := &resizeCommand{}
	c.newAPIFunc = func(ctx context.Context) (storageResizeAPI, error) {
		return c.NewStorageAPI(ctx)
	}
	return modelcmd.Wrap(c)
}

type storageResizeAPI interface {
	storageResizeResultAPI
	ResizeStorage(context.Context, string, uint64) (string, error)
	SupportsStorageResizeResults() bool
}

type resizeCommand struct {
	StorageCommandBase
	modelcmd.IAASOnlyCommand

	storageID  corestorage.ID
	size       string
	sizeMiB    uint64
	wait       bool
	out        cmd.Output
	newAPIFunc func(context.Context) (storageResizeAPI, error)
}

// Info implements Command.Info.
func (c *resizeCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    "resize-storage",
		Args:    "<storage-id> --size <size>",
		Purpose: "Records a request to grow a storage instance.",
		Doc: `
Request a new total capacity for an existing storage instance in a machine model.
The size is a target capacity, not an amount to add. Sizes use binary units;
M, G, T and P suffixes are supported, and a size without a suffix is in MiB.

The command prints the UUID of the recorded request. Recording a request does
not mean the resize has been approved or completed. Only one request can be
pending for a storage instance.

Use --wait to wait until a completed or failed result is reported. The result
is printed instead of just the UUID, and a failed resize exits unsuccessfully.
Stopping the wait does not cancel the request. Use show-storage-resize with the
request UUID to retrieve the result later.

Run any charm-specific preparation action separately before requesting a resize.
This command does not pause charm hooks or run preparation actions.
`,
		Examples: "juju resize-storage pgdata/0 --size 20G",
		SeeAlso:  []string{"show-storage-resize", "storage", "show-storage", "add-storage", "run"},
	})
}

// SetFlags implements Command.SetFlags.
func (c *resizeCommand) SetFlags(f *gnuflag.FlagSet) {
	c.StorageCommandBase.SetFlags(f)
	f.StringVar(&c.size, "size", "", "Requested total capacity (for example, 20G)")
	f.BoolVar(&c.wait, "wait", false, "Wait for the resize result")
	c.out.AddFlags(f, "yaml", output.DefaultFormatters)
}

// Init implements Command.Init.
func (c *resizeCommand) Init(args []string) error {
	if len(args) != 1 {
		return errors.New("resize-storage requires one storage ID")
	}
	var err error
	c.storageID, err = corestorage.ParseID(args[0])
	if err != nil {
		return err
	}
	if c.size == "" {
		return errors.New("--size is required")
	}
	c.sizeMiB, err = utils.ParseSize(c.size)
	if err != nil {
		return errors.Annotate(err, "invalid size")
	}
	if c.sizeMiB == 0 || c.sizeMiB > math.MaxInt64 {
		return errors.New("size must be between 1 and 9223372036854775807 MiB")
	}
	return nil
}

// Run implements Command.Run.
func (c *resizeCommand) Run(ctx *cmd.Context) error {
	api, err := c.newAPIFunc(ctx)
	if err != nil {
		return err
	}
	defer api.Close()
	if c.wait && !api.SupportsStorageResizeResults() {
		return errors.NotSupportedf("waiting for storage resize results on this version of Juju")
	}
	uuid, err := api.ResizeStorage(ctx, c.storageID.String(), c.sizeMiB)
	if err != nil {
		return err
	}
	if !c.wait {
		return c.out.Write(ctx, uuid)
	}
	ctx.Infof("Resize request %s recorded. Waiting for a result.", uuid)
	result, err := getStorageResizeResult(ctx, api, uuid, true, clock.WallClock)
	if err != nil {
		return errors.Annotatef(err, "waiting for resize request %s", uuid)
	}
	formatted, err := formatStorageResizeResult(result)
	if err != nil {
		return err
	}
	if err := c.out.Write(ctx, formatted); err != nil {
		return err
	}
	if result.Status == params.StorageResizeFailed {
		return cmd.ErrSilent
	}
	return nil
}
