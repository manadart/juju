// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"context"
	"fmt"
	"io"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"github.com/juju/names/v6"

	"github.com/juju/juju/api/client/application"
	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/cmd"
	"github.com/juju/juju/cmd/juju/config"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/core/output"
	"github.com/juju/juju/core/storage"
)

const (
	appStorageConfigSummary = "Displays or sets storage directives for an application."

	appStorageConfigDoc = `
View or update storage allocation defaults for an application.
A directive has the form <storage-name>=<size>,<pool>,<count>.
Omitted fields retain their current values. For example, data=20G updates only
capacity and preserves the pool and allocation count.

Updating a directive does not resize existing storage or run charm actions.
New units allocate storage using the updated defaults. To apply capacity targets
to an existing unit, run resize-storage <unit>, optionally with --wait.
Pool and count changes do not move, add or remove existing storage.

Directive updates currently require a machine model. Directives may be viewed
in either model type. Use --format yaml or --format json for structured output.
`
	appStorageConfigExamples = `
View all allocation defaults:
    juju application-storage postgresql

View a named directive:
    juju application-storage postgresql pgdata

Update the capacity default, then resize one unit:
    juju application-storage postgresql pgdata=20G
    juju resize-storage postgresql/0 --wait
`
)

// NewStorageCommand returns the application storage directive command.
func NewStorageCommand() cmd.Command {
	applicationStorageConfigBase := config.ConfigCommandBase{
		Resettable: false,
	}
	return modelcmd.Wrap(&applicationStorageCommand{configBase: applicationStorageConfigBase})
}

type applicationStorageCommand struct {
	modelcmd.ModelCommandBase
	configBase config.ConfigCommandBase
	out        cmd.Output

	api             StorageDirectivesAPI
	applicationName string
}

// Info returns the name and usage details for the application storage command,
// providing guidance to the user.
// It implements the [cmd.Command] interface.
func (c *applicationStorageCommand) Info() *cmd.Info {
	info := &cmd.Info{
		Name:     "application-storage",
		Args:     "<application-name> [<storage-name>[={<size>,<pool>,<count>}]] ...",
		Purpose:  appStorageConfigSummary,
		Examples: appStorageConfigExamples,
		Doc:      appStorageConfigDoc,
		SeeAlso: []string{
			"storage",
			"storage-pools",
			"resize-storage",
			"add-storage",
			"add-unit",
		},
	}

	return jujucmd.Info(info)
}

// SetFlags implements part of the cmd.Command interface.
func (c *applicationStorageCommand) SetFlags(f *gnuflag.FlagSet) {
	c.ModelCommandBase.SetFlags(f)

	c.out.AddFlags(f, "tabular", map[string]cmd.Formatter{
		"json":    cmd.FormatJson,
		"tabular": formatConfigTabular,
		"yaml":    cmd.FormatYaml,
	})
}

// Init parses the command arguments and validates the application name.
// It implements the [cmd.Command] interface.
func (c *applicationStorageCommand) Init(args []string) error {
	nArgs := len(args)
	if nArgs == 0 {
		return errors.Errorf("no application specified")
	}

	if !names.IsValidApplication(args[0]) {
		return errors.Errorf("invalid application name %q", args[0])
	}
	c.applicationName = args[0]
	return c.configBase.Init(args[1:])
}

// getAPI returns the API. This allows passing in a test configCommandAPI
// implementation.
func (c *applicationStorageCommand) getAPI(ctx context.Context) (StorageDirectivesAPI, error) {
	if c.api != nil {
		return c.api, nil
	}
	api, err := c.NewAPIRoot(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "opening API connection")
	}
	client := application.NewClient(api)
	return client, nil
}

// Run executes the main logic for the application storage command.
// It processes the configured actions by either retrieving the
// storage directive for a single storage name, setting storage
// directives, or listing all storage directives for the given application.
// It implements the [cmd.Command] interface.
func (c *applicationStorageCommand) Run(ctx *cmd.Context) error {
	client, err := c.getAPI(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	for _, action := range c.configBase.Actions {
		var err error
		switch action {
		case config.GetOne:
			err = c.getConfig(ctx, client)
		case config.SetArgs:
			err = c.setConfig(ctx, client, c.configBase.ValsToSet)
		default:
			err = c.getAllConfig(ctx, client)
		}
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

// StorageDirectivesAPI supplies allocation directive reads and updates.
type StorageDirectivesAPI interface {
	Close() error
	GetApplicationStorageWithContext(context.Context, string) (application.ApplicationStorageInfo, error)
	SetApplicationStorage(context.Context, string, map[string]storage.DirectiveUpdate) error
}

// setConfig sets the provided key/value pairs on the application.
func (c *applicationStorageCommand) setConfig(ctx context.Context, client StorageDirectivesAPI, attrs config.Attrs) error {
	updates := make(map[string]storage.DirectiveUpdate, len(attrs))
	for name, value := range attrs {
		update, err := storage.ParseDirectiveUpdate(fmt.Sprint(value))
		if err != nil {
			return errors.Annotatef(err, "parsing storage directive %q", name)
		}
		updates[name] = update
	}
	return client.SetApplicationStorage(ctx, c.applicationName, updates)
}

// getConfig writes one named storage directive.
func (c *applicationStorageCommand) getConfig(ctx *cmd.Context, client StorageDirectivesAPI) error {
	applicationStorageInfo, err := client.GetApplicationStorageWithContext(ctx, c.applicationName)
	if err != nil {
		return err
	}

	if len(c.configBase.KeysToGet) == 0 {
		return errors.New("no storage key specified")
	}

	if applicationStorageInfo.Error != nil {
		return applicationStorageInfo.Error
	}

	storeKey := c.configBase.KeysToGet[0]
	storageConsForKey, ok := applicationStorageInfo.StorageConstraints[storeKey]
	if !ok {
		return errors.NotFoundf("storage %q", storeKey)
	}

	// Retain the storage name in structured and tabular output.
	storageConsForKeyMap := map[string]storage.Directive{
		storeKey: storageConsForKey,
	}

	return c.out.Write(ctx, storageConsForKeyMap)
}

// getAllConfig returns the entire configuration for the selected application.
func (c *applicationStorageCommand) getAllConfig(ctx *cmd.Context, client StorageDirectivesAPI) error {
	applicationStorageInfo, err := client.GetApplicationStorageWithContext(ctx, c.applicationName)
	if err != nil {
		return err
	}

	if applicationStorageInfo.Error != nil {
		return applicationStorageInfo.Error
	}
	return c.out.Write(ctx, applicationStorageInfo.StorageConstraints)
}

// formatConfigTabular writes a tabular summary of config information.
func formatConfigTabular(writer io.Writer, value interface{}) error {
	configValues, ok := value.(map[string]storage.Directive)
	if !ok {
		return errors.Errorf("expected value of type %T, got %T", map[string]storage.Directive{}, value)
	}

	tw := output.TabWriter(writer)
	w := output.Wrapper{
		TabWriter: tw,
	}

	var valueNames []string
	for name := range configValues {
		valueNames = append(valueNames, name)
	}
	sort.Strings(valueNames)
	w.Println("Storage", "Pool", "Size", "Count")

	for _, name := range valueNames {
		info := configValues[name]
		w.Println(name, info.Pool, humanize.IBytes(info.Size*humanize.MiByte), info.Count)
	}

	return tw.Flush()
}
