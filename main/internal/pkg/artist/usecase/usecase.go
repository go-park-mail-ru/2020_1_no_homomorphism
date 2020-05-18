package usecase

import (
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
)

type ArtistUseCase struct {
	ArtistRepository artist.Repository
}

func (uc *ArtistUseCase) GetArtistById(aID, uID string) (models.Artist, error) {
	dbArtist, err := uc.ArtistRepository.GetArtist(aID)
	if err != nil {
		return models.Artist{}, err
	}

	if uID != "" {
		dbArtist.IsSubscribed = uc.ArtistRepository.IsSubscribed(aID, uID)
	}

	return dbArtist, nil
}

func (uc *ArtistUseCase) Subscription(aID, uID string) error {
	return uc.ArtistRepository.Subscription(aID, uID)
}

func (uc *ArtistUseCase) SubscriptionList(uID string) ([]models.ArtistSearch, error) {
	return uc.ArtistRepository.SubscriptionsList(uID)
}

func (uc *ArtistUseCase) GetBoundedArtists(start, end uint64) ([]models.Artist, error) {
	return uc.ArtistRepository.GetBoundedArtists(start, end)
}

func (uc *ArtistUseCase) GetArtistStat(id string) (models.ArtistStat, error) {
	return uc.ArtistRepository.GetArtistStat(id)
}

func (uc *ArtistUseCase) Search(text string, count uint) ([]models.ArtistSearch, error) {
	return uc.ArtistRepository.Search(text, count)
}
