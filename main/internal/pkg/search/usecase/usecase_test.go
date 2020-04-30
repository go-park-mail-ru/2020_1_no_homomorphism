package usecase

import (
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/album"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch(t *testing.T) {
	//testError := errors.New("something go wrong")

	t.Run("Create-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		artistMock := artist.NewMockRepository(ctrl)
		albumMock := album.NewMockRepository(ctrl)
		trackMock := track.NewMockRepository(ctrl)

		text := "testText"
		var count uint = 5

		artistRes := []models.ArtistSearch{
			{
				ArtistID: "234",
				Name:     "test",
				Image:    "default.png",
			},
		}
		albumRes := []models.AlbumSearch{
			{
				AlbumID:    "34234",
				AlbumName:  "testLol",
				ArtistID:   "293874",
				ArtistName: "lolkek",
				Image:      "default.png",
			},
		}
		trackRes := []models.TrackSearch{
			{
				TrackID:    "234234",
				TrackName:  "tetsTrack",
				ArtistName: "artisttest",
				ArtistID:   "234234",
				Image:      "default.png",
			},
		}

		artistMock.
			EXPECT().
			Search(text, count).
			Return(artistRes, nil)

		albumMock.
			EXPECT().
			Search(text, count).
			Return(albumRes, nil)

		trackMock.
			EXPECT().
			Search(text, count).
			Return(trackRes, nil)

		useCase := SearchUseCase{
			ArtistRepo: artistMock,
			AlbumRepo:  albumMock,
			TrackRepo:  trackMock,
		}

		res, err := useCase.Search(text, count)
		assert.NoError(t, err)

		result := models.SearchResult{
			Artists: artistRes,
			Albums:  albumRes,
			Tracks:  trackRes,
		}

		assert.Equal(t, res, result)
	})
}
