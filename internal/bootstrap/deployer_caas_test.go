// Copyright 2023 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package bootstrap

import (
	"testing"

	"github.com/juju/errors"
	"github.com/juju/tc"
	"go.uber.org/mock/gomock"

	coreapplication "github.com/juju/juju/core/application"
	corecharm "github.com/juju/juju/core/charm"
	network "github.com/juju/juju/core/network"
	unit "github.com/juju/juju/core/unit"
	"github.com/juju/juju/core/version"
	applicationservice "github.com/juju/juju/domain/application/service"
	deploymentcharm "github.com/juju/juju/domain/deployment/charm"
	"github.com/juju/juju/environs/bootstrap"
	"github.com/juju/juju/internal/uuid"
)

type deployerCAASSuite struct {
	baseSuite
	serviceManager *MockServiceManager
}

func TestDeployerCAASSuite(t *testing.T) {
	tc.Run(t, &deployerCAASSuite{})
}

func (s *deployerCAASSuite) TestValidate(c *tc.C) {
	defer s.setupMocks(c).Finish()

	cfg := s.newConfig(c)
	err := cfg.Validate()
	c.Assert(err, tc.IsNil)

	cfg = s.newConfig(c)
	cfg.ServiceManager = nil
	err = cfg.Validate()
	c.Assert(err, tc.ErrorIs, errors.NotValid)
}

func (s *deployerCAASSuite) TestControllerCharmBase(c *tc.C) {
	defer s.setupMocks(c).Finish()

	deployer := s.newDeployer(c)
	base, err := deployer.ControllerCharmBase()
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(base, tc.DeepEquals, version.DefaultSupportedLTSBase())
}

func (s *deployerCAASSuite) TestAddCAASControllerApplicationSetsApplicationPassword(c *tc.C) {
	defer s.setupMocks(c).Finish()

	cfg := s.newConfig(c)
	appID := coreapplication.UUID("controller-app-uuid")
	s.caasApplicationService.EXPECT().CreateCAASApplication(
		gomock.Any(),
		bootstrap.ControllerApplicationName,
		s.charm,
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(appID, nil)
	s.agentPasswordService.EXPECT().SetApplicationPassword(gomock.Any(), appID, cfg.UnitPassword).Return(nil)

	deployer := s.newDeployerWithConfig(c, cfg)
	origin := corecharm.Origin{
		Source:   corecharm.CharmHub,
		Type:     "charm",
		Channel:  &deploymentcharm.Channel{},
		Revision: new(1),
		Hash:     "sha-256",
		Platform: corecharm.Platform{
			Architecture: "arm64",
			OS:           "ubuntu",
			Channel:      "22.04",
		},
	}
	err := deployer.AddCAASControllerApplication(c.Context(), DeployCharmInfo{
		URL:             deploymentcharm.MustParseURL("ch:juju-controller-0"),
		Charm:           s.charm,
		Origin:          &origin,
		ArchivePath:     "path",
		ObjectStoreUUID: "1234",
	})
	c.Assert(err, tc.ErrorIsNil)
}

func (s *deployerCAASSuite) TestCompleteCAASProcess(c *tc.C) {
	defer s.setupMocks(c).Finish()

	cfg := s.newConfig(c)

	unitName := unit.Name("controller/0")

	providerAddress := network.ProviderAddresses{
		{
			MachineAddress: network.MachineAddress{
				Value: "10.0.0.1",
				Type:  network.IPv4Address,
				Scope: network.ScopeMachineLocal,
			},
		},
		{
			MachineAddress: network.MachineAddress{
				Value: "203.0.113.1",
				Type:  network.IPv4Address,
				Scope: network.ScopePublic,
			},
		},
	}

	s.caasApplicationService.EXPECT().UpdateK8sService(gomock.Any(), bootstrap.ControllerApplicationName, controllerProviderID(unitName), providerAddress).Return(nil)
	s.caasApplicationService.EXPECT().UpdateCAASUnit(gomock.Any(), unitName, applicationservice.UpdateCAASUnitParams{
		ProviderID: new("controller-0"),
	})
	s.agentPasswordService.EXPECT().SetUnitPassword(gomock.Any(), unitName, cfg.UnitPassword)

	deployer := s.newDeployerWithConfig(c, cfg)
	err := deployer.CompleteCAASProcess(c.Context())
	c.Assert(err, tc.ErrorIsNil)
}

func (s *deployerCAASSuite) newDeployer(c *tc.C) *CAASDeployer {
	return s.newDeployerWithConfig(c, s.newConfig(c))
}

func (s *deployerCAASSuite) newDeployerWithConfig(c *tc.C, cfg CAASDeployerConfig) *CAASDeployer {
	deployer, err := NewCAASDeployer(cfg)
	c.Assert(err, tc.IsNil)
	return deployer
}

func (s *deployerCAASSuite) setupMocks(c *tc.C) *gomock.Controller {
	ctrl := s.baseSuite.setupMocks(c)

	s.serviceManager = NewMockServiceManager(ctrl)

	return ctrl
}

func (s *deployerCAASSuite) newConfig(c *tc.C) CAASDeployerConfig {
	return CAASDeployerConfig{
		BaseDeployerConfig: s.baseSuite.newConfig(c),
		ApplicationService: s.caasApplicationService,
		UnitPassword:       uuid.MustNewUUID().String(),
		ServiceManager:     s.serviceManager,
	}
}
