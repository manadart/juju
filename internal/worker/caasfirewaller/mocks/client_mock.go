// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/worker/caasfirewaller (interfaces: Client,CAASFirewallerAPI,LifeGetter)
//
// Generated by this command:
//
//	mockgen -typed -package mocks -destination mocks/client_mock.go github.com/juju/juju/internal/worker/caasfirewaller Client,CAASFirewallerAPI,LifeGetter
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	charms "github.com/juju/juju/api/common/charms"
	life "github.com/juju/juju/core/life"
	watcher "github.com/juju/juju/core/watcher"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// ApplicationCharmInfo mocks base method.
func (m *MockClient) ApplicationCharmInfo(arg0 context.Context, arg1 string) (*charms.CharmInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplicationCharmInfo", arg0, arg1)
	ret0, _ := ret[0].(*charms.CharmInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplicationCharmInfo indicates an expected call of ApplicationCharmInfo.
func (mr *MockClientMockRecorder) ApplicationCharmInfo(arg0, arg1 any) *MockClientApplicationCharmInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplicationCharmInfo", reflect.TypeOf((*MockClient)(nil).ApplicationCharmInfo), arg0, arg1)
	return &MockClientApplicationCharmInfoCall{Call: call}
}

// MockClientApplicationCharmInfoCall wrap *gomock.Call
type MockClientApplicationCharmInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockClientApplicationCharmInfoCall) Return(arg0 *charms.CharmInfo, arg1 error) *MockClientApplicationCharmInfoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockClientApplicationCharmInfoCall) Do(f func(context.Context, string) (*charms.CharmInfo, error)) *MockClientApplicationCharmInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockClientApplicationCharmInfoCall) DoAndReturn(f func(context.Context, string) (*charms.CharmInfo, error)) *MockClientApplicationCharmInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsExposed mocks base method.
func (m *MockClient) IsExposed(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExposed", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExposed indicates an expected call of IsExposed.
func (mr *MockClientMockRecorder) IsExposed(arg0, arg1 any) *MockClientIsExposedCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExposed", reflect.TypeOf((*MockClient)(nil).IsExposed), arg0, arg1)
	return &MockClientIsExposedCall{Call: call}
}

// MockClientIsExposedCall wrap *gomock.Call
type MockClientIsExposedCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockClientIsExposedCall) Return(arg0 bool, arg1 error) *MockClientIsExposedCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockClientIsExposedCall) Do(f func(context.Context, string) (bool, error)) *MockClientIsExposedCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockClientIsExposedCall) DoAndReturn(f func(context.Context, string) (bool, error)) *MockClientIsExposedCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Life mocks base method.
func (m *MockClient) Life(arg0 context.Context, arg1 string) (life.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Life", arg0, arg1)
	ret0, _ := ret[0].(life.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Life indicates an expected call of Life.
func (mr *MockClientMockRecorder) Life(arg0, arg1 any) *MockClientLifeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Life", reflect.TypeOf((*MockClient)(nil).Life), arg0, arg1)
	return &MockClientLifeCall{Call: call}
}

// MockClientLifeCall wrap *gomock.Call
type MockClientLifeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockClientLifeCall) Return(arg0 life.Value, arg1 error) *MockClientLifeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockClientLifeCall) Do(f func(context.Context, string) (life.Value, error)) *MockClientLifeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockClientLifeCall) DoAndReturn(f func(context.Context, string) (life.Value, error)) *MockClientLifeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchApplication mocks base method.
func (m *MockClient) WatchApplication(arg0 context.Context, arg1 string) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchApplication", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchApplication indicates an expected call of WatchApplication.
func (mr *MockClientMockRecorder) WatchApplication(arg0, arg1 any) *MockClientWatchApplicationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchApplication", reflect.TypeOf((*MockClient)(nil).WatchApplication), arg0, arg1)
	return &MockClientWatchApplicationCall{Call: call}
}

