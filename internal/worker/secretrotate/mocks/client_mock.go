// Code generated by MockGen. DO NOT EDIT.
// Source: secretrotate.go
//
// Generated by this command:
//
//	mockgen -typed -package mocks -destination mocks/client_mock.go -source secretrotate.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	watcher "github.com/juju/juju/core/watcher"
	names "github.com/juju/names/v6"
	gomock "go.uber.org/mock/gomock"
)

// MockSecretManagerFacade is a mock of SecretManagerFacade interface.
type MockSecretManagerFacade struct {
	ctrl     *gomock.Controller
	recorder *MockSecretManagerFacadeMockRecorder
}

// MockSecretManagerFacadeMockRecorder is the mock recorder for MockSecretManagerFacade.
type MockSecretManagerFacadeMockRecorder struct {
	mock *MockSecretManagerFacade
}

// NewMockSecretManagerFacade creates a new mock instance.
func NewMockSecretManagerFacade(ctrl *gomock.Controller) *MockSecretManagerFacade {
	mock := &MockSecretManagerFacade{ctrl: ctrl}
	mock.recorder = &MockSecretManagerFacadeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretManagerFacade) EXPECT() *MockSecretManagerFacadeMockRecorder {
	return m.recorder
}

// WatchSecretsRotationChanges mocks base method.
func (m *MockSecretManagerFacade) WatchSecretsRotationChanges(ctx context.Context, ownerTags ...names.Tag) (watcher.SecretTriggerWatcher, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range ownerTags {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WatchSecretsRotationChanges", varargs...)
	ret0, _ := ret[0].(watcher.SecretTriggerWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchSecretsRotationChanges indicates an expected call of WatchSecretsRotationChanges.
func (mr *MockSecretManagerFacadeMockRecorder) WatchSecretsRotationChanges(ctx any, ownerTags ...any) *MockSecretManagerFacadeWatchSecretsRotationChangesCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, ownerTags...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchSecretsRotationChanges", reflect.TypeOf((*MockSecretManagerFacade)(nil).WatchSecretsRotationChanges), varargs...)
	return &MockSecretManagerFacadeWatchSecretsRotationChangesCall{Call: call}
}

// MockSecretManagerFacadeWatchSecretsRotationChangesCall wrap *gomock.Call
type MockSecretManagerFacadeWatchSecretsRotationChangesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSecretManagerFacadeWatchSecretsRotationChangesCall) Return(arg0 watcher.SecretTriggerWatcher, arg1 error) *MockSecretManagerFacadeWatchSecretsRotationChangesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSecretManagerFacadeWatchSecretsRotationChangesCall) Do(f func(context.Context, ...names.Tag) (watcher.SecretTriggerWatcher, error)) *MockSecretManagerFacadeWatchSecretsRotationChangesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSecretManagerFacadeWatchSecretsRotationChangesCall) DoAndReturn(f func(context.Context, ...names.Tag) (watcher.SecretTriggerWatcher, error)) *MockSecretManagerFacadeWatchSecretsRotationChangesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
