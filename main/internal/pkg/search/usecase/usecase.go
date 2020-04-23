package usecase

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/album"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
)

type SearchUseCase struct {
	ArtistRepo artist.Repository
	AlbumRepo  album.Repository
	TrackRepo  track.Repository
}

func (uc SearchUseCase) Search(text string, count uint) (models.SearchResult, error) {
	artistSearch, err := uc.ArtistRepo.Search(text, count)
	if err != nil {
		return models.SearchResult{}, fmt.Errorf("failed to search in artists")
	}

	albumSearch, err := uc.AlbumRepo.Search(text, count)
	if err != nil {
		return models.SearchResult{}, fmt.Errorf("failed to search in artists")
	}

	trackSearch, err := uc.TrackRepo.Search(text, count)
	if err != nil {
		return models.SearchResult{}, fmt.Errorf("failed to search in artists")
	}

	search := models.SearchResult{
		Artists: artistSearch,
		Albums:  albumSearch,
		Tracks:  trackSearch,
	}

	return search, nil
}
