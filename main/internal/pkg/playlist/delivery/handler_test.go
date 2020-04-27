package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/middleware"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
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
		handler := middleware.AuthMiddlewareMock(plHandler.GetUserPlaylists, true, testUser, "")

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



	t.Run("GetUserPlaylists-UseCaseError", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(plHandler.GetUserPlaylists, true, testUser, "")

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
			middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, true, testUser, ""),
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
			Return(true, nil)

		m.EXPECT().
			GetPlaylistById(testPl.Id).
			Return(testPl, nil)

		apitest.New("GetFullPlaylistById-OK").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/1234").
			Expect(t).
			Status(http.StatusOK).
			End()
	})



	t.Run("GetFullPlaylistById-NoVars", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, true, testUser, "")

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
			middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, true, testUser, ""),
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
			middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, true, testUser, ""),
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
			middleware.AuthMiddlewareMock(plHandler.GetFullPlaylistById, true, testUser, ""),
			"id",
			varsId,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m

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
}