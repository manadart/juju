// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/services (interfaces: ModelDomainServices)
//
// Generated by this command:
//
//	mockgen -typed -package mocks -destination mocks/services_mocks.go github.com/juju/juju/internal/services ModelDomainServices
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	service "github.com/juju/juju/domain/agentprovisioner/service"
	service0 "github.com/juju/juju/domain/annotation/service"
	service1 "github.com/juju/juju/domain/application/service"
	service2 "github.com/juju/juju/domain/blockcommand/service"
	service3 "github.com/juju/juju/domain/blockdevice/service"
	service4 "github.com/juju/juju/domain/cloudimagemetadata/service"
	service5 "github.com/juju/juju/domain/keymanager/service"
	service6 "github.com/juju/juju/domain/keyupdater/service"
	service7 "github.com/juju/juju/domain/machine/service"
	service8 "github.com/juju/juju/domain/model/service"
	service9 "github.com/juju/juju/domain/modelagent/service"
	service10 "github.com/juju/juju/domain/modelconfig/service"
	service11 "github.com/juju/juju/domain/modelmigration/service"
	service12 "github.com/juju/juju/domain/network/service"
	service13 "github.com/juju/juju/domain/port/service"
	service14 "github.com/juju/juju/domain/proxy/service"
	service15 "github.com/juju/juju/domain/secret/service"
	service16 "github.com/juju/juju/domain/secretbackend/service"
	service17 "github.com/juju/juju/domain/storage/service"
	stub "github.com/juju/juju/domain/stub"
	service18 "github.com/juju/juju/domain/unitstate/service"
	gomock "go.uber.org/mock/gomock"
)

// MockModelDomainServices is a mock of ModelDomainServices interface.
type MockModelDomainServices struct {
	ctrl     *gomock.Controller
	recorder *MockModelDomainServicesMockRecorder
}

// MockModelDomainServicesMockRecorder is the mock recorder for MockModelDomainServices.
type MockModelDomainServicesMockRecorder struct {
	mock *MockModelDomainServices
}

// NewMockModelDomainServices creates a new mock instance.
func NewMockModelDomainServices(ctrl *gomock.Controller) *MockModelDomainServices {
	mock := &MockModelDomainServices{ctrl: ctrl}
	mock.recorder = &MockModelDomainServicesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelDomainServices) EXPECT() *MockModelDomainServicesMockRecorder {
	return m.recorder
}

// Agent mocks base method.
func (m *MockModelDomainServices) Agent() *service9.ModelService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Agent")
	ret0, _ := ret[0].(*service9.ModelService)
	return ret0
}

// Agent indicates an expected call of Agent.
func (mr *MockModelDomainServicesMockRecorder) Agent() *MockModelDomainServicesAgentCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Agent", reflect.TypeOf((*MockModelDomainServices)(nil).Agent))
	return &MockModelDomainServicesAgentCall{Call: call}
}

// MockModelDomainServicesAgentCall wrap *gomock.Call
type MockModelDomainServicesAgentCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesAgentCall) Return(arg0 *service9.ModelService) *MockModelDomainServicesAgentCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesAgentCall) Do(f func() *service9.ModelService) *MockModelDomainServicesAgentCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesAgentCall) DoAndReturn(f func() *service9.ModelService) *MockModelDomainServicesAgentCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// AgentProvisioner mocks base method.
func (m *MockModelDomainServices) AgentProvisioner() *service.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AgentProvisioner")
	ret0, _ := ret[0].(*service.Service)
	return ret0
}

// AgentProvisioner indicates an expected call of AgentProvisioner.
func (mr *MockModelDomainServicesMockRecorder) AgentProvisioner() *MockModelDomainServicesAgentProvisionerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AgentProvisioner", reflect.TypeOf((*MockModelDomainServices)(nil).AgentProvisioner))
	return &MockModelDomainServicesAgentProvisionerCall{Call: call}
}

