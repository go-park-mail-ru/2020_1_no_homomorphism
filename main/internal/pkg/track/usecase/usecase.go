package usecase

import (
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
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
