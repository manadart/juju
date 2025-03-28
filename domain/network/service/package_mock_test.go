// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/domain/network/service (interfaces: State,Provider)
//
// Generated by this command:
//
//	mockgen -typed -package service -destination package_mock_test.go github.com/juju/juju/domain/network/service State,Provider
//

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	database "github.com/juju/juju/core/database"
	instance "github.com/juju/juju/core/instance"
	network "github.com/juju/juju/core/network"
	environs "github.com/juju/juju/environs"
	envcontext "github.com/juju/juju/environs/envcontext"
	names "github.com/juju/names/v6"
	gomock "go.uber.org/mock/gomock"
)

// MockState is a mock of State interface.
type MockState struct {
	ctrl     *gomock.Controller
	recorder *MockStateMockRecorder
}

// MockStateMockRecorder is the mock recorder for MockState.
type MockStateMockRecorder struct {
	mock *MockState
}

// NewMockState creates a new mock instance.
func NewMockState(ctrl *gomock.Controller) *MockState {
	mock := &MockState{ctrl: ctrl}
	mock.recorder = &MockStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockState) EXPECT() *MockStateMockRecorder {
	return m.recorder
}

// AddSpace mocks base method.
func (m *MockState) AddSpace(arg0 context.Context, arg1, arg2 string, arg3 network.Id, arg4 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSpace", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSpace indicates an expected call of AddSpace.
func (mr *MockStateMockRecorder) AddSpace(arg0, arg1, arg2, arg3, arg4 any) *MockStateAddSpaceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSpace", reflect.TypeOf((*MockState)(nil).AddSpace), arg0, arg1, arg2, arg3, arg4)
	return &MockStateAddSpaceCall{Call: call}
}

// MockStateAddSpaceCall wrap *gomock.Call
type MockStateAddSpaceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateAddSpaceCall) Return(arg0 error) *MockStateAddSpaceCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateAddSpaceCall) Do(f func(context.Context, string, string, network.Id, []string) error) *MockStateAddSpaceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateAddSpaceCall) DoAndReturn(f func(context.Context, string, string, network.Id, []string) error) *MockStateAddSpaceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// AddSubnet mocks base method.
func (m *MockState) AddSubnet(arg0 context.Context, arg1 network.SubnetInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSubnet", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSubnet indicates an expected call of AddSubnet.
func (mr *MockStateMockRecorder) AddSubnet(arg0, arg1 any) *MockStateAddSubnetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSubnet", reflect.TypeOf((*MockState)(nil).AddSubnet), arg0, arg1)
	return &MockStateAddSubnetCall{Call: call}
}

// MockStateAddSubnetCall wrap *gomock.Call
type MockStateAddSubnetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateAddSubnetCall) Return(arg0 error) *MockStateAddSubnetCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateAddSubnetCall) Do(f func(context.Context, network.SubnetInfo) error) *MockStateAddSubnetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateAddSubnetCall) DoAndReturn(f func(context.Context, network.SubnetInfo) error) *MockStateAddSubnetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// AllSubnetsQuery mocks base method.
func (m *MockState) AllSubnetsQuery(arg0 context.Context, arg1 database.TxnRunner) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllSubnetsQuery", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllSubnetsQuery indicates an expected call of AllSubnetsQuery.
func (mr *MockStateMockRecorder) AllSubnetsQuery(arg0, arg1 any) *MockStateAllSubnetsQueryCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllSubnetsQuery", reflect.TypeOf((*MockState)(nil).AllSubnetsQuery), arg0, arg1)
	return &MockStateAllSubnetsQueryCall{Call: call}
}

// MockStateAllSubnetsQueryCall wrap *gomock.Call
type MockStateAllSubnetsQueryCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateAllSubnetsQueryCall) Return(arg0 []string, arg1 error) *MockStateAllSubnetsQueryCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateAllSubnetsQueryCall) Do(f func(context.Context, database.TxnRunner) ([]string, error)) *MockStateAllSubnetsQueryCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateAllSubnetsQueryCall) DoAndReturn(f func(context.Context, database.TxnRunner) ([]string, error)) *MockStateAllSubnetsQueryCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DeleteSpace mocks base method.
func (m *MockState) DeleteSpace(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSpace", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSpace indicates an expected call of DeleteSpace.
func (mr *MockStateMockRecorder) DeleteSpace(arg0, arg1 any) *MockStateDeleteSpaceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSpace", reflect.TypeOf((*MockState)(nil).DeleteSpace), arg0, arg1)
	return &MockStateDeleteSpaceCall{Call: call}
}

