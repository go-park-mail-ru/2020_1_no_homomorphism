package album

import "no_homomorphism/internal/pkg/models"

type UseCase interface {
	GetUserAlbums(id string) (models.UserAlbums, error)
	GetAlbumTracks(id string) (models.AlbumTracks, error)
}
