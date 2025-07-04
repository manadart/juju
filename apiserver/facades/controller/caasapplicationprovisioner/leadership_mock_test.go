// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/core/leadership (interfaces: Revoker)
//
// Generated by this command:
//
//	mockgen -typed -package caasapplicationprovisioner_test -destination leadership_mock_test.go github.com/juju/juju/core/leadership Revoker
//

// Package caasapplicationprovisioner_test is a generated GoMock package.
package caasapplicationprovisioner_test

import (
	reflect "reflect"

	unit "github.com/juju/juju/core/unit"
	gomock "go.uber.org/mock/gomock"
)

// MockRevoker is a mock of Revoker interface.
type MockRevoker struct {
	ctrl     *gomock.Controller
	recorder *MockRevokerMockRecorder
}

// MockRevokerMockRecorder is the mock recorder for MockRevoker.
type MockRevokerMockRecorder struct {
	mock *MockRevoker
}

// NewMockRevoker creates a new mock instance.
func NewMockRevoker(ctrl *gomock.Controller) *MockRevoker {
	mock := &MockRevoker{ctrl: ctrl}
	mock.recorder = &MockRevokerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRevoker) EXPECT() *MockRevokerMockRecorder {
	return m.recorder
}

// RevokeLeadership mocks base method.
func (m *MockRevoker) RevokeLeadership(arg0 string, arg1 unit.Name) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevokeLeadership", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RevokeLeadership indicates an expected call of RevokeLeadership.
func (mr *MockRevokerMockRecorder) RevokeLeadership(arg0, arg1 any) *MockRevokerRevokeLeadershipCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeLeadership", reflect.TypeOf((*MockRevoker)(nil).RevokeLeadership), arg0, arg1)
	return &MockRevokerRevokeLeadershipCall{Call: call}
}

// MockRevokerRevokeLeadershipCall wrap *gomock.Call
type MockRevokerRevokeLeadershipCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRevokerRevokeLeadershipCall) Return(arg0 error) *MockRevokerRevokeLeadershipCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRevokerRevokeLeadershipCall) Do(f func(string, unit.Name) error) *MockRevokerRevokeLeadershipCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRevokerRevokeLeadershipCall) DoAndReturn(f func(string, unit.Name) error) *MockRevokerRevokeLeadershipCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
