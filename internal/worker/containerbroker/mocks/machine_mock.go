// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/api/agent/provisioner (interfaces: MachineProvisioner)
//
// Generated by this command:
//
//	mockgen -typed -package mocks -destination mocks/machine_mock.go github.com/juju/juju/api/agent/provisioner MachineProvisioner
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	instance "github.com/juju/juju/core/instance"
	life "github.com/juju/juju/core/life"
	semversion "github.com/juju/juju/core/semversion"
	status "github.com/juju/juju/core/status"
	watcher "github.com/juju/juju/core/watcher"
	params "github.com/juju/juju/rpc/params"
	names "github.com/juju/names/v6"
	gomock "go.uber.org/mock/gomock"
)

// MockMachineProvisioner is a mock of MachineProvisioner interface.
type MockMachineProvisioner struct {
	ctrl     *gomock.Controller
	recorder *MockMachineProvisionerMockRecorder
}

// MockMachineProvisionerMockRecorder is the mock recorder for MockMachineProvisioner.
type MockMachineProvisionerMockRecorder struct {
	mock *MockMachineProvisioner
}

// NewMockMachineProvisioner creates a new mock instance.
func NewMockMachineProvisioner(ctrl *gomock.Controller) *MockMachineProvisioner {
	mock := &MockMachineProvisioner{ctrl: ctrl}
	mock.recorder = &MockMachineProvisionerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMachineProvisioner) EXPECT() *MockMachineProvisionerMockRecorder {
	return m.recorder
}

// AvailabilityZone mocks base method.
func (m *MockMachineProvisioner) AvailabilityZone(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AvailabilityZone", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AvailabilityZone indicates an expected call of AvailabilityZone.
func (mr *MockMachineProvisionerMockRecorder) AvailabilityZone(arg0 any) *MockMachineProvisionerAvailabilityZoneCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AvailabilityZone", reflect.TypeOf((*MockMachineProvisioner)(nil).AvailabilityZone), arg0)
	return &MockMachineProvisionerAvailabilityZoneCall{Call: call}
}

// MockMachineProvisionerAvailabilityZoneCall wrap *gomock.Call
type MockMachineProvisionerAvailabilityZoneCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerAvailabilityZoneCall) Return(arg0 string, arg1 error) *MockMachineProvisionerAvailabilityZoneCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerAvailabilityZoneCall) Do(f func(context.Context) (string, error)) *MockMachineProvisionerAvailabilityZoneCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerAvailabilityZoneCall) DoAndReturn(f func(context.Context) (string, error)) *MockMachineProvisionerAvailabilityZoneCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DistributionGroup mocks base method.
func (m *MockMachineProvisioner) DistributionGroup(arg0 context.Context) ([]instance.Id, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DistributionGroup", arg0)
	ret0, _ := ret[0].([]instance.Id)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DistributionGroup indicates an expected call of DistributionGroup.
func (mr *MockMachineProvisionerMockRecorder) DistributionGroup(arg0 any) *MockMachineProvisionerDistributionGroupCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DistributionGroup", reflect.TypeOf((*MockMachineProvisioner)(nil).DistributionGroup), arg0)
	return &MockMachineProvisionerDistributionGroupCall{Call: call}
}

// MockMachineProvisionerDistributionGroupCall wrap *gomock.Call
type MockMachineProvisionerDistributionGroupCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerDistributionGroupCall) Return(arg0 []instance.Id, arg1 error) *MockMachineProvisionerDistributionGroupCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerDistributionGroupCall) Do(f func(context.Context) ([]instance.Id, error)) *MockMachineProvisionerDistributionGroupCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerDistributionGroupCall) DoAndReturn(f func(context.Context) ([]instance.Id, error)) *MockMachineProvisionerDistributionGroupCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// EnsureDead mocks base method.
func (m *MockMachineProvisioner) EnsureDead(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnsureDead", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnsureDead indicates an expected call of EnsureDead.
func (mr *MockMachineProvisionerMockRecorder) EnsureDead(arg0 any) *MockMachineProvisionerEnsureDeadCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnsureDead", reflect.TypeOf((*MockMachineProvisioner)(nil).EnsureDead), arg0)
	return &MockMachineProvisionerEnsureDeadCall{Call: call}
}

