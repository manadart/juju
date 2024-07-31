// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/worker/caasapplicationprovisioner (interfaces: CAASProvisionerFacade)
//
// Generated by this command:
//
//	mockgen -typed -package mocks -destination mocks/facade_mock.go github.com/juju/juju/internal/worker/caasapplicationprovisioner CAASProvisionerFacade
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	charms "github.com/juju/juju/api/common/charms"
	caasapplicationprovisioner "github.com/juju/juju/api/controller/caasapplicationprovisioner"
	life "github.com/juju/juju/core/life"
	resources "github.com/juju/juju/core/resources"
	status "github.com/juju/juju/core/status"
	watcher "github.com/juju/juju/core/watcher"
	params "github.com/juju/juju/rpc/params"
	gomock "go.uber.org/mock/gomock"
)

// MockCAASProvisionerFacade is a mock of CAASProvisionerFacade interface.
type MockCAASProvisionerFacade struct {
	ctrl     *gomock.Controller
	recorder *MockCAASProvisionerFacadeMockRecorder
}

// MockCAASProvisionerFacadeMockRecorder is the mock recorder for MockCAASProvisionerFacade.
type MockCAASProvisionerFacadeMockRecorder struct {
	mock *MockCAASProvisionerFacade
}

// NewMockCAASProvisionerFacade creates a new mock instance.
func NewMockCAASProvisionerFacade(ctrl *gomock.Controller) *MockCAASProvisionerFacade {
	mock := &MockCAASProvisionerFacade{ctrl: ctrl}
	mock.recorder = &MockCAASProvisionerFacadeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCAASProvisionerFacade) EXPECT() *MockCAASProvisionerFacadeMockRecorder {
	return m.recorder
}

// ApplicationCharmInfo mocks base method.
func (m *MockCAASProvisionerFacade) ApplicationCharmInfo(arg0 context.Context, arg1 string) (*charms.CharmInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplicationCharmInfo", arg0, arg1)
	ret0, _ := ret[0].(*charms.CharmInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplicationCharmInfo indicates an expected call of ApplicationCharmInfo.
func (mr *MockCAASProvisionerFacadeMockRecorder) ApplicationCharmInfo(arg0, arg1 any) *MockCAASProvisionerFacadeApplicationCharmInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplicationCharmInfo", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).ApplicationCharmInfo), arg0, arg1)
	return &MockCAASProvisionerFacadeApplicationCharmInfoCall{Call: call}
}

// MockCAASProvisionerFacadeApplicationCharmInfoCall wrap *gomock.Call
type MockCAASProvisionerFacadeApplicationCharmInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeApplicationCharmInfoCall) Return(arg0 *charms.CharmInfo, arg1 error) *MockCAASProvisionerFacadeApplicationCharmInfoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeApplicationCharmInfoCall) Do(f func(context.Context, string) (*charms.CharmInfo, error)) *MockCAASProvisionerFacadeApplicationCharmInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeApplicationCharmInfoCall) DoAndReturn(f func(context.Context, string) (*charms.CharmInfo, error)) *MockCAASProvisionerFacadeApplicationCharmInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ApplicationOCIResources mocks base method.
func (m *MockCAASProvisionerFacade) ApplicationOCIResources(arg0 context.Context, arg1 string) (map[string]resources.DockerImageDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplicationOCIResources", arg0, arg1)
	ret0, _ := ret[0].(map[string]resources.DockerImageDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplicationOCIResources indicates an expected call of ApplicationOCIResources.
func (mr *MockCAASProvisionerFacadeMockRecorder) ApplicationOCIResources(arg0, arg1 any) *MockCAASProvisionerFacadeApplicationOCIResourcesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplicationOCIResources", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).ApplicationOCIResources), arg0, arg1)
	return &MockCAASProvisionerFacadeApplicationOCIResourcesCall{Call: call}
}

// MockCAASProvisionerFacadeApplicationOCIResourcesCall wrap *gomock.Call
type MockCAASProvisionerFacadeApplicationOCIResourcesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeApplicationOCIResourcesCall) Return(arg0 map[string]resources.DockerImageDetails, arg1 error) *MockCAASProvisionerFacadeApplicationOCIResourcesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeApplicationOCIResourcesCall) Do(f func(context.Context, string) (map[string]resources.DockerImageDetails, error)) *MockCAASProvisionerFacadeApplicationOCIResourcesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeApplicationOCIResourcesCall) DoAndReturn(f func(context.Context, string) (map[string]resources.DockerImageDetails, error)) *MockCAASProvisionerFacadeApplicationOCIResourcesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CharmInfo mocks base method.
func (m *MockCAASProvisionerFacade) CharmInfo(arg0 context.Context, arg1 string) (*charms.CharmInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CharmInfo", arg0, arg1)
	ret0, _ := ret[0].(*charms.CharmInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CharmInfo indicates an expected call of CharmInfo.
func (mr *MockCAASProvisionerFacadeMockRecorder) CharmInfo(arg0, arg1 any) *MockCAASProvisionerFacadeCharmInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CharmInfo", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).CharmInfo), arg0, arg1)
	return &MockCAASProvisionerFacadeCharmInfoCall{Call: call}
}

// MockCAASProvisionerFacadeCharmInfoCall wrap *gomock.Call
type MockCAASProvisionerFacadeCharmInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeCharmInfoCall) Return(arg0 *charms.CharmInfo, arg1 error) *MockCAASProvisionerFacadeCharmInfoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeCharmInfoCall) Do(f func(context.Context, string) (*charms.CharmInfo, error)) *MockCAASProvisionerFacadeCharmInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeCharmInfoCall) DoAndReturn(f func(context.Context, string) (*charms.CharmInfo, error)) *MockCAASProvisionerFacadeCharmInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ClearApplicationResources mocks base method.
func (m *MockCAASProvisionerFacade) ClearApplicationResources(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearApplicationResources", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearApplicationResources indicates an expected call of ClearApplicationResources.
func (mr *MockCAASProvisionerFacadeMockRecorder) ClearApplicationResources(arg0, arg1 any) *MockCAASProvisionerFacadeClearApplicationResourcesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearApplicationResources", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).ClearApplicationResources), arg0, arg1)
	return &MockCAASProvisionerFacadeClearApplicationResourcesCall{Call: call}
}

