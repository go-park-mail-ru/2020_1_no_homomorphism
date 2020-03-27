package track

import "no_homomorphism/internal/pkg/models"

type Repository interface {
	GetTrackById(id string) (models.Track, error)
	GetPlaylistTracks(plId string) ([]models.Track, error)
	GetTracksByAlbumId(aId string) ([]models.Track, error)
	//GetArtistTracks(artistId uint64) ([]*models.Track, error)
}
