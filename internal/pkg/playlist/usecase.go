package playlist

import "no_homomorphism/internal/pkg/models"

type UseCase interface {
	GetUserPlaylists(id string) ([]models.Playlist, error)
	GetPlaylistById(id string) (models.Playlist, error)
	CheckAccessToPlaylist(userId string, playlistId string) (bool, error)
}
