// Code generated by MockGen. DO NOT EDIT.
// Source: watcher.go
//
// Generated by this command:
//
//	mockgen -source=watcher.go -destination=mock_watcher_test.go -package=internal
//
// Package internal is a generated GoMock package.
package internal

import (
	reflect "reflect"

	fsnotify "github.com/fsnotify/fsnotify"
	gomock "go.uber.org/mock/gomock"
)

// MockWatcher is a mock of Watcher interface.
type MockWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockWatcherMockRecorder
}

// MockWatcherMockRecorder is the mock recorder for MockWatcher.
type MockWatcherMockRecorder struct {
	mock *MockWatcher
}

// NewMockWatcher creates a new mock instance.
func NewMockWatcher(ctrl *gomock.Controller) *MockWatcher {
	mock := &MockWatcher{ctrl: ctrl}
	mock.recorder = &MockWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWatcher) EXPECT() *MockWatcherMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockWatcher) Add(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockWatcherMockRecorder) Add(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockWatcher)(nil).Add), arg0)
}

// Errors mocks base method.
func (m *MockWatcher) Errors() chan error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Errors")
	ret0, _ := ret[0].(chan error)
	return ret0
}

// Errors indicates an expected call of Errors.
func (mr *MockWatcherMockRecorder) Errors() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errors", reflect.TypeOf((*MockWatcher)(nil).Errors))
}

// Events mocks base method.
func (m *MockWatcher) Events() chan fsnotify.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Events")
	ret0, _ := ret[0].(chan fsnotify.Event)
	return ret0
}

// Events indicates an expected call of Events.
func (mr *MockWatcherMockRecorder) Events() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Events", reflect.TypeOf((*MockWatcher)(nil).Events))
}
