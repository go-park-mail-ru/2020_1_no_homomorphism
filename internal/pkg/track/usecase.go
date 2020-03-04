package track

import "no_homomorphism/internal/pkg/models"

type UseCase interface {
	GetTrackById(id uint) (*models.Track, error)
}
