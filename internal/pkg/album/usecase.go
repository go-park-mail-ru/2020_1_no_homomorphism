package album

import "no_homomorphism/internal/pkg/models"

type UseCase interface {
	GetUserAlbums(id string) ([]models.Album, error)
	GetAlbumById(id string) (models.Album, error)
}
