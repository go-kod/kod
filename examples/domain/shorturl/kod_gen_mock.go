// Code generated by MockGen. DO NOT EDIT.
// Source: examples/domain/shorturl/kod_gen_interface.go
//
// Generated by this command:
//
//	mockgen -source examples/domain/shorturl/kod_gen_interface.go -destination examples/domain/shorturl/kod_gen_mock.go -package shorturl
//

// Package shorturl is a generated GoMock package.
package shorturl

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockComponent is a mock of Component interface.
type MockComponent struct {
	ctrl     *gomock.Controller
	recorder *MockComponentMockRecorder
}

// MockComponentMockRecorder is the mock recorder for MockComponent.
type MockComponentMockRecorder struct {
	mock *MockComponent
}

// NewMockComponent creates a new mock instance.
func NewMockComponent(ctrl *gomock.Controller) *MockComponent {
	mock := &MockComponent{ctrl: ctrl}
	mock.recorder = &MockComponentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComponent) EXPECT() *MockComponentMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockComponent) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", ctx, req)
	ret0, _ := ret[0].(*GenerateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockComponentMockRecorder) Generate(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockComponent)(nil).Generate), ctx, req)
}

// Get mocks base method.
func (m *MockComponent) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, req)
	ret0, _ := ret[0].(*GetResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockComponentMockRecorder) Get(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockComponent)(nil).Get), ctx, req)
}