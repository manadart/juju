// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/domain/macaroon/service (interfaces: State)
//
// Generated by this command:
//
//	mockgen -typed -package service -destination package_mock_test.go github.com/juju/juju/domain/macaroon/service State
//

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"
	time "time"

	bakery "github.com/go-macaroon-bakery/macaroon-bakery/v3/bakery"
	macaroon "github.com/juju/juju/domain/macaroon"
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

// FindLatestKey mocks base method.
func (m *MockState) FindLatestKey(arg0 context.Context, arg1, arg2, arg3, arg4 time.Time) (macaroon.RootKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindLatestKey", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(macaroon.RootKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindLatestKey indicates an expected call of FindLatestKey.
func (mr *MockStateMockRecorder) FindLatestKey(arg0, arg1, arg2, arg3, arg4 any) *MockStateFindLatestKeyCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindLatestKey", reflect.TypeOf((*MockState)(nil).FindLatestKey), arg0, arg1, arg2, arg3, arg4)
	return &MockStateFindLatestKeyCall{Call: call}
}

// MockStateFindLatestKeyCall wrap *gomock.Call
type MockStateFindLatestKeyCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateFindLatestKeyCall) Return(arg0 macaroon.RootKey, arg1 error) *MockStateFindLatestKeyCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateFindLatestKeyCall) Do(f func(context.Context, time.Time, time.Time, time.Time, time.Time) (macaroon.RootKey, error)) *MockStateFindLatestKeyCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateFindLatestKeyCall) DoAndReturn(f func(context.Context, time.Time, time.Time, time.Time, time.Time) (macaroon.RootKey, error)) *MockStateFindLatestKeyCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetExternalUsersThirdPartyKey mocks base method.
func (m *MockState) GetExternalUsersThirdPartyKey(arg0 context.Context) (*bakery.KeyPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExternalUsersThirdPartyKey", arg0)
	ret0, _ := ret[0].(*bakery.KeyPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExternalUsersThirdPartyKey indicates an expected call of GetExternalUsersThirdPartyKey.
func (mr *MockStateMockRecorder) GetExternalUsersThirdPartyKey(arg0 any) *MockStateGetExternalUsersThirdPartyKeyCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExternalUsersThirdPartyKey", reflect.TypeOf((*MockState)(nil).GetExternalUsersThirdPartyKey), arg0)
	return &MockStateGetExternalUsersThirdPartyKeyCall{Call: call}
}

// MockStateGetExternalUsersThirdPartyKeyCall wrap *gomock.Call
type MockStateGetExternalUsersThirdPartyKeyCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateGetExternalUsersThirdPartyKeyCall) Return(arg0 *bakery.KeyPair, arg1 error) *MockStateGetExternalUsersThirdPartyKeyCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateGetExternalUsersThirdPartyKeyCall) Do(f func(context.Context) (*bakery.KeyPair, error)) *MockStateGetExternalUsersThirdPartyKeyCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateGetExternalUsersThirdPartyKeyCall) DoAndReturn(f func(context.Context) (*bakery.KeyPair, error)) *MockStateGetExternalUsersThirdPartyKeyCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetKey mocks base method.
func (m *MockState) GetKey(arg0 context.Context, arg1 []byte, arg2 time.Time) (macaroon.RootKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKey", arg0, arg1, arg2)
	ret0, _ := ret[0].(macaroon.RootKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKey indicates an expected call of GetKey.
func (mr *MockStateMockRecorder) GetKey(arg0, arg1, arg2 any) *MockStateGetKeyCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKey", reflect.TypeOf((*MockState)(nil).GetKey), arg0, arg1, arg2)
	return &MockStateGetKeyCall{Call: call}
}

// MockStateGetKeyCall wrap *gomock.Call
type MockStateGetKeyCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateGetKeyCall) Return(arg0 macaroon.RootKey, arg1 error) *MockStateGetKeyCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateGetKeyCall) Do(f func(context.Context, []byte, time.Time) (macaroon.RootKey, error)) *MockStateGetKeyCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateGetKeyCall) DoAndReturn(f func(context.Context, []byte, time.Time) (macaroon.RootKey, error)) *MockStateGetKeyCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetLocalUsersKey mocks base method.
func (m *MockState) GetLocalUsersKey(arg0 context.Context) (*bakery.KeyPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocalUsersKey", arg0)
	ret0, _ := ret[0].(*bakery.KeyPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLocalUsersKey indicates an expected call of GetLocalUsersKey.
func (mr *MockStateMockRecorder) GetLocalUsersKey(arg0 any) *MockStateGetLocalUsersKeyCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocalUsersKey", reflect.TypeOf((*MockState)(nil).GetLocalUsersKey), arg0)
	return &MockStateGetLocalUsersKeyCall{Call: call}
}

// MockStateGetLocalUsersKeyCall wrap *gomock.Call
type MockStateGetLocalUsersKeyCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateGetLocalUsersKeyCall) Return(arg0 *bakery.KeyPair, arg1 error) *MockStateGetLocalUsersKeyCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateGetLocalUsersKeyCall) Do(f func(context.Context) (*bakery.KeyPair, error)) *MockStateGetLocalUsersKeyCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateGetLocalUsersKeyCall) DoAndReturn(f func(context.Context) (*bakery.KeyPair, error)) *MockStateGetLocalUsersKeyCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetLocalUsersThirdPartyKey mocks base method.
func (m *MockState) GetLocalUsersThirdPartyKey(arg0 context.Context) (*bakery.KeyPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocalUsersThirdPartyKey", arg0)
	ret0, _ := ret[0].(*bakery.KeyPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLocalUsersThirdPartyKey indicates an expected call of GetLocalUsersThirdPartyKey.
func (mr *MockStateMockRecorder) GetLocalUsersThirdPartyKey(arg0 any) *MockStateGetLocalUsersThirdPartyKeyCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocalUsersThirdPartyKey", reflect.TypeOf((*MockState)(nil).GetLocalUsersThirdPartyKey), arg0)
	return &MockStateGetLocalUsersThirdPartyKeyCall{Call: call}
}

// MockStateGetLocalUsersThirdPartyKeyCall wrap *gomock.Call
type MockStateGetLocalUsersThirdPartyKeyCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateGetLocalUsersThirdPartyKeyCall) Return(arg0 *bakery.KeyPair, arg1 error) *MockStateGetLocalUsersThirdPartyKeyCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateGetLocalUsersThirdPartyKeyCall) Do(f func(context.Context) (*bakery.KeyPair, error)) *MockStateGetLocalUsersThirdPartyKeyCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateGetLocalUsersThirdPartyKeyCall) DoAndReturn(f func(context.Context) (*bakery.KeyPair, error)) *MockStateGetLocalUsersThirdPartyKeyCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetOffersThirdPartyKey mocks base method.
func (m *MockState) GetOffersThirdPartyKey(arg0 context.Context) (*bakery.KeyPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOffersThirdPartyKey", arg0)
	ret0, _ := ret[0].(*bakery.KeyPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOffersThirdPartyKey indicates an expected call of GetOffersThirdPartyKey.
func (mr *MockStateMockRecorder) GetOffersThirdPartyKey(arg0 any) *MockStateGetOffersThirdPartyKeyCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOffersThirdPartyKey", reflect.TypeOf((*MockState)(nil).GetOffersThirdPartyKey), arg0)
	return &MockStateGetOffersThirdPartyKeyCall{Call: call}
}

// MockStateGetOffersThirdPartyKeyCall wrap *gomock.Call
type MockStateGetOffersThirdPartyKeyCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateGetOffersThirdPartyKeyCall) Return(arg0 *bakery.KeyPair, arg1 error) *MockStateGetOffersThirdPartyKeyCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateGetOffersThirdPartyKeyCall) Do(f func(context.Context) (*bakery.KeyPair, error)) *MockStateGetOffersThirdPartyKeyCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateGetOffersThirdPartyKeyCall) DoAndReturn(f func(context.Context) (*bakery.KeyPair, error)) *MockStateGetOffersThirdPartyKeyCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// InitialiseBakeryConfig mocks base method.
func (m *MockState) InitialiseBakeryConfig(arg0 context.Context, arg1, arg2, arg3, arg4 *bakery.KeyPair) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitialiseBakeryConfig", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitialiseBakeryConfig indicates an expected call of InitialiseBakeryConfig.
func (mr *MockStateMockRecorder) InitialiseBakeryConfig(arg0, arg1, arg2, arg3, arg4 any) *MockStateInitialiseBakeryConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitialiseBakeryConfig", reflect.TypeOf((*MockState)(nil).InitialiseBakeryConfig), arg0, arg1, arg2, arg3, arg4)
	return &MockStateInitialiseBakeryConfigCall{Call: call}
}

