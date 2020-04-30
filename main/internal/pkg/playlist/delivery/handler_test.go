package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/middleware"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
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

var testTrack = models.Track{
	Id:       "123",
	Name:     "name",
	Artist:   "Mc",
	Duration: 243,
	Image:    "/asdw/sdaasd/asd",
	Link:     "http://kek.lol.ru/test.pm3",
}

var testTrack2 = models.Track{
	Id:       "123",
	Name:     "sdfa",
	Artist:   "Mc",
	Duration: 1234,
	Image:    "/asdw/asdf/asd",
	Link:     "http://kekasdf.lol.ru/teasdfst.pm3",
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
	t.Run("GetUserPlaylists-NoAuth", func(t *testing.T) {
		apitest.New("GetUserPlaylists-NoAuth").
			Handler(http.HandlerFunc(plHandler.GetUserPlaylists)).
			Method("Get").
			URL("/api/v1/users/playlists").
			Expect(t).
			Status(http.StatusInternalServerError).
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

func TestGetBoundedPlaylistTracks(t *testing.T) {
	t.Run("GetBoundedPlaylistTracks-OK", func(t *testing.T) {
		id := "1234"
		start := "0"
		end := "50"

		boundedVars := middleware.BoundedVars(
			middleware.AuthMiddlewareMock(plHandler.GetBoundedPlaylistTracks, true, testUser, ""),
			plHandler.Log,
		)

		handler := middleware.SetTripleVars(
			middleware.AuthMiddlewareMock(boundedVars, true, testUser, ""),
			id, start, end,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m
		plHandler.TrackUC = tr

		testPl := models.Playlist{
			Id:     id,
			Name:   "KekLol",
			Image:  "/static/img",
			UserId: testUser.Id,
		}

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, testPl.Id).
			Return(true, nil)

		var startUint uint64 = 0
		var endUint uint64 = 50

		tracks := []models.Track{testTrack, testTrack2}
		output := models.PlaylistTracksArray{
			Id:     id,
			Tracks: tracks,
		}

		jsonData, err := json.Marshal(output)
		assert.NoError(t, err)

		tr.EXPECT().
			GetBoundedTracksByPlaylistId(id, startUint, endUint).
			Return(tracks, nil)

		apitest.New("GetBoundedPlaylistTracks-OK").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/1234/0/50").
			Expect(t).
			Body(string(jsonData)).
			Status(http.StatusOK).
			End()
	})
	t.Run("GetBoundedPlaylistTracks-NoVars", func(t *testing.T) {
		apitest.New("GetBoundedPlaylistTracks-NoVars").
			Handler(http.HandlerFunc(plHandler.GetBoundedPlaylistTracks)).
			Method("Get").
			URL("/api/v1/playlists").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
	t.Run("GetBoundedPlaylistTracks-CheckAccessErr", func(t *testing.T) {
		id := "1234"
		start := "0"
		end := "50"

		boundedVars := middleware.BoundedVars(
			middleware.AuthMiddlewareMock(plHandler.GetBoundedPlaylistTracks, true, testUser, ""),
			plHandler.Log,
		)

		handler := middleware.SetTripleVars(
			middleware.AuthMiddlewareMock(boundedVars, true, testUser, ""),
			id, start, end,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m
		plHandler.TrackUC = tr

		testPl := models.Playlist{
			Id:     id,
			Name:   "KekLol",
			Image:  "/static/img",
			UserId: testUser.Id,
		}

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, testPl.Id).
			Return(false, nil)

		apitest.New("GetBoundedPlaylistTracks-CheckAccessErr").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/1234/0/50").
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})
	t.Run("GetBoundedPlaylistTracks-Error", func(t *testing.T) {
		id := "1234"
		start := "0"
		end := "50"

		boundedVars := middleware.BoundedVars(
			middleware.AuthMiddlewareMock(plHandler.GetBoundedPlaylistTracks, true, testUser, ""),
			plHandler.Log,
		)

		handler := middleware.SetTripleVars(
			middleware.AuthMiddlewareMock(boundedVars, true, testUser, ""),
			id, start, end,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m
		plHandler.TrackUC = tr

		testPl := models.Playlist{
			Id:     id,
			Name:   "KekLol",
			Image:  "/static/img",
			UserId: testUser.Id,
		}

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, testPl.Id).
			Return(true, nil)

		var startUint uint64 = 0
		var endUint uint64 = 50

		testError := errors.New("testError")

		tr.EXPECT().
			GetBoundedTracksByPlaylistId(id, startUint, endUint).
			Return([]models.Track{}, testError)

		apitest.New("GetBoundedPlaylistTracks-Error").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/1234/0/50").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestCreatePlaylist(t *testing.T) { //todo check errors
	t.Run("CreatePlaylist-OK", func(t *testing.T) {
		name := "testName"
		handler := middleware.SetMuxVars(
			middleware.AuthMiddlewareMock(plHandler.CreatePlaylist, true, testUser, ""),
			"name",
			name,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m

		plID := "12314"

		m.EXPECT().
			CreatePlaylist(name, testUser.Id).
			Return(plID, nil)

		apitest.New("CreatePlaylist-OK").
			Handler(handler).
			Method("Post").
			URL("/api/v1/playlists/testName").
			Expect(t).
			Status(http.StatusCreated).
			End()
	})
	t.Run("CreatePlaylist-NoVars", func(t *testing.T) {
		apitest.New("CreatePlaylist-NoVars").
			Handler(http.HandlerFunc(plHandler.CreatePlaylist)).
			Method("Post").
			URL("/api/v1/playlists/testName").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("CreatePlaylist-NoAuth", func(t *testing.T) {
		name := "testName"
		handler := middleware.SetMuxVars(
			plHandler.CreatePlaylist,
			"name",
			name,
		)

		apitest.New("CreatePlaylist-NoAuth").
			Handler(handler).
			Method("Post").
			URL("/api/v1/playlists/testName").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}

func TestAddTrackToPlaylist(t *testing.T) { //todo check errors, body check
	t.Run("AddTrackToPlaylist-OK", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(plHandler.AddTrackToPlaylist, true, testUser, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m
		plHandler.TrackUC = tr

		plTracks := models.PlaylistTracks{
			PlaylistID: "12341",
			TrackID:    "23",
			Index:      "5",
			Image:      "/static/default",
		}

		jsonData, err := json.Marshal(plTracks)
		assert.NoError(t, err)

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, plTracks.PlaylistID).
			Return(true, nil)

		m.EXPECT().
			AddTrackToPlaylist(plTracks).
			Return(nil)

		apitest.New("AddTrackToPlaylist-OK").
			Handler(handler).
			Method("Post").
			URL("/api/v1/playlists/tracks").
			Body(string(jsonData)).
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	t.Run("AddTrackToPlaylist-FailedToUnmarshall", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(plHandler.AddTrackToPlaylist, true, testUser, "")

		apitest.New("AddTrackToPlaylist-FailedToUnmarshall").
			Handler(handler).
			Method("Post").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("AddTrackToPlaylist-NoAccess", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(plHandler.AddTrackToPlaylist, true, testUser, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m
		plHandler.TrackUC = tr

		plTracks := models.PlaylistTracks{
			PlaylistID: "12341",
			TrackID:    "23",
			Index:      "5",
			Image:      "/static/default",
		}

		jsonData, err := json.Marshal(plTracks)
		assert.NoError(t, err)

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, plTracks.PlaylistID).
			Return(false, nil)

		apitest.New("AddTrackToPlaylist-NoAccess").
			Handler(handler).
			Method("Post").
			URL("/api/v1/playlists/tracks").
			Body(string(jsonData)).
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})
	t.Run("AddTrackToPlaylist-Error", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(plHandler.AddTrackToPlaylist, true, testUser, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)
		tr := track.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m
		plHandler.TrackUC = tr

		plTracks := models.PlaylistTracks{
			PlaylistID: "12341",
			TrackID:    "23",
			Index:      "5",
			Image:      "/static/default",
		}

		jsonData, err := json.Marshal(plTracks)
		assert.NoError(t, err)

		testError := errors.New("testError")

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, plTracks.PlaylistID).
			Return(true, nil)

		m.EXPECT().
			AddTrackToPlaylist(plTracks).
			Return(testError)

		apitest.New("AddTrackToPlaylist-Error").
			Handler(handler).
			Method("Post").
			URL("/api/v1/playlists/tracks").
			Body(string(jsonData)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestGetPlaylistsIDByTrack(t *testing.T) { //todo check errors, body check
	t.Run("GetPlaylistsIDByTrack-OK", func(t *testing.T) {
		id := "1234"
		handler := middleware.SetMuxVars(
			middleware.AuthMiddlewareMock(plHandler.GetPlaylistsIDByTrack, true, testUser, ""),
			"id",
			id,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m

		IDs := []string{"1", "2", "3"}

		m.EXPECT().
			GetUserPlaylistsIdByTrack(testUser.Id, id).
			Return(IDs, nil)

		jsonData, err := json.Marshal(models.PlaylistsID{IDs: IDs})
		assert.NoError(t, err)

		apitest.New("GetPlaylistsIDByTrack-OK").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Body(string(jsonData)).
			Status(http.StatusOK).
			End()
	})
	t.Run("GetPlaylistsIDByTrack-NoVars", func(t *testing.T) {
		apitest.New("GetPlaylistsIDByTrack-NoVars").
			Handler(http.HandlerFunc(plHandler.GetPlaylistsIDByTrack)).
			Method("Get").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("GetPlaylistsIDByTrack-NoAuth", func(t *testing.T) {
		id := "1234"
		handler := middleware.SetMuxVars(
			plHandler.GetPlaylistsIDByTrack,
			"id",
			id,
		)
		apitest.New("GetPlaylistsIDByTrack-NoAuth").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
	t.Run("GetPlaylistsIDByTrack-Error", func(t *testing.T) {
		id := "1234"
		handler := middleware.SetMuxVars(
			middleware.AuthMiddlewareMock(plHandler.GetPlaylistsIDByTrack, true, testUser, ""),
			"id",
			id,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m

		testError := errors.New("testError")

		m.EXPECT().
			GetUserPlaylistsIdByTrack(testUser.Id, id).
			Return([]string{}, testError)

		apitest.New("GetPlaylistsIDByTrack-Error").
			Handler(handler).
			Method("Get").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}

func TestDeletePlaylist(t *testing.T) { //todo check errors, body check
	t.Run("DeletePlaylist-OK", func(t *testing.T) {
		id := "1234"
		handler := middleware.SetMuxVars(
			middleware.AuthMiddlewareMock(plHandler.DeletePlaylist, true, testUser, ""),
			"id",
			id,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, id).
			Return(true, nil)

		m.EXPECT().
			DeletePlaylist(id).
			Return(nil)

		apitest.New("DeletePlaylist-OK").
			Handler(handler).
			Method("Delete").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	t.Run("DeletePlaylist-NoVars", func(t *testing.T) {
		apitest.New("DeletePlaylist-NoVars").
			Handler(http.HandlerFunc(plHandler.DeletePlaylist)).
			Method("Delete").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("DeletePlaylist-NoAccess", func(t *testing.T) {
		id := "1234"
		handler := middleware.SetMuxVars(
			plHandler.DeletePlaylist,
			"id",
			id,
		)

		apitest.New("DeletePlaylist-NoAccess").
			Handler(handler).
			Method("Delete").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
	t.Run("DeletePlaylist-Error", func(t *testing.T) {
		id := "1234"
		handler := middleware.SetMuxVars(
			middleware.AuthMiddlewareMock(plHandler.DeletePlaylist, true, testUser, ""),
			"id",
			id,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, id).
			Return(true, nil)

		testError := errors.New("testError")

		m.EXPECT().
			DeletePlaylist(id).
			Return(testError)

		apitest.New("DeletePlaylist-Error").
			Handler(handler).
			Method("Delete").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestDeleteTrackFromPlaylist(t *testing.T) { //todo check errors, body check
	t.Run("DeleteTrackFromPlaylist-OK", func(t *testing.T) {
		pl := middleware.VarsPair{
			Key:   "playlist",
			Value: "1234",
		}
		tr := middleware.VarsPair{
			Key:   "track",
			Value: "23523",
		}

		handler := middleware.SetUnlimitedVars(
			middleware.AuthMiddlewareMock(plHandler.DeleteTrackFromPlaylist, true, testUser, ""),
			pl, tr,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, pl.Value).
			Return(true, nil)

		m.EXPECT().
			DeleteTrackFromPlaylist(pl.Value, tr.Value).
			Return(nil)

		apitest.New("DeleteTrackFromPlaylist-OK").
			Handler(handler).
			Method("Delete").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	t.Run("DeleteTrackFromPlaylist-NoPlVar", func(t *testing.T) {
		tr := middleware.VarsPair{
			Key:   "track",
			Value: "23523",
		}

		handler := middleware.SetUnlimitedVars(
			middleware.AuthMiddlewareMock(plHandler.DeleteTrackFromPlaylist, true, testUser, ""),
			tr,
		)

		apitest.New("DeleteTrackFromPlaylist-NoPlVar").
			Handler(handler).
			Method("Delete").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("DeleteTrackFromPlaylist-NoTrVar", func(t *testing.T) {
		pl := middleware.VarsPair{
			Key:   "playlist",
			Value: "23523",
		}

		handler := middleware.SetUnlimitedVars(
			middleware.AuthMiddlewareMock(plHandler.DeleteTrackFromPlaylist, true, testUser, ""),
			pl,
		)

		apitest.New("DeleteTrackFromPlaylist-NoTrVar").
			Handler(handler).
			Method("Delete").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("DeleteTrackFromPlaylist-NoAccess", func(t *testing.T) {
		pl := middleware.VarsPair{
			Key:   "playlist",
			Value: "1234",
		}
		tr := middleware.VarsPair{
			Key:   "track",
			Value: "23523",
		}

		handler := middleware.SetUnlimitedVars(
			plHandler.DeleteTrackFromPlaylist,
			pl, tr,
		)

		apitest.New("DeleteTrackFromPlaylist-NoAccess").
			Handler(handler).
			Method("Delete").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
	t.Run("DeleteTrackFromPlaylist-Error", func(t *testing.T) {
		pl := middleware.VarsPair{
			Key:   "playlist",
			Value: "1234",
		}
		tr := middleware.VarsPair{
			Key:   "track",
			Value: "23523",
		}

		handler := middleware.SetUnlimitedVars(
			middleware.AuthMiddlewareMock(plHandler.DeleteTrackFromPlaylist, true, testUser, ""),
			pl, tr,
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := playlist.NewMockUseCase(ctrl)

		plHandler.PlaylistUC = m

		m.EXPECT().
			CheckAccessToPlaylist(testUser.Id, pl.Value).
			Return(true, nil)

		testError := errors.New("testError")

		m.EXPECT().
			DeleteTrackFromPlaylist(pl.Value, tr.Value).
			Return(testError)

		apitest.New("DeleteTrackFromPlaylist-Error").
			Handler(handler).
			Method("Delete").
			URL("/api/v1/playlists/tracks").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}