// MockStateDeleteSpaceCall wrap *gomock.Call
type MockStateDeleteSpaceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateDeleteSpaceCall) Return(arg0 error) *MockStateDeleteSpaceCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateDeleteSpaceCall) Do(f func(context.Context, string) error) *MockStateDeleteSpaceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateDeleteSpaceCall) DoAndReturn(f func(context.Context, string) error) *MockStateDeleteSpaceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DeleteSubnet mocks base method.
func (m *MockState) DeleteSubnet(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSubnet", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSubnet indicates an expected call of DeleteSubnet.
func (mr *MockStateMockRecorder) DeleteSubnet(arg0, arg1 any) *MockStateDeleteSubnetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubnet", reflect.TypeOf((*MockState)(nil).DeleteSubnet), arg0, arg1)
	return &MockStateDeleteSubnetCall{Call: call}
}

// MockStateDeleteSubnetCall wrap *gomock.Call
type MockStateDeleteSubnetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateDeleteSubnetCall) Return(arg0 error) *MockStateDeleteSubnetCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateDeleteSubnetCall) Do(f func(context.Context, string) error) *MockStateDeleteSubnetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateDeleteSubnetCall) DoAndReturn(f func(context.Context, string) error) *MockStateDeleteSubnetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetAllSpaces mocks base method.
func (m *MockState) GetAllSpaces(arg0 context.Context) (network.SpaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSpaces", arg0)
	ret0, _ := ret[0].(network.SpaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSpaces indicates an expected call of GetAllSpaces.
func (mr *MockStateMockRecorder) GetAllSpaces(arg0 any) *MockStateGetAllSpacesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSpaces", reflect.TypeOf((*MockState)(nil).GetAllSpaces), arg0)
	return &MockStateGetAllSpacesCall{Call: call}
}

// MockStateGetAllSpacesCall wrap *gomock.Call
type MockStateGetAllSpacesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateGetAllSpacesCall) Return(arg0 network.SpaceInfos, arg1 error) *MockStateGetAllSpacesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateGetAllSpacesCall) Do(f func(context.Context) (network.SpaceInfos, error)) *MockStateGetAllSpacesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateGetAllSpacesCall) DoAndReturn(f func(context.Context) (network.SpaceInfos, error)) *MockStateGetAllSpacesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetAllSubnets mocks base method.
func (m *MockState) GetAllSubnets(arg0 context.Context) (network.SubnetInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSubnets", arg0)
	ret0, _ := ret[0].(network.SubnetInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSubnets indicates an expected call of GetAllSubnets.
func (mr *MockStateMockRecorder) GetAllSubnets(arg0 any) *MockStateGetAllSubnetsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSubnets", reflect.TypeOf((*MockState)(nil).GetAllSubnets), arg0)
	return &MockStateGetAllSubnetsCall{Call: call}
}

// MockStateGetAllSubnetsCall wrap *gomock.Call
type MockStateGetAllSubnetsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateGetAllSubnetsCall) Return(arg0 network.SubnetInfos, arg1 error) *MockStateGetAllSubnetsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateGetAllSubnetsCall) Do(f func(context.Context) (network.SubnetInfos, error)) *MockStateGetAllSubnetsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateGetAllSubnetsCall) DoAndReturn(f func(context.Context) (network.SubnetInfos, error)) *MockStateGetAllSubnetsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetSpace mocks base method.
func (m *MockState) GetSpace(arg0 context.Context, arg1 string) (*network.SpaceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpace", arg0, arg1)
	ret0, _ := ret[0].(*network.SpaceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpace indicates an expected call of GetSpace.
func (mr *MockStateMockRecorder) GetSpace(arg0, arg1 any) *MockStateGetSpaceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpace", reflect.TypeOf((*MockState)(nil).GetSpace), arg0, arg1)
	return &MockStateGetSpaceCall{Call: call}
}

// MockStateGetSpaceCall wrap *gomock.Call
type MockStateGetSpaceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateGetSpaceCall) Return(arg0 *network.SpaceInfo, arg1 error) *MockStateGetSpaceCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateGetSpaceCall) Do(f func(context.Context, string) (*network.SpaceInfo, error)) *MockStateGetSpaceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateGetSpaceCall) DoAndReturn(f func(context.Context, string) (*network.SpaceInfo, error)) *MockStateGetSpaceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetSpaceByName mocks base method.
func (m *MockState) GetSpaceByName(arg0 context.Context, arg1 string) (*network.SpaceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpaceByName", arg0, arg1)
	ret0, _ := ret[0].(*network.SpaceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpaceByName indicates an expected call of GetSpaceByName.
func (mr *MockStateMockRecorder) GetSpaceByName(arg0, arg1 any) *MockStateGetSpaceByNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpaceByName", reflect.TypeOf((*MockState)(nil).GetSpaceByName), arg0, arg1)
	return &MockStateGetSpaceByNameCall{Call: call}
}

// MockStateGetSpaceByNameCall wrap *gomock.Call
type MockStateGetSpaceByNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateGetSpaceByNameCall) Return(arg0 *network.SpaceInfo, arg1 error) *MockStateGetSpaceByNameCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateGetSpaceByNameCall) Do(f func(context.Context, string) (*network.SpaceInfo, error)) *MockStateGetSpaceByNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateGetSpaceByNameCall) DoAndReturn(f func(context.Context, string) (*network.SpaceInfo, error)) *MockStateGetSpaceByNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetSubnet mocks base method.
func (m *MockState) GetSubnet(arg0 context.Context, arg1 string) (*network.SubnetInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubnet", arg0, arg1)
	ret0, _ := ret[0].(*network.SubnetInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubnet indicates an expected call of GetSubnet.
func (mr *MockStateMockRecorder) GetSubnet(arg0, arg1 any) *MockStateGetSubnetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubnet", reflect.TypeOf((*MockState)(nil).GetSubnet), arg0, arg1)
	return &MockStateGetSubnetCall{Call: call}
}

// MockStateGetSubnetCall wrap *gomock.Call
type MockStateGetSubnetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateGetSubnetCall) Return(arg0 *network.SubnetInfo, arg1 error) *MockStateGetSubnetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateGetSubnetCall) Do(f func(context.Context, string) (*network.SubnetInfo, error)) *MockStateGetSubnetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateGetSubnetCall) DoAndReturn(f func(context.Context, string) (*network.SubnetInfo, error)) *MockStateGetSubnetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetSubnetsByCIDR mocks base method.
func (m *MockState) GetSubnetsByCIDR(arg0 context.Context, arg1 ...string) (network.SubnetInfos, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSubnetsByCIDR", varargs...)
	ret0, _ := ret[0].(network.SubnetInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubnetsByCIDR indicates an expected call of GetSubnetsByCIDR.
func (mr *MockStateMockRecorder) GetSubnetsByCIDR(arg0 any, arg1 ...any) *MockStateGetSubnetsByCIDRCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubnetsByCIDR", reflect.TypeOf((*MockState)(nil).GetSubnetsByCIDR), varargs...)
	return &MockStateGetSubnetsByCIDRCall{Call: call}
}

// MockStateGetSubnetsByCIDRCall wrap *gomock.Call
type MockStateGetSubnetsByCIDRCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateGetSubnetsByCIDRCall) Return(arg0 network.SubnetInfos, arg1 error) *MockStateGetSubnetsByCIDRCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateGetSubnetsByCIDRCall) Do(f func(context.Context, ...string) (network.SubnetInfos, error)) *MockStateGetSubnetsByCIDRCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateGetSubnetsByCIDRCall) DoAndReturn(f func(context.Context, ...string) (network.SubnetInfos, error)) *MockStateGetSubnetsByCIDRCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// NamespaceForWatchSubnet mocks base method.
func (m *MockState) NamespaceForWatchSubnet() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamespaceForWatchSubnet")
	ret0, _ := ret[0].(string)
	return ret0
}