// MockModelDomainServicesAgentProvisionerCall wrap *gomock.Call
type MockModelDomainServicesAgentProvisionerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesAgentProvisionerCall) Return(arg0 *service.Service) *MockModelDomainServicesAgentProvisionerCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesAgentProvisionerCall) Do(f func() *service.Service) *MockModelDomainServicesAgentProvisionerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesAgentProvisionerCall) DoAndReturn(f func() *service.Service) *MockModelDomainServicesAgentProvisionerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Annotation mocks base method.
func (m *MockModelDomainServices) Annotation() *service0.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Annotation")
	ret0, _ := ret[0].(*service0.Service)
	return ret0
}

// Annotation indicates an expected call of Annotation.
func (mr *MockModelDomainServicesMockRecorder) Annotation() *MockModelDomainServicesAnnotationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Annotation", reflect.TypeOf((*MockModelDomainServices)(nil).Annotation))
	return &MockModelDomainServicesAnnotationCall{Call: call}
}

// MockModelDomainServicesAnnotationCall wrap *gomock.Call
type MockModelDomainServicesAnnotationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesAnnotationCall) Return(arg0 *service0.Service) *MockModelDomainServicesAnnotationCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesAnnotationCall) Do(f func() *service0.Service) *MockModelDomainServicesAnnotationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesAnnotationCall) DoAndReturn(f func() *service0.Service) *MockModelDomainServicesAnnotationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Application mocks base method.
func (m *MockModelDomainServices) Application() *service1.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Application")
	ret0, _ := ret[0].(*service1.WatchableService)
	return ret0
}

// Application indicates an expected call of Application.
func (mr *MockModelDomainServicesMockRecorder) Application() *MockModelDomainServicesApplicationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Application", reflect.TypeOf((*MockModelDomainServices)(nil).Application))
	return &MockModelDomainServicesApplicationCall{Call: call}
}

// MockModelDomainServicesApplicationCall wrap *gomock.Call
type MockModelDomainServicesApplicationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesApplicationCall) Return(arg0 *service1.WatchableService) *MockModelDomainServicesApplicationCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesApplicationCall) Do(f func() *service1.WatchableService) *MockModelDomainServicesApplicationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesApplicationCall) DoAndReturn(f func() *service1.WatchableService) *MockModelDomainServicesApplicationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BlockCommand mocks base method.
func (m *MockModelDomainServices) BlockCommand() *service2.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockCommand")
	ret0, _ := ret[0].(*service2.Service)
	return ret0
}

// BlockCommand indicates an expected call of BlockCommand.
func (mr *MockModelDomainServicesMockRecorder) BlockCommand() *MockModelDomainServicesBlockCommandCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockCommand", reflect.TypeOf((*MockModelDomainServices)(nil).BlockCommand))
	return &MockModelDomainServicesBlockCommandCall{Call: call}
}

// MockModelDomainServicesBlockCommandCall wrap *gomock.Call
type MockModelDomainServicesBlockCommandCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesBlockCommandCall) Return(arg0 *service2.Service) *MockModelDomainServicesBlockCommandCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesBlockCommandCall) Do(f func() *service2.Service) *MockModelDomainServicesBlockCommandCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesBlockCommandCall) DoAndReturn(f func() *service2.Service) *MockModelDomainServicesBlockCommandCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BlockDevice mocks base method.
func (m *MockModelDomainServices) BlockDevice() *service3.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockDevice")
	ret0, _ := ret[0].(*service3.WatchableService)
	return ret0
}

// BlockDevice indicates an expected call of BlockDevice.
func (mr *MockModelDomainServicesMockRecorder) BlockDevice() *MockModelDomainServicesBlockDeviceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockDevice", reflect.TypeOf((*MockModelDomainServices)(nil).BlockDevice))
	return &MockModelDomainServicesBlockDeviceCall{Call: call}
}

