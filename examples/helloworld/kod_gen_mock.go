//go:build !ignoreKodGen

// Code generated by MockGen. DO NOT EDIT.
// Source: examples/helloworld/kod_gen_interface.go
//
// Generated by this command:
//
//	mockgen -source examples/helloworld/kod_gen_interface.go -destination examples/helloworld/kod_gen_mock.go -package helloworld -typed -build_constraint !ignoreKodGen
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
	isgomock struct{}
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
func (mr *MockHelloWorldMockRecorder) SayHello(ctx any) *MockHelloWorldSayHelloCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SayHello", reflect.TypeOf((*MockHelloWorld)(nil).SayHello), ctx)
	return &MockHelloWorldSayHelloCall{Call: call}
}

// MockHelloWorldSayHelloCall wrap *gomock.Call
type MockHelloWorldSayHelloCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockHelloWorldSayHelloCall) Return() *MockHelloWorldSayHelloCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockHelloWorldSayHelloCall) Do(f func(context.Context)) *MockHelloWorldSayHelloCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockHelloWorldSayHelloCall) DoAndReturn(f func(context.Context)) *MockHelloWorldSayHelloCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockHelloWorldLazy is a mock of HelloWorldLazy interface.
type MockHelloWorldLazy struct {
	ctrl     *gomock.Controller
	recorder *MockHelloWorldLazyMockRecorder
	isgomock struct{}
}

// MockHelloWorldLazyMockRecorder is the mock recorder for MockHelloWorldLazy.
type MockHelloWorldLazyMockRecorder struct {
	mock *MockHelloWorldLazy
}

// NewMockHelloWorldLazy creates a new mock instance.
func NewMockHelloWorldLazy(ctrl *gomock.Controller) *MockHelloWorldLazy {
	mock := &MockHelloWorldLazy{ctrl: ctrl}
	mock.recorder = &MockHelloWorldLazyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHelloWorldLazy) EXPECT() *MockHelloWorldLazyMockRecorder {
	return m.recorder
}

// SayHello mocks base method.
func (m *MockHelloWorldLazy) SayHello(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SayHello", arg0)
}

// SayHello indicates an expected call of SayHello.
func (mr *MockHelloWorldLazyMockRecorder) SayHello(arg0 any) *MockHelloWorldLazySayHelloCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SayHello", reflect.TypeOf((*MockHelloWorldLazy)(nil).SayHello), arg0)
	return &MockHelloWorldLazySayHelloCall{Call: call}
}

// MockHelloWorldLazySayHelloCall wrap *gomock.Call
type MockHelloWorldLazySayHelloCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockHelloWorldLazySayHelloCall) Return() *MockHelloWorldLazySayHelloCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockHelloWorldLazySayHelloCall) Do(f func(context.Context)) *MockHelloWorldLazySayHelloCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockHelloWorldLazySayHelloCall) DoAndReturn(f func(context.Context)) *MockHelloWorldLazySayHelloCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockHelloWorldInterceptor is a mock of HelloWorldInterceptor interface.
type MockHelloWorldInterceptor struct {
	ctrl     *gomock.Controller
	recorder *MockHelloWorldInterceptorMockRecorder
	isgomock struct{}
}

// MockHelloWorldInterceptorMockRecorder is the mock recorder for MockHelloWorldInterceptor.
type MockHelloWorldInterceptorMockRecorder struct {
	mock *MockHelloWorldInterceptor
}

// NewMockHelloWorldInterceptor creates a new mock instance.
func NewMockHelloWorldInterceptor(ctrl *gomock.Controller) *MockHelloWorldInterceptor {
	mock := &MockHelloWorldInterceptor{ctrl: ctrl}
	mock.recorder = &MockHelloWorldInterceptorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHelloWorldInterceptor) EXPECT() *MockHelloWorldInterceptorMockRecorder {
	return m.recorder
}

// SayHello mocks base method.
func (m *MockHelloWorldInterceptor) SayHello(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SayHello", arg0)
}

// SayHello indicates an expected call of SayHello.
func (mr *MockHelloWorldInterceptorMockRecorder) SayHello(arg0 any) *MockHelloWorldInterceptorSayHelloCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SayHello", reflect.TypeOf((*MockHelloWorldInterceptor)(nil).SayHello), arg0)
	return &MockHelloWorldInterceptorSayHelloCall{Call: call}
}

// MockHelloWorldInterceptorSayHelloCall wrap *gomock.Call
type MockHelloWorldInterceptorSayHelloCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockHelloWorldInterceptorSayHelloCall) Return() *MockHelloWorldInterceptorSayHelloCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockHelloWorldInterceptorSayHelloCall) Do(f func(context.Context)) *MockHelloWorldInterceptorSayHelloCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockHelloWorldInterceptorSayHelloCall) DoAndReturn(f func(context.Context)) *MockHelloWorldInterceptorSayHelloCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
