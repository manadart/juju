// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/prometheus/client_golang/prometheus (interfaces: Counter,Gauge)
//
// Generated by this command:
//
//	mockgen -typed -package mocks -destination mocks/metrics_mock.go github.com/prometheus/client_golang/prometheus Counter,Gauge
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	prometheus "github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	gomock "go.uber.org/mock/gomock"
)

// MockCounter is a mock of Counter interface.
type MockCounter struct {
	ctrl     *gomock.Controller
	recorder *MockCounterMockRecorder
}

// MockCounterMockRecorder is the mock recorder for MockCounter.
type MockCounterMockRecorder struct {
	mock *MockCounter
}

// NewMockCounter creates a new mock instance.
func NewMockCounter(ctrl *gomock.Controller) *MockCounter {
	mock := &MockCounter{ctrl: ctrl}
	mock.recorder = &MockCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCounter) EXPECT() *MockCounterMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockCounter) Add(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0)
}

// Add indicates an expected call of Add.
func (mr *MockCounterMockRecorder) Add(arg0 any) *MockCounterAddCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCounter)(nil).Add), arg0)
	return &MockCounterAddCall{Call: call}
}

// MockCounterAddCall wrap *gomock.Call
type MockCounterAddCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCounterAddCall) Return() *MockCounterAddCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCounterAddCall) Do(f func(float64)) *MockCounterAddCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCounterAddCall) DoAndReturn(f func(float64)) *MockCounterAddCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Collect mocks base method.
func (m *MockCounter) Collect(arg0 chan<- prometheus.Metric) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Collect", arg0)
}

// Collect indicates an expected call of Collect.
func (mr *MockCounterMockRecorder) Collect(arg0 any) *MockCounterCollectCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collect", reflect.TypeOf((*MockCounter)(nil).Collect), arg0)
	return &MockCounterCollectCall{Call: call}
}

// MockCounterCollectCall wrap *gomock.Call
type MockCounterCollectCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCounterCollectCall) Return() *MockCounterCollectCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCounterCollectCall) Do(f func(chan<- prometheus.Metric)) *MockCounterCollectCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCounterCollectCall) DoAndReturn(f func(chan<- prometheus.Metric)) *MockCounterCollectCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Desc mocks base method.
func (m *MockCounter) Desc() *prometheus.Desc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Desc")
	ret0, _ := ret[0].(*prometheus.Desc)
	return ret0
}

// Desc indicates an expected call of Desc.
func (mr *MockCounterMockRecorder) Desc() *MockCounterDescCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Desc", reflect.TypeOf((*MockCounter)(nil).Desc))
	return &MockCounterDescCall{Call: call}
}

// MockCounterDescCall wrap *gomock.Call
type MockCounterDescCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCounterDescCall) Return(arg0 *prometheus.Desc) *MockCounterDescCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCounterDescCall) Do(f func() *prometheus.Desc) *MockCounterDescCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCounterDescCall) DoAndReturn(f func() *prometheus.Desc) *MockCounterDescCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Describe mocks base method.
func (m *MockCounter) Describe(arg0 chan<- *prometheus.Desc) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Describe", arg0)
}

// Describe indicates an expected call of Describe.
func (mr *MockCounterMockRecorder) Describe(arg0 any) *MockCounterDescribeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Describe", reflect.TypeOf((*MockCounter)(nil).Describe), arg0)
	return &MockCounterDescribeCall{Call: call}
}

// MockCounterDescribeCall wrap *gomock.Call
type MockCounterDescribeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCounterDescribeCall) Return() *MockCounterDescribeCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCounterDescribeCall) Do(f func(chan<- *prometheus.Desc)) *MockCounterDescribeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCounterDescribeCall) DoAndReturn(f func(chan<- *prometheus.Desc)) *MockCounterDescribeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Inc mocks base method.
func (m *MockCounter) Inc() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Inc")
}

// Inc indicates an expected call of Inc.
func (mr *MockCounterMockRecorder) Inc() *MockCounterIncCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inc", reflect.TypeOf((*MockCounter)(nil).Inc))
	return &MockCounterIncCall{Call: call}
}

