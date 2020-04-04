package artist

import "no_homomorphism/internal/pkg/models"

type Repository interface {
	GetArtist(id string) (models.Artist, error)
	GetBoundedArtists(start, end uint64) ([]models.Artist, error)
}