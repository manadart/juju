// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package scriptletcharm

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/juju/errors"
	"github.com/juju/names/v6"
	"gopkg.in/yaml.v2"

	"github.com/juju/juju/api/base"
	apiscriptletcharm "github.com/juju/juju/api/client/scriptletcharm"
	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/cmd"
	"github.com/juju/juju/cmd/juju/block"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/rpc/params"
)

const deployDoc = `
Registers a scriptlet charm with the current model.

The command records the supplied Starform scriptlet source as raw text for the
application name. It does not yet create a normal Juju application, provision
units, or run the scriptlet.
`

const deployExamples = `
    juju deploy-scriptlet-charm ./scriptlet
    juju deploy-scriptlet-charm ./scriptlet/scriptlets/hooks.star model-info
`

// NewDeployCommand returns a command to register scriptlet charms.
func NewDeployCommand() cmd.Command {
	command := &deployCommand{}
	command.newClient = func(caller base.APICallCloser) ScriptletCharmAPI {
		return apiscriptletcharm.NewClient(caller)
	}
	return modelcmd.Wrap(command)
}

// ScriptletCharmAPI defines the client methods required by the command.
type ScriptletCharmAPI interface {
	Register(ctx context.Context, args params.RegisterScriptletCharmArgs) error
	Close() error
}

type deployCommand struct {
	modelcmd.ModelCommandBase

	scriptletPath   string
	applicationName string
	newClient       func(base.APICallCloser) ScriptletCharmAPI
}

// Info implements Command.Info.
func (c *deployCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:     "deploy-scriptlet-charm",
		Args:     "<scriptlet charm or .star file> [<application name>]",
		Purpose:  "Registers a scriptlet charm.",
		Doc:      deployDoc,
		Examples: deployExamples,
		SeeAlso: []string{
			"deploy",
		},
	})
}

// Init implements Command.Init.
func (c *deployCommand) Init(args []string) error {
	if len(args) < 1 {
		return errors.New("deploy-scriptlet-charm requires a scriptlet charm or .star file")
	}
	c.scriptletPath = args[0]
	if len(args) > 1 {
		c.applicationName = args[1]
		if err := names.ValidateApplicationName(c.applicationName); err != nil {
			return errors.Trace(err)
		}
	}
	if len(args) < 3 {
		return nil
	}
	return cmd.CheckEmpty(args[2:])
}

// Run implements Command.Run.
func (c *deployCommand) Run(ctx *cmd.Context) error {
	playClintEastwood(ctx, ctx.Stdout)

	registerArgs, defaultApplicationName, err := readScriptlet(ctx.AbsPath(c.scriptletPath))
	if err != nil {
		return errors.Trace(err)
	}
	if strings.TrimSpace(registerArgs.Scriptlet) == "" {
		return errors.New("scriptlet file is empty")
	}

	applicationName := c.applicationName
	if applicationName == "" {
		applicationName = defaultApplicationName
	}
	if applicationName == "" {
		return errors.New("application name not specified and metadata.yaml has no name")
	}
	if err := names.ValidateApplicationName(applicationName); err != nil {
		return errors.Trace(err)
	}
	registerArgs.ApplicationName = applicationName

	apiRoot, err := c.NewAPIRoot(ctx)
	if err != nil {
		return errors.Trace(err)
	}
	client := c.newClient(apiRoot)
	defer client.Close()

	err = client.Register(ctx, registerArgs)
	return block.ProcessBlockedError(err, block.BlockChange)
}

type scriptletConfig struct {
	Sources []string `yaml:"sources"`
}

type charmMetadataRelation struct {
	Interface string `yaml:"interface"`
	Scope     string `yaml:"scope"`
	Optional  bool   `yaml:"optional"`
	Limit     int    `yaml:"limit"`
}

type charmMetadata struct {
	Name     string                           `yaml:"name"`
	Provides map[string]charmMetadataRelation `yaml:"provides"`
	Requires map[string]charmMetadataRelation `yaml:"requires"`
	Peers    map[string]charmMetadataRelation `yaml:"peers"`
}

