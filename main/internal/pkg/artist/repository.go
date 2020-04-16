package artist

import "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"

type Repository interface {
	GetArtist(id string) (models.Artist, error)
	GetBoundedArtists(start, end uint64) ([]models.Artist, error)
	GetArtistStat(id string) (models.ArtistStat, error)
}
