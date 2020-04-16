package album

import "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"

type UseCase interface {
	GetUserAlbums(id string) ([]models.Album, error)
	GetAlbumById(id string) (models.Album, error)
	GetBoundedAlbumsByArtistId(id string, start uint64, end uint64) ([]models.Album, error)
}
