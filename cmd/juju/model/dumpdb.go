// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package model

import (
	"context"

	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"github.com/juju/names/v6"

	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/core/output"
	"github.com/juju/juju/internal/cmd"
)

// NewDumpDBCommand returns a fully constructed dump-db command.
func NewDumpDBCommand() cmd.Command {
	return modelcmd.Wrap(&dumpDBCommand{})
}

type dumpDBCommand struct {
	modelcmd.ModelCommandBase
	out cmd.Output
	api DumpDBAPI
}

const dumpDBHelpDoc = `
dump-db returns all that is stored in the database for the specified model.

Examples:

    juju dump-db
    juju dump-db -m mymodel

See also:
    models
`

// Info implements Command.
func (c *dumpDBCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    "dump-db",
		Purpose: "Displays the mongo documents for of the model.",
		Doc:     dumpDBHelpDoc,
	})
}

// SetFlags implements Command.
func (c *dumpDBCommand) SetFlags(f *gnuflag.FlagSet) {
	c.ModelCommandBase.SetFlags(f)
	c.out.AddFlags(f, "yaml", output.DefaultFormatters)
}

// Init implements Command.
func (c *dumpDBCommand) Init(args []string) error {
	return cmd.CheckEmpty(args)
}

// DumpDBAPI specifies the used function calls of the ModelManager.
type DumpDBAPI interface {
	Close() error
	DumpModelDB(context.Context, names.ModelTag) (map[string]interface{}, error)
}

func (c *dumpDBCommand) getAPI(ctx context.Context) (DumpDBAPI, error) {
	if c.api != nil {
		return c.api, nil
	}
	return c.ModelCommandBase.NewModelManagerAPIClient(ctx)
}

// Run implements Command.
func (c *dumpDBCommand) Run(ctx *cmd.Context) error {
	client, err := c.getAPI(ctx)
	if err != nil {
		return errors.Trace(err)
	}
	defer client.Close()

	_, modelDetails, err := c.ModelCommandBase.ModelDetails(ctx)
	if err != nil {
		return errors.Annotate(err, "getting model details")
	}

	modelTag := names.NewModelTag(modelDetails.ModelUUID)
	results, err := client.DumpModelDB(ctx, modelTag)
	if err != nil {
		return err
	}

	return c.out.Write(ctx, results)
}