// MockClientWatchApplicationCall wrap *gomock.Call
type MockClientWatchApplicationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockClientWatchApplicationCall) Return(arg0 watcher.Watcher[struct{}], arg1 error) *MockClientWatchApplicationCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockClientWatchApplicationCall) Do(f func(context.Context, string) (watcher.Watcher[struct{}], error)) *MockClientWatchApplicationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockClientWatchApplicationCall) DoAndReturn(f func(context.Context, string) (watcher.Watcher[struct{}], error)) *MockClientWatchApplicationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchApplications mocks base method.
func (m *MockClient) WatchApplications(arg0 context.Context) (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchApplications", arg0)
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchApplications indicates an expected call of WatchApplications.
func (mr *MockClientMockRecorder) WatchApplications(arg0 any) *MockClientWatchApplicationsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchApplications", reflect.TypeOf((*MockClient)(nil).WatchApplications), arg0)
	return &MockClientWatchApplicationsCall{Call: call}
}

// MockClientWatchApplicationsCall wrap *gomock.Call
type MockClientWatchApplicationsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockClientWatchApplicationsCall) Return(arg0 watcher.Watcher[[]string], arg1 error) *MockClientWatchApplicationsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockClientWatchApplicationsCall) Do(f func(context.Context) (watcher.Watcher[[]string], error)) *MockClientWatchApplicationsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockClientWatchApplicationsCall) DoAndReturn(f func(context.Context) (watcher.Watcher[[]string], error)) *MockClientWatchApplicationsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockCAASFirewallerAPI is a mock of CAASFirewallerAPI interface.
type MockCAASFirewallerAPI struct {
	ctrl     *gomock.Controller
	recorder *MockCAASFirewallerAPIMockRecorder
}

// MockCAASFirewallerAPIMockRecorder is the mock recorder for MockCAASFirewallerAPI.
type MockCAASFirewallerAPIMockRecorder struct {
	mock *MockCAASFirewallerAPI
}

// NewMockCAASFirewallerAPI creates a new mock instance.
func NewMockCAASFirewallerAPI(ctrl *gomock.Controller) *MockCAASFirewallerAPI {
	mock := &MockCAASFirewallerAPI{ctrl: ctrl}
	mock.recorder = &MockCAASFirewallerAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCAASFirewallerAPI) EXPECT() *MockCAASFirewallerAPIMockRecorder {
	return m.recorder
}

// ApplicationCharmInfo mocks base method.
func (m *MockCAASFirewallerAPI) ApplicationCharmInfo(arg0 context.Context, arg1 string) (*charms.CharmInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplicationCharmInfo", arg0, arg1)
	ret0, _ := ret[0].(*charms.CharmInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplicationCharmInfo indicates an expected call of ApplicationCharmInfo.
func (mr *MockCAASFirewallerAPIMockRecorder) ApplicationCharmInfo(arg0, arg1 any) *MockCAASFirewallerAPIApplicationCharmInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplicationCharmInfo", reflect.TypeOf((*MockCAASFirewallerAPI)(nil).ApplicationCharmInfo), arg0, arg1)
	return &MockCAASFirewallerAPIApplicationCharmInfoCall{Call: call}
}

// MockCAASFirewallerAPIApplicationCharmInfoCall wrap *gomock.Call
type MockCAASFirewallerAPIApplicationCharmInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASFirewallerAPIApplicationCharmInfoCall) Return(arg0 *charms.CharmInfo, arg1 error) *MockCAASFirewallerAPIApplicationCharmInfoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASFirewallerAPIApplicationCharmInfoCall) Do(f func(context.Context, string) (*charms.CharmInfo, error)) *MockCAASFirewallerAPIApplicationCharmInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASFirewallerAPIApplicationCharmInfoCall) DoAndReturn(f func(context.Context, string) (*charms.CharmInfo, error)) *MockCAASFirewallerAPIApplicationCharmInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsExposed mocks base method.
func (m *MockCAASFirewallerAPI) IsExposed(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExposed", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExposed indicates an expected call of IsExposed.
func (mr *MockCAASFirewallerAPIMockRecorder) IsExposed(arg0, arg1 any) *MockCAASFirewallerAPIIsExposedCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExposed", reflect.TypeOf((*MockCAASFirewallerAPI)(nil).IsExposed), arg0, arg1)
	return &MockCAASFirewallerAPIIsExposedCall{Call: call}
}

