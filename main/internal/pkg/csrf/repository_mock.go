// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package csrf is a generated GoMock package.
package csrf

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

// Add mocks base method
func (m *MockRepository) Add(token string, expire int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", token, expire)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockRepositoryMockRecorder) Add(token, expire interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockRepository)(nil).Add), token, expire)
}

// Check mocks base method
func (m *MockRepository) Check(token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check
func (mr *MockRepositoryMockRecorder) Check(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockRepository)(nil).Check), token)
}
