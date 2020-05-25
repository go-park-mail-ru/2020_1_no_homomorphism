package album

import "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"

type UseCase interface {
	GetUserAlbums(id string) ([]models.Album, error)
	GetAlbumById(aID, uID string) (models.Album, error)
	GetBoundedAlbumsByArtistId(id string, start uint64, end uint64) ([]models.Album, error)
	Search(text string, count uint) ([]models.AlbumSearch, error)
	RateAlbum(aID, uID string) error
	GetNewestReleases(uID string, begin, end int) ([]models.NewestReleases, error)
}
