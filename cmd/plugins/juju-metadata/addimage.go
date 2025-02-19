// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"context"

	"github.com/juju/errors"
	"github.com/juju/gnuflag"

	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/core/arch"
	corebase "github.com/juju/juju/core/base"
	"github.com/juju/juju/environs/simplestreams"
	"github.com/juju/juju/internal/cmd"
	"github.com/juju/juju/rpc/params"
)

func newAddImageMetadataCommand() cmd.Command {
	return modelcmd.Wrap(&addImageMetadataCommand{})
}

const addImageCommandDoc = `
Add image metadata to Juju model.

Image metadata properties vary between providers. Consequently, some properties
are optional for this command but they may still be needed by your provider.

Adding an image for a specific base can be done via --base. --base can be 
specified using the OS name and the version of the OS, separated by @. For 
example, --base ubuntu@22.04.

Valid values for --stream are released, testing, proposed and devel. The image
stream used by Juju can be configured with 'juju model-config'.
`

// addImageMetadataCommand stores image metadata in Juju environment.
type addImageMetadataCommand struct {
	cloudImageMetadataCommandBase

	ImageId         string
	Region          string
	Base            string
	Version         string
	Arch            string
	VirtType        string
	RootStorageType string
	RootStorageSize uint64
	Stream          string
}

// Init implements Command.Init.
func (c *addImageMetadataCommand) Init(args []string) (err error) {
	if len(args) == 0 {
		return errors.New("image id must be supplied when adding image metadata")
	}
	if len(args) != 1 {
		return errors.New("only one image id can be supplied as an argument to this command")
	}
	c.ImageId = args[0]
	return nil
}

// Info implements Command.Info.
func (c *addImageMetadataCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    "add-image",
		Purpose: "adds image metadata to model",
		Doc:     addImageCommandDoc,
		SeeAlso: []string{
			"delete-image",
			"list-images",
			"model-config",
		},
	})
}

// SetFlags implements Command.SetFlags.
func (c *addImageMetadataCommand) SetFlags(f *gnuflag.FlagSet) {
	c.cloudImageMetadataCommandBase.SetFlags(f)

	f.StringVar(&c.Region, "region", "", "image cloud region")
	f.StringVar(&c.Base, "base", "", "image base")
	f.StringVar(&c.Arch, "arch", arch.AMD64, "image architecture")
	f.StringVar(&c.VirtType, "virt-type", "", "image metadata virtualisation type")
	f.StringVar(&c.RootStorageType, "storage-type", "", "image metadata root storage type")
	f.Uint64Var(&c.RootStorageSize, "storage-size", 0, "image metadata root storage size")
	f.StringVar(&c.Stream, "stream", "released", "image metadata stream")
}

// Run implements Command.Run.
func (c *addImageMetadataCommand) Run(ctx *cmd.Context) error {
	var (
		base corebase.Base
		err  error
	)
	if c.Base != "" {
		if base, err = corebase.ParseBaseFromString(c.Base); err != nil {
			return errors.Trace(err)
		}
	}

	api, err := getImageMetadataAddAPI(c, ctx)
	if err != nil {
		return err
	}
	defer api.Close()

	m := c.constructMetadataParam(base)
	if err := api.Save(ctx, []params.CloudImageMetadata{m}); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// MetadataAddAPI defines the API methods that add image metadata command uses.
type MetadataAddAPI interface {
	Close() error
	Save(ctx context.Context, metadata []params.CloudImageMetadata) error
}

var getImageMetadataAddAPI = (*addImageMetadataCommand).getImageMetadataAddAPI

func (c *addImageMetadataCommand) getImageMetadataAddAPI(ctx context.Context) (MetadataAddAPI, error) {
	return c.NewImageMetadataAPI(ctx)
}

// constructMetadataParam returns cloud image metadata as a param.
func (c *addImageMetadataCommand) constructMetadataParam(base corebase.Base) params.CloudImageMetadata {
	info := params.CloudImageMetadata{
		ImageId: c.ImageId,
		Region:  c.Region,
		// TODO (stickupkid): Allow passing the channel risk through to the API
		// to target an image. Currently limited to track only.
		Version:         base.Channel.Track,
		Arch:            c.Arch,
		VirtType:        c.VirtType,
		RootStorageType: c.RootStorageType,
		Stream:          c.Stream,
		Source:          "custom",
		Priority:        simplestreams.CUSTOM_CLOUD_DATA,
	}
	if c.RootStorageSize != 0 {
		info.RootStorageSize = &c.RootStorageSize
	}
	return info
}
