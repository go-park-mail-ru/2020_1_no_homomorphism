package playlist

import (
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"io"
)

type UseCase interface {
	GetUserPlaylists(id string) ([]models.Playlist, error)
	GetPlaylistById(id string) (models.Playlist, error)
	CreatePlaylist(name string, uID string) (plID string, err error)
	CheckAccessToPlaylist(userId string, playlistId string, isStrict bool) (bool, error)
	AddTrackToPlaylist(plTracks models.PlaylistTracks) error
	GetUserPlaylistsIdByTrack(userID, trackID string) ([]string, error)
	DeleteTrackFromPlaylist(plID, trackID string) error
	DeletePlaylist(plID string) error
	ChangePrivacy(plID string) error
	AddSharedPlaylist(plID string, uID string) (string, error)
	UpdateAvatar(plID string, file io.Reader, fileType string) (string, error)
	Update(id string, name string) error
}
