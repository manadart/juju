// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/worker/instancemutater (interfaces: MutaterContext)
//
// Generated by this command:
//
//	mockgen -typed -package mocks -destination mocks/mutatercontext_mock.go github.com/juju/juju/internal/worker/instancemutater MutaterContext
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	instancemutater "github.com/juju/juju/api/agent/instancemutater"
	environs "github.com/juju/juju/environs"
	instancemutater0 "github.com/juju/juju/internal/worker/instancemutater"
	names "github.com/juju/names/v6"
	worker "github.com/juju/worker/v4"
	gomock "go.uber.org/mock/gomock"
)

// MockMutaterContext is a mock of MutaterContext interface.
type MockMutaterContext struct {
	ctrl     *gomock.Controller
	recorder *MockMutaterContextMockRecorder
}

// MockMutaterContextMockRecorder is the mock recorder for MockMutaterContext.
type MockMutaterContextMockRecorder struct {
	mock *MockMutaterContext
}

// NewMockMutaterContext creates a new mock instance.
func NewMockMutaterContext(ctrl *gomock.Controller) *MockMutaterContext {
	mock := &MockMutaterContext{ctrl: ctrl}
	mock.recorder = &MockMutaterContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMutaterContext) EXPECT() *MockMutaterContextMockRecorder {
	return m.recorder
}

// KillWithError mocks base method.
func (m *MockMutaterContext) KillWithError(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "KillWithError", arg0)
}

// KillWithError indicates an expected call of KillWithError.
func (mr *MockMutaterContextMockRecorder) KillWithError(arg0 any) *MockMutaterContextKillWithErrorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KillWithError", reflect.TypeOf((*MockMutaterContext)(nil).KillWithError), arg0)
	return &MockMutaterContextKillWithErrorCall{Call: call}
}

// MockMutaterContextKillWithErrorCall wrap *gomock.Call
type MockMutaterContextKillWithErrorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMutaterContextKillWithErrorCall) Return() *MockMutaterContextKillWithErrorCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMutaterContextKillWithErrorCall) Do(f func(error)) *MockMutaterContextKillWithErrorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMutaterContextKillWithErrorCall) DoAndReturn(f func(error)) *MockMutaterContextKillWithErrorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// add mocks base method.
func (m *MockMutaterContext) add(arg0 worker.Worker) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// add indicates an expected call of add.
func (mr *MockMutaterContextMockRecorder) add(arg0 any) *MockMutaterContextaddCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "add", reflect.TypeOf((*MockMutaterContext)(nil).add), arg0)
	return &MockMutaterContextaddCall{Call: call}
}

// MockMutaterContextaddCall wrap *gomock.Call
type MockMutaterContextaddCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMutaterContextaddCall) Return(arg0 error) *MockMutaterContextaddCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMutaterContextaddCall) Do(f func(worker.Worker) error) *MockMutaterContextaddCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMutaterContextaddCall) DoAndReturn(f func(worker.Worker) error) *MockMutaterContextaddCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// dying mocks base method.
func (m *MockMutaterContext) dying() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "dying")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// dying indicates an expected call of dying.
func (mr *MockMutaterContextMockRecorder) dying() *MockMutaterContextdyingCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "dying", reflect.TypeOf((*MockMutaterContext)(nil).dying))
	return &MockMutaterContextdyingCall{Call: call}
}

// MockMutaterContextdyingCall wrap *gomock.Call
type MockMutaterContextdyingCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMutaterContextdyingCall) Return(arg0 <-chan struct{}) *MockMutaterContextdyingCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMutaterContextdyingCall) Do(f func() <-chan struct{}) *MockMutaterContextdyingCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMutaterContextdyingCall) DoAndReturn(f func() <-chan struct{}) *MockMutaterContextdyingCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// errDying mocks base method.
func (m *MockMutaterContext) errDying() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "errDying")
	ret0, _ := ret[0].(error)
	return ret0
}

// errDying indicates an expected call of errDying.
func (mr *MockMutaterContextMockRecorder) errDying() *MockMutaterContexterrDyingCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "errDying", reflect.TypeOf((*MockMutaterContext)(nil).errDying))
	return &MockMutaterContexterrDyingCall{Call: call}
}

// MockMutaterContexterrDyingCall wrap *gomock.Call
type MockMutaterContexterrDyingCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMutaterContexterrDyingCall) Return(arg0 error) *MockMutaterContexterrDyingCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMutaterContexterrDyingCall) Do(f func() error) *MockMutaterContexterrDyingCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMutaterContexterrDyingCall) DoAndReturn(f func() error) *MockMutaterContexterrDyingCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// getBroker mocks base method.
func (m *MockMutaterContext) getBroker() environs.LXDProfiler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getBroker")
	ret0, _ := ret[0].(environs.LXDProfiler)
	return ret0
}

