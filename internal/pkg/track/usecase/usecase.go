package usecase

import (
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/track"
)

type TrackUseCase struct {
	Repository track.Repository
}

func (uc *TrackUseCase) GetTrackById(id uint) (*models.Track, error) {
	return uc.Repository.GetTrackById(id)
}
