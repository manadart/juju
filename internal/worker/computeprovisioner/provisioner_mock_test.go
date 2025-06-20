// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/worker/computeprovisioner (interfaces: ControllerAPI,MachinesAPI,MachineService)
//
// Generated by this command:
//
//	mockgen -typed -package computeprovisioner_test -destination provisioner_mock_test.go github.com/juju/juju/internal/worker/computeprovisioner ControllerAPI,MachinesAPI,MachineService
//

// Package computeprovisioner_test is a generated GoMock package.
package computeprovisioner_test

import (
	context "context"
	reflect "reflect"

	provisioner "github.com/juju/juju/api/agent/provisioner"
	controller "github.com/juju/juju/controller"
	instance "github.com/juju/juju/core/instance"
	machine "github.com/juju/juju/core/machine"
	watcher "github.com/juju/juju/core/watcher"
	config "github.com/juju/juju/environs/config"
	params "github.com/juju/juju/rpc/params"
	names "github.com/juju/names/v6"
	gomock "go.uber.org/mock/gomock"
)

// MockControllerAPI is a mock of ControllerAPI interface.
type MockControllerAPI struct {
	ctrl     *gomock.Controller
	recorder *MockControllerAPIMockRecorder
}

// MockControllerAPIMockRecorder is the mock recorder for MockControllerAPI.
type MockControllerAPIMockRecorder struct {
	mock *MockControllerAPI
}

// NewMockControllerAPI creates a new mock instance.
func NewMockControllerAPI(ctrl *gomock.Controller) *MockControllerAPI {
	mock := &MockControllerAPI{ctrl: ctrl}
	mock.recorder = &MockControllerAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockControllerAPI) EXPECT() *MockControllerAPIMockRecorder {
	return m.recorder
}

// APIAddresses mocks base method.
func (m *MockControllerAPI) APIAddresses(arg0 context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "APIAddresses", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// APIAddresses indicates an expected call of APIAddresses.
func (mr *MockControllerAPIMockRecorder) APIAddresses(arg0 any) *MockControllerAPIAPIAddressesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APIAddresses", reflect.TypeOf((*MockControllerAPI)(nil).APIAddresses), arg0)
	return &MockControllerAPIAPIAddressesCall{Call: call}
}

// MockControllerAPIAPIAddressesCall wrap *gomock.Call
type MockControllerAPIAPIAddressesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockControllerAPIAPIAddressesCall) Return(arg0 []string, arg1 error) *MockControllerAPIAPIAddressesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockControllerAPIAPIAddressesCall) Do(f func(context.Context) ([]string, error)) *MockControllerAPIAPIAddressesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockControllerAPIAPIAddressesCall) DoAndReturn(f func(context.Context) ([]string, error)) *MockControllerAPIAPIAddressesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CACert mocks base method.
func (m *MockControllerAPI) CACert(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CACert", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CACert indicates an expected call of CACert.
func (mr *MockControllerAPIMockRecorder) CACert(arg0 any) *MockControllerAPICACertCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CACert", reflect.TypeOf((*MockControllerAPI)(nil).CACert), arg0)
	return &MockControllerAPICACertCall{Call: call}
}

// MockControllerAPICACertCall wrap *gomock.Call
type MockControllerAPICACertCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockControllerAPICACertCall) Return(arg0 string, arg1 error) *MockControllerAPICACertCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockControllerAPICACertCall) Do(f func(context.Context) (string, error)) *MockControllerAPICACertCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockControllerAPICACertCall) DoAndReturn(f func(context.Context) (string, error)) *MockControllerAPICACertCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ControllerConfig mocks base method.
func (m *MockControllerAPI) ControllerConfig(arg0 context.Context) (controller.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerConfig", arg0)
	ret0, _ := ret[0].(controller.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerConfig indicates an expected call of ControllerConfig.
func (mr *MockControllerAPIMockRecorder) ControllerConfig(arg0 any) *MockControllerAPIControllerConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerConfig", reflect.TypeOf((*MockControllerAPI)(nil).ControllerConfig), arg0)
	return &MockControllerAPIControllerConfigCall{Call: call}
}

// MockControllerAPIControllerConfigCall wrap *gomock.Call
type MockControllerAPIControllerConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockControllerAPIControllerConfigCall) Return(arg0 controller.Config, arg1 error) *MockControllerAPIControllerConfigCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockControllerAPIControllerConfigCall) Do(f func(context.Context) (controller.Config, error)) *MockControllerAPIControllerConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockControllerAPIControllerConfigCall) DoAndReturn(f func(context.Context) (controller.Config, error)) *MockControllerAPIControllerConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ModelConfig mocks base method.
func (m *MockControllerAPI) ModelConfig(arg0 context.Context) (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelConfig", arg0)
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelConfig indicates an expected call of ModelConfig.
func (mr *MockControllerAPIMockRecorder) ModelConfig(arg0 any) *MockControllerAPIModelConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelConfig", reflect.TypeOf((*MockControllerAPI)(nil).ModelConfig), arg0)
	return &MockControllerAPIModelConfigCall{Call: call}
}