// MockCAASFirewallerAPIIsExposedCall wrap *gomock.Call
type MockCAASFirewallerAPIIsExposedCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASFirewallerAPIIsExposedCall) Return(arg0 bool, arg1 error) *MockCAASFirewallerAPIIsExposedCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASFirewallerAPIIsExposedCall) Do(f func(context.Context, string) (bool, error)) *MockCAASFirewallerAPIIsExposedCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASFirewallerAPIIsExposedCall) DoAndReturn(f func(context.Context, string) (bool, error)) *MockCAASFirewallerAPIIsExposedCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchApplication mocks base method.
func (m *MockCAASFirewallerAPI) WatchApplication(arg0 context.Context, arg1 string) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchApplication", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchApplication indicates an expected call of WatchApplication.
func (mr *MockCAASFirewallerAPIMockRecorder) WatchApplication(arg0, arg1 any) *MockCAASFirewallerAPIWatchApplicationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchApplication", reflect.TypeOf((*MockCAASFirewallerAPI)(nil).WatchApplication), arg0, arg1)
	return &MockCAASFirewallerAPIWatchApplicationCall{Call: call}
}

// MockCAASFirewallerAPIWatchApplicationCall wrap *gomock.Call
type MockCAASFirewallerAPIWatchApplicationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASFirewallerAPIWatchApplicationCall) Return(arg0 watcher.Watcher[struct{}], arg1 error) *MockCAASFirewallerAPIWatchApplicationCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASFirewallerAPIWatchApplicationCall) Do(f func(context.Context, string) (watcher.Watcher[struct{}], error)) *MockCAASFirewallerAPIWatchApplicationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASFirewallerAPIWatchApplicationCall) DoAndReturn(f func(context.Context, string) (watcher.Watcher[struct{}], error)) *MockCAASFirewallerAPIWatchApplicationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchApplications mocks base method.
func (m *MockCAASFirewallerAPI) WatchApplications(arg0 context.Context) (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchApplications", arg0)
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchApplications indicates an expected call of WatchApplications.
func (mr *MockCAASFirewallerAPIMockRecorder) WatchApplications(arg0 any) *MockCAASFirewallerAPIWatchApplicationsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchApplications", reflect.TypeOf((*MockCAASFirewallerAPI)(nil).WatchApplications), arg0)
	return &MockCAASFirewallerAPIWatchApplicationsCall{Call: call}
}

// MockCAASFirewallerAPIWatchApplicationsCall wrap *gomock.Call
type MockCAASFirewallerAPIWatchApplicationsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCAASFirewallerAPIWatchApplicationsCall) Return(arg0 watcher.Watcher[[]string], arg1 error) *MockCAASFirewallerAPIWatchApplicationsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCAASFirewallerAPIWatchApplicationsCall) Do(f func(context.Context) (watcher.Watcher[[]string], error)) *MockCAASFirewallerAPIWatchApplicationsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCAASFirewallerAPIWatchApplicationsCall) DoAndReturn(f func(context.Context) (watcher.Watcher[[]string], error)) *MockCAASFirewallerAPIWatchApplicationsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockLifeGetter is a mock of LifeGetter interface.
type MockLifeGetter struct {
	ctrl     *gomock.Controller
	recorder *MockLifeGetterMockRecorder
}

// MockLifeGetterMockRecorder is the mock recorder for MockLifeGetter.
type MockLifeGetterMockRecorder struct {
	mock *MockLifeGetter
}

// NewMockLifeGetter creates a new mock instance.
func NewMockLifeGetter(ctrl *gomock.Controller) *MockLifeGetter {
	mock := &MockLifeGetter{ctrl: ctrl}
	mock.recorder = &MockLifeGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLifeGetter) EXPECT() *MockLifeGetterMockRecorder {
	return m.recorder
}

// Life mocks base method.
func (m *MockLifeGetter) Life(arg0 context.Context, arg1 string) (life.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Life", arg0, arg1)
	ret0, _ := ret[0].(life.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Life indicates an expected call of Life.
func (mr *MockLifeGetterMockRecorder) Life(arg0, arg1 any) *MockLifeGetterLifeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Life", reflect.TypeOf((*MockLifeGetter)(nil).Life), arg0, arg1)
	return &MockLifeGetterLifeCall{Call: call}
}

// MockLifeGetterLifeCall wrap *gomock.Call
type MockLifeGetterLifeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockLifeGetterLifeCall) Return(arg0 life.Value, arg1 error) *MockLifeGetterLifeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockLifeGetterLifeCall) Do(f func(context.Context, string) (life.Value, error)) *MockLifeGetterLifeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockLifeGetterLifeCall) DoAndReturn(f func(context.Context, string) (life.Value, error)) *MockLifeGetterLifeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
