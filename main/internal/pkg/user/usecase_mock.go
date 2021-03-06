// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package user is a generated GoMock package.
package user

import (
	models "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
)

// MockUseCase is a mock of UseCase interface
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockUseCase) Create(user models.User) (SameUserExists, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", user)
	ret0, _ := ret[0].(SameUserExists)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockUseCaseMockRecorder) Create(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUseCase)(nil).Create), user)
}

// Update mocks base method
func (m *MockUseCase) Update(user models.User, input models.UserSettings) (SameUserExists, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", user, input)
	ret0, _ := ret[0].(SameUserExists)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockUseCaseMockRecorder) Update(user, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUseCase)(nil).Update), user, input)
}

// Login mocks base method
func (m *MockUseCase) Login(input models.UserSignIn) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", input)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login
func (mr *MockUseCaseMockRecorder) Login(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUseCase)(nil).Login), input)
}

// UpdateAvatar mocks base method
func (m *MockUseCase) UpdateAvatar(user models.User, file io.Reader, fileType string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAvatar", user, file, fileType)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAvatar indicates an expected call of UpdateAvatar
func (mr *MockUseCaseMockRecorder) UpdateAvatar(user, file, fileType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAvatar", reflect.TypeOf((*MockUseCase)(nil).UpdateAvatar), user, file, fileType)
}

// GetUserByLogin mocks base method
func (m *MockUseCase) GetUserByLogin(login string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", login)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin
func (mr *MockUseCaseMockRecorder) GetUserByLogin(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockUseCase)(nil).GetUserByLogin), login)
}

// GetProfileByLogin mocks base method
func (m *MockUseCase) GetProfileByLogin(login string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileByLogin", login)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileByLogin indicates an expected call of GetProfileByLogin
func (mr *MockUseCaseMockRecorder) GetProfileByLogin(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileByLogin", reflect.TypeOf((*MockUseCase)(nil).GetProfileByLogin), login)
}

// GetOutputUserData mocks base method
func (m *MockUseCase) GetOutputUserData(user models.User) models.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOutputUserData", user)
	ret0, _ := ret[0].(models.User)
	return ret0
}

// GetOutputUserData indicates an expected call of GetOutputUserData
func (mr *MockUseCaseMockRecorder) GetOutputUserData(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOutputUserData", reflect.TypeOf((*MockUseCase)(nil).GetOutputUserData), user)
}

// CheckUserPassword mocks base method
func (m *MockUseCase) CheckUserPassword(userPassword, InputPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserPassword", userPassword, InputPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckUserPassword indicates an expected call of CheckUserPassword
func (mr *MockUseCaseMockRecorder) CheckUserPassword(userPassword, InputPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserPassword", reflect.TypeOf((*MockUseCase)(nil).CheckUserPassword), userPassword, InputPassword)
}

// GetUserStat mocks base method
func (m *MockUseCase) GetUserStat(id string) (models.UserStat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserStat", id)
	ret0, _ := ret[0].(models.UserStat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserStat indicates an expected call of GetUserStat
func (mr *MockUseCaseMockRecorder) GetUserStat(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserStat", reflect.TypeOf((*MockUseCase)(nil).GetUserStat), id)
}
