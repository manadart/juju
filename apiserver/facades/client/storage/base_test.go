// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage_test

import (
	"context"

	"github.com/juju/names/v6"
	"github.com/juju/tc"
	"go.uber.org/mock/gomock"

	"github.com/juju/juju/apiserver/facades/client/storage"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/core/machine"
	coremodel "github.com/juju/juju/core/model"
	modeltesting "github.com/juju/juju/core/model/testing"
	"github.com/juju/juju/core/unit"
	coretesting "github.com/juju/juju/internal/testing"
	"github.com/juju/juju/internal/uuid"
)

type baseStorageSuite struct {
	coretesting.BaseSuite

	authorizer apiservertesting.FakeAuthorizer

	controllerUUID string
	modelUUID      coremodel.UUID

	api     *storage.StorageAPI
	apiCaas *storage.StorageAPI

	unitTag    names.UnitTag
	machineTag names.MachineTag

	blockDeviceService  *storage.MockBlockDeviceService
	blockCommandService *storage.MockBlockCommandService
	storageService      *storage.MockStorageService
	applicationService  *storage.MockApplicationService
	removalService      *storage.MockRemovalService

	poolsInUse []string
}

func (s *baseStorageSuite) setupMocks(c *tc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.unitTag = names.NewUnitTag("mysql/0")
	s.machineTag = names.NewMachineTag("1234")

	s.authorizer = apiservertesting.FakeAuthorizer{Tag: names.NewUserTag("admin"), Controller: true}

	s.blockDeviceService = storage.NewMockBlockDeviceService(ctrl)
	s.blockCommandService = storage.NewMockBlockCommandService(ctrl)
	s.storageService = storage.NewMockStorageService(ctrl)
	s.applicationService = storage.NewMockApplicationService(ctrl)
	s.removalService = storage.NewMockRemovalService(ctrl)

	s.applicationService.EXPECT().GetUnitMachineName(
		gomock.Any(), unit.Name("mysql/0"),
	).DoAndReturn(func(ctx context.Context, u unit.Name) (machine.Name, error) {
		c.Assert(u.String(), tc.Equals, s.unitTag.Id())
		return machine.Name(s.machineTag.Id()), nil
	}).AnyTimes()

	s.poolsInUse = []string{}

	s.controllerUUID = uuid.MustNewUUID().String()
	s.modelUUID = modeltesting.GenModelUUID(c)
	s.api = storage.NewStorageAPI(
		s.controllerUUID, s.modelUUID,
		s.blockDeviceService,
		s.storageService, s.applicationService,
		s.authorizer, s.blockCommandService)
	s.apiCaas = storage.NewStorageAPI(
		s.controllerUUID, s.modelUUID,
		s.blockDeviceService,
		s.storageService, s.applicationService,
		s.authorizer, s.blockCommandService)

	return ctrl
}
