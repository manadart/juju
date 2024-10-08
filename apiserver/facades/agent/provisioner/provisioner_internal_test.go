// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provisioner

import (
	"context"

	jc "github.com/juju/testing/checkers"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/container"
	"github.com/juju/juju/core/containermanager"
	"github.com/juju/juju/core/instance"
	modeltesting "github.com/juju/juju/core/model/testing"
	"github.com/juju/juju/rpc/params"
)

// This file contains new provisioner tests which use mocked dependencies, and
// it lives inside the 'provisioner' package (not a separate '_test' package),
// so that we don't have to bend over backwards to test provisioner logic.
// As we refactor the provisioner to use services, we should gradually delete
// the old provisioner tests and write equivalent tests in this file.

type provisionerSuite struct {
	agentProvisionerService *MockAgentProvisionerService
	keyUpdaterService       *MockKeyUpdaterService
}

var _ = gc.Suite(&provisionerSuite{})

func (s *provisionerSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)
	s.agentProvisionerService = NewMockAgentProvisionerService(ctrl)
	s.keyUpdaterService = NewMockKeyUpdaterService(ctrl)
	return ctrl
}

func (s *provisionerSuite) TestContainerManagerConfig(c *gc.C) {
	defer s.setupMocks(c).Finish()
	api := &ProvisionerAPI{
		agentProvisionerService: s.agentProvisionerService,
	}

	modelID := modeltesting.GenModelUUID(c)
	s.agentProvisionerService.EXPECT().ContainerManagerConfigForType(gomock.Any(), instance.LXD).
		Return(containermanager.Config{
			ModelID:                  modelID,
			ImageMetadataURL:         "https://images.linuxcontainers.org/",
			ImageStream:              "released",
			LXDSnapChannel:           "5.0/stable",
			MetadataDefaultsDisabled: true,
		}, nil)
	s.agentProvisionerService.EXPECT().ContainerNetworkingMethod(gomock.Any()).
		Return(containermanager.NetworkingMethodProvider, nil)

	containerManagerConfig, err := api.ContainerManagerConfig(context.Background(),
		params.ContainerManagerConfigParams{Type: instance.LXD},
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(containerManagerConfig, jc.DeepEquals, params.ContainerManagerConfig{ManagerConfig: map[string]string{
		"model-uuid":                                 modelID.String(),
		"lxd-snap-channel":                           "5.0/stable",
		"container-image-metadata-url":               "https://images.linuxcontainers.org/",
		"container-image-metadata-defaults-disabled": "true",
		"container-image-stream":                     "released",
		"container-networking-method":                "provider",
	}})
}

func (s *provisionerSuite) TestContainerConfig(c *gc.C) {
	defer s.setupMocks(c).Finish()
	api := &ProvisionerAPI{
		agentProvisionerService: s.agentProvisionerService,
		keyUpdaterService:       s.keyUpdaterService,
	}

	s.agentProvisionerService.EXPECT().ContainerConfig(gomock.Any()).Return(container.Config{
		EnableOSRefreshUpdate:      true,
		EnableOSUpgrade:            false,
		ProviderType:               "ec2",
		SSLHostnameVerification:    true,
		SnapStoreAssertions:        "trust us",
		SnapStoreProxyID:           "42",
		SnapStoreProxyURL:          "http://snap-store-proxy",
		AptMirror:                  "http://my.archive.ubuntu.com",
		ContainerInheritProperties: "ca-certs,apt-primary",
	}, nil)
	s.keyUpdaterService.EXPECT().GetInitialAuthorisedKeysForContainer(gomock.Any()).Return([]string{
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN8h8XBpjS9aBUG5cdoSWubs7wT2Lc/BEZIUQCqoaOZR juju-client-key",
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN8h8XBpjS9aBUG5cdoSWubs7wT2Lc/BEZIUQCqoaOZR juju-system-key",
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN8h8XBpjS9aBUG5cdoSWubs7wT2Lc/BEZIUQCqoaOZR foo-key",
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN8h8XBpjS9aBUG5cdoSWubs7wT2Lc/BEZIUQCqoaOZR bar-key",
	}, nil)

	containerManagerConfig, err := api.ContainerConfig(context.Background())
	c.Assert(err, jc.ErrorIsNil)
	c.Check(containerManagerConfig, jc.DeepEquals, params.ContainerConfig{
		UpdateBehavior: &params.UpdateBehavior{
			EnableOSRefreshUpdate: true,
			EnableOSUpgrade:       false,
		},
		ProviderType:               "ec2",
		SSLHostnameVerification:    true,
		SnapStoreAssertions:        "trust us",
		SnapStoreProxyID:           "42",
		SnapStoreProxyURL:          "http://snap-store-proxy",
		AptMirror:                  "http://my.archive.ubuntu.com",
		ContainerInheritProperties: "ca-certs,apt-primary",
		AuthorizedKeys: `
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN8h8XBpjS9aBUG5cdoSWubs7wT2Lc/BEZIUQCqoaOZR juju-client-key
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN8h8XBpjS9aBUG5cdoSWubs7wT2Lc/BEZIUQCqoaOZR juju-system-key
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN8h8XBpjS9aBUG5cdoSWubs7wT2Lc/BEZIUQCqoaOZR foo-key
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN8h8XBpjS9aBUG5cdoSWubs7wT2Lc/BEZIUQCqoaOZR bar-key
`[1:],
	})
}
