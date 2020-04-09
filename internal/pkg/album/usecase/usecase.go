package usecase

import (
	"no_homomorphism/internal/pkg/album"
	"no_homomorphism/internal/pkg/models"
)

type AlbumUseCase struct {
	AlbumRepository album.Repository
}

func (uc AlbumUseCase) GetUserAlbums(id string) ([]models.Album, error) {
	return uc.AlbumRepository.GetUserAlbums(id)
}

func (uc AlbumUseCase) GetAlbumById(id string) (models.Album, error) {
	return uc.AlbumRepository.GetAlbumById(id)
}

func (uc AlbumUseCase) GetBoundedAlbumsByArtistId(id string, start uint64, end uint64) ([]models.Album, error) {
	return uc.AlbumRepository.GetBoundedAlbumsByArtistId(id, start, end)
}