func readScriptlet(path string) (params.RegisterScriptletCharmArgs, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return params.RegisterScriptletCharmArgs{}, "", errors.Annotatef(err, "checking scriptlet path %q", path)
	}
	if info.IsDir() {
		return readScriptletCharmDir(path)
	}
	if filepath.Ext(path) != ".star" {
		return params.RegisterScriptletCharmArgs{}, "", errors.Errorf("scriptlet file %q must have .star extension", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return params.RegisterScriptletCharmArgs{}, "", errors.Annotatef(err, "reading scriptlet file %q", path)
	}
	return params.RegisterScriptletCharmArgs{Scriptlet: string(data)}, "", nil
}

func readScriptletCharmDir(path string) (params.RegisterScriptletCharmArgs, string, error) {
	metadata, err := readCharmMetadata(path)
	if err != nil {
		return params.RegisterScriptletCharmArgs{}, "", errors.Trace(err)
	}

	config, err := readScriptletConfig(path)
	if err != nil {
		return params.RegisterScriptletCharmArgs{}, "", errors.Trace(err)
	}
	if len(config.Sources) != 1 {
		return params.RegisterScriptletCharmArgs{}, "", errors.Errorf("expected exactly one scriptlet source, got %d", len(config.Sources))
	}

	source := filepath.Clean(config.Sources[0])
	if filepath.IsAbs(source) || source == ".." || strings.HasPrefix(source, ".."+string(os.PathSeparator)) {
		return params.RegisterScriptletCharmArgs{}, "", errors.Errorf("scriptlet source %q escapes charm directory", config.Sources[0])
	}
	if filepath.Ext(source) != ".star" {
		return params.RegisterScriptletCharmArgs{}, "", errors.Errorf("scriptlet source %q must have .star extension", config.Sources[0])
	}

	scriptletPath := filepath.Join(path, source)
	data, err := os.ReadFile(scriptletPath)
	if err != nil {
		return params.RegisterScriptletCharmArgs{}, "", errors.Annotatef(err, "reading scriptlet file %q", scriptletPath)
	}

	args := params.RegisterScriptletCharmArgs{
		Scriptlet: string(data),
		Relations: encodeMetadataRelations(metadata),
	}
	return args, metadata.Name, nil
}

func encodeMetadataRelations(metadata charmMetadata) []params.ScriptletRelation {
	var relations []params.ScriptletRelation
	for name, r := range metadata.Provides {
		relations = append(relations, params.ScriptletRelation{
			Name: name, Role: "provider", Interface: r.Interface,
			Scope: r.Scope, Optional: r.Optional, Limit: r.Limit,
		})
	}
	for name, r := range metadata.Requires {
		relations = append(relations, params.ScriptletRelation{
			Name: name, Role: "requirer", Interface: r.Interface,
			Scope: r.Scope, Optional: r.Optional, Limit: r.Limit,
		})
	}
	for name, r := range metadata.Peers {
		relations = append(relations, params.ScriptletRelation{
			Name: name, Role: "peer", Interface: r.Interface,
			Scope: r.Scope, Optional: r.Optional, Limit: r.Limit,
		})
	}
	sort.Slice(relations, func(i, j int) bool { return relations[i].Name < relations[j].Name })
	return relations
}

func readCharmMetadata(path string) (charmMetadata, error) {
	data, err := os.ReadFile(filepath.Join(path, "metadata.yaml"))
	if err != nil {
		return charmMetadata{}, errors.Annotate(err, "reading metadata.yaml")
	}
	var metadata charmMetadata
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return charmMetadata{}, errors.Annotate(err, "parsing metadata.yaml")
	}
	return metadata, nil
}

func readScriptletConfig(path string) (scriptletConfig, error) {
	data, err := os.ReadFile(filepath.Join(path, "scriptlet.yaml"))
	if err != nil {
		return scriptletConfig{}, errors.Annotate(err, "reading scriptlet.yaml")
	}

	var config scriptletConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return scriptletConfig{}, errors.Annotate(err, "parsing scriptlet.yaml")
	}
	return config, nil
}
