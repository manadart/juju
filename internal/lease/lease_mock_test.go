// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/core/lease (interfaces: Secretary)
//
// Generated by this command:
//
//	mockgen -package lease -destination lease_mock_test.go github.com/juju/juju/core/lease Secretary
//

// Package lease is a generated GoMock package.
package lease

import (
	reflect "reflect"
	time "time"

	lease "github.com/juju/juju/core/lease"
	gomock "go.uber.org/mock/gomock"
)

// MockSecretary is a mock of Secretary interface.
type MockSecretary struct {
	ctrl     *gomock.Controller
	recorder *MockSecretaryMockRecorder
}

// MockSecretaryMockRecorder is the mock recorder for MockSecretary.
type MockSecretaryMockRecorder struct {
	mock *MockSecretary
}

// NewMockSecretary creates a new mock instance.
func NewMockSecretary(ctrl *gomock.Controller) *MockSecretary {
	mock := &MockSecretary{ctrl: ctrl}
	mock.recorder = &MockSecretaryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretary) EXPECT() *MockSecretaryMockRecorder {
	return m.recorder
}

// CheckDuration mocks base method.
func (m *MockSecretary) CheckDuration(arg0 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDuration", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckDuration indicates an expected call of CheckDuration.
func (mr *MockSecretaryMockRecorder) CheckDuration(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDuration", reflect.TypeOf((*MockSecretary)(nil).CheckDuration), arg0)
}

// CheckHolder mocks base method.
func (m *MockSecretary) CheckHolder(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckHolder", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckHolder indicates an expected call of CheckHolder.
func (mr *MockSecretaryMockRecorder) CheckHolder(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckHolder", reflect.TypeOf((*MockSecretary)(nil).CheckHolder), arg0)
}

// CheckLease mocks base method.
func (m *MockSecretary) CheckLease(arg0 lease.Key) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckLease", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckLease indicates an expected call of CheckLease.
func (mr *MockSecretaryMockRecorder) CheckLease(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLease", reflect.TypeOf((*MockSecretary)(nil).CheckLease), arg0)
}