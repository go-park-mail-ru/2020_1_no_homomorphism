package artist

import "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"

type UseCase interface {
	GetArtistById(aID, uID string) (models.Artist, error)
	GetBoundedArtists(start, end uint64) ([]models.Artist, error)
	GetArtistStat(id string) (models.ArtistStat, error)
	Search(text string, count uint) ([]models.ArtistSearch, error)
	Subscription(aID, uID string) error
	SubscriptionList(uID string) ([]models.ArtistSearch, error)
}
