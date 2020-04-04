package track

import "no_homomorphism/internal/pkg/models"

type UseCase interface {
	GetTrackById(id string) (models.Track, error)
	GetTracksByAlbumId(id string) ([]models.Track, error)
	GetTracksByPlaylistId(id string) ([]models.Track, error)
	GetBoundedTracksByArtistId(id string, start uint64, end uint64) ([]models.Track, error)
	GetBoundedTracksByAlbumId(aId string, start, end uint64) ([]models.Track, error)
	GetBoundedTracksByPlaylistId(plId string, start, end uint64) ([]models.Track, error)
}
