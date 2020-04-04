package usecase

import (
	"no_homomorphism/internal/pkg/artist"
	"no_homomorphism/internal/pkg/models"
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
