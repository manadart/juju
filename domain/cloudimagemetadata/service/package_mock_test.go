// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/domain/cloudimagemetadata/service (interfaces: State)
//
// Generated by this command:
//
//	mockgen -typed -package service -destination package_mock_test.go github.com/juju/juju/domain/cloudimagemetadata/service State
//

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	set "github.com/juju/collections/set"
	cloudimagemetadata "github.com/juju/juju/domain/cloudimagemetadata"
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

// AllCloudImageMetadata mocks base method.
func (m *MockState) AllCloudImageMetadata(arg0 context.Context) ([]cloudimagemetadata.Metadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllCloudImageMetadata", arg0)
	ret0, _ := ret[0].([]cloudimagemetadata.Metadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllCloudImageMetadata indicates an expected call of AllCloudImageMetadata.
func (mr *MockStateMockRecorder) AllCloudImageMetadata(arg0 any) *MockStateAllCloudImageMetadataCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllCloudImageMetadata", reflect.TypeOf((*MockState)(nil).AllCloudImageMetadata), arg0)
	return &MockStateAllCloudImageMetadataCall{Call: call}
}

// MockStateAllCloudImageMetadataCall wrap *gomock.Call
type MockStateAllCloudImageMetadataCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateAllCloudImageMetadataCall) Return(arg0 []cloudimagemetadata.Metadata, arg1 error) *MockStateAllCloudImageMetadataCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateAllCloudImageMetadataCall) Do(f func(context.Context) ([]cloudimagemetadata.Metadata, error)) *MockStateAllCloudImageMetadataCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateAllCloudImageMetadataCall) DoAndReturn(f func(context.Context) ([]cloudimagemetadata.Metadata, error)) *MockStateAllCloudImageMetadataCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DeleteMetadataWithImageID mocks base method.
func (m *MockState) DeleteMetadataWithImageID(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMetadataWithImageID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMetadataWithImageID indicates an expected call of DeleteMetadataWithImageID.
func (mr *MockStateMockRecorder) DeleteMetadataWithImageID(arg0, arg1 any) *MockStateDeleteMetadataWithImageIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMetadataWithImageID", reflect.TypeOf((*MockState)(nil).DeleteMetadataWithImageID), arg0, arg1)
	return &MockStateDeleteMetadataWithImageIDCall{Call: call}
}

// MockStateDeleteMetadataWithImageIDCall wrap *gomock.Call
type MockStateDeleteMetadataWithImageIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateDeleteMetadataWithImageIDCall) Return(arg0 error) *MockStateDeleteMetadataWithImageIDCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateDeleteMetadataWithImageIDCall) Do(f func(context.Context, string) error) *MockStateDeleteMetadataWithImageIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateDeleteMetadataWithImageIDCall) DoAndReturn(f func(context.Context, string) error) *MockStateDeleteMetadataWithImageIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// FindMetadata mocks base method.
func (m *MockState) FindMetadata(arg0 context.Context, arg1 cloudimagemetadata.MetadataFilter) ([]cloudimagemetadata.Metadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMetadata", arg0, arg1)
	ret0, _ := ret[0].([]cloudimagemetadata.Metadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMetadata indicates an expected call of FindMetadata.
func (mr *MockStateMockRecorder) FindMetadata(arg0, arg1 any) *MockStateFindMetadataCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMetadata", reflect.TypeOf((*MockState)(nil).FindMetadata), arg0, arg1)
	return &MockStateFindMetadataCall{Call: call}
}

// MockStateFindMetadataCall wrap *gomock.Call
type MockStateFindMetadataCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateFindMetadataCall) Return(arg0 []cloudimagemetadata.Metadata, arg1 error) *MockStateFindMetadataCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateFindMetadataCall) Do(f func(context.Context, cloudimagemetadata.MetadataFilter) ([]cloudimagemetadata.Metadata, error)) *MockStateFindMetadataCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateFindMetadataCall) DoAndReturn(f func(context.Context, cloudimagemetadata.MetadataFilter) ([]cloudimagemetadata.Metadata, error)) *MockStateFindMetadataCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SaveMetadata mocks base method.
func (m *MockState) SaveMetadata(arg0 context.Context, arg1 []cloudimagemetadata.Metadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMetadata", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveMetadata indicates an expected call of SaveMetadata.
func (mr *MockStateMockRecorder) SaveMetadata(arg0, arg1 any) *MockStateSaveMetadataCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMetadata", reflect.TypeOf((*MockState)(nil).SaveMetadata), arg0, arg1)
	return &MockStateSaveMetadataCall{Call: call}
}

// MockStateSaveMetadataCall wrap *gomock.Call
type MockStateSaveMetadataCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateSaveMetadataCall) Return(arg0 error) *MockStateSaveMetadataCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateSaveMetadataCall) Do(f func(context.Context, []cloudimagemetadata.Metadata) error) *MockStateSaveMetadataCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateSaveMetadataCall) DoAndReturn(f func(context.Context, []cloudimagemetadata.Metadata) error) *MockStateSaveMetadataCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SupportedArchitectures mocks base method.
func (m *MockState) SupportedArchitectures(arg0 context.Context) set.Strings {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportedArchitectures", arg0)
	ret0, _ := ret[0].(set.Strings)
	return ret0
}

// SupportedArchitectures indicates an expected call of SupportedArchitectures.
func (mr *MockStateMockRecorder) SupportedArchitectures(arg0 any) *MockStateSupportedArchitecturesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportedArchitectures", reflect.TypeOf((*MockState)(nil).SupportedArchitectures), arg0)
	return &MockStateSupportedArchitecturesCall{Call: call}
}

// MockStateSupportedArchitecturesCall wrap *gomock.Call
type MockStateSupportedArchitecturesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateSupportedArchitecturesCall) Return(arg0 set.Strings) *MockStateSupportedArchitecturesCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateSupportedArchitecturesCall) Do(f func(context.Context) set.Strings) *MockStateSupportedArchitecturesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateSupportedArchitecturesCall) DoAndReturn(f func(context.Context) set.Strings) *MockStateSupportedArchitecturesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