// MockCAASProvisionerFacadeClearApplicationResourcesCall wrap *gomock.Call
type MockCAASProvisionerFacadeClearApplicationResourcesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeClearApplicationResourcesCall) Return(arg0 error) *MockCAASProvisionerFacadeClearApplicationResourcesCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeClearApplicationResourcesCall) Do(f func(context.Context, string) error) *MockCAASProvisionerFacadeClearApplicationResourcesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeClearApplicationResourcesCall) DoAndReturn(f func(context.Context, string) error) *MockCAASProvisionerFacadeClearApplicationResourcesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DestroyUnits mocks base method.
func (m *MockCAASProvisionerFacade) DestroyUnits(arg0 context.Context, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyUnits", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DestroyUnits indicates an expected call of DestroyUnits.
func (mr *MockCAASProvisionerFacadeMockRecorder) DestroyUnits(arg0, arg1 any) *MockCAASProvisionerFacadeDestroyUnitsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyUnits", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).DestroyUnits), arg0, arg1)
	return &MockCAASProvisionerFacadeDestroyUnitsCall{Call: call}
}

// MockCAASProvisionerFacadeDestroyUnitsCall wrap *gomock.Call
type MockCAASProvisionerFacadeDestroyUnitsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeDestroyUnitsCall) Return(arg0 error) *MockCAASProvisionerFacadeDestroyUnitsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeDestroyUnitsCall) Do(f func(context.Context, []string) error) *MockCAASProvisionerFacadeDestroyUnitsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeDestroyUnitsCall) DoAndReturn(f func(context.Context, []string) error) *MockCAASProvisionerFacadeDestroyUnitsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Life mocks base method.
func (m *MockCAASProvisionerFacade) Life(arg0 context.Context, arg1 string) (life.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Life", arg0, arg1)
	ret0, _ := ret[0].(life.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Life indicates an expected call of Life.
func (mr *MockCAASProvisionerFacadeMockRecorder) Life(arg0, arg1 any) *MockCAASProvisionerFacadeLifeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Life", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).Life), arg0, arg1)
	return &MockCAASProvisionerFacadeLifeCall{Call: call}
}

// MockCAASProvisionerFacadeLifeCall wrap *gomock.Call
type MockCAASProvisionerFacadeLifeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeLifeCall) Return(arg0 life.Value, arg1 error) *MockCAASProvisionerFacadeLifeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeLifeCall) Do(f func(context.Context, string) (life.Value, error)) *MockCAASProvisionerFacadeLifeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeLifeCall) DoAndReturn(f func(context.Context, string) (life.Value, error)) *MockCAASProvisionerFacadeLifeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ProvisionerConfig mocks base method.
func (m *MockCAASProvisionerFacade) ProvisionerConfig(arg0 context.Context) (params.CAASApplicationProvisionerConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProvisionerConfig", arg0)
	ret0, _ := ret[0].(params.CAASApplicationProvisionerConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProvisionerConfig indicates an expected call of ProvisionerConfig.
func (mr *MockCAASProvisionerFacadeMockRecorder) ProvisionerConfig(arg0 any) *MockCAASProvisionerFacadeProvisionerConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProvisionerConfig", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).ProvisionerConfig), arg0)
	return &MockCAASProvisionerFacadeProvisionerConfigCall{Call: call}
}

// MockCAASProvisionerFacadeProvisionerConfigCall wrap *gomock.Call
type MockCAASProvisionerFacadeProvisionerConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeProvisionerConfigCall) Return(arg0 params.CAASApplicationProvisionerConfig, arg1 error) *MockCAASProvisionerFacadeProvisionerConfigCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeProvisionerConfigCall) Do(f func(context.Context) (params.CAASApplicationProvisionerConfig, error)) *MockCAASProvisionerFacadeProvisionerConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeProvisionerConfigCall) DoAndReturn(f func(context.Context) (params.CAASApplicationProvisionerConfig, error)) *MockCAASProvisionerFacadeProvisionerConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ProvisioningInfo mocks base method.
func (m *MockCAASProvisionerFacade) ProvisioningInfo(arg0 context.Context, arg1 string) (caasapplicationprovisioner.ProvisioningInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProvisioningInfo", arg0, arg1)
	ret0, _ := ret[0].(caasapplicationprovisioner.ProvisioningInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProvisioningInfo indicates an expected call of ProvisioningInfo.
func (mr *MockCAASProvisionerFacadeMockRecorder) ProvisioningInfo(arg0, arg1 any) *MockCAASProvisionerFacadeProvisioningInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProvisioningInfo", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).ProvisioningInfo), arg0, arg1)
	return &MockCAASProvisionerFacadeProvisioningInfoCall{Call: call}
}

// MockCAASProvisionerFacadeProvisioningInfoCall wrap *gomock.Call
type MockCAASProvisionerFacadeProvisioningInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeProvisioningInfoCall) Return(arg0 caasapplicationprovisioner.ProvisioningInfo, arg1 error) *MockCAASProvisionerFacadeProvisioningInfoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeProvisioningInfoCall) Do(f func(context.Context, string) (caasapplicationprovisioner.ProvisioningInfo, error)) *MockCAASProvisionerFacadeProvisioningInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeProvisioningInfoCall) DoAndReturn(f func(context.Context, string) (caasapplicationprovisioner.ProvisioningInfo, error)) *MockCAASProvisionerFacadeProvisioningInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ProvisioningState mocks base method.
func (m *MockCAASProvisionerFacade) ProvisioningState(arg0 context.Context, arg1 string) (*params.CAASApplicationProvisioningState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProvisioningState", arg0, arg1)
	ret0, _ := ret[0].(*params.CAASApplicationProvisioningState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProvisioningState indicates an expected call of ProvisioningState.
func (mr *MockCAASProvisionerFacadeMockRecorder) ProvisioningState(arg0, arg1 any) *MockCAASProvisionerFacadeProvisioningStateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProvisioningState", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).ProvisioningState), arg0, arg1)
	return &MockCAASProvisionerFacadeProvisioningStateCall{Call: call}
}

// MockCAASProvisionerFacadeProvisioningStateCall wrap *gomock.Call
type MockCAASProvisionerFacadeProvisioningStateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeProvisioningStateCall) Return(arg0 *params.CAASApplicationProvisioningState, arg1 error) *MockCAASProvisionerFacadeProvisioningStateCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeProvisioningStateCall) Do(f func(context.Context, string) (*params.CAASApplicationProvisioningState, error)) *MockCAASProvisionerFacadeProvisioningStateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeProvisioningStateCall) DoAndReturn(f func(context.Context, string) (*params.CAASApplicationProvisioningState, error)) *MockCAASProvisionerFacadeProvisioningStateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RemoveUnit mocks base method.
func (m *MockCAASProvisionerFacade) RemoveUnit(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUnit", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUnit indicates an expected call of RemoveUnit.
func (mr *MockCAASProvisionerFacadeMockRecorder) RemoveUnit(arg0, arg1 any) *MockCAASProvisionerFacadeRemoveUnitCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUnit", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).RemoveUnit), arg0, arg1)
	return &MockCAASProvisionerFacadeRemoveUnitCall{Call: call}
}

