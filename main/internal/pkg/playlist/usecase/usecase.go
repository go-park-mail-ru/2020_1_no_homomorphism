package usecase

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist"
)

type PlaylistUseCase struct {
	PlRepository playlist.Repository
}

func (uc PlaylistUseCase) GetUserPlaylists(id string) ([]models.Playlist, error) {
	return uc.PlRepository.GetUserPlaylists(id)
}

func (uc PlaylistUseCase) GetPlaylistById(id string) (models.Playlist, error) {
	return uc.PlRepository.GetPlaylistById(id)
}

func (uc PlaylistUseCase) CreatePlaylist(name string, uID string) (plID string, err error) {
	return uc.PlRepository.CreatePlaylist(name, uID)
}

func (uc PlaylistUseCase) AddTrackToPlaylist(plTracks models.PlaylistTracks) error {
	return uc.PlRepository.AddTrackToPlaylist(plTracks)
}

func (uc PlaylistUseCase) GetUserPlaylistsIdByTrack(userID, trackID string) ([]string, error) {
	return uc.PlRepository.GetUserPlaylistsIdByTrack(userID, trackID)
}

func (uc PlaylistUseCase) DeletePlaylist(plID string) error {
	return uc.PlRepository.DeletePlaylist(plID)
}

func (uc PlaylistUseCase) DeleteTrackFromPlaylist(plID, trackID string) error {
	return uc.PlRepository.DeleteTrackFromPlaylist(plID, trackID)
}

func (uc PlaylistUseCase) CheckAccessToPlaylist(userId string, playlistId string) (bool, error) {
	pl, err := uc.PlRepository.GetPlaylistById(playlistId)
	if err != nil {
		return false, fmt.Errorf("cant get playlist: %v", err)
	}
	return pl.UserId == userId, nil
}
