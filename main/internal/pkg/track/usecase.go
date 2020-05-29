package track

import "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"

type UseCase interface {
	GetTrackById(id string) (models.Track, error)
	GetBoundedTracksByArtistId(id string, start, end uint64, uID string) ([]models.Track, error)
	GetBoundedTracksByAlbumId(aId string, start, end uint64, uID string) ([]models.Track, error)
	GetBoundedTracksByPlaylistId(plId string, start, end uint64, uID string) ([]models.Track, error)
	Search(text string, count uint) ([]models.TrackSearch, error)
	GetUserTracks(uID string) ([]models.Track, error)
	RateTrack(uID, tID string) error
	IsLikedByUser(uID string, tID string) (bool, error)
}
