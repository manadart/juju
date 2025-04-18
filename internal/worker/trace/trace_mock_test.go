// Code generated by MockGen. DO NOT EDIT.
// Source: go.opentelemetry.io/otel/trace (interfaces: Span)
//
// Generated by this command:
//
//	mockgen -typed -package trace -destination trace_mock_test.go go.opentelemetry.io/otel/trace Span
//

// Package trace is a generated GoMock package.
package trace

import (
	reflect "reflect"

	attribute "go.opentelemetry.io/otel/attribute"
	codes "go.opentelemetry.io/otel/codes"
	trace "go.opentelemetry.io/otel/trace"
	gomock "go.uber.org/mock/gomock"
)

// MockSpan is a mock of Span interface.
type MockSpan struct {
	ctrl     *gomock.Controller
	recorder *MockSpanMockRecorder
}

// MockSpanMockRecorder is the mock recorder for MockSpan.
type MockSpanMockRecorder struct {
	mock *MockSpan
}

// NewMockSpan creates a new mock instance.
func NewMockSpan(ctrl *gomock.Controller) *MockSpan {
	mock := &MockSpan{ctrl: ctrl}
	mock.recorder = &MockSpanMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpan) EXPECT() *MockSpanMockRecorder {
	return m.recorder
}

// AddEvent mocks base method.
func (m *MockSpan) AddEvent(arg0 string, arg1 ...trace.EventOption) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "AddEvent", varargs...)
}

// AddEvent indicates an expected call of AddEvent.
func (mr *MockSpanMockRecorder) AddEvent(arg0 any, arg1 ...any) *MockSpanAddEventCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvent", reflect.TypeOf((*MockSpan)(nil).AddEvent), varargs...)
	return &MockSpanAddEventCall{Call: call}
}

// MockSpanAddEventCall wrap *gomock.Call
type MockSpanAddEventCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSpanAddEventCall) Return() *MockSpanAddEventCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSpanAddEventCall) Do(f func(string, ...trace.EventOption)) *MockSpanAddEventCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSpanAddEventCall) DoAndReturn(f func(string, ...trace.EventOption)) *MockSpanAddEventCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// AddLink mocks base method.
func (m *MockSpan) AddLink(arg0 trace.Link) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddLink", arg0)
}

// AddLink indicates an expected call of AddLink.
func (mr *MockSpanMockRecorder) AddLink(arg0 any) *MockSpanAddLinkCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLink", reflect.TypeOf((*MockSpan)(nil).AddLink), arg0)
	return &MockSpanAddLinkCall{Call: call}
}

// MockSpanAddLinkCall wrap *gomock.Call
type MockSpanAddLinkCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSpanAddLinkCall) Return() *MockSpanAddLinkCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSpanAddLinkCall) Do(f func(trace.Link)) *MockSpanAddLinkCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSpanAddLinkCall) DoAndReturn(f func(trace.Link)) *MockSpanAddLinkCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// End mocks base method.
func (m *MockSpan) End(arg0 ...trace.SpanEndOption) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "End", varargs...)
}

// End indicates an expected call of End.
func (mr *MockSpanMockRecorder) End(arg0 ...any) *MockSpanEndCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "End", reflect.TypeOf((*MockSpan)(nil).End), arg0...)
	return &MockSpanEndCall{Call: call}
}

// MockSpanEndCall wrap *gomock.Call
type MockSpanEndCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSpanEndCall) Return() *MockSpanEndCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSpanEndCall) Do(f func(...trace.SpanEndOption)) *MockSpanEndCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSpanEndCall) DoAndReturn(f func(...trace.SpanEndOption)) *MockSpanEndCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsRecording mocks base method.
func (m *MockSpan) IsRecording() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRecording")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRecording indicates an expected call of IsRecording.
func (mr *MockSpanMockRecorder) IsRecording() *MockSpanIsRecordingCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRecording", reflect.TypeOf((*MockSpan)(nil).IsRecording))
	return &MockSpanIsRecordingCall{Call: call}
}

// MockSpanIsRecordingCall wrap *gomock.Call
type MockSpanIsRecordingCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSpanIsRecordingCall) Return(arg0 bool) *MockSpanIsRecordingCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSpanIsRecordingCall) Do(f func() bool) *MockSpanIsRecordingCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSpanIsRecordingCall) DoAndReturn(f func() bool) *MockSpanIsRecordingCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RecordError mocks base method.
func (m *MockSpan) RecordError(arg0 error, arg1 ...trace.EventOption) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "RecordError", varargs...)
}

// RecordError indicates an expected call of RecordError.
func (mr *MockSpanMockRecorder) RecordError(arg0 any, arg1 ...any) *MockSpanRecordErrorCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordError", reflect.TypeOf((*MockSpan)(nil).RecordError), varargs...)
	return &MockSpanRecordErrorCall{Call: call}
}

