package usecase

import (
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
)

type ArtistUseCase struct {
	ArtistRepository artist.Repository
}

func (uc *ArtistUseCase) GetArtistById(id string) (models.Artist, error) {
	return uc.ArtistRepository.GetArtist(id)
}

func (uc *ArtistUseCase) GetBoundedArtists(start, end uint64) ([]models.Artist, error) {
	return uc.ArtistRepository.GetBoundedArtists(start, end)
}

func (uc *ArtistUseCase) GetArtistStat(id string) (models.ArtistStat, error) {
	return uc.ArtistRepository.GetArtistStat(id)
}
