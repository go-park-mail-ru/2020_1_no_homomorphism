package track

import "no_homomorphism/internal/pkg/models"

type UseCase interface {
	GetTrackById(id string) (*models.Track, error)
}
