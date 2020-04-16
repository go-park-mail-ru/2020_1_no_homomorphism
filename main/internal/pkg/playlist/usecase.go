package playlist

import "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"

type UseCase interface {
	GetUserPlaylists(id string) ([]models.Playlist, error)
	GetPlaylistById(id string) (models.Playlist, error)
	CreatePlaylist(name string, uID string) (plID string, err error)
	CheckAccessToPlaylist(userId string, playlistId string) (bool, error)
	AddTrackToPlaylist(plTracks models.PlaylistTracks) error
}