// MockModelDomainServicesBlockDeviceCall wrap *gomock.Call
type MockModelDomainServicesBlockDeviceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesBlockDeviceCall) Return(arg0 *service3.WatchableService) *MockModelDomainServicesBlockDeviceCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesBlockDeviceCall) Do(f func() *service3.WatchableService) *MockModelDomainServicesBlockDeviceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesBlockDeviceCall) DoAndReturn(f func() *service3.WatchableService) *MockModelDomainServicesBlockDeviceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CloudImageMetadata mocks base method.
func (m *MockModelDomainServices) CloudImageMetadata() *service4.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudImageMetadata")
	ret0, _ := ret[0].(*service4.Service)
	return ret0
}

// CloudImageMetadata indicates an expected call of CloudImageMetadata.
func (mr *MockModelDomainServicesMockRecorder) CloudImageMetadata() *MockModelDomainServicesCloudImageMetadataCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudImageMetadata", reflect.TypeOf((*MockModelDomainServices)(nil).CloudImageMetadata))
	return &MockModelDomainServicesCloudImageMetadataCall{Call: call}
}

// MockModelDomainServicesCloudImageMetadataCall wrap *gomock.Call
type MockModelDomainServicesCloudImageMetadataCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesCloudImageMetadataCall) Return(arg0 *service4.Service) *MockModelDomainServicesCloudImageMetadataCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesCloudImageMetadataCall) Do(f func() *service4.Service) *MockModelDomainServicesCloudImageMetadataCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesCloudImageMetadataCall) DoAndReturn(f func() *service4.Service) *MockModelDomainServicesCloudImageMetadataCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Config mocks base method.
func (m *MockModelDomainServices) Config() *service10.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*service10.WatchableService)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockModelDomainServicesMockRecorder) Config() *MockModelDomainServicesConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockModelDomainServices)(nil).Config))
	return &MockModelDomainServicesConfigCall{Call: call}
}

// MockModelDomainServicesConfigCall wrap *gomock.Call
type MockModelDomainServicesConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesConfigCall) Return(arg0 *service10.WatchableService) *MockModelDomainServicesConfigCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesConfigCall) Do(f func() *service10.WatchableService) *MockModelDomainServicesConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesConfigCall) DoAndReturn(f func() *service10.WatchableService) *MockModelDomainServicesConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// KeyManager mocks base method.
func (m *MockModelDomainServices) KeyManager() *service5.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeyManager")
	ret0, _ := ret[0].(*service5.Service)
	return ret0
}

// KeyManager indicates an expected call of KeyManager.
func (mr *MockModelDomainServicesMockRecorder) KeyManager() *MockModelDomainServicesKeyManagerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeyManager", reflect.TypeOf((*MockModelDomainServices)(nil).KeyManager))
	return &MockModelDomainServicesKeyManagerCall{Call: call}
}

// MockModelDomainServicesKeyManagerCall wrap *gomock.Call
type MockModelDomainServicesKeyManagerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesKeyManagerCall) Return(arg0 *service5.Service) *MockModelDomainServicesKeyManagerCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesKeyManagerCall) Do(f func() *service5.Service) *MockModelDomainServicesKeyManagerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesKeyManagerCall) DoAndReturn(f func() *service5.Service) *MockModelDomainServicesKeyManagerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// KeyManagerWithImporter mocks base method.
func (m *MockModelDomainServices) KeyManagerWithImporter() *service5.ImporterService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeyManagerWithImporter")
	ret0, _ := ret[0].(*service5.ImporterService)
	return ret0
}

// KeyManagerWithImporter indicates an expected call of KeyManagerWithImporter.
func (mr *MockModelDomainServicesMockRecorder) KeyManagerWithImporter() *MockModelDomainServicesKeyManagerWithImporterCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeyManagerWithImporter", reflect.TypeOf((*MockModelDomainServices)(nil).KeyManagerWithImporter))
	return &MockModelDomainServicesKeyManagerWithImporterCall{Call: call}
}

