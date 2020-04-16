package usecase

import (
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/album"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
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
