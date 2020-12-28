package track

import "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"

type Repository interface {
	GetTrackById(id string) (models.Track, error)
	GetBoundedTracksByPlaylistId(plId string, start, end uint64) ([]models.Track, error)
	GetBoundedTracksByAlbumId(aId string, start, end uint64) ([]models.Track, error)
	GetBoundedTracksByArtistId(id string, start, end uint64) ([]models.Track, error)
	GetAllTracks() ([]models.Track, error)
	Search(text string, count uint) ([]models.TrackSearch, error)
	GetUserTracks(uID string) ([]models.Track, error)
	GetUserLikedTracksIDs(uID string) ([]int64, error)
	RateTrack(uID string, tID string) error
	IsLikedByUser(uID string, tID string) (bool, error)
}