// MockModelDomainServicesKeyManagerWithImporterCall wrap *gomock.Call
type MockModelDomainServicesKeyManagerWithImporterCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesKeyManagerWithImporterCall) Return(arg0 *service5.ImporterService) *MockModelDomainServicesKeyManagerWithImporterCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesKeyManagerWithImporterCall) Do(f func() *service5.ImporterService) *MockModelDomainServicesKeyManagerWithImporterCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesKeyManagerWithImporterCall) DoAndReturn(f func() *service5.ImporterService) *MockModelDomainServicesKeyManagerWithImporterCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// KeyUpdater mocks base method.
func (m *MockModelDomainServices) KeyUpdater() *service6.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeyUpdater")
	ret0, _ := ret[0].(*service6.WatchableService)
	return ret0
}

// KeyUpdater indicates an expected call of KeyUpdater.
func (mr *MockModelDomainServicesMockRecorder) KeyUpdater() *MockModelDomainServicesKeyUpdaterCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeyUpdater", reflect.TypeOf((*MockModelDomainServices)(nil).KeyUpdater))
	return &MockModelDomainServicesKeyUpdaterCall{Call: call}
}

// MockModelDomainServicesKeyUpdaterCall wrap *gomock.Call
type MockModelDomainServicesKeyUpdaterCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesKeyUpdaterCall) Return(arg0 *service6.WatchableService) *MockModelDomainServicesKeyUpdaterCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesKeyUpdaterCall) Do(f func() *service6.WatchableService) *MockModelDomainServicesKeyUpdaterCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesKeyUpdaterCall) DoAndReturn(f func() *service6.WatchableService) *MockModelDomainServicesKeyUpdaterCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Machine mocks base method.
func (m *MockModelDomainServices) Machine() *service7.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Machine")
	ret0, _ := ret[0].(*service7.WatchableService)
	return ret0
}

// Machine indicates an expected call of Machine.
func (mr *MockModelDomainServicesMockRecorder) Machine() *MockModelDomainServicesMachineCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Machine", reflect.TypeOf((*MockModelDomainServices)(nil).Machine))
	return &MockModelDomainServicesMachineCall{Call: call}
}

// MockModelDomainServicesMachineCall wrap *gomock.Call
type MockModelDomainServicesMachineCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesMachineCall) Return(arg0 *service7.WatchableService) *MockModelDomainServicesMachineCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesMachineCall) Do(f func() *service7.WatchableService) *MockModelDomainServicesMachineCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesMachineCall) DoAndReturn(f func() *service7.WatchableService) *MockModelDomainServicesMachineCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ModelInfo mocks base method.
func (m *MockModelDomainServices) ModelInfo() *service8.ModelService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelInfo")
	ret0, _ := ret[0].(*service8.ModelService)
	return ret0
}

// ModelInfo indicates an expected call of ModelInfo.
func (mr *MockModelDomainServicesMockRecorder) ModelInfo() *MockModelDomainServicesModelInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelInfo", reflect.TypeOf((*MockModelDomainServices)(nil).ModelInfo))
	return &MockModelDomainServicesModelInfoCall{Call: call}
}

// MockModelDomainServicesModelInfoCall wrap *gomock.Call
type MockModelDomainServicesModelInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesModelInfoCall) Return(arg0 *service8.ModelService) *MockModelDomainServicesModelInfoCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesModelInfoCall) Do(f func() *service8.ModelService) *MockModelDomainServicesModelInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesModelInfoCall) DoAndReturn(f func() *service8.ModelService) *MockModelDomainServicesModelInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ModelMigration mocks base method.
func (m *MockModelDomainServices) ModelMigration() *service11.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelMigration")
	ret0, _ := ret[0].(*service11.Service)
	return ret0
}

// ModelMigration indicates an expected call of ModelMigration.
func (mr *MockModelDomainServicesMockRecorder) ModelMigration() *MockModelDomainServicesModelMigrationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelMigration", reflect.TypeOf((*MockModelDomainServices)(nil).ModelMigration))
	return &MockModelDomainServicesModelMigrationCall{Call: call}
}

