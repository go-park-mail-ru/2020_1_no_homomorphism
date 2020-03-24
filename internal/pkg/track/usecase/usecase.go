package usecase

import (
	"fmt"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/track"
	"strconv"
)

type TrackUseCase struct {
	Repository track.Repository
}

func (uc TrackUseCase) GetTrackById(id string) (models.Track, error) {
	tId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return models.Track{}, fmt.Errorf("failed to convert id: %e", err)
	}
	return uc.Repository.GetTrackById(tId)
}
