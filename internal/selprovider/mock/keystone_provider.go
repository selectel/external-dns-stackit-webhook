// Code generated by MockGen. DO NOT EDIT.
// Source: ./config.go
//
// Generated by this command:
//
//	mockgen -destination=./mock/keystone_provider.go -source=./config.go KeystoneProvider
//

// Package mock_selprovider is a generated GoMock package.
package mock_selprovider

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockKeystoneProvider is a mock of KeystoneProvider interface.
type MockKeystoneProvider struct {
	ctrl     *gomock.Controller
	recorder *MockKeystoneProviderMockRecorder
	isgomock struct{}
}

// MockKeystoneProviderMockRecorder is the mock recorder for MockKeystoneProvider.
type MockKeystoneProviderMockRecorder struct {
	mock *MockKeystoneProvider
}

// NewMockKeystoneProvider creates a new mock instance.
func NewMockKeystoneProvider(ctrl *gomock.Controller) *MockKeystoneProvider {
	mock := &MockKeystoneProvider{ctrl: ctrl}
	mock.recorder = &MockKeystoneProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeystoneProvider) EXPECT() *MockKeystoneProviderMockRecorder {
	return m.recorder
}

// GetToken mocks base method.
func (m *MockKeystoneProvider) GetToken() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetToken")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetToken indicates an expected call of GetToken.
func (mr *MockKeystoneProviderMockRecorder) GetToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetToken", reflect.TypeOf((*MockKeystoneProvider)(nil).GetToken))
}
