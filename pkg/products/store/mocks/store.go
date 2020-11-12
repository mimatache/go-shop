// Code generated by MockGen. DO NOT EDIT.
// Source: ./store.go

// Package mock_store is a generated GoMock package.
package mock_store

import (
	gomock "github.com/golang/mock/gomock"
	store "github.com/mimatache/go-shop/internal/store"
	store0 "github.com/mimatache/go-shop/pkg/products/store"
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

// Debugw mocks base method
func (m *Mocklogger) Debugw(msg string, keysAndValues ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{msg}
	for _, a := range keysAndValues {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debugw", varargs...)
}

// Debugw indicates an expected call of Debugw
func (mr *MockloggerMockRecorder) Debugw(msg interface{}, keysAndValues ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{msg}, keysAndValues...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debugw", reflect.TypeOf((*Mocklogger)(nil).Debugw), varargs...)
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

// WriteAndBlock mocks base method
func (m *MockUnderlyingStore) WriteAndBlock(table string, value ...interface{}) (store.Transaction, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{table}
	for _, a := range value {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WriteAndBlock", varargs...)
	ret0, _ := ret[0].(store.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteAndBlock indicates an expected call of WriteAndBlock
func (mr *MockUnderlyingStoreMockRecorder) WriteAndBlock(table interface{}, value ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{table}, value...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAndBlock", reflect.TypeOf((*MockUnderlyingStore)(nil).WriteAndBlock), varargs...)
}

// MockProductStore is a mock of ProductStore interface
type MockProductStore struct {
	ctrl     *gomock.Controller
	recorder *MockProductStoreMockRecorder
}

// MockProductStoreMockRecorder is the mock recorder for MockProductStore
type MockProductStoreMockRecorder struct {
	mock *MockProductStore
}

// NewMockProductStore creates a new mock instance
func NewMockProductStore(ctrl *gomock.Controller) *MockProductStore {
	mock := &MockProductStore{ctrl: ctrl}
	mock.recorder = &MockProductStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProductStore) EXPECT() *MockProductStoreMockRecorder {
	return m.recorder
}

// GetProductByID mocks base method
func (m *MockProductStore) GetProductByID(ID uint) (*store0.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductByID", ID)
	ret0, _ := ret[0].(*store0.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductByID indicates an expected call of GetProductByID
func (mr *MockProductStoreMockRecorder) GetProductByID(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductByID", reflect.TypeOf((*MockProductStore)(nil).GetProductByID), ID)
}

// SetProducts mocks base method
func (m *MockProductStore) SetProducts(products ...*store0.Product) (*store0.ProductTransaction, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range products {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SetProducts", varargs...)
	ret0, _ := ret[0].(*store0.ProductTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetProducts indicates an expected call of SetProducts
func (mr *MockProductStoreMockRecorder) SetProducts(products ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProducts", reflect.TypeOf((*MockProductStore)(nil).SetProducts), products...)
}
