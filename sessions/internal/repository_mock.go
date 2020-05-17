// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package session is a generated GoMock package.
package session

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockRepository) Create(sId, value string, expire uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", sId, value, expire)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockRepositoryMockRecorder) Create(sId, value, expire interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), sId, value, expire)
}

// Delete mocks base method
func (m *MockRepository) Delete(sId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", sId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRepositoryMockRecorder) Delete(sId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), sId)
}

// GetLoginBySessionID mocks base method
func (m *MockRepository) GetLoginBySessionID(sId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoginBySessionID", sId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoginBySessionID indicates an expected call of GetLoginBySessionID
func (mr *MockRepositoryMockRecorder) GetLoginBySessionID(sId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoginBySessionID", reflect.TypeOf((*MockRepository)(nil).GetLoginBySessionID), sId)
}