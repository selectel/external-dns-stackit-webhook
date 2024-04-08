// Code generated by MockGen. DO NOT EDIT.
// Source: ./http_middleware.go
//
// Generated by this command:
//
//	mockgen -destination=./mock/http_middleware.go -source=./http_middleware.go HttpApiMetrics
//

// Package mock_metrics is a generated GoMock package.
package mock_metrics

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockHttpApiMetrics is a mock of HttpApiMetrics interface.
type MockHttpApiMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockHttpApiMetricsMockRecorder
}

// MockHttpApiMetricsMockRecorder is the mock recorder for MockHttpApiMetrics.
type MockHttpApiMetricsMockRecorder struct {
	mock *MockHttpApiMetrics
}

// NewMockHttpApiMetrics creates a new mock instance.
func NewMockHttpApiMetrics(ctrl *gomock.Controller) *MockHttpApiMetrics {
	mock := &MockHttpApiMetrics{ctrl: ctrl}
	mock.recorder = &MockHttpApiMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpApiMetrics) EXPECT() *MockHttpApiMetricsMockRecorder {
	return m.recorder
}

// Collect400TotalRequests mocks base method.
func (m *MockHttpApiMetrics) Collect400TotalRequests() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Collect400TotalRequests")
}

// Collect400TotalRequests indicates an expected call of Collect400TotalRequests.
func (mr *MockHttpApiMetricsMockRecorder) Collect400TotalRequests() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collect400TotalRequests", reflect.TypeOf((*MockHttpApiMetrics)(nil).Collect400TotalRequests))
}

// Collect500TotalRequests mocks base method.
func (m *MockHttpApiMetrics) Collect500TotalRequests() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Collect500TotalRequests")
}

// Collect500TotalRequests indicates an expected call of Collect500TotalRequests.
func (mr *MockHttpApiMetricsMockRecorder) Collect500TotalRequests() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collect500TotalRequests", reflect.TypeOf((*MockHttpApiMetrics)(nil).Collect500TotalRequests))
}

// CollectRequest mocks base method.
func (m *MockHttpApiMetrics) CollectRequest(method, path string, statusCode int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CollectRequest", method, path, statusCode)
}

// CollectRequest indicates an expected call of CollectRequest.
func (mr *MockHttpApiMetricsMockRecorder) CollectRequest(method, path, statusCode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectRequest", reflect.TypeOf((*MockHttpApiMetrics)(nil).CollectRequest), method, path, statusCode)
}

// CollectRequestContentLength mocks base method.
func (m *MockHttpApiMetrics) CollectRequestContentLength(method, path string, contentLength float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CollectRequestContentLength", method, path, contentLength)
}

// CollectRequestContentLength indicates an expected call of CollectRequestContentLength.
func (mr *MockHttpApiMetricsMockRecorder) CollectRequestContentLength(method, path, contentLength any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectRequestContentLength", reflect.TypeOf((*MockHttpApiMetrics)(nil).CollectRequestContentLength), method, path, contentLength)
}

// CollectRequestDuration mocks base method.
func (m *MockHttpApiMetrics) CollectRequestDuration(method, path string, duration float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CollectRequestDuration", method, path, duration)
}

// CollectRequestDuration indicates an expected call of CollectRequestDuration.
func (mr *MockHttpApiMetricsMockRecorder) CollectRequestDuration(method, path, duration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectRequestDuration", reflect.TypeOf((*MockHttpApiMetrics)(nil).CollectRequestDuration), method, path, duration)
}

// CollectRequestResponseSize mocks base method.
func (m *MockHttpApiMetrics) CollectRequestResponseSize(method, path string, contentLength float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CollectRequestResponseSize", method, path, contentLength)
}

// CollectRequestResponseSize indicates an expected call of CollectRequestResponseSize.
func (mr *MockHttpApiMetricsMockRecorder) CollectRequestResponseSize(method, path, contentLength any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectRequestResponseSize", reflect.TypeOf((*MockHttpApiMetrics)(nil).CollectRequestResponseSize), method, path, contentLength)
}

// CollectTotalRequests mocks base method.
func (m *MockHttpApiMetrics) CollectTotalRequests() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CollectTotalRequests")
}

// CollectTotalRequests indicates an expected call of CollectTotalRequests.
func (mr *MockHttpApiMetricsMockRecorder) CollectTotalRequests() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectTotalRequests", reflect.TypeOf((*MockHttpApiMetrics)(nil).CollectTotalRequests))
}