// MockMachineProvisionerEnsureDeadCall wrap *gomock.Call
type MockMachineProvisionerEnsureDeadCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerEnsureDeadCall) Return(arg0 error) *MockMachineProvisionerEnsureDeadCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerEnsureDeadCall) Do(f func(context.Context) error) *MockMachineProvisionerEnsureDeadCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerEnsureDeadCall) DoAndReturn(f func(context.Context) error) *MockMachineProvisionerEnsureDeadCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Id mocks base method.
func (m *MockMachineProvisioner) Id() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(string)
	return ret0
}

// Id indicates an expected call of Id.
func (mr *MockMachineProvisionerMockRecorder) Id() *MockMachineProvisionerIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*MockMachineProvisioner)(nil).Id))
	return &MockMachineProvisionerIdCall{Call: call}
}

// MockMachineProvisionerIdCall wrap *gomock.Call
type MockMachineProvisionerIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerIdCall) Return(arg0 string) *MockMachineProvisionerIdCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerIdCall) Do(f func() string) *MockMachineProvisionerIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerIdCall) DoAndReturn(f func() string) *MockMachineProvisionerIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// InstanceId mocks base method.
func (m *MockMachineProvisioner) InstanceId(arg0 context.Context) (instance.Id, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceId", arg0)
	ret0, _ := ret[0].(instance.Id)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstanceId indicates an expected call of InstanceId.
func (mr *MockMachineProvisionerMockRecorder) InstanceId(arg0 any) *MockMachineProvisionerInstanceIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceId", reflect.TypeOf((*MockMachineProvisioner)(nil).InstanceId), arg0)
	return &MockMachineProvisionerInstanceIdCall{Call: call}
}

// MockMachineProvisionerInstanceIdCall wrap *gomock.Call
type MockMachineProvisionerInstanceIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerInstanceIdCall) Return(arg0 instance.Id, arg1 error) *MockMachineProvisionerInstanceIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerInstanceIdCall) Do(f func(context.Context) (instance.Id, error)) *MockMachineProvisionerInstanceIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerInstanceIdCall) DoAndReturn(f func(context.Context) (instance.Id, error)) *MockMachineProvisionerInstanceIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// InstanceStatus mocks base method.
func (m *MockMachineProvisioner) InstanceStatus(arg0 context.Context) (status.Status, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceStatus", arg0)
	ret0, _ := ret[0].(status.Status)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// InstanceStatus indicates an expected call of InstanceStatus.
func (mr *MockMachineProvisionerMockRecorder) InstanceStatus(arg0 any) *MockMachineProvisionerInstanceStatusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceStatus", reflect.TypeOf((*MockMachineProvisioner)(nil).InstanceStatus), arg0)
	return &MockMachineProvisionerInstanceStatusCall{Call: call}
}

// MockMachineProvisionerInstanceStatusCall wrap *gomock.Call
type MockMachineProvisionerInstanceStatusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerInstanceStatusCall) Return(arg0 status.Status, arg1 string, arg2 error) *MockMachineProvisionerInstanceStatusCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerInstanceStatusCall) Do(f func(context.Context) (status.Status, string, error)) *MockMachineProvisionerInstanceStatusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerInstanceStatusCall) DoAndReturn(f func(context.Context) (status.Status, string, error)) *MockMachineProvisionerInstanceStatusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// KeepInstance mocks base method.
func (m *MockMachineProvisioner) KeepInstance(arg0 context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeepInstance", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KeepInstance indicates an expected call of KeepInstance.
func (mr *MockMachineProvisionerMockRecorder) KeepInstance(arg0 any) *MockMachineProvisionerKeepInstanceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeepInstance", reflect.TypeOf((*MockMachineProvisioner)(nil).KeepInstance), arg0)
	return &MockMachineProvisionerKeepInstanceCall{Call: call}
}

// MockMachineProvisionerKeepInstanceCall wrap *gomock.Call
type MockMachineProvisionerKeepInstanceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerKeepInstanceCall) Return(arg0 bool, arg1 error) *MockMachineProvisionerKeepInstanceCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerKeepInstanceCall) Do(f func(context.Context) (bool, error)) *MockMachineProvisionerKeepInstanceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerKeepInstanceCall) DoAndReturn(f func(context.Context) (bool, error)) *MockMachineProvisionerKeepInstanceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Life mocks base method.
func (m *MockMachineProvisioner) Life() life.Value {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Life")
	ret0, _ := ret[0].(life.Value)
	return ret0
}

// Life indicates an expected call of Life.
func (mr *MockMachineProvisionerMockRecorder) Life() *MockMachineProvisionerLifeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Life", reflect.TypeOf((*MockMachineProvisioner)(nil).Life))
	return &MockMachineProvisionerLifeCall{Call: call}
}

