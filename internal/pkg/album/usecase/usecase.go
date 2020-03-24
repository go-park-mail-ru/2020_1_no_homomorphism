package usecase

import (
	"fmt"
	"no_homomorphism/internal/pkg/album"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/track"
	"strconv"
)

type AlbumUseCase struct {
	AlbumRepository album.Repository
	TrackRepository track.Repository
}

func (uc AlbumUseCase) GetUserAlbums(id string) (models.UserAlbums, error) {
	uId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return models.UserAlbums{}, fmt.Errorf("failed to convert id: %s", err)
	}
	albums, err := uc.AlbumRepository.GetUserAlbums(uId)
	if err != nil {
		return models.UserAlbums{}, fmt.Errorf("failed to get from album repo: %s", err)
	}
	return models.UserAlbums{
		Count:  len(albums),
		Albums: albums,
	}, nil
}

func (uc AlbumUseCase) GetAlbumTracks(id string) (models.AlbumTracks, error) {
	aId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return models.AlbumTracks{}, fmt.Errorf("failed to convert id: %s", err)
	}
	dbAlbum, err := uc.AlbumRepository.GetAlbumById(aId)
	if err != nil {
		return models.AlbumTracks{}, fmt.Errorf("failed to get from dbAlbum repo: %s", err)
	}
	tracks, err := uc.TrackRepository.GetAlbumTracks(aId)
	if err != nil {
		return models.AlbumTracks{}, fmt.Errorf("failed to get from track repo: %s", err)
	}

	return models.AlbumTracks{
		Album:  dbAlbum,
		Count:  len(tracks),
		Tracks: tracks,
	}, nil
}