// MockCAASProvisionerFacadeRemoveUnitCall wrap *gomock.Call
type MockCAASProvisionerFacadeRemoveUnitCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeRemoveUnitCall) Return(arg0 error) *MockCAASProvisionerFacadeRemoveUnitCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeRemoveUnitCall) Do(f func(context.Context, string) error) *MockCAASProvisionerFacadeRemoveUnitCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeRemoveUnitCall) DoAndReturn(f func(context.Context, string) error) *MockCAASProvisionerFacadeRemoveUnitCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetOperatorStatus mocks base method.
func (m *MockCAASProvisionerFacade) SetOperatorStatus(arg0 context.Context, arg1 string, arg2 status.Status, arg3 string, arg4 map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOperatorStatus", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetOperatorStatus indicates an expected call of SetOperatorStatus.
func (mr *MockCAASProvisionerFacadeMockRecorder) SetOperatorStatus(arg0, arg1, arg2, arg3, arg4 any) *MockCAASProvisionerFacadeSetOperatorStatusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOperatorStatus", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).SetOperatorStatus), arg0, arg1, arg2, arg3, arg4)
	return &MockCAASProvisionerFacadeSetOperatorStatusCall{Call: call}
}

// MockCAASProvisionerFacadeSetOperatorStatusCall wrap *gomock.Call
type MockCAASProvisionerFacadeSetOperatorStatusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeSetOperatorStatusCall) Return(arg0 error) *MockCAASProvisionerFacadeSetOperatorStatusCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeSetOperatorStatusCall) Do(f func(context.Context, string, status.Status, string, map[string]any) error) *MockCAASProvisionerFacadeSetOperatorStatusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeSetOperatorStatusCall) DoAndReturn(f func(context.Context, string, status.Status, string, map[string]any) error) *MockCAASProvisionerFacadeSetOperatorStatusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetPassword mocks base method.
func (m *MockCAASProvisionerFacade) SetPassword(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPassword indicates an expected call of SetPassword.
func (mr *MockCAASProvisionerFacadeMockRecorder) SetPassword(arg0, arg1, arg2 any) *MockCAASProvisionerFacadeSetPasswordCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPassword", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).SetPassword), arg0, arg1, arg2)
	return &MockCAASProvisionerFacadeSetPasswordCall{Call: call}
}