// getBroker indicates an expected call of getBroker.
func (mr *MockMutaterContextMockRecorder) getBroker() *MockMutaterContextgetBrokerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getBroker", reflect.TypeOf((*MockMutaterContext)(nil).getBroker))
	return &MockMutaterContextgetBrokerCall{Call: call}
}

// MockMutaterContextgetBrokerCall wrap *gomock.Call
type MockMutaterContextgetBrokerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMutaterContextgetBrokerCall) Return(arg0 environs.LXDProfiler) *MockMutaterContextgetBrokerCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMutaterContextgetBrokerCall) Do(f func() environs.LXDProfiler) *MockMutaterContextgetBrokerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMutaterContextgetBrokerCall) DoAndReturn(f func() environs.LXDProfiler) *MockMutaterContextgetBrokerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// getMachine mocks base method.
func (m *MockMutaterContext) getMachine(arg0 context.Context, arg1 names.MachineTag) (instancemutater.MutaterMachine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getMachine", arg0, arg1)
	ret0, _ := ret[0].(instancemutater.MutaterMachine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getMachine indicates an expected call of getMachine.
func (mr *MockMutaterContextMockRecorder) getMachine(arg0, arg1 any) *MockMutaterContextgetMachineCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getMachine", reflect.TypeOf((*MockMutaterContext)(nil).getMachine), arg0, arg1)
	return &MockMutaterContextgetMachineCall{Call: call}
}

// MockMutaterContextgetMachineCall wrap *gomock.Call
type MockMutaterContextgetMachineCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMutaterContextgetMachineCall) Return(arg0 instancemutater.MutaterMachine, arg1 error) *MockMutaterContextgetMachineCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMutaterContextgetMachineCall) Do(f func(context.Context, names.MachineTag) (instancemutater.MutaterMachine, error)) *MockMutaterContextgetMachineCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMutaterContextgetMachineCall) DoAndReturn(f func(context.Context, names.MachineTag) (instancemutater.MutaterMachine, error)) *MockMutaterContextgetMachineCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// getRequiredLXDProfiles mocks base method.
func (m *MockMutaterContext) getRequiredLXDProfiles(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getRequiredLXDProfiles", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// getRequiredLXDProfiles indicates an expected call of getRequiredLXDProfiles.
func (mr *MockMutaterContextMockRecorder) getRequiredLXDProfiles(arg0 any) *MockMutaterContextgetRequiredLXDProfilesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getRequiredLXDProfiles", reflect.TypeOf((*MockMutaterContext)(nil).getRequiredLXDProfiles), arg0)
	return &MockMutaterContextgetRequiredLXDProfilesCall{Call: call}
}

// MockMutaterContextgetRequiredLXDProfilesCall wrap *gomock.Call
type MockMutaterContextgetRequiredLXDProfilesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMutaterContextgetRequiredLXDProfilesCall) Return(arg0 []string) *MockMutaterContextgetRequiredLXDProfilesCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMutaterContextgetRequiredLXDProfilesCall) Do(f func(string) []string) *MockMutaterContextgetRequiredLXDProfilesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMutaterContextgetRequiredLXDProfilesCall) DoAndReturn(f func(string) []string) *MockMutaterContextgetRequiredLXDProfilesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// newMachineContext mocks base method.
func (m *MockMutaterContext) newMachineContext() instancemutater0.MachineContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "newMachineContext")
	ret0, _ := ret[0].(instancemutater0.MachineContext)
	return ret0
}

// newMachineContext indicates an expected call of newMachineContext.
func (mr *MockMutaterContextMockRecorder) newMachineContext() *MockMutaterContextnewMachineContextCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "newMachineContext", reflect.TypeOf((*MockMutaterContext)(nil).newMachineContext))
	return &MockMutaterContextnewMachineContextCall{Call: call}
}

// MockMutaterContextnewMachineContextCall wrap *gomock.Call
type MockMutaterContextnewMachineContextCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMutaterContextnewMachineContextCall) Return(arg0 instancemutater0.MachineContext) *MockMutaterContextnewMachineContextCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMutaterContextnewMachineContextCall) Do(f func() instancemutater0.MachineContext) *MockMutaterContextnewMachineContextCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMutaterContextnewMachineContextCall) DoAndReturn(f func() instancemutater0.MachineContext) *MockMutaterContextnewMachineContextCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
