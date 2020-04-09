package track

import "no_homomorphism/internal/pkg/models"

type Repository interface {
	GetTrackById(id string) (models.Track, error)
	GetBoundedTracksByPlaylistId(plId string, start, end uint64) ([]models.Track, error)
	GetBoundedTracksByAlbumId(aId string, start, end uint64) ([]models.Track, error)
	GetBoundedTracksByArtistId(id string, start, end uint64) ([]models.Track, error)
}