// MockControllerAPIModelConfigCall wrap *gomock.Call
type MockControllerAPIModelConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockControllerAPIModelConfigCall) Return(arg0 *config.Config, arg1 error) *MockControllerAPIModelConfigCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockControllerAPIModelConfigCall) Do(f func(context.Context) (*config.Config, error)) *MockControllerAPIModelConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockControllerAPIModelConfigCall) DoAndReturn(f func(context.Context) (*config.Config, error)) *MockControllerAPIModelConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ModelUUID mocks base method.
func (m *MockControllerAPI) ModelUUID(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelUUID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelUUID indicates an expected call of ModelUUID.
func (mr *MockControllerAPIMockRecorder) ModelUUID(arg0 any) *MockControllerAPIModelUUIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelUUID", reflect.TypeOf((*MockControllerAPI)(nil).ModelUUID), arg0)
	return &MockControllerAPIModelUUIDCall{Call: call}
}

// MockControllerAPIModelUUIDCall wrap *gomock.Call
type MockControllerAPIModelUUIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockControllerAPIModelUUIDCall) Return(arg0 string, arg1 error) *MockControllerAPIModelUUIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockControllerAPIModelUUIDCall) Do(f func(context.Context) (string, error)) *MockControllerAPIModelUUIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockControllerAPIModelUUIDCall) DoAndReturn(f func(context.Context) (string, error)) *MockControllerAPIModelUUIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchForModelConfigChanges mocks base method.
func (m *MockControllerAPI) WatchForModelConfigChanges(arg0 context.Context) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchForModelConfigChanges", arg0)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchForModelConfigChanges indicates an expected call of WatchForModelConfigChanges.
func (mr *MockControllerAPIMockRecorder) WatchForModelConfigChanges(arg0 any) *MockControllerAPIWatchForModelConfigChangesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchForModelConfigChanges", reflect.TypeOf((*MockControllerAPI)(nil).WatchForModelConfigChanges), arg0)
	return &MockControllerAPIWatchForModelConfigChangesCall{Call: call}
}

// MockControllerAPIWatchForModelConfigChangesCall wrap *gomock.Call
type MockControllerAPIWatchForModelConfigChangesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockControllerAPIWatchForModelConfigChangesCall) Return(arg0 watcher.Watcher[struct{}], arg1 error) *MockControllerAPIWatchForModelConfigChangesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockControllerAPIWatchForModelConfigChangesCall) Do(f func(context.Context) (watcher.Watcher[struct{}], error)) *MockControllerAPIWatchForModelConfigChangesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockControllerAPIWatchForModelConfigChangesCall) DoAndReturn(f func(context.Context) (watcher.Watcher[struct{}], error)) *MockControllerAPIWatchForModelConfigChangesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockMachinesAPI is a mock of MachinesAPI interface.
type MockMachinesAPI struct {
	ctrl     *gomock.Controller
	recorder *MockMachinesAPIMockRecorder
}

// MockMachinesAPIMockRecorder is the mock recorder for MockMachinesAPI.
type MockMachinesAPIMockRecorder struct {
	mock *MockMachinesAPI
}

// NewMockMachinesAPI creates a new mock instance.
func NewMockMachinesAPI(ctrl *gomock.Controller) *MockMachinesAPI {
	mock := &MockMachinesAPI{ctrl: ctrl}
	mock.recorder = &MockMachinesAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMachinesAPI) EXPECT() *MockMachinesAPIMockRecorder {
	return m.recorder
}

