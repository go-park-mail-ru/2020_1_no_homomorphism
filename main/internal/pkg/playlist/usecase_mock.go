// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package playlist is a generated GoMock package.
package playlist

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

// GetUserPlaylists mocks base method
func (m *MockUseCase) GetUserPlaylists(id string) ([]models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPlaylists", id)
	ret0, _ := ret[0].([]models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPlaylists indicates an expected call of GetUserPlaylists
func (mr *MockUseCaseMockRecorder) GetUserPlaylists(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPlaylists", reflect.TypeOf((*MockUseCase)(nil).GetUserPlaylists), id)
}

// GetPlaylistById mocks base method
func (m *MockUseCase) GetPlaylistById(id string) (models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlaylistById", id)
	ret0, _ := ret[0].(models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlaylistById indicates an expected call of GetPlaylistById
func (mr *MockUseCaseMockRecorder) GetPlaylistById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlaylistById", reflect.TypeOf((*MockUseCase)(nil).GetPlaylistById), id)
}

// CreatePlaylist mocks base method
func (m *MockUseCase) CreatePlaylist(name, uID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePlaylist", name, uID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePlaylist indicates an expected call of CreatePlaylist
func (mr *MockUseCaseMockRecorder) CreatePlaylist(name, uID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePlaylist", reflect.TypeOf((*MockUseCase)(nil).CreatePlaylist), name, uID)
}

// CheckAccessToPlaylist mocks base method
func (m *MockUseCase) CheckAccessToPlaylist(userId, playlistId string, isStrict bool) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAccessToPlaylist", userId, playlistId, isStrict)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAccessToPlaylist indicates an expected call of CheckAccessToPlaylist
func (mr *MockUseCaseMockRecorder) CheckAccessToPlaylist(userId, playlistId, isStrict interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAccessToPlaylist", reflect.TypeOf((*MockUseCase)(nil).CheckAccessToPlaylist), userId, playlistId, isStrict)
}

// AddTrackToPlaylist mocks base method
func (m *MockUseCase) AddTrackToPlaylist(plTracks models.PlaylistTracks) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTrackToPlaylist", plTracks)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTrackToPlaylist indicates an expected call of AddTrackToPlaylist
func (mr *MockUseCaseMockRecorder) AddTrackToPlaylist(plTracks interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTrackToPlaylist", reflect.TypeOf((*MockUseCase)(nil).AddTrackToPlaylist), plTracks)
}

// GetUserPlaylistsIdByTrack mocks base method
func (m *MockUseCase) GetUserPlaylistsIdByTrack(userID, trackID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPlaylistsIdByTrack", userID, trackID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPlaylistsIdByTrack indicates an expected call of GetUserPlaylistsIdByTrack
func (mr *MockUseCaseMockRecorder) GetUserPlaylistsIdByTrack(userID, trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPlaylistsIdByTrack", reflect.TypeOf((*MockUseCase)(nil).GetUserPlaylistsIdByTrack), userID, trackID)
}

// DeleteTrackFromPlaylist mocks base method
func (m *MockUseCase) DeleteTrackFromPlaylist(plID, trackID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrackFromPlaylist", plID, trackID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTrackFromPlaylist indicates an expected call of DeleteTrackFromPlaylist
func (mr *MockUseCaseMockRecorder) DeleteTrackFromPlaylist(plID, trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrackFromPlaylist", reflect.TypeOf((*MockUseCase)(nil).DeleteTrackFromPlaylist), plID, trackID)
}

// DeletePlaylist mocks base method
func (m *MockUseCase) DeletePlaylist(plID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlaylist", plID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlaylist indicates an expected call of DeletePlaylist
func (mr *MockUseCaseMockRecorder) DeletePlaylist(plID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlaylist", reflect.TypeOf((*MockUseCase)(nil).DeletePlaylist), plID)
}

// ChangePrivacy mocks base method
func (m *MockUseCase) ChangePrivacy(plID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePrivacy", plID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePrivacy indicates an expected call of ChangePrivacy
func (mr *MockUseCaseMockRecorder) ChangePrivacy(plID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePrivacy", reflect.TypeOf((*MockUseCase)(nil).ChangePrivacy), plID)
}

// AddSharedPlaylist mocks base method
func (m *MockUseCase) AddSharedPlaylist(plID, uID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSharedPlaylist", plID, uID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSharedPlaylist indicates an expected call of AddSharedPlaylist
func (mr *MockUseCaseMockRecorder) AddSharedPlaylist(plID, uID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSharedPlaylist", reflect.TypeOf((*MockUseCase)(nil).AddSharedPlaylist), plID, uID)
}