// MockModelDomainServicesModelMigrationCall wrap *gomock.Call
type MockModelDomainServicesModelMigrationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesModelMigrationCall) Return(arg0 *service11.Service) *MockModelDomainServicesModelMigrationCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesModelMigrationCall) Do(f func() *service11.Service) *MockModelDomainServicesModelMigrationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesModelMigrationCall) DoAndReturn(f func() *service11.Service) *MockModelDomainServicesModelMigrationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ModelSecretBackend mocks base method.
func (m *MockModelDomainServices) ModelSecretBackend() *service16.ModelSecretBackendService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelSecretBackend")
	ret0, _ := ret[0].(*service16.ModelSecretBackendService)
	return ret0
}

// ModelSecretBackend indicates an expected call of ModelSecretBackend.
func (mr *MockModelDomainServicesMockRecorder) ModelSecretBackend() *MockModelDomainServicesModelSecretBackendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelSecretBackend", reflect.TypeOf((*MockModelDomainServices)(nil).ModelSecretBackend))
	return &MockModelDomainServicesModelSecretBackendCall{Call: call}
}

// MockModelDomainServicesModelSecretBackendCall wrap *gomock.Call
type MockModelDomainServicesModelSecretBackendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesModelSecretBackendCall) Return(arg0 *service16.ModelSecretBackendService) *MockModelDomainServicesModelSecretBackendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesModelSecretBackendCall) Do(f func() *service16.ModelSecretBackendService) *MockModelDomainServicesModelSecretBackendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesModelSecretBackendCall) DoAndReturn(f func() *service16.ModelSecretBackendService) *MockModelDomainServicesModelSecretBackendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Network mocks base method.
func (m *MockModelDomainServices) Network() *service12.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Network")
	ret0, _ := ret[0].(*service12.WatchableService)
	return ret0
}

// Network indicates an expected call of Network.
func (mr *MockModelDomainServicesMockRecorder) Network() *MockModelDomainServicesNetworkCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Network", reflect.TypeOf((*MockModelDomainServices)(nil).Network))
	return &MockModelDomainServicesNetworkCall{Call: call}
}

// MockModelDomainServicesNetworkCall wrap *gomock.Call
type MockModelDomainServicesNetworkCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesNetworkCall) Return(arg0 *service12.WatchableService) *MockModelDomainServicesNetworkCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesNetworkCall) Do(f func() *service12.WatchableService) *MockModelDomainServicesNetworkCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesNetworkCall) DoAndReturn(f func() *service12.WatchableService) *MockModelDomainServicesNetworkCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Port mocks base method.
func (m *MockModelDomainServices) Port() *service13.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Port")
	ret0, _ := ret[0].(*service13.WatchableService)
	return ret0
}

// Port indicates an expected call of Port.
func (mr *MockModelDomainServicesMockRecorder) Port() *MockModelDomainServicesPortCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Port", reflect.TypeOf((*MockModelDomainServices)(nil).Port))
	return &MockModelDomainServicesPortCall{Call: call}
}

// MockModelDomainServicesPortCall wrap *gomock.Call
type MockModelDomainServicesPortCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesPortCall) Return(arg0 *service13.WatchableService) *MockModelDomainServicesPortCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesPortCall) Do(f func() *service13.WatchableService) *MockModelDomainServicesPortCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesPortCall) DoAndReturn(f func() *service13.WatchableService) *MockModelDomainServicesPortCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Proxy mocks base method.
func (m *MockModelDomainServices) Proxy() *service14.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Proxy")
	ret0, _ := ret[0].(*service14.Service)
	return ret0
}

// Proxy indicates an expected call of Proxy.
func (mr *MockModelDomainServicesMockRecorder) Proxy() *MockModelDomainServicesProxyCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Proxy", reflect.TypeOf((*MockModelDomainServices)(nil).Proxy))
	return &MockModelDomainServicesProxyCall{Call: call}
}