// MockCAASProvisionerFacadeSetPasswordCall wrap *gomock.Call
type MockCAASProvisionerFacadeSetPasswordCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeSetPasswordCall) Return(arg0 error) *MockCAASProvisionerFacadeSetPasswordCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeSetPasswordCall) Do(f func(context.Context, string, string) error) *MockCAASProvisionerFacadeSetPasswordCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeSetPasswordCall) DoAndReturn(f func(context.Context, string, string) error) *MockCAASProvisionerFacadeSetPasswordCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetProvisioningState mocks base method.
func (m *MockCAASProvisionerFacade) SetProvisioningState(arg0 context.Context, arg1 string, arg2 params.CAASApplicationProvisioningState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetProvisioningState", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetProvisioningState indicates an expected call of SetProvisioningState.
func (mr *MockCAASProvisionerFacadeMockRecorder) SetProvisioningState(arg0, arg1, arg2 any) *MockCAASProvisionerFacadeSetProvisioningStateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProvisioningState", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).SetProvisioningState), arg0, arg1, arg2)
	return &MockCAASProvisionerFacadeSetProvisioningStateCall{Call: call}
}

// MockCAASProvisionerFacadeSetProvisioningStateCall wrap *gomock.Call
type MockCAASProvisionerFacadeSetProvisioningStateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeSetProvisioningStateCall) Return(arg0 error) *MockCAASProvisionerFacadeSetProvisioningStateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeSetProvisioningStateCall) Do(f func(context.Context, string, params.CAASApplicationProvisioningState) error) *MockCAASProvisionerFacadeSetProvisioningStateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeSetProvisioningStateCall) DoAndReturn(f func(context.Context, string, params.CAASApplicationProvisioningState) error) *MockCAASProvisionerFacadeSetProvisioningStateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Units mocks base method.
func (m *MockCAASProvisionerFacade) Units(arg0 context.Context, arg1 string) ([]params.CAASUnit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Units", arg0, arg1)
	ret0, _ := ret[0].([]params.CAASUnit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Units indicates an expected call of Units.
func (mr *MockCAASProvisionerFacadeMockRecorder) Units(arg0, arg1 any) *MockCAASProvisionerFacadeUnitsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Units", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).Units), arg0, arg1)
	return &MockCAASProvisionerFacadeUnitsCall{Call: call}
}

// MockCAASProvisionerFacadeUnitsCall wrap *gomock.Call
type MockCAASProvisionerFacadeUnitsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeUnitsCall) Return(arg0 []params.CAASUnit, arg1 error) *MockCAASProvisionerFacadeUnitsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeUnitsCall) Do(f func(context.Context, string) ([]params.CAASUnit, error)) *MockCAASProvisionerFacadeUnitsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeUnitsCall) DoAndReturn(f func(context.Context, string) ([]params.CAASUnit, error)) *MockCAASProvisionerFacadeUnitsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateUnits mocks base method.
func (m *MockCAASProvisionerFacade) UpdateUnits(arg0 context.Context, arg1 params.UpdateApplicationUnits) (*params.UpdateApplicationUnitsInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUnits", arg0, arg1)
	ret0, _ := ret[0].(*params.UpdateApplicationUnitsInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUnits indicates an expected call of UpdateUnits.
func (mr *MockCAASProvisionerFacadeMockRecorder) UpdateUnits(arg0, arg1 any) *MockCAASProvisionerFacadeUpdateUnitsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUnits", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).UpdateUnits), arg0, arg1)
	return &MockCAASProvisionerFacadeUpdateUnitsCall{Call: call}
}

// MockCAASProvisionerFacadeUpdateUnitsCall wrap *gomock.Call
type MockCAASProvisionerFacadeUpdateUnitsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeUpdateUnitsCall) Return(arg0 *params.UpdateApplicationUnitsInfo, arg1 error) *MockCAASProvisionerFacadeUpdateUnitsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeUpdateUnitsCall) Do(f func(context.Context, params.UpdateApplicationUnits) (*params.UpdateApplicationUnitsInfo, error)) *MockCAASProvisionerFacadeUpdateUnitsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeUpdateUnitsCall) DoAndReturn(f func(context.Context, params.UpdateApplicationUnits) (*params.UpdateApplicationUnitsInfo, error)) *MockCAASProvisionerFacadeUpdateUnitsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchApplication mocks base method.
func (m *MockCAASProvisionerFacade) WatchApplication(arg0 context.Context, arg1 string) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchApplication", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchApplication indicates an expected call of WatchApplication.
func (mr *MockCAASProvisionerFacadeMockRecorder) WatchApplication(arg0, arg1 any) *MockCAASProvisionerFacadeWatchApplicationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchApplication", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).WatchApplication), arg0, arg1)
	return &MockCAASProvisionerFacadeWatchApplicationCall{Call: call}
}