// MockCounterIncCall wrap *gomock.Call
type MockCounterIncCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCounterIncCall) Return() *MockCounterIncCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCounterIncCall) Do(f func()) *MockCounterIncCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCounterIncCall) DoAndReturn(f func()) *MockCounterIncCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Write mocks base method.
func (m *MockCounter) Write(arg0 *io_prometheus_client.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write.
func (mr *MockCounterMockRecorder) Write(arg0 any) *MockCounterWriteCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockCounter)(nil).Write), arg0)
	return &MockCounterWriteCall{Call: call}
}

// MockCounterWriteCall wrap *gomock.Call
type MockCounterWriteCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCounterWriteCall) Return(arg0 error) *MockCounterWriteCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCounterWriteCall) Do(f func(*io_prometheus_client.Metric) error) *MockCounterWriteCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCounterWriteCall) DoAndReturn(f func(*io_prometheus_client.Metric) error) *MockCounterWriteCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockGauge is a mock of Gauge interface.
type MockGauge struct {
	ctrl     *gomock.Controller
	recorder *MockGaugeMockRecorder
}

// MockGaugeMockRecorder is the mock recorder for MockGauge.
type MockGaugeMockRecorder struct {
	mock *MockGauge
}

// NewMockGauge creates a new mock instance.
func NewMockGauge(ctrl *gomock.Controller) *MockGauge {
	mock := &MockGauge{ctrl: ctrl}
	mock.recorder = &MockGaugeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGauge) EXPECT() *MockGaugeMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockGauge) Add(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0)
}

// Add indicates an expected call of Add.
func (mr *MockGaugeMockRecorder) Add(arg0 any) *MockGaugeAddCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockGauge)(nil).Add), arg0)
	return &MockGaugeAddCall{Call: call}
}

// MockGaugeAddCall wrap *gomock.Call
type MockGaugeAddCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGaugeAddCall) Return() *MockGaugeAddCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGaugeAddCall) Do(f func(float64)) *MockGaugeAddCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGaugeAddCall) DoAndReturn(f func(float64)) *MockGaugeAddCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Collect mocks base method.
func (m *MockGauge) Collect(arg0 chan<- prometheus.Metric) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Collect", arg0)
}

// Collect indicates an expected call of Collect.
func (mr *MockGaugeMockRecorder) Collect(arg0 any) *MockGaugeCollectCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collect", reflect.TypeOf((*MockGauge)(nil).Collect), arg0)
	return &MockGaugeCollectCall{Call: call}
}

// MockGaugeCollectCall wrap *gomock.Call
type MockGaugeCollectCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGaugeCollectCall) Return() *MockGaugeCollectCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGaugeCollectCall) Do(f func(chan<- prometheus.Metric)) *MockGaugeCollectCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGaugeCollectCall) DoAndReturn(f func(chan<- prometheus.Metric)) *MockGaugeCollectCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Dec mocks base method.
func (m *MockGauge) Dec() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dec")
}

// Dec indicates an expected call of Dec.
func (mr *MockGaugeMockRecorder) Dec() *MockGaugeDecCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dec", reflect.TypeOf((*MockGauge)(nil).Dec))
	return &MockGaugeDecCall{Call: call}
}

// MockGaugeDecCall wrap *gomock.Call
type MockGaugeDecCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGaugeDecCall) Return() *MockGaugeDecCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGaugeDecCall) Do(f func()) *MockGaugeDecCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGaugeDecCall) DoAndReturn(f func()) *MockGaugeDecCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Desc mocks base method.
func (m *MockGauge) Desc() *prometheus.Desc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Desc")
	ret0, _ := ret[0].(*prometheus.Desc)
	return ret0
}

// Desc indicates an expected call of Desc.
func (mr *MockGaugeMockRecorder) Desc() *MockGaugeDescCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Desc", reflect.TypeOf((*MockGauge)(nil).Desc))
	return &MockGaugeDescCall{Call: call}
}

// MockGaugeDescCall wrap *gomock.Call
type MockGaugeDescCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGaugeDescCall) Return(arg0 *prometheus.Desc) *MockGaugeDescCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGaugeDescCall) Do(f func() *prometheus.Desc) *MockGaugeDescCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGaugeDescCall) DoAndReturn(f func() *prometheus.Desc) *MockGaugeDescCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Describe mocks base method.
func (m *MockGauge) Describe(arg0 chan<- *prometheus.Desc) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Describe", arg0)
}