// MockMachineProvisionerLifeCall wrap *gomock.Call
type MockMachineProvisionerLifeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerLifeCall) Return(arg0 life.Value) *MockMachineProvisionerLifeCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerLifeCall) Do(f func() life.Value) *MockMachineProvisionerLifeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerLifeCall) DoAndReturn(f func() life.Value) *MockMachineProvisionerLifeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MachineTag mocks base method.
func (m *MockMachineProvisioner) MachineTag() names.MachineTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MachineTag")
	ret0, _ := ret[0].(names.MachineTag)
	return ret0
}

// MachineTag indicates an expected call of MachineTag.
func (mr *MockMachineProvisionerMockRecorder) MachineTag() *MockMachineProvisionerMachineTagCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MachineTag", reflect.TypeOf((*MockMachineProvisioner)(nil).MachineTag))
	return &MockMachineProvisionerMachineTagCall{Call: call}
}

// MockMachineProvisionerMachineTagCall wrap *gomock.Call
type MockMachineProvisionerMachineTagCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerMachineTagCall) Return(arg0 names.MachineTag) *MockMachineProvisionerMachineTagCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerMachineTagCall) Do(f func() names.MachineTag) *MockMachineProvisionerMachineTagCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerMachineTagCall) DoAndReturn(f func() names.MachineTag) *MockMachineProvisionerMachineTagCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MarkForRemoval mocks base method.
func (m *MockMachineProvisioner) MarkForRemoval(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkForRemoval", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkForRemoval indicates an expected call of MarkForRemoval.
func (mr *MockMachineProvisionerMockRecorder) MarkForRemoval(arg0 any) *MockMachineProvisionerMarkForRemovalCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkForRemoval", reflect.TypeOf((*MockMachineProvisioner)(nil).MarkForRemoval), arg0)
	return &MockMachineProvisionerMarkForRemovalCall{Call: call}
}

// MockMachineProvisionerMarkForRemovalCall wrap *gomock.Call
type MockMachineProvisionerMarkForRemovalCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerMarkForRemovalCall) Return(arg0 error) *MockMachineProvisionerMarkForRemovalCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerMarkForRemovalCall) Do(f func(context.Context) error) *MockMachineProvisionerMarkForRemovalCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerMarkForRemovalCall) DoAndReturn(f func(context.Context) error) *MockMachineProvisionerMarkForRemovalCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ModelAgentVersion mocks base method.
func (m *MockMachineProvisioner) ModelAgentVersion(arg0 context.Context) (*semversion.Number, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelAgentVersion", arg0)
	ret0, _ := ret[0].(*semversion.Number)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelAgentVersion indicates an expected call of ModelAgentVersion.
func (mr *MockMachineProvisionerMockRecorder) ModelAgentVersion(arg0 any) *MockMachineProvisionerModelAgentVersionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelAgentVersion", reflect.TypeOf((*MockMachineProvisioner)(nil).ModelAgentVersion), arg0)
	return &MockMachineProvisionerModelAgentVersionCall{Call: call}
}

// MockMachineProvisionerModelAgentVersionCall wrap *gomock.Call
type MockMachineProvisionerModelAgentVersionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerModelAgentVersionCall) Return(arg0 *semversion.Number, arg1 error) *MockMachineProvisionerModelAgentVersionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerModelAgentVersionCall) Do(f func(context.Context) (*semversion.Number, error)) *MockMachineProvisionerModelAgentVersionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerModelAgentVersionCall) DoAndReturn(f func(context.Context) (*semversion.Number, error)) *MockMachineProvisionerModelAgentVersionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Refresh mocks base method.
func (m *MockMachineProvisioner) Refresh(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh.
func (mr *MockMachineProvisionerMockRecorder) Refresh(arg0 any) *MockMachineProvisionerRefreshCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockMachineProvisioner)(nil).Refresh), arg0)
	return &MockMachineProvisionerRefreshCall{Call: call}
}

// MockMachineProvisionerRefreshCall wrap *gomock.Call
type MockMachineProvisionerRefreshCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerRefreshCall) Return(arg0 error) *MockMachineProvisionerRefreshCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerRefreshCall) Do(f func(context.Context) error) *MockMachineProvisionerRefreshCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerRefreshCall) DoAndReturn(f func(context.Context) error) *MockMachineProvisionerRefreshCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Remove mocks base method.
func (m *MockMachineProvisioner) Remove(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockMachineProvisionerMockRecorder) Remove(arg0 any) *MockMachineProvisionerRemoveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockMachineProvisioner)(nil).Remove), arg0)
	return &MockMachineProvisionerRemoveCall{Call: call}
}