// MockSpanRecordErrorCall wrap *gomock.Call
type MockSpanRecordErrorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSpanRecordErrorCall) Return() *MockSpanRecordErrorCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSpanRecordErrorCall) Do(f func(error, ...trace.EventOption)) *MockSpanRecordErrorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSpanRecordErrorCall) DoAndReturn(f func(error, ...trace.EventOption)) *MockSpanRecordErrorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetAttributes mocks base method.
func (m *MockSpan) SetAttributes(arg0 ...attribute.KeyValue) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "SetAttributes", varargs...)
}

// SetAttributes indicates an expected call of SetAttributes.
func (mr *MockSpanMockRecorder) SetAttributes(arg0 ...any) *MockSpanSetAttributesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAttributes", reflect.TypeOf((*MockSpan)(nil).SetAttributes), arg0...)
	return &MockSpanSetAttributesCall{Call: call}
}

// MockSpanSetAttributesCall wrap *gomock.Call
type MockSpanSetAttributesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSpanSetAttributesCall) Return() *MockSpanSetAttributesCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSpanSetAttributesCall) Do(f func(...attribute.KeyValue)) *MockSpanSetAttributesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSpanSetAttributesCall) DoAndReturn(f func(...attribute.KeyValue)) *MockSpanSetAttributesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetName mocks base method.
func (m *MockSpan) SetName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetName", arg0)
}

// SetName indicates an expected call of SetName.
func (mr *MockSpanMockRecorder) SetName(arg0 any) *MockSpanSetNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetName", reflect.TypeOf((*MockSpan)(nil).SetName), arg0)
	return &MockSpanSetNameCall{Call: call}
}

// MockSpanSetNameCall wrap *gomock.Call
type MockSpanSetNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSpanSetNameCall) Return() *MockSpanSetNameCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSpanSetNameCall) Do(f func(string)) *MockSpanSetNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSpanSetNameCall) DoAndReturn(f func(string)) *MockSpanSetNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetStatus mocks base method.
func (m *MockSpan) SetStatus(arg0 codes.Code, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStatus", arg0, arg1)
}

// SetStatus indicates an expected call of SetStatus.
func (mr *MockSpanMockRecorder) SetStatus(arg0, arg1 any) *MockSpanSetStatusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatus", reflect.TypeOf((*MockSpan)(nil).SetStatus), arg0, arg1)
	return &MockSpanSetStatusCall{Call: call}
}

// MockSpanSetStatusCall wrap *gomock.Call
type MockSpanSetStatusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSpanSetStatusCall) Return() *MockSpanSetStatusCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSpanSetStatusCall) Do(f func(codes.Code, string)) *MockSpanSetStatusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSpanSetStatusCall) DoAndReturn(f func(codes.Code, string)) *MockSpanSetStatusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SpanContext mocks base method.
func (m *MockSpan) SpanContext() trace.SpanContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpanContext")
	ret0, _ := ret[0].(trace.SpanContext)
	return ret0
}

// SpanContext indicates an expected call of SpanContext.
func (mr *MockSpanMockRecorder) SpanContext() *MockSpanSpanContextCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpanContext", reflect.TypeOf((*MockSpan)(nil).SpanContext))
	return &MockSpanSpanContextCall{Call: call}
}

// MockSpanSpanContextCall wrap *gomock.Call
type MockSpanSpanContextCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSpanSpanContextCall) Return(arg0 trace.SpanContext) *MockSpanSpanContextCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSpanSpanContextCall) Do(f func() trace.SpanContext) *MockSpanSpanContextCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSpanSpanContextCall) DoAndReturn(f func() trace.SpanContext) *MockSpanSpanContextCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// TracerProvider mocks base method.
func (m *MockSpan) TracerProvider() trace.TracerProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TracerProvider")
	ret0, _ := ret[0].(trace.TracerProvider)
	return ret0
}

// TracerProvider indicates an expected call of TracerProvider.
func (mr *MockSpanMockRecorder) TracerProvider() *MockSpanTracerProviderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TracerProvider", reflect.TypeOf((*MockSpan)(nil).TracerProvider))
	return &MockSpanTracerProviderCall{Call: call}
}

// MockSpanTracerProviderCall wrap *gomock.Call
type MockSpanTracerProviderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSpanTracerProviderCall) Return(arg0 trace.TracerProvider) *MockSpanTracerProviderCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSpanTracerProviderCall) Do(f func() trace.TracerProvider) *MockSpanTracerProviderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSpanTracerProviderCall) DoAndReturn(f func() trace.TracerProvider) *MockSpanTracerProviderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// span mocks base method.
func (m *MockSpan) span() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "span")
}

// span indicates an expected call of span.
func (mr *MockSpanMockRecorder) span() *MockSpanspanCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "span", reflect.TypeOf((*MockSpan)(nil).span))
	return &MockSpanspanCall{Call: call}
}

// MockSpanspanCall wrap *gomock.Call
type MockSpanspanCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSpanspanCall) Return() *MockSpanspanCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSpanspanCall) Do(f func()) *MockSpanspanCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSpanspanCall) DoAndReturn(f func()) *MockSpanspanCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
