package album

import "no_homomorphism/internal/pkg/models"

type Repository interface {
	GetUserAlbums(uId string) ([]models.Album, error)
	GetAlbumById(aId string) (models.Album, error)
}
