// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	define "github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	gomock "github.com/golang/mock/gomock"
)

// MockBkApiClient is a mock of BkApiClient interface.
type MockBkApiClient struct {
	ctrl     *gomock.Controller
	recorder *MockBkApiClientMockRecorder
}

// MockBkApiClientMockRecorder is the mock recorder for MockBkApiClient.
type MockBkApiClientMockRecorder struct {
	mock *MockBkApiClient
}

// NewMockBkApiClient creates a new mock instance.
func NewMockBkApiClient(ctrl *gomock.Controller) *MockBkApiClient {
	mock := &MockBkApiClient{ctrl: ctrl}
	mock.recorder = &MockBkApiClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBkApiClient) EXPECT() *MockBkApiClientMockRecorder {
	return m.recorder
}

// NewOperation mocks base method.
func (m *MockBkApiClient) NewOperation(config define.OperationConfig, opts ...define.OperationOption) define.Operation {
	m.ctrl.T.Helper()
	varargs := []interface{}{config}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewOperation", varargs...)
	ret0, _ := ret[0].(define.Operation)
	return ret0
}

// NewOperation indicates an expected call of NewOperation.
func (mr *MockBkApiClientMockRecorder) NewOperation(config interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{config}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewOperation", reflect.TypeOf((*MockBkApiClient)(nil).NewOperation), varargs...)
}

// MockBkApiClientOption is a mock of BkApiClientOption interface.
type MockBkApiClientOption struct {
	ctrl     *gomock.Controller
	recorder *MockBkApiClientOptionMockRecorder
}

// MockBkApiClientOptionMockRecorder is the mock recorder for MockBkApiClientOption.
type MockBkApiClientOptionMockRecorder struct {
	mock *MockBkApiClientOption
}

// NewMockBkApiClientOption creates a new mock instance.
func NewMockBkApiClientOption(ctrl *gomock.Controller) *MockBkApiClientOption {
	mock := &MockBkApiClientOption{ctrl: ctrl}
	mock.recorder = &MockBkApiClientOptionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBkApiClientOption) EXPECT() *MockBkApiClientOptionMockRecorder {
	return m.recorder
}

// ApplyTo mocks base method.
func (m *MockBkApiClientOption) ApplyTo(arg0 define.BkApiClient) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ApplyTo", arg0)
}

// ApplyTo indicates an expected call of ApplyTo.
func (mr *MockBkApiClientOptionMockRecorder) ApplyTo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyTo", reflect.TypeOf((*MockBkApiClientOption)(nil).ApplyTo), arg0)
}
