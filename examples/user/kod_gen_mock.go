// Code generated by MockGen. DO NOT EDIT.
// Source: examples/user/kod_gen_interface.go
//
// Generated by this command:
//
//	mockgen -source examples/user/kod_gen_interface.go -destination examples/user/kod_gen_mock.go -package user
//

// Package user is a generated GoMock package.
package user

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

// Auth mocks base method.
func (m *MockComponent) Auth(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth", ctx, req)
	ret0, _ := ret[0].(*AuthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Auth indicates an expected call of Auth.
func (mr *MockComponentMockRecorder) Auth(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockComponent)(nil).Auth), ctx, req)
}

// DeRegister mocks base method.
func (m *MockComponent) DeRegister(ctx context.Context, req *DeRegisterRequest) (*DeRegisterResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeRegister", ctx, req)
	ret0, _ := ret[0].(*DeRegisterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeRegister indicates an expected call of DeRegister.
func (mr *MockComponentMockRecorder) DeRegister(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeRegister", reflect.TypeOf((*MockComponent)(nil).DeRegister), ctx, req)
}

// Login mocks base method.
func (m *MockComponent) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, req)
	ret0, _ := ret[0].(*LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockComponentMockRecorder) Login(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockComponent)(nil).Login), ctx, req)
}

// Register mocks base method.
func (m *MockComponent) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, req)
	ret0, _ := ret[0].(*RegisterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockComponentMockRecorder) Register(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockComponent)(nil).Register), ctx, req)
}