// MockStateInitialiseBakeryConfigCall wrap *gomock.Call
type MockStateInitialiseBakeryConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateInitialiseBakeryConfigCall) Return(arg0 error) *MockStateInitialiseBakeryConfigCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateInitialiseBakeryConfigCall) Do(f func(context.Context, *bakery.KeyPair, *bakery.KeyPair, *bakery.KeyPair, *bakery.KeyPair) error) *MockStateInitialiseBakeryConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateInitialiseBakeryConfigCall) DoAndReturn(f func(context.Context, *bakery.KeyPair, *bakery.KeyPair, *bakery.KeyPair, *bakery.KeyPair) error) *MockStateInitialiseBakeryConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// InsertKey mocks base method.
func (m *MockState) InsertKey(arg0 context.Context, arg1 macaroon.RootKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertKey", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertKey indicates an expected call of InsertKey.
func (mr *MockStateMockRecorder) InsertKey(arg0, arg1 any) *MockStateInsertKeyCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertKey", reflect.TypeOf((*MockState)(nil).InsertKey), arg0, arg1)
	return &MockStateInsertKeyCall{Call: call}
}

// MockStateInsertKeyCall wrap *gomock.Call
type MockStateInsertKeyCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateInsertKeyCall) Return(arg0 error) *MockStateInsertKeyCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateInsertKeyCall) Do(f func(context.Context, macaroon.RootKey) error) *MockStateInsertKeyCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateInsertKeyCall) DoAndReturn(f func(context.Context, macaroon.RootKey) error) *MockStateInsertKeyCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
