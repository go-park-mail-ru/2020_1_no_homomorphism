package artist

import "no_homomorphism/internal/pkg/models"

type UseCase interface {
	GetArtistById(id string) (models.Artist, error)
	GetBoundedArtists(start, end uint64) ([]models.Artist, error)
	GetArtistStat(id string) (models.ArtistStat, error)
}
