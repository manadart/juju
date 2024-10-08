// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/agent/keyupdater (interfaces: KeyUpdaterService)
//
// Generated by this command:
//
//	mockgen -typed -package keyupdater -destination service_mock_test.go github.com/juju/juju/apiserver/facades/agent/keyupdater KeyUpdaterService
//

// Package keyupdater is a generated GoMock package.
package keyupdater

import (
	context "context"
	reflect "reflect"

	machine "github.com/juju/juju/core/machine"
	watcher "github.com/juju/juju/core/watcher"
	gomock "go.uber.org/mock/gomock"
)

// MockKeyUpdaterService is a mock of KeyUpdaterService interface.
type MockKeyUpdaterService struct {
	ctrl     *gomock.Controller
	recorder *MockKeyUpdaterServiceMockRecorder
}

// MockKeyUpdaterServiceMockRecorder is the mock recorder for MockKeyUpdaterService.
type MockKeyUpdaterServiceMockRecorder struct {
	mock *MockKeyUpdaterService
}

// NewMockKeyUpdaterService creates a new mock instance.
func NewMockKeyUpdaterService(ctrl *gomock.Controller) *MockKeyUpdaterService {
	mock := &MockKeyUpdaterService{ctrl: ctrl}
	mock.recorder = &MockKeyUpdaterServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyUpdaterService) EXPECT() *MockKeyUpdaterServiceMockRecorder {
	return m.recorder
}

// GetAuthorisedKeysForMachine mocks base method.
func (m *MockKeyUpdaterService) GetAuthorisedKeysForMachine(arg0 context.Context, arg1 machine.Name) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthorisedKeysForMachine", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuthorisedKeysForMachine indicates an expected call of GetAuthorisedKeysForMachine.
func (mr *MockKeyUpdaterServiceMockRecorder) GetAuthorisedKeysForMachine(arg0, arg1 any) *MockKeyUpdaterServiceGetAuthorisedKeysForMachineCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthorisedKeysForMachine", reflect.TypeOf((*MockKeyUpdaterService)(nil).GetAuthorisedKeysForMachine), arg0, arg1)
	return &MockKeyUpdaterServiceGetAuthorisedKeysForMachineCall{Call: call}
}

// MockKeyUpdaterServiceGetAuthorisedKeysForMachineCall wrap *gomock.Call
type MockKeyUpdaterServiceGetAuthorisedKeysForMachineCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKeyUpdaterServiceGetAuthorisedKeysForMachineCall) Return(arg0 []string, arg1 error) *MockKeyUpdaterServiceGetAuthorisedKeysForMachineCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKeyUpdaterServiceGetAuthorisedKeysForMachineCall) Do(f func(context.Context, machine.Name) ([]string, error)) *MockKeyUpdaterServiceGetAuthorisedKeysForMachineCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKeyUpdaterServiceGetAuthorisedKeysForMachineCall) DoAndReturn(f func(context.Context, machine.Name) ([]string, error)) *MockKeyUpdaterServiceGetAuthorisedKeysForMachineCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchAuthorisedKeysForMachine mocks base method.
func (m *MockKeyUpdaterService) WatchAuthorisedKeysForMachine(arg0 context.Context, arg1 machine.Name) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchAuthorisedKeysForMachine", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchAuthorisedKeysForMachine indicates an expected call of WatchAuthorisedKeysForMachine.
func (mr *MockKeyUpdaterServiceMockRecorder) WatchAuthorisedKeysForMachine(arg0, arg1 any) *MockKeyUpdaterServiceWatchAuthorisedKeysForMachineCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchAuthorisedKeysForMachine", reflect.TypeOf((*MockKeyUpdaterService)(nil).WatchAuthorisedKeysForMachine), arg0, arg1)
	return &MockKeyUpdaterServiceWatchAuthorisedKeysForMachineCall{Call: call}
}

// MockKeyUpdaterServiceWatchAuthorisedKeysForMachineCall wrap *gomock.Call
type MockKeyUpdaterServiceWatchAuthorisedKeysForMachineCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKeyUpdaterServiceWatchAuthorisedKeysForMachineCall) Return(arg0 watcher.Watcher[struct{}], arg1 error) *MockKeyUpdaterServiceWatchAuthorisedKeysForMachineCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKeyUpdaterServiceWatchAuthorisedKeysForMachineCall) Do(f func(context.Context, machine.Name) (watcher.Watcher[struct{}], error)) *MockKeyUpdaterServiceWatchAuthorisedKeysForMachineCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKeyUpdaterServiceWatchAuthorisedKeysForMachineCall) DoAndReturn(f func(context.Context, machine.Name) (watcher.Watcher[struct{}], error)) *MockKeyUpdaterServiceWatchAuthorisedKeysForMachineCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