// Machines mocks base method.
func (m *MockMachinesAPI) Machines(arg0 context.Context, arg1 ...names.MachineTag) ([]provisioner.MachineResult, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Machines", varargs...)
	ret0, _ := ret[0].([]provisioner.MachineResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Machines indicates an expected call of Machines.
func (mr *MockMachinesAPIMockRecorder) Machines(arg0 any, arg1 ...any) *MockMachinesAPIMachinesCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Machines", reflect.TypeOf((*MockMachinesAPI)(nil).Machines), varargs...)
	return &MockMachinesAPIMachinesCall{Call: call}
}

// MockMachinesAPIMachinesCall wrap *gomock.Call
type MockMachinesAPIMachinesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachinesAPIMachinesCall) Return(arg0 []provisioner.MachineResult, arg1 error) *MockMachinesAPIMachinesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachinesAPIMachinesCall) Do(f func(context.Context, ...names.MachineTag) ([]provisioner.MachineResult, error)) *MockMachinesAPIMachinesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachinesAPIMachinesCall) DoAndReturn(f func(context.Context, ...names.MachineTag) ([]provisioner.MachineResult, error)) *MockMachinesAPIMachinesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MachinesWithTransientErrors mocks base method.
func (m *MockMachinesAPI) MachinesWithTransientErrors(arg0 context.Context) ([]provisioner.MachineStatusResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MachinesWithTransientErrors", arg0)
	ret0, _ := ret[0].([]provisioner.MachineStatusResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MachinesWithTransientErrors indicates an expected call of MachinesWithTransientErrors.
func (mr *MockMachinesAPIMockRecorder) MachinesWithTransientErrors(arg0 any) *MockMachinesAPIMachinesWithTransientErrorsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MachinesWithTransientErrors", reflect.TypeOf((*MockMachinesAPI)(nil).MachinesWithTransientErrors), arg0)
	return &MockMachinesAPIMachinesWithTransientErrorsCall{Call: call}
}

// MockMachinesAPIMachinesWithTransientErrorsCall wrap *gomock.Call
type MockMachinesAPIMachinesWithTransientErrorsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachinesAPIMachinesWithTransientErrorsCall) Return(arg0 []provisioner.MachineStatusResult, arg1 error) *MockMachinesAPIMachinesWithTransientErrorsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachinesAPIMachinesWithTransientErrorsCall) Do(f func(context.Context) ([]provisioner.MachineStatusResult, error)) *MockMachinesAPIMachinesWithTransientErrorsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachinesAPIMachinesWithTransientErrorsCall) DoAndReturn(f func(context.Context) ([]provisioner.MachineStatusResult, error)) *MockMachinesAPIMachinesWithTransientErrorsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ProvisioningInfo mocks base method.
func (m *MockMachinesAPI) ProvisioningInfo(arg0 context.Context, arg1 []names.MachineTag) (params.ProvisioningInfoResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProvisioningInfo", arg0, arg1)
	ret0, _ := ret[0].(params.ProvisioningInfoResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProvisioningInfo indicates an expected call of ProvisioningInfo.
func (mr *MockMachinesAPIMockRecorder) ProvisioningInfo(arg0, arg1 any) *MockMachinesAPIProvisioningInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProvisioningInfo", reflect.TypeOf((*MockMachinesAPI)(nil).ProvisioningInfo), arg0, arg1)
	return &MockMachinesAPIProvisioningInfoCall{Call: call}
}

// MockMachinesAPIProvisioningInfoCall wrap *gomock.Call
type MockMachinesAPIProvisioningInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachinesAPIProvisioningInfoCall) Return(arg0 params.ProvisioningInfoResults, arg1 error) *MockMachinesAPIProvisioningInfoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachinesAPIProvisioningInfoCall) Do(f func(context.Context, []names.MachineTag) (params.ProvisioningInfoResults, error)) *MockMachinesAPIProvisioningInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachinesAPIProvisioningInfoCall) DoAndReturn(f func(context.Context, []names.MachineTag) (params.ProvisioningInfoResults, error)) *MockMachinesAPIProvisioningInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchMachineErrorRetry mocks base method.
func (m *MockMachinesAPI) WatchMachineErrorRetry(arg0 context.Context) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchMachineErrorRetry", arg0)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchMachineErrorRetry indicates an expected call of WatchMachineErrorRetry.
func (mr *MockMachinesAPIMockRecorder) WatchMachineErrorRetry(arg0 any) *MockMachinesAPIWatchMachineErrorRetryCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchMachineErrorRetry", reflect.TypeOf((*MockMachinesAPI)(nil).WatchMachineErrorRetry), arg0)
	return &MockMachinesAPIWatchMachineErrorRetryCall{Call: call}
}

// MockMachinesAPIWatchMachineErrorRetryCall wrap *gomock.Call
type MockMachinesAPIWatchMachineErrorRetryCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachinesAPIWatchMachineErrorRetryCall) Return(arg0 watcher.Watcher[struct{}], arg1 error) *MockMachinesAPIWatchMachineErrorRetryCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachinesAPIWatchMachineErrorRetryCall) Do(f func(context.Context) (watcher.Watcher[struct{}], error)) *MockMachinesAPIWatchMachineErrorRetryCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachinesAPIWatchMachineErrorRetryCall) DoAndReturn(f func(context.Context) (watcher.Watcher[struct{}], error)) *MockMachinesAPIWatchMachineErrorRetryCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchModelMachines mocks base method.
func (m *MockMachinesAPI) WatchModelMachines(arg0 context.Context) (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchModelMachines", arg0)
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchModelMachines indicates an expected call of WatchModelMachines.
func (mr *MockMachinesAPIMockRecorder) WatchModelMachines(arg0 any) *MockMachinesAPIWatchModelMachinesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchModelMachines", reflect.TypeOf((*MockMachinesAPI)(nil).WatchModelMachines), arg0)
	return &MockMachinesAPIWatchModelMachinesCall{Call: call}
}

