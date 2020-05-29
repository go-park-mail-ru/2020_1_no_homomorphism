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

func (uc AlbumUseCase) GetAlbumById(aID, uID string) (models.Album, error) {
	dbAlbum, err := uc.AlbumRepository.GetAlbumById(aID)
	if err != nil {
		return models.Album{}, err
	}
	if uID != "" {
		dbAlbum.IsLiked = uc.AlbumRepository.CheckLike(aID, uID)
	}
	return dbAlbum, nil
}

func (uc AlbumUseCase) GetBoundedAlbumsByArtistId(id string, start uint64, end uint64) ([]models.Album, error) {
	return uc.AlbumRepository.GetBoundedAlbumsByArtistId(id, start, end)
}

func (uc AlbumUseCase) Search(text string, count uint) ([]models.AlbumSearch, error) {
	return uc.AlbumRepository.Search(text, count)
}

func (uc AlbumUseCase) RateAlbum(aID, uID string) error {
	return uc.AlbumRepository.RateAlbum(aID, uID)
}

func (uc AlbumUseCase) GetNewestReleases(uID string, begin, end int) ([]models.NewestReleases, error) {
	return uc.AlbumRepository.GetNewestReleases(uID, begin, end)
}
func (uc AlbumUseCase) GetWorldNews() ([]models.Album, error) {
	return uc.AlbumRepository.GetWorldNews()
}