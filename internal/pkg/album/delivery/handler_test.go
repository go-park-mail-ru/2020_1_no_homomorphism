package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"no_homomorphism/internal/pkg/album"
	"no_homomorphism/internal/pkg/middleware"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/track"
	"no_homomorphism/pkg/logger"
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
		trueAuthPreHandle := middleware.AuthMiddlewareMock(albumHandlers.GetUserAlbums, true, testUser)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		albums := []models.Album{
			{
				Id:         "1",
				Name:       "KekLol",
				Image:      "/static/img",
				ArtistName: "TestArtist",
				ArtistId:   "42",
			},
			{
				Id:         "2",
				Name:       "testName",
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

		albumStr := `{"artist_id":"%s", "artist_name":"%s", "id":"%s", "image":"%s", "name":"%s"}`
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
				albums[1].ArtistId,
				albums[1].ArtistName,
				albums[1].Id,
				albums[1].Image,
				albums[1].Name,
			)).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserAlbums-NoAuth", func(t *testing.T) {
		trueAuthPreHandle := middleware.AuthMiddlewareMock(albumHandlers.GetUserAlbums, false, models.User{})

		apitest.New("GetUserAlbums-NoAuth").
			Handler(trueAuthPreHandle).
			Method("Get").
			Cookie("session_id", "randomSessionIdValueForTesting").
			URL("/users/albums").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("GetUserAlbums-error", func(t *testing.T) {
		trueAuthPreHandle := middleware.AuthMiddlewareMock(albumHandlers.GetUserAlbums, true, testUser)

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
			ArtistName: "TestArtist",
			ArtistId:   "42",
		}

		tracks := []models.Track{
			{
				Id:       "123",
				Name:     "name",
				Artist:   "Mc",
				Duration: 243,
				Image:    "/asdw/sdaasd/asd",
				Link:     "http://kek.lol.ru/test.pm3",
			},
			{
				Id:       "4312",
				Name:     "name2",
				Artist:   "NotMc",
				Duration: 232,
				Image:    "/ss/23/fgbg",
				Link:     "https://kek.test.ru/test.pm3",
			},
		}

		tracksMarshal, err := json.Marshal(tracks)
		assert.NoError(t, err)

		albumMarshal, err := json.Marshal(albumModel)
		assert.NoError(t, err)

		m := album.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		m.EXPECT().
			GetAlbumById(idVal).
			Return(albumModel, nil)

		tr.EXPECT().
			GetTracksByAlbumId(idVal).
			Return(tracks, nil)

		albumHandlers.AlbumUC = m
		albumHandlers.TrackUC = tr

		str := fmt.Sprintf(`{"album":%s,"count":%d,"tracks":%s}`, albumMarshal, len(tracks), tracksMarshal)

		apitest.New("GetFullAlbum-OK").
			Handler(trueAuthPreHandle).
			Method("Get").
			URL("/api/v1/albums/12").
			Expect(t).
			Body(str).
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

	t.Run("GetFullAlbum-GetTracksError", func(t *testing.T) {
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
			ArtistName: "TestArtist",
			ArtistId:   "42",
		}

		m := album.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		m.EXPECT().
			GetAlbumById(idVal).
			Return(albumModel, nil)

		testError := errors.New("test error")

		tr.EXPECT().
			GetTracksByAlbumId(idVal).
			Return([]models.Track{}, testError)

		albumHandlers.AlbumUC = m
		albumHandlers.TrackUC = tr

		apitest.New("GetFullAlbum-GetTracksError").
			Handler(trueAuthPreHandle).
			Method("Get").
			URL("/api/v1/albums/12").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}