// NamespaceForWatchSubnet indicates an expected call of NamespaceForWatchSubnet.
func (mr *MockStateMockRecorder) NamespaceForWatchSubnet() *MockStateNamespaceForWatchSubnetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamespaceForWatchSubnet", reflect.TypeOf((*MockState)(nil).NamespaceForWatchSubnet))
	return &MockStateNamespaceForWatchSubnetCall{Call: call}
}

// MockStateNamespaceForWatchSubnetCall wrap *gomock.Call
type MockStateNamespaceForWatchSubnetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateNamespaceForWatchSubnetCall) Return(arg0 string) *MockStateNamespaceForWatchSubnetCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateNamespaceForWatchSubnetCall) Do(f func() string) *MockStateNamespaceForWatchSubnetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateNamespaceForWatchSubnetCall) DoAndReturn(f func() string) *MockStateNamespaceForWatchSubnetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateSpace mocks base method.
func (m *MockState) UpdateSpace(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSpace", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSpace indicates an expected call of UpdateSpace.
func (mr *MockStateMockRecorder) UpdateSpace(arg0, arg1, arg2 any) *MockStateUpdateSpaceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSpace", reflect.TypeOf((*MockState)(nil).UpdateSpace), arg0, arg1, arg2)
	return &MockStateUpdateSpaceCall{Call: call}
}

