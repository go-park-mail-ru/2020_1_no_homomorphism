package album

import "no_homomorphism/internal/pkg/models"

type Repository interface {
	GetUserAlbums(uId uint64) ([]models.AlbumWithArtist, error)
	GetAlbumById(aId uint64) (models.Album, error)
}
