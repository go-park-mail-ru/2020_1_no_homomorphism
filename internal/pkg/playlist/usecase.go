package playlist

import "no_homomorphism/internal/pkg/models"

type UseCase interface {
	GetUserPlaylists(id string) (models.UserPlaylists, error)
	GetPlaylistWithTracks(id string) (models.PlaylistTracks, error)
}
