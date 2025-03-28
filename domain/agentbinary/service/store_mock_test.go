// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/domain/agentbinary/service (interfaces: State)
//
// Generated by this command:
//
//	mockgen -typed -package service -destination store_mock_test.go github.com/juju/juju/domain/agentbinary/service State
//

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	agentbinary "github.com/juju/juju/domain/agentbinary"
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

// Add mocks base method.
func (m *MockState) Add(arg0 context.Context, arg1 agentbinary.Metadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockStateMockRecorder) Add(arg0, arg1 any) *MockStateAddCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockState)(nil).Add), arg0, arg1)
	return &MockStateAddCall{Call: call}
}

// MockStateAddCall wrap *gomock.Call
type MockStateAddCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateAddCall) Return(arg0 error) *MockStateAddCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateAddCall) Do(f func(context.Context, agentbinary.Metadata) error) *MockStateAddCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateAddCall) DoAndReturn(f func(context.Context, agentbinary.Metadata) error) *MockStateAddCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