// MockCAASProvisionerFacadeWatchApplicationCall wrap *gomock.Call
type MockCAASProvisionerFacadeWatchApplicationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeWatchApplicationCall) Return(arg0 watcher.Watcher[struct{}], arg1 error) *MockCAASProvisionerFacadeWatchApplicationCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeWatchApplicationCall) Do(f func(context.Context, string) (watcher.Watcher[struct{}], error)) *MockCAASProvisionerFacadeWatchApplicationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeWatchApplicationCall) DoAndReturn(f func(context.Context, string) (watcher.Watcher[struct{}], error)) *MockCAASProvisionerFacadeWatchApplicationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchApplications mocks base method.
func (m *MockCAASProvisionerFacade) WatchApplications(arg0 context.Context) (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchApplications", arg0)
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchApplications indicates an expected call of WatchApplications.
func (mr *MockCAASProvisionerFacadeMockRecorder) WatchApplications(arg0 any) *MockCAASProvisionerFacadeWatchApplicationsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchApplications", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).WatchApplications), arg0)
	return &MockCAASProvisionerFacadeWatchApplicationsCall{Call: call}
}

// MockCAASProvisionerFacadeWatchApplicationsCall wrap *gomock.Call
type MockCAASProvisionerFacadeWatchApplicationsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeWatchApplicationsCall) Return(arg0 watcher.Watcher[[]string], arg1 error) *MockCAASProvisionerFacadeWatchApplicationsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeWatchApplicationsCall) Do(f func(context.Context) (watcher.Watcher[[]string], error)) *MockCAASProvisionerFacadeWatchApplicationsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeWatchApplicationsCall) DoAndReturn(f func(context.Context) (watcher.Watcher[[]string], error)) *MockCAASProvisionerFacadeWatchApplicationsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchProvisioningInfo mocks base method.
func (m *MockCAASProvisionerFacade) WatchProvisioningInfo(arg0 context.Context, arg1 string) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchProvisioningInfo", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchProvisioningInfo indicates an expected call of WatchProvisioningInfo.
func (mr *MockCAASProvisionerFacadeMockRecorder) WatchProvisioningInfo(arg0, arg1 any) *MockCAASProvisionerFacadeWatchProvisioningInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchProvisioningInfo", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).WatchProvisioningInfo), arg0, arg1)
	return &MockCAASProvisionerFacadeWatchProvisioningInfoCall{Call: call}
}