// Describe indicates an expected call of Describe.
func (mr *MockGaugeMockRecorder) Describe(arg0 any) *MockGaugeDescribeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Describe", reflect.TypeOf((*MockGauge)(nil).Describe), arg0)
	return &MockGaugeDescribeCall{Call: call}
}

// MockGaugeDescribeCall wrap *gomock.Call
type MockGaugeDescribeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGaugeDescribeCall) Return() *MockGaugeDescribeCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGaugeDescribeCall) Do(f func(chan<- *prometheus.Desc)) *MockGaugeDescribeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGaugeDescribeCall) DoAndReturn(f func(chan<- *prometheus.Desc)) *MockGaugeDescribeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Inc mocks base method.
func (m *MockGauge) Inc() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Inc")
}

// Inc indicates an expected call of Inc.
func (mr *MockGaugeMockRecorder) Inc() *MockGaugeIncCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inc", reflect.TypeOf((*MockGauge)(nil).Inc))
	return &MockGaugeIncCall{Call: call}
}

// MockGaugeIncCall wrap *gomock.Call
type MockGaugeIncCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGaugeIncCall) Return() *MockGaugeIncCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGaugeIncCall) Do(f func()) *MockGaugeIncCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGaugeIncCall) DoAndReturn(f func()) *MockGaugeIncCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Set mocks base method.
func (m *MockGauge) Set(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", arg0)
}

// Set indicates an expected call of Set.
func (mr *MockGaugeMockRecorder) Set(arg0 any) *MockGaugeSetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockGauge)(nil).Set), arg0)
	return &MockGaugeSetCall{Call: call}
}

// MockGaugeSetCall wrap *gomock.Call
type MockGaugeSetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGaugeSetCall) Return() *MockGaugeSetCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGaugeSetCall) Do(f func(float64)) *MockGaugeSetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGaugeSetCall) DoAndReturn(f func(float64)) *MockGaugeSetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetToCurrentTime mocks base method.
func (m *MockGauge) SetToCurrentTime() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetToCurrentTime")
}

// SetToCurrentTime indicates an expected call of SetToCurrentTime.
func (mr *MockGaugeMockRecorder) SetToCurrentTime() *MockGaugeSetToCurrentTimeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetToCurrentTime", reflect.TypeOf((*MockGauge)(nil).SetToCurrentTime))
	return &MockGaugeSetToCurrentTimeCall{Call: call}
}

// MockGaugeSetToCurrentTimeCall wrap *gomock.Call
type MockGaugeSetToCurrentTimeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGaugeSetToCurrentTimeCall) Return() *MockGaugeSetToCurrentTimeCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGaugeSetToCurrentTimeCall) Do(f func()) *MockGaugeSetToCurrentTimeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGaugeSetToCurrentTimeCall) DoAndReturn(f func()) *MockGaugeSetToCurrentTimeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Sub mocks base method.
func (m *MockGauge) Sub(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Sub", arg0)
}

// Sub indicates an expected call of Sub.
func (mr *MockGaugeMockRecorder) Sub(arg0 any) *MockGaugeSubCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sub", reflect.TypeOf((*MockGauge)(nil).Sub), arg0)
	return &MockGaugeSubCall{Call: call}
}

// MockGaugeSubCall wrap *gomock.Call
type MockGaugeSubCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGaugeSubCall) Return() *MockGaugeSubCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGaugeSubCall) Do(f func(float64)) *MockGaugeSubCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGaugeSubCall) DoAndReturn(f func(float64)) *MockGaugeSubCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Write mocks base method.
func (m *MockGauge) Write(arg0 *io_prometheus_client.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write.
func (mr *MockGaugeMockRecorder) Write(arg0 any) *MockGaugeWriteCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockGauge)(nil).Write), arg0)
	return &MockGaugeWriteCall{Call: call}
}

// MockGaugeWriteCall wrap *gomock.Call
type MockGaugeWriteCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGaugeWriteCall) Return(arg0 error) *MockGaugeWriteCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGaugeWriteCall) Do(f func(*io_prometheus_client.Metric) error) *MockGaugeWriteCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGaugeWriteCall) DoAndReturn(f func(*io_prometheus_client.Metric) error) *MockGaugeWriteCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
