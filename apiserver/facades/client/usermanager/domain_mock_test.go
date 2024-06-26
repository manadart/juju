// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/usermanager (interfaces: AccessService)
//
// Generated by this command:
//
//	mockgen -typed -package usermanager_test -destination domain_mock_test.go github.com/juju/juju/apiserver/facades/client/usermanager AccessService
//

// Package usermanager_test is a generated GoMock package.
package usermanager_test

import (
	context "context"
	reflect "reflect"

	permission "github.com/juju/juju/core/permission"
	user "github.com/juju/juju/core/user"
	service "github.com/juju/juju/domain/access/service"
	auth "github.com/juju/juju/internal/auth"
	gomock "go.uber.org/mock/gomock"
)

// MockAccessService is a mock of AccessService interface.
type MockAccessService struct {
	ctrl     *gomock.Controller
	recorder *MockAccessServiceMockRecorder
}

// MockAccessServiceMockRecorder is the mock recorder for MockAccessService.
type MockAccessServiceMockRecorder struct {
	mock *MockAccessService
}

// NewMockAccessService creates a new mock instance.
func NewMockAccessService(ctrl *gomock.Controller) *MockAccessService {
	mock := &MockAccessService{ctrl: ctrl}
	mock.recorder = &MockAccessServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccessService) EXPECT() *MockAccessServiceMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *MockAccessService) AddUser(arg0 context.Context, arg1 service.AddUserArg) (user.UUID, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1)
	ret0, _ := ret[0].(user.UUID)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AddUser indicates an expected call of AddUser.
func (mr *MockAccessServiceMockRecorder) AddUser(arg0, arg1 any) *MockAccessServiceAddUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockAccessService)(nil).AddUser), arg0, arg1)
	return &MockAccessServiceAddUserCall{Call: call}
}

// MockAccessServiceAddUserCall wrap *gomock.Call
type MockAccessServiceAddUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceAddUserCall) Return(arg0 user.UUID, arg1 []byte, arg2 error) *MockAccessServiceAddUserCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceAddUserCall) Do(f func(context.Context, service.AddUserArg) (user.UUID, []byte, error)) *MockAccessServiceAddUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceAddUserCall) DoAndReturn(f func(context.Context, service.AddUserArg) (user.UUID, []byte, error)) *MockAccessServiceAddUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DisableUserAuthentication mocks base method.
func (m *MockAccessService) DisableUserAuthentication(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableUserAuthentication", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DisableUserAuthentication indicates an expected call of DisableUserAuthentication.
func (mr *MockAccessServiceMockRecorder) DisableUserAuthentication(arg0, arg1 any) *MockAccessServiceDisableUserAuthenticationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableUserAuthentication", reflect.TypeOf((*MockAccessService)(nil).DisableUserAuthentication), arg0, arg1)
	return &MockAccessServiceDisableUserAuthenticationCall{Call: call}
}

// MockAccessServiceDisableUserAuthenticationCall wrap *gomock.Call
type MockAccessServiceDisableUserAuthenticationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceDisableUserAuthenticationCall) Return(arg0 error) *MockAccessServiceDisableUserAuthenticationCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceDisableUserAuthenticationCall) Do(f func(context.Context, string) error) *MockAccessServiceDisableUserAuthenticationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceDisableUserAuthenticationCall) DoAndReturn(f func(context.Context, string) error) *MockAccessServiceDisableUserAuthenticationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// EnableUserAuthentication mocks base method.
func (m *MockAccessService) EnableUserAuthentication(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableUserAuthentication", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnableUserAuthentication indicates an expected call of EnableUserAuthentication.
func (mr *MockAccessServiceMockRecorder) EnableUserAuthentication(arg0, arg1 any) *MockAccessServiceEnableUserAuthenticationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableUserAuthentication", reflect.TypeOf((*MockAccessService)(nil).EnableUserAuthentication), arg0, arg1)
	return &MockAccessServiceEnableUserAuthenticationCall{Call: call}
}

// MockAccessServiceEnableUserAuthenticationCall wrap *gomock.Call
type MockAccessServiceEnableUserAuthenticationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceEnableUserAuthenticationCall) Return(arg0 error) *MockAccessServiceEnableUserAuthenticationCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceEnableUserAuthenticationCall) Do(f func(context.Context, string) error) *MockAccessServiceEnableUserAuthenticationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceEnableUserAuthenticationCall) DoAndReturn(f func(context.Context, string) error) *MockAccessServiceEnableUserAuthenticationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetAllUsers mocks base method.
func (m *MockAccessService) GetAllUsers(arg0 context.Context) ([]user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers", arg0)
	ret0, _ := ret[0].([]user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockAccessServiceMockRecorder) GetAllUsers(arg0 any) *MockAccessServiceGetAllUsersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockAccessService)(nil).GetAllUsers), arg0)
	return &MockAccessServiceGetAllUsersCall{Call: call}
}

