package usecase

import (
	"fmt"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/playlist"
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

func (uc PlaylistUseCase) CheckAccessToPlaylist(userId string, playlistId string) (bool, error) {
	pl, err := uc.PlRepository.GetPlaylistById(playlistId)
	if err != nil {
		return false, fmt.Errorf("cant get playlist: %v", err)
	}
	return pl.UserId == userId, nil
}
