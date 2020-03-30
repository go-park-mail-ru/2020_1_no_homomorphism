// Code generated by MockGen. DO NOT EDIT.
// Source: delivery.go

// Package session is a generated GoMock package.
package session

import (
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	http "net/http"
	models "no_homomorphism/internal/pkg/models"
	reflect "reflect"
)

// MockDelivery is a mock of Delivery interface
type MockDelivery struct {
	ctrl     *gomock.Controller
	recorder *MockDeliveryMockRecorder
}

// MockDeliveryMockRecorder is the mock recorder for MockDelivery
type MockDeliveryMockRecorder struct {
	mock *MockDelivery
}

// NewMockDelivery creates a new mock instance
func NewMockDelivery(ctrl *gomock.Controller) *MockDelivery {
	mock := &MockDelivery{ctrl: ctrl}
	mock.recorder = &MockDeliveryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDelivery) EXPECT() *MockDeliveryMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockDelivery) Create(user models.User) (http.Cookie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", user)
	ret0, _ := ret[0].(http.Cookie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockDeliveryMockRecorder) Create(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDelivery)(nil).Create), user)
}

// Delete mocks base method
func (m *MockDelivery) Delete(sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockDeliveryMockRecorder) Delete(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDelivery)(nil).Delete), sessionID)
}

// GetLoginBySessionID mocks base method
func (m *MockDelivery) GetLoginBySessionID(sessionID uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoginBySessionID", sessionID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoginBySessionID indicates an expected call of GetLoginBySessionID
func (mr *MockDeliveryMockRecorder) GetLoginBySessionID(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoginBySessionID", reflect.TypeOf((*MockDelivery)(nil).GetLoginBySessionID), sessionID)
}