// MockAccessServiceGetAllUsersCall wrap *gomock.Call
type MockAccessServiceGetAllUsersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceGetAllUsersCall) Return(arg0 []user.User, arg1 error) *MockAccessServiceGetAllUsersCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceGetAllUsersCall) Do(f func(context.Context) ([]user.User, error)) *MockAccessServiceGetAllUsersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceGetAllUsersCall) DoAndReturn(f func(context.Context) ([]user.User, error)) *MockAccessServiceGetAllUsersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetUserByName mocks base method.
func (m *MockAccessService) GetUserByName(arg0 context.Context, arg1 string) (user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByName", arg0, arg1)
	ret0, _ := ret[0].(user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByName indicates an expected call of GetUserByName.
func (mr *MockAccessServiceMockRecorder) GetUserByName(arg0, arg1 any) *MockAccessServiceGetUserByNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByName", reflect.TypeOf((*MockAccessService)(nil).GetUserByName), arg0, arg1)
	return &MockAccessServiceGetUserByNameCall{Call: call}
}

// MockAccessServiceGetUserByNameCall wrap *gomock.Call
type MockAccessServiceGetUserByNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceGetUserByNameCall) Return(arg0 user.User, arg1 error) *MockAccessServiceGetUserByNameCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceGetUserByNameCall) Do(f func(context.Context, string) (user.User, error)) *MockAccessServiceGetUserByNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceGetUserByNameCall) DoAndReturn(f func(context.Context, string) (user.User, error)) *MockAccessServiceGetUserByNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ReadUserAccessForTarget mocks base method.
func (m *MockAccessService) ReadUserAccessForTarget(arg0 context.Context, arg1 string, arg2 permission.ID) (permission.UserAccess, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadUserAccessForTarget", arg0, arg1, arg2)
	ret0, _ := ret[0].(permission.UserAccess)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadUserAccessForTarget indicates an expected call of ReadUserAccessForTarget.
func (mr *MockAccessServiceMockRecorder) ReadUserAccessForTarget(arg0, arg1, arg2 any) *MockAccessServiceReadUserAccessForTargetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadUserAccessForTarget", reflect.TypeOf((*MockAccessService)(nil).ReadUserAccessForTarget), arg0, arg1, arg2)
	return &MockAccessServiceReadUserAccessForTargetCall{Call: call}
}

// MockAccessServiceReadUserAccessForTargetCall wrap *gomock.Call
type MockAccessServiceReadUserAccessForTargetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceReadUserAccessForTargetCall) Return(arg0 permission.UserAccess, arg1 error) *MockAccessServiceReadUserAccessForTargetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceReadUserAccessForTargetCall) Do(f func(context.Context, string, permission.ID) (permission.UserAccess, error)) *MockAccessServiceReadUserAccessForTargetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceReadUserAccessForTargetCall) DoAndReturn(f func(context.Context, string, permission.ID) (permission.UserAccess, error)) *MockAccessServiceReadUserAccessForTargetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RemoveUser mocks base method.
func (m *MockAccessService) RemoveUser(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUser indicates an expected call of RemoveUser.
func (mr *MockAccessServiceMockRecorder) RemoveUser(arg0, arg1 any) *MockAccessServiceRemoveUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUser", reflect.TypeOf((*MockAccessService)(nil).RemoveUser), arg0, arg1)
	return &MockAccessServiceRemoveUserCall{Call: call}
}

// MockAccessServiceRemoveUserCall wrap *gomock.Call
type MockAccessServiceRemoveUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceRemoveUserCall) Return(arg0 error) *MockAccessServiceRemoveUserCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceRemoveUserCall) Do(f func(context.Context, string) error) *MockAccessServiceRemoveUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceRemoveUserCall) DoAndReturn(f func(context.Context, string) error) *MockAccessServiceRemoveUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ResetPassword mocks base method.
func (m *MockAccessService) ResetPassword(arg0 context.Context, arg1 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *MockAccessServiceMockRecorder) ResetPassword(arg0, arg1 any) *MockAccessServiceResetPasswordCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockAccessService)(nil).ResetPassword), arg0, arg1)
	return &MockAccessServiceResetPasswordCall{Call: call}
}

// MockAccessServiceResetPasswordCall wrap *gomock.Call
type MockAccessServiceResetPasswordCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceResetPasswordCall) Return(arg0 []byte, arg1 error) *MockAccessServiceResetPasswordCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceResetPasswordCall) Do(f func(context.Context, string) ([]byte, error)) *MockAccessServiceResetPasswordCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceResetPasswordCall) DoAndReturn(f func(context.Context, string) ([]byte, error)) *MockAccessServiceResetPasswordCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetPassword mocks base method.
func (m *MockAccessService) SetPassword(arg0 context.Context, arg1 string, arg2 auth.Password) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPassword indicates an expected call of SetPassword.
func (mr *MockAccessServiceMockRecorder) SetPassword(arg0, arg1, arg2 any) *MockAccessServiceSetPasswordCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPassword", reflect.TypeOf((*MockAccessService)(nil).SetPassword), arg0, arg1, arg2)
	return &MockAccessServiceSetPasswordCall{Call: call}
}

// MockAccessServiceSetPasswordCall wrap *gomock.Call
type MockAccessServiceSetPasswordCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceSetPasswordCall) Return(arg0 error) *MockAccessServiceSetPasswordCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceSetPasswordCall) Do(f func(context.Context, string, auth.Password) error) *MockAccessServiceSetPasswordCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceSetPasswordCall) DoAndReturn(f func(context.Context, string, auth.Password) error) *MockAccessServiceSetPasswordCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
