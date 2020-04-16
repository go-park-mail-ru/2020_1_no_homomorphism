package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/album"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/middleware"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"strings"
	"testing"
)

var albumHandlers AlbumHandler

var testUser = models.User{
	Id:       "1234",
	Password: "76453647fvd",
	Name:     "TestName",
	Login:    "nnnagibator",
	Sex:      "Man",
	Image:    "/static/avatar/default.png",
	Email:    "klsJDLKfj@mail.ru",
}

func init() {
	albumHandlers.Log = logger.NewLogger(os.Stdout)
}

func TestGetUserAlbums(t *testing.T) {
	t.Run("GetUserAlbums-OK", func(t *testing.T) {
		trueAuthPreHandle := middleware.AuthMiddlewareMock(albumHandlers.GetUserAlbums, true, testUser, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		albums := []models.Album{
			{
				Id:         "1",
				Name:       "KekLol",
				Image:      "/static/img",
				Release:    "23-01-1999",
				ArtistName: "TestArtist",
				ArtistId:   "42",
			},
			{
				Id:         "2",
				Name:       "testName",
				Release:    "23-05-2010",
				Image:      "/static/img/test",
				ArtistName: "Mc Test",
				ArtistId:   "74",
			},
		}
		m := album.NewMockUseCase(ctrl)

		m.EXPECT().
			GetUserAlbums(testUser.Id).
			Return(albums, nil)

		albumHandlers.AlbumUC = m

		albumStr := `{"artist_id":"%s", "artist_name":"%s", "id":"%s", "image":"%s", "name":"%s", "release":"%s"}`
		bodyStr := strings.Join([]string{albumStr, albumStr}, ",")
		bodyStr = `{"albums":[` + bodyStr + `]}`

		apitest.New("GetUserAlbums-OK").
			Handler(trueAuthPreHandle).
			Method("Get").
			Cookie("session_id", "randomSessionIdValueForTesting").
			URL("/users/albums").
			Expect(t).
			Body(fmt.Sprintf(bodyStr,
				albums[0].ArtistId,
				albums[0].ArtistName,
				albums[0].Id,
				albums[0].Image,
				albums[0].Name,
				albums[0].Release,
				albums[1].ArtistId,
				albums[1].ArtistName,
				albums[1].Id,
				albums[1].Image,
				albums[1].Name,
				albums[1].Release,
			)).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserAlbums-error", func(t *testing.T) {
		trueAuthPreHandle := middleware.AuthMiddlewareMock(albumHandlers.GetUserAlbums, true, testUser, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := album.NewMockUseCase(ctrl)

		testError := errors.New("test error")

		m.EXPECT().
			GetUserAlbums(testUser.Id).
			Return([]models.Album{}, testError)

		albumHandlers.AlbumUC = m

		apitest.New("GetUserAlbums-error").
			Handler(trueAuthPreHandle).
			Method("Get").
			Cookie("session_id", "randomSessionIdValueForTesting").
			URL("/users/albums").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestGetFullAlbum(t *testing.T) {
	t.Run("GetFullAlbum-OK", func(t *testing.T) {
		idVal := "12"
		trueAuthPreHandle := middleware.SetMuxVars(
			albumHandlers.GetFullAlbum,
			"id",
			idVal,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		albumModel := models.Album{
			Id:         "1",
			Name:       "KekLol",
			Image:      "/static/img",
			Release:    "23-01-1999",
			ArtistName: "TestArtist",
			ArtistId:   "42",
		}

		albumMarshal, err := json.Marshal(albumModel)
		assert.NoError(t, err)

		m := album.NewMockUseCase(ctrl)

		m.EXPECT().
			GetAlbumById(idVal).
			Return(albumModel, nil)

		albumHandlers.AlbumUC = m

		apitest.New("GetFullAlbum-OK").
			Handler(trueAuthPreHandle).
			Method("Get").
			URL("/api/v1/albums/12").
			Expect(t).
			Body(string(albumMarshal)).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetFullAlbum-NoVars", func(t *testing.T) {

		apitest.New("GetFullAlbum-NoVars").
			Handler(http.HandlerFunc(albumHandlers.GetFullAlbum)).
			Method("Get").
			URL("/api/v1/albums/12").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetFullAlbum-GetAlbumError", func(t *testing.T) {
		idVal := "12"
		trueAuthPreHandle := middleware.SetMuxVars(
			albumHandlers.GetFullAlbum,
			"id",
			idVal,
		)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := album.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		testError := errors.New("test error")

		m.EXPECT().
			GetAlbumById(idVal).
			Return(models.Album{}, testError)

		albumHandlers.AlbumUC = m
		albumHandlers.TrackUC = tr

		apitest.New("GetFullAlbum-GetAlbumError").
			Handler(trueAuthPreHandle).
			Method("Get").
			URL("/api/v1/albums/12").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestGetBoundedAlbumsByArtistId(t *testing.T) {
	t.Run("GetBoundedAlbumsByArtistId-OK", func(t *testing.T) {
		artistId := "1231"
		start := "0"
		end := "2"

		boundedVars := middleware.GetBoundedVars(albumHandlers.GetBoundedAlbumsByArtistId, albumHandlers.Log)
		vars := middleware.SetTripleVars(boundedVars, artistId, start, end)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		albums := []models.Album{
			{
				Id:         "1",
				Name:       "KekLol",
				Image:      "/static/img",
				Release:    "23-01-1999",
				ArtistName: "TestArtist",
				ArtistId:   artistId,
			},
			{
				Id:         "2",
				Name:       "testName",
				Release:    "23-05-2010",
				Image:      "/static/img/test",
				ArtistName: "Mc Test",
				ArtistId:   artistId,
			},
		}

		m := album.NewMockUseCase(ctrl)

		var startUint uint64 = 0
		var endUint uint64 = 2

		m.EXPECT().
			GetBoundedAlbumsByArtistId(artistId, startUint, endUint).
			Return(albums, nil)

		albumHandlers.AlbumUC = m

		albumsMarshal, err := json.Marshal(albums)
		assert.NoError(t, err)

		str := `{"id":"1231","albums":` + string(albumsMarshal) + `}`

		apitest.New("GetBoundedAlbumsByArtistId-OK").
			Handler(vars).
			Method("Get").
			Expect(t).
			Body(str).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetBoundedAlbumsByArtistId-NoVars", func(t *testing.T) {

		apitest.New("GetBoundedAlbumsByArtistId-NoVars").
			Handler(http.HandlerFunc(albumHandlers.GetBoundedAlbumsByArtistId)).
			Method("Get").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetBoundedAlbumsByArtistId-UseCaseError", func(t *testing.T) {
		artistId := "1231"
		start := "0"
		end := "2"

		boundedVars := middleware.GetBoundedVars(albumHandlers.GetBoundedAlbumsByArtistId, albumHandlers.Log)
		vars := middleware.SetTripleVars(boundedVars, artistId, start, end)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := album.NewMockUseCase(ctrl)

		var startUint uint64 = 0
		var endUint uint64 = 2

		testError := errors.New("test error")

		m.EXPECT().
			GetBoundedAlbumsByArtistId(artistId, startUint, endUint).
			Return([]models.Album{}, testError)

		albumHandlers.AlbumUC = m

		apitest.New("GetBoundedAlbumsByArtistId-UseCaseError").
			Handler(vars).
			Method("Get").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}
