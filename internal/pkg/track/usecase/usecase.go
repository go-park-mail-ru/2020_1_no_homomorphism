package usecase

import (
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/track"
)

type TrackUseCase struct {
	Repository track.Repository
}

func (uc TrackUseCase) GetTrackById(id string) (models.Track, error) {
	return uc.Repository.GetTrackById(id)
}

func (uc TrackUseCase) GetTracksByAlbumId(id string) ([]models.Track, error) {
	return uc.Repository.GetTracksByAlbumId(id)
}

func (uc TrackUseCase) GetTracksByPlaylistId(id string) ([]models.Track, error) {
	return uc.Repository.GetPlaylistTracks(id)
}