// MockMachineProvisionerRemoveCall wrap *gomock.Call
type MockMachineProvisionerRemoveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerRemoveCall) Return(arg0 error) *MockMachineProvisionerRemoveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerRemoveCall) Do(f func(context.Context) error) *MockMachineProvisionerRemoveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerRemoveCall) DoAndReturn(f func(context.Context) error) *MockMachineProvisionerRemoveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetCharmProfiles mocks base method.
func (m *MockMachineProvisioner) SetCharmProfiles(arg0 context.Context, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCharmProfiles", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCharmProfiles indicates an expected call of SetCharmProfiles.
func (mr *MockMachineProvisionerMockRecorder) SetCharmProfiles(arg0, arg1 any) *MockMachineProvisionerSetCharmProfilesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCharmProfiles", reflect.TypeOf((*MockMachineProvisioner)(nil).SetCharmProfiles), arg0, arg1)
	return &MockMachineProvisionerSetCharmProfilesCall{Call: call}
}

// MockMachineProvisionerSetCharmProfilesCall wrap *gomock.Call
type MockMachineProvisionerSetCharmProfilesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerSetCharmProfilesCall) Return(arg0 error) *MockMachineProvisionerSetCharmProfilesCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerSetCharmProfilesCall) Do(f func(context.Context, []string) error) *MockMachineProvisionerSetCharmProfilesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerSetCharmProfilesCall) DoAndReturn(f func(context.Context, []string) error) *MockMachineProvisionerSetCharmProfilesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetInstanceInfo mocks base method.
func (m *MockMachineProvisioner) SetInstanceInfo(arg0 context.Context, arg1 instance.Id, arg2, arg3 string, arg4 *instance.HardwareCharacteristics, arg5 []params.NetworkConfig, arg6 []params.Volume, arg7 map[string]params.VolumeAttachmentInfo, arg8 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetInstanceInfo", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetInstanceInfo indicates an expected call of SetInstanceInfo.
func (mr *MockMachineProvisionerMockRecorder) SetInstanceInfo(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8 any) *MockMachineProvisionerSetInstanceInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInstanceInfo", reflect.TypeOf((*MockMachineProvisioner)(nil).SetInstanceInfo), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
	return &MockMachineProvisionerSetInstanceInfoCall{Call: call}
}

// MockMachineProvisionerSetInstanceInfoCall wrap *gomock.Call
type MockMachineProvisionerSetInstanceInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerSetInstanceInfoCall) Return(arg0 error) *MockMachineProvisionerSetInstanceInfoCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerSetInstanceInfoCall) Do(f func(context.Context, instance.Id, string, string, *instance.HardwareCharacteristics, []params.NetworkConfig, []params.Volume, map[string]params.VolumeAttachmentInfo, []string) error) *MockMachineProvisionerSetInstanceInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerSetInstanceInfoCall) DoAndReturn(f func(context.Context, instance.Id, string, string, *instance.HardwareCharacteristics, []params.NetworkConfig, []params.Volume, map[string]params.VolumeAttachmentInfo, []string) error) *MockMachineProvisionerSetInstanceInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetInstanceStatus mocks base method.
func (m *MockMachineProvisioner) SetInstanceStatus(arg0 context.Context, arg1 status.Status, arg2 string, arg3 map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetInstanceStatus", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetInstanceStatus indicates an expected call of SetInstanceStatus.
func (mr *MockMachineProvisionerMockRecorder) SetInstanceStatus(arg0, arg1, arg2, arg3 any) *MockMachineProvisionerSetInstanceStatusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInstanceStatus", reflect.TypeOf((*MockMachineProvisioner)(nil).SetInstanceStatus), arg0, arg1, arg2, arg3)
	return &MockMachineProvisionerSetInstanceStatusCall{Call: call}
}

// MockMachineProvisionerSetInstanceStatusCall wrap *gomock.Call
type MockMachineProvisionerSetInstanceStatusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerSetInstanceStatusCall) Return(arg0 error) *MockMachineProvisionerSetInstanceStatusCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerSetInstanceStatusCall) Do(f func(context.Context, status.Status, string, map[string]any) error) *MockMachineProvisionerSetInstanceStatusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerSetInstanceStatusCall) DoAndReturn(f func(context.Context, status.Status, string, map[string]any) error) *MockMachineProvisionerSetInstanceStatusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetPassword mocks base method.
func (m *MockMachineProvisioner) SetPassword(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPassword", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPassword indicates an expected call of SetPassword.
func (mr *MockMachineProvisionerMockRecorder) SetPassword(arg0, arg1 any) *MockMachineProvisionerSetPasswordCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPassword", reflect.TypeOf((*MockMachineProvisioner)(nil).SetPassword), arg0, arg1)
	return &MockMachineProvisionerSetPasswordCall{Call: call}
}

