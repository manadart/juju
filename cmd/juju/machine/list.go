// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package machine

import (
	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/internal/cmd"
)

var usageListMachinesSummary = `
Lists machines in a model.`[1:]

var usageListMachinesDetails = `
By default, the tabular format is used.
The following sections are included: ID, STATE, DNS, INS-ID, SERIES, AZ
Note: AZ above is the cloud region's availability zone.

`

const usageListMachinesExamples = `
     juju machines
`

// NewListMachineCommand returns a command that lists the machines in a model.
func NewListMachinesCommand() cmd.Command {
	return modelcmd.Wrap(newListMachinesCommand(nil))
}

func newListMachinesCommand(api statusAPI) *listMachinesCommand {
	listCmd := &listMachinesCommand{}
	listCmd.defaultFormat = "tabular"
	listCmd.api = api
	return listCmd
}

// listMachineCommand holds information about machines in a model.
type listMachinesCommand struct {
	baselistMachinesCommand
}

// Info implements Command.Info.
func (c *listMachinesCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:     "machines",
		Purpose:  usageListMachinesSummary,
		Doc:      usageListMachinesDetails,
		Aliases:  []string{"list-machines"},
		Examples: usageListMachinesExamples,
		SeeAlso: []string{
			"status",
		},
	})
}

// Init ensures the machines Command does not take arguments.
func (c *listMachinesCommand) Init(args []string) error {
	return cmd.CheckEmpty(args)
}