// MockStateUpdateSpaceCall wrap *gomock.Call
type MockStateUpdateSpaceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateUpdateSpaceCall) Return(arg0 error) *MockStateUpdateSpaceCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateUpdateSpaceCall) Do(f func(context.Context, string, string) error) *MockStateUpdateSpaceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateUpdateSpaceCall) DoAndReturn(f func(context.Context, string, string) error) *MockStateUpdateSpaceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateSubnet mocks base method.
func (m *MockState) UpdateSubnet(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSubnet", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSubnet indicates an expected call of UpdateSubnet.
func (mr *MockStateMockRecorder) UpdateSubnet(arg0, arg1, arg2 any) *MockStateUpdateSubnetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSubnet", reflect.TypeOf((*MockState)(nil).UpdateSubnet), arg0, arg1, arg2)
	return &MockStateUpdateSubnetCall{Call: call}
}

// MockStateUpdateSubnetCall wrap *gomock.Call
type MockStateUpdateSubnetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateUpdateSubnetCall) Return(arg0 error) *MockStateUpdateSubnetCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateUpdateSubnetCall) Do(f func(context.Context, string, string) error) *MockStateUpdateSubnetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateUpdateSubnetCall) DoAndReturn(f func(context.Context, string, string) error) *MockStateUpdateSubnetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpsertSubnets mocks base method.
func (m *MockState) UpsertSubnets(arg0 context.Context, arg1 []network.SubnetInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertSubnets", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertSubnets indicates an expected call of UpsertSubnets.
func (mr *MockStateMockRecorder) UpsertSubnets(arg0, arg1 any) *MockStateUpsertSubnetsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertSubnets", reflect.TypeOf((*MockState)(nil).UpsertSubnets), arg0, arg1)
	return &MockStateUpsertSubnetsCall{Call: call}
}

// MockStateUpsertSubnetsCall wrap *gomock.Call
type MockStateUpsertSubnetsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateUpsertSubnetsCall) Return(arg0 error) *MockStateUpsertSubnetsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateUpsertSubnetsCall) Do(f func(context.Context, []network.SubnetInfo) error) *MockStateUpsertSubnetsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateUpsertSubnetsCall) DoAndReturn(f func(context.Context, []network.SubnetInfo) error) *MockStateUpsertSubnetsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockProvider is a mock of Provider interface.
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance.
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// AllocateContainerAddresses mocks base method.
func (m *MockProvider) AllocateContainerAddresses(arg0 context.Context, arg1 instance.Id, arg2 names.MachineTag, arg3 network.InterfaceInfos) (network.InterfaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllocateContainerAddresses", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(network.InterfaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllocateContainerAddresses indicates an expected call of AllocateContainerAddresses.
func (mr *MockProviderMockRecorder) AllocateContainerAddresses(arg0, arg1, arg2, arg3 any) *MockProviderAllocateContainerAddressesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllocateContainerAddresses", reflect.TypeOf((*MockProvider)(nil).AllocateContainerAddresses), arg0, arg1, arg2, arg3)
	return &MockProviderAllocateContainerAddressesCall{Call: call}
}

