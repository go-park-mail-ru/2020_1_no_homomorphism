package track

import "no_homomorphism/internal/pkg/models"

type Repository interface {
	GetTrackById(id uint) (*models.Track, error)
}
