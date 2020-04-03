package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"no_homomorphism/internal/pkg/middleware"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/playlist"
	"no_homomorphism/internal/pkg/track"
	"no_homomorphism/pkg/logger"
	"os"
	"testing"
)

var plHandler PlaylistHandler

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
	plHandler.Log = logger.NewLogger(os.Stdout)
}

func TestGetUserPlaylists(t *testing.T) {
	t.Run("GetUserPlaylists-OK", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(plHandler.GetUserPlaylists, true, testUser)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		plHandler.PlaylistUC = m

		playlists := []models.Playlist{
			{
				Id:     "1",
				Name:   "KekLol",
				Image:  "/static/img",
				UserId: testUser.Id,
			},
			{
				Id:     "2",
				Name:   "testName",
				Image:  "/static/img/test",
				UserId: testUser.Id,
			},
		}
		plMarshal, err := json.Marshal(playlists)
		assert.NoError(t, err)
		resBody := fmt.Sprintf(`{"playlists":%s}`, string(plMarshal))

		m.EXPECT().
			GetUserPlaylists(testUser.Id).
			Return(playlists, nil)

		apitest.New("GetUserPlaylists-OK").
			Handler(handler).
			Method("Get").
			URL("/api/v1/users/playlists").
			Expect(t).
			Body(resBody).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserPlaylists-NoAuth", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(plHandler.GetUserPlaylists, false, models.User{})

		apitest.New("GetUserPlaylists-NoAuth").
			Handler(handler).
			Method("Get").
			URL("/api/v1/users/playlists").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("GetUserPlaylists-UseCaseError", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(plHandler.GetUserPlaylists, true, testUser)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		plHandler.PlaylistUC = m

		testError := errors.New("test error")

		m.EXPECT().
			GetUserPlaylists(testUser.Id).
			Return([]models.Playlist{}, testError)

		apitest.New("GetUserPlaylists-UseCaseError").
			Handler(handler).
			Method("Get").
			URL("/api/v1/users/playlists").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestGetFullPlaylistById(t *testing.T) {
	t.Run("GetFullPlaylistById-OK", func(t *testing.T) {
		varsId := "1234"
		handler := middleware.SetMuxVars(
			middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, true, testUser),
			"id",
			varsId,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m
		plHandler.TrackUC = tr

		testPl := models.Playlist{
			Id:     varsId,
			Name:   "KekLol",
			Image:  "/static/img",
			UserId: testUser.Id,
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

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, testPl.Id).
			Return(true, nil)

		m.EXPECT().
			GetPlaylistById(testPl.Id).
			Return(testPl, nil)

		tr.EXPECT().
			GetTracksByPlaylistId(testPl.Id).
			Return(tracks, nil)

		apitest.New("GetFullPlaylistById-OK").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/1234").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetFullPlaylistById-NoAuth", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, false, models.User{})

		apitest.New("GetFullPlaylistById-NoAuth").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/1234").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("GetFullPlaylistById-NoVars", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, true, testUser)

		apitest.New("GetFullPlaylistById-NoVars").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetFullPlaylistById-CheckAccessError", func(t *testing.T) {
		varsId := "1234"
		handler := middleware.SetMuxVars(
			middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, true, testUser),
			"id",
			varsId,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m
		plHandler.TrackUC = tr

		testError := errors.New("test error")

		testPl := models.Playlist{
			Id:     varsId,
			Name:   "KekLol",
			Image:  "/static/img",
			UserId: testUser.Id,
		}

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, testPl.Id).
			Return(false, testError)

		apitest.New("GetFullPlaylistById-CheckAccessError").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/1234").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetFullPlaylistById-NoAccess", func(t *testing.T) {
		varsId := "1234"
		handler := middleware.SetMuxVars(
			middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, true, testUser),
			"id",
			varsId,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m
		plHandler.TrackUC = tr

		testPl := models.Playlist{
			Id:     varsId,
			Name:   "KekLol",
			Image:  "/static/img",
			UserId: testUser.Id,
		}

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, testPl.Id).
			Return(false, nil)

		apitest.New("GetFullPlaylistById-NoAccess").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/1234").
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("GetFullPlaylistById-GetPlError", func(t *testing.T) {
		varsId := "1234"
		handler := middleware.SetMuxVars(
			middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, true, testUser),
			"id",
			varsId,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m
		plHandler.TrackUC = tr

		testPl := models.Playlist{
			Id:     varsId,
			Name:   "KekLol",
			Image:  "/static/img",
			UserId: testUser.Id,
		}

		testError := errors.New("test error")

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, testPl.Id).
			Return(true, nil)

		m.EXPECT().
			GetPlaylistById(testPl.Id).
			Return(models.Playlist{}, testError)

		apitest.New("GetFullPlaylistById-GetPlError").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/1234").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetFullPlaylistById-GetTracksError", func(t *testing.T) {
		varsId := "1234"
		handler := middleware.SetMuxVars(
			middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, true, testUser),
			"id",
			varsId,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m
		plHandler.TrackUC = tr

		testPl := models.Playlist{
			Id:     varsId,
			Name:   "KekLol",
			Image:  "/static/img",
			UserId: testUser.Id,
		}

		testError := errors.New("test error")

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, testPl.Id).
			Return(true, nil)

		m.EXPECT().
			GetPlaylistById(testPl.Id).
			Return(testPl, nil)

		tr.EXPECT().
			GetTracksByPlaylistId(testPl.Id).
			Return([]models.Track{}, testError)

		apitest.New("GetFullPlaylistById-GetTracksError").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/1234").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}
