package track

import "no_homomorphism/internal/pkg/models"

type Repository interface {
	GetTrackById(id uint64) (models.Track, error)
	GetPlaylistTracks(plId uint64) ([]models.Track, error)
	GetAlbumTracks(aId uint64) ([]models.Track, error)
	//GetArtistTracks(artistId uint64) ([]*models.Track, error)
}