// MockProviderAllocateContainerAddressesCall wrap *gomock.Call
type MockProviderAllocateContainerAddressesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProviderAllocateContainerAddressesCall) Return(arg0 network.InterfaceInfos, arg1 error) *MockProviderAllocateContainerAddressesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProviderAllocateContainerAddressesCall) Do(f func(context.Context, instance.Id, names.MachineTag, network.InterfaceInfos) (network.InterfaceInfos, error)) *MockProviderAllocateContainerAddressesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProviderAllocateContainerAddressesCall) DoAndReturn(f func(context.Context, instance.Id, names.MachineTag, network.InterfaceInfos) (network.InterfaceInfos, error)) *MockProviderAllocateContainerAddressesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// NetworkInterfaces mocks base method.
func (m *MockProvider) NetworkInterfaces(arg0 envcontext.ProviderCallContext, arg1 []instance.Id) ([]network.InterfaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetworkInterfaces", arg0, arg1)
	ret0, _ := ret[0].([]network.InterfaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NetworkInterfaces indicates an expected call of NetworkInterfaces.
func (mr *MockProviderMockRecorder) NetworkInterfaces(arg0, arg1 any) *MockProviderNetworkInterfacesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetworkInterfaces", reflect.TypeOf((*MockProvider)(nil).NetworkInterfaces), arg0, arg1)
	return &MockProviderNetworkInterfacesCall{Call: call}
}

// MockProviderNetworkInterfacesCall wrap *gomock.Call
type MockProviderNetworkInterfacesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProviderNetworkInterfacesCall) Return(arg0 []network.InterfaceInfos, arg1 error) *MockProviderNetworkInterfacesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProviderNetworkInterfacesCall) Do(f func(envcontext.ProviderCallContext, []instance.Id) ([]network.InterfaceInfos, error)) *MockProviderNetworkInterfacesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProviderNetworkInterfacesCall) DoAndReturn(f func(envcontext.ProviderCallContext, []instance.Id) ([]network.InterfaceInfos, error)) *MockProviderNetworkInterfacesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ProviderSpaceInfo mocks base method.
func (m *MockProvider) ProviderSpaceInfo(arg0 envcontext.ProviderCallContext, arg1 *network.SpaceInfo) (*environs.ProviderSpaceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProviderSpaceInfo", arg0, arg1)
	ret0, _ := ret[0].(*environs.ProviderSpaceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProviderSpaceInfo indicates an expected call of ProviderSpaceInfo.
func (mr *MockProviderMockRecorder) ProviderSpaceInfo(arg0, arg1 any) *MockProviderProviderSpaceInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProviderSpaceInfo", reflect.TypeOf((*MockProvider)(nil).ProviderSpaceInfo), arg0, arg1)
	return &MockProviderProviderSpaceInfoCall{Call: call}
}

// MockProviderProviderSpaceInfoCall wrap *gomock.Call
type MockProviderProviderSpaceInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProviderProviderSpaceInfoCall) Return(arg0 *environs.ProviderSpaceInfo, arg1 error) *MockProviderProviderSpaceInfoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProviderProviderSpaceInfoCall) Do(f func(envcontext.ProviderCallContext, *network.SpaceInfo) (*environs.ProviderSpaceInfo, error)) *MockProviderProviderSpaceInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProviderProviderSpaceInfoCall) DoAndReturn(f func(envcontext.ProviderCallContext, *network.SpaceInfo) (*environs.ProviderSpaceInfo, error)) *MockProviderProviderSpaceInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ReleaseContainerAddresses mocks base method.
func (m *MockProvider) ReleaseContainerAddresses(arg0 envcontext.ProviderCallContext, arg1 []network.ProviderInterfaceInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleaseContainerAddresses", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReleaseContainerAddresses indicates an expected call of ReleaseContainerAddresses.
func (mr *MockProviderMockRecorder) ReleaseContainerAddresses(arg0, arg1 any) *MockProviderReleaseContainerAddressesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseContainerAddresses", reflect.TypeOf((*MockProvider)(nil).ReleaseContainerAddresses), arg0, arg1)
	return &MockProviderReleaseContainerAddressesCall{Call: call}
}

