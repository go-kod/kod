// Code generated by MockGen. DO NOT EDIT.
// Source: examples/helloworld/kod_gen_interface.go
//
// Generated by this command:
//
//	mockgen -source examples/helloworld/kod_gen_interface.go -destination examples/helloworld/kod_gen_mock.go -package helloworld
//

// Package helloworld is a generated GoMock package.
package helloworld

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockHelloWorld is a mock of HelloWorld interface.
type MockHelloWorld struct {
	ctrl     *gomock.Controller
	recorder *MockHelloWorldMockRecorder
}

// MockHelloWorldMockRecorder is the mock recorder for MockHelloWorld.
type MockHelloWorldMockRecorder struct {
	mock *MockHelloWorld
}

// NewMockHelloWorld creates a new mock instance.
func NewMockHelloWorld(ctrl *gomock.Controller) *MockHelloWorld {
	mock := &MockHelloWorld{ctrl: ctrl}
	mock.recorder = &MockHelloWorldMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHelloWorld) EXPECT() *MockHelloWorldMockRecorder {
	return m.recorder
}

// SayHello mocks base method.
func (m *MockHelloWorld) SayHello(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SayHello", ctx)
}

// SayHello indicates an expected call of SayHello.
func (mr *MockHelloWorldMockRecorder) SayHello(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SayHello", reflect.TypeOf((*MockHelloWorld)(nil).SayHello), ctx)
}

// MockHelloBob is a mock of HelloBob interface.
type MockHelloBob struct {
	ctrl     *gomock.Controller
	recorder *MockHelloBobMockRecorder
}

// MockHelloBobMockRecorder is the mock recorder for MockHelloBob.
type MockHelloBobMockRecorder struct {
	mock *MockHelloBob
}

// NewMockHelloBob creates a new mock instance.
func NewMockHelloBob(ctrl *gomock.Controller) *MockHelloBob {
	mock := &MockHelloBob{ctrl: ctrl}
	mock.recorder = &MockHelloBobMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHelloBob) EXPECT() *MockHelloBobMockRecorder {
	return m.recorder
}

// SayHello mocks base method.
func (m *MockHelloBob) SayHello(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SayHello", ctx)
}

// SayHello indicates an expected call of SayHello.
func (mr *MockHelloBobMockRecorder) SayHello(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SayHello", reflect.TypeOf((*MockHelloBob)(nil).SayHello), ctx)
}
