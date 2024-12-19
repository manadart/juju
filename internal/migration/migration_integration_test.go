// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package migration_test

import (
	"context"

	"github.com/juju/clock"
	jc "github.com/juju/testing/checkers"
	gomock "go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	model "github.com/juju/juju/core/model"
	"github.com/juju/juju/core/modelmigration"
	corestorage "github.com/juju/juju/core/storage"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	"github.com/juju/juju/internal/migration"
	"github.com/juju/juju/internal/storage"
	"github.com/juju/juju/internal/storage/provider"
	jujutesting "github.com/juju/juju/internal/testing"
	"github.com/juju/juju/state"
)

type ExportImportSuite struct {
	controllerConfigService *MockControllerConfigService
	domainServices          *MockDomainServices
	domainServicesGetter    *MockDomainServicesGetter
	objectStoreGetter       *MockModelObjectStoreGetter
}

var _ = gc.Suite(&ExportImportSuite{})

func (s *ExportImportSuite) SetUpSuite(c *gc.C) {
	c.Skip(`
TODO tlm: We are skipping these tests as they are currently relying heavily on
mocks for how the importer is working internally. Now that we are trying to test
model migration into DQlite we have hit the problem where this can no longer be
an isolation suite and needs a full database. This is due to the fact that the
Setup call to the import operations construct their own services and they're not
something that can be injected as "mocks" from this layer.

I have added this to the risk register for 4.0 and will discuss further with
Simon. For the moment these tests can't continue as is due to DB dependencies
that are needed now.
`)
}

func (s *ExportImportSuite) exportImport(c *gc.C, leaders map[string]string) {
	bytes := []byte(modelYaml)
	st := &state.State{}
	m := &state.Model{}
	controller := &fakeImporter{st: st, m: m}
	scope := func(model.UUID) modelmigration.Scope { return modelmigration.NewScope(nil, nil, nil) }
	importer := migration.NewModelImporter(
		controller, scope, s.controllerConfigService, s.domainServicesGetter,
		corestorage.ConstModelStorageRegistry(func() storage.ProviderRegistry {
			return provider.CommonStorageProviders()
		}),
		s.objectStoreGetter,
		loggertesting.WrapCheckLog(c),
		clock.WallClock,
	)
	gotM, gotSt, err := importer.ImportModel(context.Background(), bytes)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(controller.model.Tag().Id(), gc.Equals, "bd3fae18-5ea1-4bc5-8837-45400cf1f8f6")
	c.Assert(gotM, gc.Equals, m)
	c.Assert(gotSt, gc.Equals, st)
}

func (s *ExportImportSuite) TestExportImportModel(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.exportImport(c, map[string]string{})
}

func (s *ExportImportSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.controllerConfigService = NewMockControllerConfigService(ctrl)
	s.controllerConfigService.EXPECT().ControllerConfig(gomock.Any()).Return(jujutesting.FakeControllerConfig(), nil).AnyTimes()

	s.domainServices = NewMockDomainServices(ctrl)
	s.domainServices.EXPECT().Cloud().Return(nil).AnyTimes()
	s.domainServices.EXPECT().Credential().Return(nil).AnyTimes()
	s.domainServices.EXPECT().Machine().Return(nil)
	s.domainServices.EXPECT().Application().Return(nil)
	s.domainServicesGetter = NewMockDomainServicesGetter(ctrl)
	s.domainServicesGetter.EXPECT().ServicesForModel(model.UUID("bd3fae18-5ea1-4bc5-8837-45400cf1f8f6")).Return(s.domainServices)
	s.objectStoreGetter = NewMockModelObjectStoreGetter(ctrl)

	return ctrl
}