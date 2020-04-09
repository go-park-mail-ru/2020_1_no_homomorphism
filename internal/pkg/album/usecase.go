package album

import "no_homomorphism/internal/pkg/models"

type UseCase interface {
	GetUserAlbums(id string) ([]models.Album, error)
	GetAlbumById(id string) (models.Album, error)
	GetBoundedAlbumsByArtistId(id string, start uint64, end uint64) ([]models.Album, error)
}
