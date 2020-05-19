// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package album is a generated GoMock package.
package album

import (
	models "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
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

// GetUserAlbums mocks base method
func (m *MockUseCase) GetUserAlbums(id string) ([]models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserAlbums", id)
	ret0, _ := ret[0].([]models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserAlbums indicates an expected call of GetUserAlbums
func (mr *MockUseCaseMockRecorder) GetUserAlbums(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserAlbums", reflect.TypeOf((*MockUseCase)(nil).GetUserAlbums), id)
}

// GetAlbumById mocks base method
func (m *MockUseCase) GetAlbumById(aID, uID string) (models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlbumById", aID, uID)
	ret0, _ := ret[0].(models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlbumById indicates an expected call of GetAlbumById
func (mr *MockUseCaseMockRecorder) GetAlbumById(aID, uID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlbumById", reflect.TypeOf((*MockUseCase)(nil).GetAlbumById), aID, uID)
}

// GetBoundedAlbumsByArtistId mocks base method
func (m *MockUseCase) GetBoundedAlbumsByArtistId(id string, start, end uint64) ([]models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoundedAlbumsByArtistId", id, start, end)
	ret0, _ := ret[0].([]models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBoundedAlbumsByArtistId indicates an expected call of GetBoundedAlbumsByArtistId
func (mr *MockUseCaseMockRecorder) GetBoundedAlbumsByArtistId(id, start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoundedAlbumsByArtistId", reflect.TypeOf((*MockUseCase)(nil).GetBoundedAlbumsByArtistId), id, start, end)
}

// Search mocks base method
func (m *MockUseCase) Search(text string, count uint) ([]models.AlbumSearch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", text, count)
	ret0, _ := ret[0].([]models.AlbumSearch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockUseCaseMockRecorder) Search(text, count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockUseCase)(nil).Search), text, count)
}

// RateAlbum mocks base method
func (m *MockUseCase) RateAlbum(aID, uID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RateAlbum", aID, uID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RateAlbum indicates an expected call of RateAlbum
func (mr *MockUseCaseMockRecorder) RateAlbum(aID, uID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RateAlbum", reflect.TypeOf((*MockUseCase)(nil).RateAlbum), aID, uID)
}
