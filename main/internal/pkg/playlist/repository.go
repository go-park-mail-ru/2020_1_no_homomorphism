package playlist

import "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"

type Repository interface {
	GetUserPlaylists(uId string) ([]models.Playlist, error)
	GetPlaylistById(pId string) (models.Playlist, error)
	CreatePlaylist(name string, uID string) (plID string, err error)
}