// MockCAASProvisionerFacadeWatchProvisioningInfoCall wrap *gomock.Call
type MockCAASProvisionerFacadeWatchProvisioningInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeWatchProvisioningInfoCall) Return(arg0 watcher.Watcher[struct{}], arg1 error) *MockCAASProvisionerFacadeWatchProvisioningInfoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeWatchProvisioningInfoCall) Do(f func(context.Context, string) (watcher.Watcher[struct{}], error)) *MockCAASProvisionerFacadeWatchProvisioningInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeWatchProvisioningInfoCall) DoAndReturn(f func(context.Context, string) (watcher.Watcher[struct{}], error)) *MockCAASProvisionerFacadeWatchProvisioningInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchUnits mocks base method.
func (m *MockCAASProvisionerFacade) WatchUnits(arg0 context.Context, arg1 string) (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchUnits", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchUnits indicates an expected call of WatchUnits.
func (mr *MockCAASProvisionerFacadeMockRecorder) WatchUnits(arg0, arg1 any) *MockCAASProvisionerFacadeWatchUnitsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchUnits", reflect.TypeOf((*MockCAASProvisionerFacade)(nil).WatchUnits), arg0, arg1)
	return &MockCAASProvisionerFacadeWatchUnitsCall{Call: call}
}

// MockCAASProvisionerFacadeWatchUnitsCall wrap *gomock.Call
type MockCAASProvisionerFacadeWatchUnitsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASProvisionerFacadeWatchUnitsCall) Return(arg0 watcher.Watcher[[]string], arg1 error) *MockCAASProvisionerFacadeWatchUnitsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASProvisionerFacadeWatchUnitsCall) Do(f func(context.Context, string) (watcher.Watcher[[]string], error)) *MockCAASProvisionerFacadeWatchUnitsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASProvisionerFacadeWatchUnitsCall) DoAndReturn(f func(context.Context, string) (watcher.Watcher[[]string], error)) *MockCAASProvisionerFacadeWatchUnitsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}