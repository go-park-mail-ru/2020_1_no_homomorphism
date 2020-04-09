package playlist

import "no_homomorphism/internal/pkg/models"

type Repository interface {
	GetUserPlaylists(uId string) ([]models.Playlist, error)
	GetPlaylistById(pId string) (models.Playlist, error)
}
