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

func (uc TrackUseCase) GetBoundedTracksByArtistId(id string, start, end uint64) ([]models.Track, error) {
	return uc.Repository.GetBoundedTracksByArtistId(id, start, end)
}

func (uc TrackUseCase) GetBoundedTracksByAlbumId(id string, start, end uint64) ([]models.Track, error) {
	return uc.Repository.GetBoundedTracksByAlbumId(id, start, end)
}

func (uc TrackUseCase) GetBoundedTracksByPlaylistId(plId string, start, end uint64) ([]models.Track, error) {
	return uc.Repository.GetBoundedTracksByPlaylistId(plId, start, end)
}
