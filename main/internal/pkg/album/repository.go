package album

import "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"

type Repository interface {
	GetUserAlbums(uId string) ([]models.Album, error)
	GetAlbumById(aId string) (models.Album, error)
	GetBoundedAlbumsByArtistId(id string, start, end uint64) ([]models.Album, error)
	Search(text string, count uint) ([]models.AlbumSearch, error)
	RateAlbum(aID, uID string) error
	CheckLike(aID, uID string) bool
	GetNewestReleases(uID string, begin int, end int) ([]models.NewestReleases, error)
}
