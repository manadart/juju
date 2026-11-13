// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"github.com/juju/names/v6"

	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/cmd"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/core/output"
	"github.com/juju/juju/rpc/params"
)

// NewResizeCommand returns a command applying capacity directives to a unit.
func NewResizeCommand() cmd.Command {
	c := &resizeCommand{}
	c.newAPIFunc = func(ctx context.Context) (storageResizeAPI, error) { return c.NewStorageAPI(ctx) }
	return modelcmd.Wrap(c)
}

type storageResizeAPI interface {
	storageResizeResultAPI
	ResizeUnitStorage(context.Context, string) ([]params.UnitStorageResizeResult, error)
}

type resizeCommand struct {
	StorageCommandBase
	modelcmd.IAASOnlyCommand
	unitName   string
	wait       bool
	out        cmd.Output
	newAPIFunc func(context.Context) (storageResizeAPI, error)
}

// Info implements Command.Info.
func (c *resizeCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name: "resize-storage", Args: "<unit>",
		Purpose: "Grows a unit's storage to the application capacity targets.",
		Doc: `
Apply the current application storage directives to a unit's existing storage.
Use application-storage first to set the desired capacity for each storage name.
The controller stages requests for all instances that need growth, with targets
fixed when submitted. Later directive edits do not change those requests.

Instances already large enough are reported as already-satisfied. Matching
pending requests are reused. Conflicting pending targets reject the submission.
Pool and count changes do not move, add or remove existing storage. Storage is
never shrunk. This prototype supports machine models.

Output lists each integer request ID and its storage name and instance ID.
Each request is independent; there is no parent operation. YAML and JSON output
map storage IDs to request IDs, target capacities and statuses.
Use --wait to wait for all requests and report individual results. Any failed
request makes the command exit unsuccessfully. Stopping the wait does not cancel
requests. Use show-storage-resize <request-id> to retrieve a result later.

Run any charm-specific preparation action separately before requesting a resize.
`,
		Examples: "juju application-storage postgresql pgdata=20G\njuju resize-storage postgresql/0 --wait",
		SeeAlso:  []string{"application-storage", "show-storage-resize", "storage", "show-storage", "run"},
	})
}

// SetFlags implements Command.SetFlags.
func (c *resizeCommand) SetFlags(f *gnuflag.FlagSet) {
	c.StorageCommandBase.SetFlags(f)
	f.BoolVar(&c.wait, "wait", false, "Wait for all resize results")
	formatters := maps.Clone(output.DefaultFormatters)
	formatters["plain"] = c.formatRequests
	c.out.AddFlags(f, "plain", formatters)
}

// Init implements Command.Init.
func (c *resizeCommand) Init(args []string) error {
	if len(args) != 1 {
		return errors.New("resize-storage requires one unit name")
	}
	if !names.IsValidUnit(args[0]) {
		return errors.NotValidf("unit name %q", args[0])
	}
	c.unitName = args[0]
	return nil
}

// Run implements Command.Run.
func (c *resizeCommand) Run(ctx *cmd.Context) error {
	api, err := c.newAPIFunc(ctx)
	if err != nil {
		return err
	}
	defer api.Close()
	results, err := api.ResizeUnitStorage(ctx, c.unitName)
	if err != nil {
		return err
	}
	out := make(map[string]unitStorageResizeOutput, len(results))
	for _, result := range results {
		tag, err := names.ParseStorageTag(result.StorageTag)
		if err != nil {
			return err
		}
		status := params.StorageResizePending
		if result.AlreadySatisfied {
			status = "already-satisfied"
		}
		out[tag.Id()] = unitStorageResizeOutput{ID: result.ID, RequestedSizeMiB: result.RequestedSizeMiB, Status: status}
		if c.wait && result.ID != "" {
			name, _, _ := strings.Cut(tag.Id(), "/")
			ctx.Infof("  - request %s for %s (%s)", result.ID, name, tag.Id())
		}
	}
	if !c.wait {
		return c.out.Write(ctx, out)
	}
	// Poll all requests so a pending member never hides another member's result.
	for {
		pending, failed := false, false
		for storageID, result := range out {
			if result.Status == params.StorageResizePending {
				details, err := api.StorageResizeDetailsByID(ctx, result.ID)
				if err != nil {
					if writeErr := c.out.Write(ctx, out); writeErr != nil {
						return writeErr
					}
					return errors.Annotatef(err, "waiting for resize request %s", result.ID)
				}
				switch details.Status {
				case params.StorageResizePending, params.StorageResizeCompleted, params.StorageResizeFailed:
				default:
					return errors.Errorf("unknown resize status %q", details.Status)
				}
				if details.Status != params.StorageResizePending {
					ctx.Infof("Request %s for %s: %s (%s)", result.ID, storageID, details.Status, details.Message)
				}
				result.Status, result.Message = details.Status, details.Message
				out[storageID] = result
			}
			pending = pending || result.Status == params.StorageResizePending
			failed = failed || result.Status == params.StorageResizeFailed
		}
		if !pending {
			if err := c.out.Write(ctx, out); err != nil {
				return err
			}
			if failed {
				return cmd.ErrSilent
			}
			return nil
		}
		timer := time.NewTimer(2 * time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			if err := c.out.Write(ctx, out); err != nil {
				return err
			}
			return ctx.Err()
		case <-timer.C:
		}
	}
}

func (c *resizeCommand) formatRequests(writer io.Writer, value any) error {
	results, ok := value.(map[string]unitStorageResizeOutput)
	if !ok {
		return errors.Errorf("unexpected resize result type %T", value)
	}
	var buf bytes.Buffer
	storageIDs := make([]string, 0, len(results))
	requests, pending := 0, false
	for storageID, result := range results {
		storageIDs = append(storageIDs, storageID)
		if result.ID != "" {
			requests++
		}
		pending = pending || result.Status == params.StorageResizePending
	}
	slices.Sort(storageIDs)
	if requests == 0 {
		fmt.Fprintf(&buf, "No storage resize requests needed for %s\n", c.unitName)
	} else {
		fmt.Fprintf(&buf, "Storage resize requests for %s:\n", c.unitName)
	}
	for _, storageID := range storageIDs {
		result := results[storageID]
		if result.ID == "" {
			fmt.Fprintf(&buf, "  - storage %s already satisfies %d MiB\n", storageID, result.RequestedSizeMiB)
			continue
		}
		name, _, _ := strings.Cut(storageID, "/")
		fmt.Fprintf(&buf, "  - request %s for %s (%s)", result.ID, name, storageID)
		if result.Status != params.StorageResizePending {
			fmt.Fprintf(&buf, ": %s", result.Status)
		}
		fmt.Fprintln(&buf)
		if result.Message != "" {
			fmt.Fprintf(&buf, "    %s\n", result.Message)
		}
	}
	if pending {
		fmt.Fprintln(&buf, "Check request status with 'juju show-storage-resize <id>'")
	}
	_, err := writer.Write(buf.Bytes())
	return err
}

type unitStorageResizeOutput struct {
	ID               string `json:"id,omitempty" yaml:"id,omitempty"`
	RequestedSizeMiB uint64 `json:"requested-size-mib" yaml:"requested-size-mib"`
	Status           string `json:"status" yaml:"status"`
	Message          string `json:"message,omitempty" yaml:"message,omitempty"`
}