// MockMachineProvisionerSetPasswordCall wrap *gomock.Call
type MockMachineProvisionerSetPasswordCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerSetPasswordCall) Return(arg0 error) *MockMachineProvisionerSetPasswordCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerSetPasswordCall) Do(f func(context.Context, string) error) *MockMachineProvisionerSetPasswordCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerSetPasswordCall) DoAndReturn(f func(context.Context, string) error) *MockMachineProvisionerSetPasswordCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetStatus mocks base method.
func (m *MockMachineProvisioner) SetStatus(arg0 context.Context, arg1 status.Status, arg2 string, arg3 map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStatus", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStatus indicates an expected call of SetStatus.
func (mr *MockMachineProvisionerMockRecorder) SetStatus(arg0, arg1, arg2, arg3 any) *MockMachineProvisionerSetStatusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatus", reflect.TypeOf((*MockMachineProvisioner)(nil).SetStatus), arg0, arg1, arg2, arg3)
	return &MockMachineProvisionerSetStatusCall{Call: call}
}

// MockMachineProvisionerSetStatusCall wrap *gomock.Call
type MockMachineProvisionerSetStatusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerSetStatusCall) Return(arg0 error) *MockMachineProvisionerSetStatusCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerSetStatusCall) Do(f func(context.Context, status.Status, string, map[string]any) error) *MockMachineProvisionerSetStatusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerSetStatusCall) DoAndReturn(f func(context.Context, status.Status, string, map[string]any) error) *MockMachineProvisionerSetStatusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetSupportedContainers mocks base method.
func (m *MockMachineProvisioner) SetSupportedContainers(arg0 context.Context, arg1 ...instance.ContainerType) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SetSupportedContainers", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSupportedContainers indicates an expected call of SetSupportedContainers.
func (mr *MockMachineProvisionerMockRecorder) SetSupportedContainers(arg0 any, arg1 ...any) *MockMachineProvisionerSetSupportedContainersCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSupportedContainers", reflect.TypeOf((*MockMachineProvisioner)(nil).SetSupportedContainers), varargs...)
	return &MockMachineProvisionerSetSupportedContainersCall{Call: call}
}

// MockMachineProvisionerSetSupportedContainersCall wrap *gomock.Call
type MockMachineProvisionerSetSupportedContainersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerSetSupportedContainersCall) Return(arg0 error) *MockMachineProvisionerSetSupportedContainersCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerSetSupportedContainersCall) Do(f func(context.Context, ...instance.ContainerType) error) *MockMachineProvisionerSetSupportedContainersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerSetSupportedContainersCall) DoAndReturn(f func(context.Context, ...instance.ContainerType) error) *MockMachineProvisionerSetSupportedContainersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Status mocks base method.
func (m *MockMachineProvisioner) Status(arg0 context.Context) (status.Status, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", arg0)
	ret0, _ := ret[0].(status.Status)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Status indicates an expected call of Status.
func (mr *MockMachineProvisionerMockRecorder) Status(arg0 any) *MockMachineProvisionerStatusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockMachineProvisioner)(nil).Status), arg0)
	return &MockMachineProvisionerStatusCall{Call: call}
}

// MockMachineProvisionerStatusCall wrap *gomock.Call
type MockMachineProvisionerStatusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerStatusCall) Return(arg0 status.Status, arg1 string, arg2 error) *MockMachineProvisionerStatusCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerStatusCall) Do(f func(context.Context) (status.Status, string, error)) *MockMachineProvisionerStatusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerStatusCall) DoAndReturn(f func(context.Context) (status.Status, string, error)) *MockMachineProvisionerStatusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// String mocks base method.
func (m *MockMachineProvisioner) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockMachineProvisionerMockRecorder) String() *MockMachineProvisionerStringCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockMachineProvisioner)(nil).String))
	return &MockMachineProvisionerStringCall{Call: call}
}