// MockModelDomainServicesProxyCall wrap *gomock.Call
type MockModelDomainServicesProxyCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesProxyCall) Return(arg0 *service14.Service) *MockModelDomainServicesProxyCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesProxyCall) Do(f func() *service14.Service) *MockModelDomainServicesProxyCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesProxyCall) DoAndReturn(f func() *service14.Service) *MockModelDomainServicesProxyCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Secret mocks base method.
func (m *MockModelDomainServices) Secret(arg0 service15.SecretServiceParams) *service15.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Secret", arg0)
	ret0, _ := ret[0].(*service15.WatchableService)
	return ret0
}

// Secret indicates an expected call of Secret.
func (mr *MockModelDomainServicesMockRecorder) Secret(arg0 any) *MockModelDomainServicesSecretCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Secret", reflect.TypeOf((*MockModelDomainServices)(nil).Secret), arg0)
	return &MockModelDomainServicesSecretCall{Call: call}
}

// MockModelDomainServicesSecretCall wrap *gomock.Call
type MockModelDomainServicesSecretCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesSecretCall) Return(arg0 *service15.WatchableService) *MockModelDomainServicesSecretCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesSecretCall) Do(f func(service15.SecretServiceParams) *service15.WatchableService) *MockModelDomainServicesSecretCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesSecretCall) DoAndReturn(f func(service15.SecretServiceParams) *service15.WatchableService) *MockModelDomainServicesSecretCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Storage mocks base method.
func (m *MockModelDomainServices) Storage() *service17.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Storage")
	ret0, _ := ret[0].(*service17.Service)
	return ret0
}

// Storage indicates an expected call of Storage.
func (mr *MockModelDomainServicesMockRecorder) Storage() *MockModelDomainServicesStorageCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Storage", reflect.TypeOf((*MockModelDomainServices)(nil).Storage))
	return &MockModelDomainServicesStorageCall{Call: call}
}

// MockModelDomainServicesStorageCall wrap *gomock.Call
type MockModelDomainServicesStorageCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesStorageCall) Return(arg0 *service17.Service) *MockModelDomainServicesStorageCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesStorageCall) Do(f func() *service17.Service) *MockModelDomainServicesStorageCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesStorageCall) DoAndReturn(f func() *service17.Service) *MockModelDomainServicesStorageCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Stub mocks base method.
func (m *MockModelDomainServices) Stub() *stub.StubService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stub")
	ret0, _ := ret[0].(*stub.StubService)
	return ret0
}

// Stub indicates an expected call of Stub.
func (mr *MockModelDomainServicesMockRecorder) Stub() *MockModelDomainServicesStubCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stub", reflect.TypeOf((*MockModelDomainServices)(nil).Stub))
	return &MockModelDomainServicesStubCall{Call: call}
}

// MockModelDomainServicesStubCall wrap *gomock.Call
type MockModelDomainServicesStubCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesStubCall) Return(arg0 *stub.StubService) *MockModelDomainServicesStubCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesStubCall) Do(f func() *stub.StubService) *MockModelDomainServicesStubCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesStubCall) DoAndReturn(f func() *stub.StubService) *MockModelDomainServicesStubCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UnitState mocks base method.
func (m *MockModelDomainServices) UnitState() *service18.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnitState")
	ret0, _ := ret[0].(*service18.Service)
	return ret0
}

// UnitState indicates an expected call of UnitState.
func (mr *MockModelDomainServicesMockRecorder) UnitState() *MockModelDomainServicesUnitStateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnitState", reflect.TypeOf((*MockModelDomainServices)(nil).UnitState))
	return &MockModelDomainServicesUnitStateCall{Call: call}
}

// MockModelDomainServicesUnitStateCall wrap *gomock.Call
type MockModelDomainServicesUnitStateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelDomainServicesUnitStateCall) Return(arg0 *service18.Service) *MockModelDomainServicesUnitStateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelDomainServicesUnitStateCall) Do(f func() *service18.Service) *MockModelDomainServicesUnitStateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelDomainServicesUnitStateCall) DoAndReturn(f func() *service18.Service) *MockModelDomainServicesUnitStateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}