// MockProviderReleaseContainerAddressesCall wrap *gomock.Call
type MockProviderReleaseContainerAddressesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProviderReleaseContainerAddressesCall) Return(arg0 error) *MockProviderReleaseContainerAddressesCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProviderReleaseContainerAddressesCall) Do(f func(envcontext.ProviderCallContext, []network.ProviderInterfaceInfo) error) *MockProviderReleaseContainerAddressesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProviderReleaseContainerAddressesCall) DoAndReturn(f func(envcontext.ProviderCallContext, []network.ProviderInterfaceInfo) error) *MockProviderReleaseContainerAddressesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Spaces mocks base method.
func (m *MockProvider) Spaces(arg0 envcontext.ProviderCallContext) (network.SpaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Spaces", arg0)
	ret0, _ := ret[0].(network.SpaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Spaces indicates an expected call of Spaces.
func (mr *MockProviderMockRecorder) Spaces(arg0 any) *MockProviderSpacesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Spaces", reflect.TypeOf((*MockProvider)(nil).Spaces), arg0)
	return &MockProviderSpacesCall{Call: call}
}

// MockProviderSpacesCall wrap *gomock.Call
type MockProviderSpacesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProviderSpacesCall) Return(arg0 network.SpaceInfos, arg1 error) *MockProviderSpacesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProviderSpacesCall) Do(f func(envcontext.ProviderCallContext) (network.SpaceInfos, error)) *MockProviderSpacesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProviderSpacesCall) DoAndReturn(f func(envcontext.ProviderCallContext) (network.SpaceInfos, error)) *MockProviderSpacesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Subnets mocks base method.
func (m *MockProvider) Subnets(arg0 envcontext.ProviderCallContext, arg1 instance.Id, arg2 []network.Id) ([]network.SubnetInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subnets", arg0, arg1, arg2)
	ret0, _ := ret[0].([]network.SubnetInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subnets indicates an expected call of Subnets.
func (mr *MockProviderMockRecorder) Subnets(arg0, arg1, arg2 any) *MockProviderSubnetsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subnets", reflect.TypeOf((*MockProvider)(nil).Subnets), arg0, arg1, arg2)
	return &MockProviderSubnetsCall{Call: call}
}

// MockProviderSubnetsCall wrap *gomock.Call
type MockProviderSubnetsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProviderSubnetsCall) Return(arg0 []network.SubnetInfo, arg1 error) *MockProviderSubnetsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProviderSubnetsCall) Do(f func(envcontext.ProviderCallContext, instance.Id, []network.Id) ([]network.SubnetInfo, error)) *MockProviderSubnetsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProviderSubnetsCall) DoAndReturn(f func(envcontext.ProviderCallContext, instance.Id, []network.Id) ([]network.SubnetInfo, error)) *MockProviderSubnetsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SupportsContainerAddresses mocks base method.
func (m *MockProvider) SupportsContainerAddresses(arg0 context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportsContainerAddresses", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SupportsContainerAddresses indicates an expected call of SupportsContainerAddresses.
func (mr *MockProviderMockRecorder) SupportsContainerAddresses(arg0 any) *MockProviderSupportsContainerAddressesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportsContainerAddresses", reflect.TypeOf((*MockProvider)(nil).SupportsContainerAddresses), arg0)
	return &MockProviderSupportsContainerAddressesCall{Call: call}
}

// MockProviderSupportsContainerAddressesCall wrap *gomock.Call
type MockProviderSupportsContainerAddressesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProviderSupportsContainerAddressesCall) Return(arg0 bool, arg1 error) *MockProviderSupportsContainerAddressesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProviderSupportsContainerAddressesCall) Do(f func(context.Context) (bool, error)) *MockProviderSupportsContainerAddressesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProviderSupportsContainerAddressesCall) DoAndReturn(f func(context.Context) (bool, error)) *MockProviderSupportsContainerAddressesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SupportsSpaceDiscovery mocks base method.
func (m *MockProvider) SupportsSpaceDiscovery() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportsSpaceDiscovery")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SupportsSpaceDiscovery indicates an expected call of SupportsSpaceDiscovery.
func (mr *MockProviderMockRecorder) SupportsSpaceDiscovery() *MockProviderSupportsSpaceDiscoveryCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportsSpaceDiscovery", reflect.TypeOf((*MockProvider)(nil).SupportsSpaceDiscovery))
	return &MockProviderSupportsSpaceDiscoveryCall{Call: call}
}

// MockProviderSupportsSpaceDiscoveryCall wrap *gomock.Call
type MockProviderSupportsSpaceDiscoveryCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProviderSupportsSpaceDiscoveryCall) Return(arg0 bool, arg1 error) *MockProviderSupportsSpaceDiscoveryCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProviderSupportsSpaceDiscoveryCall) Do(f func() (bool, error)) *MockProviderSupportsSpaceDiscoveryCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProviderSupportsSpaceDiscoveryCall) DoAndReturn(f func() (bool, error)) *MockProviderSupportsSpaceDiscoveryCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SupportsSpaces mocks base method.
func (m *MockProvider) SupportsSpaces() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportsSpaces")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SupportsSpaces indicates an expected call of SupportsSpaces.
func (mr *MockProviderMockRecorder) SupportsSpaces() *MockProviderSupportsSpacesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportsSpaces", reflect.TypeOf((*MockProvider)(nil).SupportsSpaces))
	return &MockProviderSupportsSpacesCall{Call: call}
}

// MockProviderSupportsSpacesCall wrap *gomock.Call
type MockProviderSupportsSpacesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProviderSupportsSpacesCall) Return(arg0 bool, arg1 error) *MockProviderSupportsSpacesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProviderSupportsSpacesCall) Do(f func() (bool, error)) *MockProviderSupportsSpacesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProviderSupportsSpacesCall) DoAndReturn(f func() (bool, error)) *MockProviderSupportsSpacesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
