package usecase

import (
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
