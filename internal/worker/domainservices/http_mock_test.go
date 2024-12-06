// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/core/http (interfaces: HTTPClientGetter,HTTPClient)
//
// Generated by this command:
//
//	mockgen -typed -package domainservices -destination http_mock_test.go github.com/juju/juju/core/http HTTPClientGetter,HTTPClient
//

// Package domainservices is a generated GoMock package.
package domainservices

import (
	context "context"
	http0 "net/http"
	reflect "reflect"

	http "github.com/juju/juju/core/http"
	gomock "go.uber.org/mock/gomock"
)

// MockHTTPClientGetter is a mock of HTTPClientGetter interface.
type MockHTTPClientGetter struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPClientGetterMockRecorder
}

// MockHTTPClientGetterMockRecorder is the mock recorder for MockHTTPClientGetter.
type MockHTTPClientGetterMockRecorder struct {
	mock *MockHTTPClientGetter
}

// NewMockHTTPClientGetter creates a new mock instance.
func NewMockHTTPClientGetter(ctrl *gomock.Controller) *MockHTTPClientGetter {
	mock := &MockHTTPClientGetter{ctrl: ctrl}
	mock.recorder = &MockHTTPClientGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHTTPClientGetter) EXPECT() *MockHTTPClientGetterMockRecorder {
	return m.recorder
}

// GetHTTPClient mocks base method.
func (m *MockHTTPClientGetter) GetHTTPClient(arg0 context.Context, arg1 http.Purpose) (http.HTTPClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHTTPClient", arg0, arg1)
	ret0, _ := ret[0].(http.HTTPClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHTTPClient indicates an expected call of GetHTTPClient.
func (mr *MockHTTPClientGetterMockRecorder) GetHTTPClient(arg0, arg1 any) *MockHTTPClientGetterGetHTTPClientCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHTTPClient", reflect.TypeOf((*MockHTTPClientGetter)(nil).GetHTTPClient), arg0, arg1)
	return &MockHTTPClientGetterGetHTTPClientCall{Call: call}
}

// MockHTTPClientGetterGetHTTPClientCall wrap *gomock.Call
type MockHTTPClientGetterGetHTTPClientCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockHTTPClientGetterGetHTTPClientCall) Return(arg0 http.HTTPClient, arg1 error) *MockHTTPClientGetterGetHTTPClientCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockHTTPClientGetterGetHTTPClientCall) Do(f func(context.Context, http.Purpose) (http.HTTPClient, error)) *MockHTTPClientGetterGetHTTPClientCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockHTTPClientGetterGetHTTPClientCall) DoAndReturn(f func(context.Context, http.Purpose) (http.HTTPClient, error)) *MockHTTPClientGetterGetHTTPClientCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockHTTPClient is a mock of HTTPClient interface.
type MockHTTPClient struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPClientMockRecorder
}

// MockHTTPClientMockRecorder is the mock recorder for MockHTTPClient.
type MockHTTPClientMockRecorder struct {
	mock *MockHTTPClient
}

// NewMockHTTPClient creates a new mock instance.
func NewMockHTTPClient(ctrl *gomock.Controller) *MockHTTPClient {
	mock := &MockHTTPClient{ctrl: ctrl}
	mock.recorder = &MockHTTPClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHTTPClient) EXPECT() *MockHTTPClientMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockHTTPClient) Do(arg0 *http0.Request) (*http0.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(*http0.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockHTTPClientMockRecorder) Do(arg0 any) *MockHTTPClientDoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockHTTPClient)(nil).Do), arg0)
	return &MockHTTPClientDoCall{Call: call}
}

// MockHTTPClientDoCall wrap *gomock.Call
type MockHTTPClientDoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockHTTPClientDoCall) Return(arg0 *http0.Response, arg1 error) *MockHTTPClientDoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockHTTPClientDoCall) Do(f func(*http0.Request) (*http0.Response, error)) *MockHTTPClientDoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockHTTPClientDoCall) DoAndReturn(f func(*http0.Request) (*http0.Response, error)) *MockHTTPClientDoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}