// MockMachinesAPIWatchModelMachinesCall wrap *gomock.Call
type MockMachinesAPIWatchModelMachinesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachinesAPIWatchModelMachinesCall) Return(arg0 watcher.Watcher[[]string], arg1 error) *MockMachinesAPIWatchModelMachinesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachinesAPIWatchModelMachinesCall) Do(f func(context.Context) (watcher.Watcher[[]string], error)) *MockMachinesAPIWatchModelMachinesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachinesAPIWatchModelMachinesCall) DoAndReturn(f func(context.Context) (watcher.Watcher[[]string], error)) *MockMachinesAPIWatchModelMachinesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockMachineService is a mock of MachineService interface.
type MockMachineService struct {
	ctrl     *gomock.Controller
	recorder *MockMachineServiceMockRecorder
}

// MockMachineServiceMockRecorder is the mock recorder for MockMachineService.
type MockMachineServiceMockRecorder struct {
	mock *MockMachineService
}

// NewMockMachineService creates a new mock instance.
func NewMockMachineService(ctrl *gomock.Controller) *MockMachineService {
	mock := &MockMachineService{ctrl: ctrl}
	mock.recorder = &MockMachineServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMachineService) EXPECT() *MockMachineServiceMockRecorder {
	return m.recorder
}

// GetMachineUUID mocks base method.
func (m *MockMachineService) GetMachineUUID(arg0 context.Context, arg1 machine.Name) (machine.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMachineUUID", arg0, arg1)
	ret0, _ := ret[0].(machine.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMachineUUID indicates an expected call of GetMachineUUID.
func (mr *MockMachineServiceMockRecorder) GetMachineUUID(arg0, arg1 any) *MockMachineServiceGetMachineUUIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMachineUUID", reflect.TypeOf((*MockMachineService)(nil).GetMachineUUID), arg0, arg1)
	return &MockMachineServiceGetMachineUUIDCall{Call: call}
}

// MockMachineServiceGetMachineUUIDCall wrap *gomock.Call
type MockMachineServiceGetMachineUUIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineServiceGetMachineUUIDCall) Return(arg0 machine.UUID, arg1 error) *MockMachineServiceGetMachineUUIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineServiceGetMachineUUIDCall) Do(f func(context.Context, machine.Name) (machine.UUID, error)) *MockMachineServiceGetMachineUUIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineServiceGetMachineUUIDCall) DoAndReturn(f func(context.Context, machine.Name) (machine.UUID, error)) *MockMachineServiceGetMachineUUIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetMachineCloudInstance mocks base method.
func (m *MockMachineService) SetMachineCloudInstance(arg0 context.Context, arg1 machine.UUID, arg2 instance.Id, arg3, arg4 string, arg5 *instance.HardwareCharacteristics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMachineCloudInstance", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetMachineCloudInstance indicates an expected call of SetMachineCloudInstance.
func (mr *MockMachineServiceMockRecorder) SetMachineCloudInstance(arg0, arg1, arg2, arg3, arg4, arg5 any) *MockMachineServiceSetMachineCloudInstanceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMachineCloudInstance", reflect.TypeOf((*MockMachineService)(nil).SetMachineCloudInstance), arg0, arg1, arg2, arg3, arg4, arg5)
	return &MockMachineServiceSetMachineCloudInstanceCall{Call: call}
}

// MockMachineServiceSetMachineCloudInstanceCall wrap *gomock.Call
type MockMachineServiceSetMachineCloudInstanceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineServiceSetMachineCloudInstanceCall) Return(arg0 error) *MockMachineServiceSetMachineCloudInstanceCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineServiceSetMachineCloudInstanceCall) Do(f func(context.Context, machine.UUID, instance.Id, string, string, *instance.HardwareCharacteristics) error) *MockMachineServiceSetMachineCloudInstanceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineServiceSetMachineCloudInstanceCall) DoAndReturn(f func(context.Context, machine.UUID, instance.Id, string, string, *instance.HardwareCharacteristics) error) *MockMachineServiceSetMachineCloudInstanceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
