package delivery

import (
	"encoding/json"
	"errors"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/middleware"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/search"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"strconv"
	"testing"
)

var searchHandler SearchHandler

func init() {
	searchHandler.Log = logger.NewLogger(os.Stdout)
}

func TestSearch(t *testing.T) {
	t.Run("Search-OK", func(t *testing.T) {
		someText := "RandomSearchRequest"
		count := 5

		paramText := middleware.VarsPair{Key: "text", Value: someText}
		paramCount := middleware.VarsPair{Key: "count", Value: strconv.Itoa(count)}

		handler := middleware.SetUnlimitedVars(searchHandler.Search, paramText, paramCount)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := search.NewMockUseCase(ctrl)

		searchRes := models.SearchResult{
			Artists: []models.ArtistSearch{
				{
					ArtistID: "123",
					Name:     "Test",
					Image:    "/img/default.png",
				},
			},
			Albums: []models.AlbumSearch{
				{
					AlbumID:    "1234",
					AlbumName:  "name",
					ArtistID:   "32423",
					ArtistName: "rfwe",
					Image:      "/default.png",
				},
			},
			Tracks: []models.TrackSearch{
				{
					TrackID:    "123",
					TrackName:  "track",
					ArtistName: "asd",
					ArtistID:   "12513",
					Image:      "track.png",
				},
			},
		}

		jsonData, err := json.Marshal(searchRes)
		assert.NoError(t, err)

		m.EXPECT().
			Search(someText, uint(count)).
			Return(searchRes, nil)

		searchHandler.SearchUC = m

		apitest.New("Search-OK").
			Handler(handler).
			Method("Get").
			URL("/media/RandomSearchRequest/5").
			Expect(t).
			Body(string(jsonData)).
			Status(http.StatusOK).
			End()
	})
	t.Run("Search-NoCount", func(t *testing.T) {
		someText := "RandomSearchRequest"

		paramText := middleware.VarsPair{Key: "text", Value: someText}

		handler := middleware.SetUnlimitedVars(searchHandler.Search, paramText)

		apitest.New("Search-NoCount").
			Handler(handler).
			Method("Get").
			URL("/media/RandomSearchRequest/5").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Search-NoText", func(t *testing.T) {
		count := 5

		paramCount := middleware.VarsPair{Key: "count", Value: strconv.Itoa(count)}
		handler := middleware.SetUnlimitedVars(searchHandler.Search, paramCount)

		apitest.New("Search-NoText").
			Handler(handler).
			Method("Get").
			URL("/media/RandomSearchRequest/5").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Search-BadCount", func(t *testing.T) {
		someText := "RandomSearchRequest"

		paramText := middleware.VarsPair{Key: "text", Value: someText}
		paramCount := middleware.VarsPair{Key: "count", Value: "NotAnInteger"}

		handler := middleware.SetUnlimitedVars(searchHandler.Search, paramText, paramCount)

		apitest.New("Search-BadCount").
			Handler(handler).
			Method("Get").
			URL("/media/RandomSearchRequest/5").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Search-SearchError", func(t *testing.T) {
		someText := "RandomSearchRequest"
		count := 5

		paramText := middleware.VarsPair{Key: "text", Value: someText}
		paramCount := middleware.VarsPair{Key: "count", Value: strconv.Itoa(count)}

		handler := middleware.SetUnlimitedVars(searchHandler.Search, paramText, paramCount)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := search.NewMockUseCase(ctrl)

		testError := errors.New("testError")

		m.EXPECT().
			Search(someText, uint(count)).
			Return(models.SearchResult{}, testError)

		searchHandler.SearchUC = m

		apitest.New("Search-SearchError").
			Handler(handler).
			Method("Get").
			URL("/media/RandomSearchRequest/5").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

}
