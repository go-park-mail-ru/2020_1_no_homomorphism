package playlist

import "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"

type Repository interface {
	GetUserPlaylists(uId string) ([]models.Playlist, error)
	GetPlaylistById(pId string) (models.Playlist, error)
	CreatePlaylist(name string, uID string) (plID string, err error)
	AddTrackToPlaylist(plTracks models.PlaylistTracks) error
	GetUserPlaylistsIdByTrack(userID, trackID string) ([]string, error)
	DeleteTrackFromPlaylist(plID, trackID string) error
	DeletePlaylist(plID string) error
	ChangePrivacy(plID string) error
	GetAllPlaylistTracks(plID string) ([]models.PlaylistTracks, error)
	UpdateAvatar(playlist models.Playlist, avatarDir string, fileType string) (string, error)
	Update(id string, name string) error
}
