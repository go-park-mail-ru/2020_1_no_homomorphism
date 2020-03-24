package playlist

import "no_homomorphism/internal/pkg/models"

type Repository interface {
	GetUserPlaylists(uId uint64) ([]models.Playlist, error)
	GetPlaylistById(pId uint64) (models.Playlist, error)
}
