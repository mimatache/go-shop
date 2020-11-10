// Code generated by MockGen. DO NOT EDIT.
// Source: ./store.go

// Package mock_store is a generated GoMock package.
package mock_store

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// Mocklogger is a mock of logger interface
type Mocklogger struct {
	ctrl     *gomock.Controller
	recorder *MockloggerMockRecorder
}

// MockloggerMockRecorder is the mock recorder for Mocklogger
type MockloggerMockRecorder struct {
	mock *Mocklogger
}

// NewMocklogger creates a new mock instance
func NewMocklogger(ctrl *gomock.Controller) *Mocklogger {
	mock := &Mocklogger{ctrl: ctrl}
	mock.recorder = &MockloggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mocklogger) EXPECT() *MockloggerMockRecorder {
	return m.recorder
}

// Infof mocks base method
func (m *Mocklogger) Infof(msg string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{msg}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof
func (mr *MockloggerMockRecorder) Infof(msg interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{msg}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*Mocklogger)(nil).Infof), varargs...)
}

// Debugf mocks base method
func (m *Mocklogger) Debugf(msg string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{msg}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debugf", varargs...)
}

// Debugf indicates an expected call of Debugf
func (mr *MockloggerMockRecorder) Debugf(msg interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{msg}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debugf", reflect.TypeOf((*Mocklogger)(nil).Debugf), varargs...)
}

// MockUnderlyingStore is a mock of UnderlyingStore interface
type MockUnderlyingStore struct {
	ctrl     *gomock.Controller
	recorder *MockUnderlyingStoreMockRecorder
}

// MockUnderlyingStoreMockRecorder is the mock recorder for MockUnderlyingStore
type MockUnderlyingStoreMockRecorder struct {
	mock *MockUnderlyingStore
}

// NewMockUnderlyingStore creates a new mock instance
func NewMockUnderlyingStore(ctrl *gomock.Controller) *MockUnderlyingStore {
	mock := &MockUnderlyingStore{ctrl: ctrl}
	mock.recorder = &MockUnderlyingStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUnderlyingStore) EXPECT() *MockUnderlyingStoreMockRecorder {
	return m.recorder
}

// Read mocks base method
func (m *MockUnderlyingStore) Read(table, key string, value interface{}) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", table, key, value)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockUnderlyingStoreMockRecorder) Read(table, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockUnderlyingStore)(nil).Read), table, key, value)
}

// Write mocks base method
func (m *MockUnderlyingStore) Write(table string, value ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{table}
	for _, a := range value {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Write", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write
func (mr *MockUnderlyingStoreMockRecorder) Write(table interface{}, value ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{table}, value...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockUnderlyingStore)(nil).Write), varargs...)
}

// MockUserStore is a mock of UserStore interface
type MockUserStore struct {
	ctrl     *gomock.Controller
	recorder *MockUserStoreMockRecorder
}

// MockUserStoreMockRecorder is the mock recorder for MockUserStore
type MockUserStoreMockRecorder struct {
	mock *MockUserStore
}

// NewMockUserStore creates a new mock instance
func NewMockUserStore(ctrl *gomock.Controller) *MockUserStore {
	mock := &MockUserStore{ctrl: ctrl}
	mock.recorder = &MockUserStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserStore) EXPECT() *MockUserStoreMockRecorder {
	return m.recorder
}

// GetPasswordFor mocks base method
func (m *MockUserStore) GetPasswordFor(name string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPasswordFor", name)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPasswordFor indicates an expected call of GetPasswordFor
func (mr *MockUserStoreMockRecorder) GetPasswordFor(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPasswordFor", reflect.TypeOf((*MockUserStore)(nil).GetPasswordFor), name)
}