// MockMachineProvisionerStringCall wrap *gomock.Call
type MockMachineProvisionerStringCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerStringCall) Return(arg0 string) *MockMachineProvisionerStringCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerStringCall) Do(f func() string) *MockMachineProvisionerStringCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerStringCall) DoAndReturn(f func() string) *MockMachineProvisionerStringCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SupportedContainers mocks base method.
func (m *MockMachineProvisioner) SupportedContainers(arg0 context.Context) ([]instance.ContainerType, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportedContainers", arg0)
	ret0, _ := ret[0].([]instance.ContainerType)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SupportedContainers indicates an expected call of SupportedContainers.
func (mr *MockMachineProvisionerMockRecorder) SupportedContainers(arg0 any) *MockMachineProvisionerSupportedContainersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportedContainers", reflect.TypeOf((*MockMachineProvisioner)(nil).SupportedContainers), arg0)
	return &MockMachineProvisionerSupportedContainersCall{Call: call}
}

// MockMachineProvisionerSupportedContainersCall wrap *gomock.Call
type MockMachineProvisionerSupportedContainersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerSupportedContainersCall) Return(arg0 []instance.ContainerType, arg1 bool, arg2 error) *MockMachineProvisionerSupportedContainersCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerSupportedContainersCall) Do(f func(context.Context) ([]instance.ContainerType, bool, error)) *MockMachineProvisionerSupportedContainersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerSupportedContainersCall) DoAndReturn(f func(context.Context) ([]instance.ContainerType, bool, error)) *MockMachineProvisionerSupportedContainersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SupportsNoContainers mocks base method.
func (m *MockMachineProvisioner) SupportsNoContainers(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportsNoContainers", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SupportsNoContainers indicates an expected call of SupportsNoContainers.
func (mr *MockMachineProvisionerMockRecorder) SupportsNoContainers(arg0 any) *MockMachineProvisionerSupportsNoContainersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportsNoContainers", reflect.TypeOf((*MockMachineProvisioner)(nil).SupportsNoContainers), arg0)
	return &MockMachineProvisionerSupportsNoContainersCall{Call: call}
}

// MockMachineProvisionerSupportsNoContainersCall wrap *gomock.Call
type MockMachineProvisionerSupportsNoContainersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerSupportsNoContainersCall) Return(arg0 error) *MockMachineProvisionerSupportsNoContainersCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerSupportsNoContainersCall) Do(f func(context.Context) error) *MockMachineProvisionerSupportsNoContainersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerSupportsNoContainersCall) DoAndReturn(f func(context.Context) error) *MockMachineProvisionerSupportsNoContainersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Tag mocks base method.
func (m *MockMachineProvisioner) Tag() names.Tag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tag")
	ret0, _ := ret[0].(names.Tag)
	return ret0
}

// Tag indicates an expected call of Tag.
func (mr *MockMachineProvisionerMockRecorder) Tag() *MockMachineProvisionerTagCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockMachineProvisioner)(nil).Tag))
	return &MockMachineProvisionerTagCall{Call: call}
}

// MockMachineProvisionerTagCall wrap *gomock.Call
type MockMachineProvisionerTagCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerTagCall) Return(arg0 names.Tag) *MockMachineProvisionerTagCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerTagCall) Do(f func() names.Tag) *MockMachineProvisionerTagCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerTagCall) DoAndReturn(f func() names.Tag) *MockMachineProvisionerTagCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchContainers mocks base method.
func (m *MockMachineProvisioner) WatchContainers(arg0 context.Context, arg1 instance.ContainerType) (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchContainers", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchContainers indicates an expected call of WatchContainers.
func (mr *MockMachineProvisionerMockRecorder) WatchContainers(arg0, arg1 any) *MockMachineProvisionerWatchContainersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchContainers", reflect.TypeOf((*MockMachineProvisioner)(nil).WatchContainers), arg0, arg1)
	return &MockMachineProvisionerWatchContainersCall{Call: call}
}

// MockMachineProvisionerWatchContainersCall wrap *gomock.Call
type MockMachineProvisionerWatchContainersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineProvisionerWatchContainersCall) Return(arg0 watcher.Watcher[[]string], arg1 error) *MockMachineProvisionerWatchContainersCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineProvisionerWatchContainersCall) Do(f func(context.Context, instance.ContainerType) (watcher.Watcher[[]string], error)) *MockMachineProvisionerWatchContainersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineProvisionerWatchContainersCall) DoAndReturn(f func(context.Context, instance.ContainerType) (watcher.Watcher[[]string], error)) *MockMachineProvisionerWatchContainersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
