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
	"no_homomorphism/internal/pkg/track"
	"no_homomorphism/pkg/logger"
	"os"
	"testing"
)

var trackHandler TrackHandler

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
	trackHandler.Log = logger.NewLogger(os.Stdout)
}

func TestGetTrack(t *testing.T) {
	t.Run("GetTrack-OK", func(t *testing.T) {
		tId := "2123"
		handler := middleware.SetMuxVars(trackHandler.GetTrack, "id", tId)

		testTrack := models.Track{
			Id:       "123",
			Name:     "name",
			Artist:   "Mc",
			Duration: 243,
			Image:    "/asdw/sdaasd/asd",
			Link:     "http://kek.lol.ru/test.pm3",
		}

		trackMarshal, err := json.Marshal(testTrack)
		assert.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := track.NewMockUseCase(ctrl)

		trackHandler.TrackUC = m

		m.EXPECT().
			GetTrackById(tId).
			Return(testTrack, nil)

		apitest.New("GetTrack-OK").
			Handler(handler).
			Method("Get").
			URL("/api/v1/tracks/2123").
			Expect(t).
			Body(string(trackMarshal)).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetTrack-NoMux", func(t *testing.T) {
		apitest.New("GetTrack-NoMux").
			Handler(http.HandlerFunc(trackHandler.GetTrack)).
			Method("Get").
			URL("/api/v1/tracks/2123").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetTrack-UseCaseError", func(t *testing.T) {
		tId := "2123"
		handler := middleware.SetMuxVars(trackHandler.GetTrack, "id", tId)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := track.NewMockUseCase(ctrl)
		trackHandler.TrackUC = m

		testError := errors.New("testError")

		m.EXPECT().
			GetTrackById(tId).
			Return(models.Track{}, testError)

		apitest.New("GetTrack-UseCaseError").
			Handler(handler).
			Method("Get").
			URL("/api/v1/tracks/2123").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestGetBoundedArtistTracks(t *testing.T) {
	t.Run("GetBoundedArtistTracks-OK", func(t *testing.T) {
		artistId := "231"
		start := "0"
		end := "2"

		boundedVars := middleware.GetBoundedVars(trackHandler.GetBoundedArtistTracks, trackHandler.Log)
		vars := middleware.SetTripleVars(boundedVars, artistId, start, end)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := track.NewMockUseCase(ctrl)
		trackHandler.TrackUC = m

		tracksArray := []models.Track{testTrack, testTrack2}
		tracksMarshal, err := json.Marshal(tracksArray)
		assert.NoError(t, err)

		var startUint uint64 = 0
		var endUint uint64 = 2

		m.EXPECT().
			GetBoundedTracksByArtistId(artistId, startUint, endUint).
			Return(tracksArray, nil)

		str := fmt.Sprintf(`{"id":"%v","tracks":%v}`, artistId, string(tracksMarshal))

		apitest.New("GetBoundedArtistTracks-OK").
			Handler(vars).
			Method("Get").
			Expect(t).
			Body(str).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetBoundedArtistTracks-NoVars", func(t *testing.T) {
		apitest.New("GetBoundedArtistTracks-OK").
			Handler(http.HandlerFunc(trackHandler.GetBoundedArtistTracks)).
			Method("Get").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetBoundedArtistTracks-UseCaseError", func(t *testing.T) {
		artistId := "231"
		start := "0"
		end := "2"

		boundedVars := middleware.GetBoundedVars(trackHandler.GetBoundedArtistTracks, trackHandler.Log)
		vars := middleware.SetTripleVars(boundedVars, artistId, start, end)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := track.NewMockUseCase(ctrl)
		trackHandler.TrackUC = m

		var startUint uint64 = 0
		var endUint uint64 = 2

		testError := errors.New("testError")

		m.EXPECT().
			GetBoundedTracksByArtistId(artistId, startUint, endUint).
			Return([]models.Track{}, testError)


		apitest.New("GetBoundedArtistTracks-UseCaseError").
			Handler(vars).
			Method("Get").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}


func TestGetBoundedAlbumTracks(t *testing.T) {
	t.Run("GetBoundedAlbumTracks-OK", func(t *testing.T) {
		artistId := "231"
		start := "0"
		end := "2"

		boundedVars := middleware.GetBoundedVars(trackHandler.GetBoundedAlbumTracks, trackHandler.Log)
		vars := middleware.SetTripleVars(boundedVars, artistId, start, end)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := track.NewMockUseCase(ctrl)
		trackHandler.TrackUC = m

		tracksArray := []models.Track{testTrack, testTrack2}
		tracksMarshal, err := json.Marshal(tracksArray)
		assert.NoError(t, err)

		var startUint uint64 = 0
		var endUint uint64 = 2

		m.EXPECT().
			GetBoundedTracksByAlbumId(artistId, startUint, endUint).
			Return(tracksArray, nil)

		str := fmt.Sprintf(`{"id":"%v","tracks":%v}`, artistId, string(tracksMarshal))

		apitest.New("GetBoundedAlbumTracks-OK").
			Handler(vars).
			Method("Get").
			Expect(t).
			Body(str).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetBoundedAlbumTracks-NoVars", func(t *testing.T) {
		apitest.New("GetBoundedAlbumTracks-OK").
			Handler(http.HandlerFunc(trackHandler.GetBoundedAlbumTracks)).
			Method("Get").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetBoundedAlbumTracks-UseCaseError", func(t *testing.T) {
		artistId := "231"
		start := "0"
		end := "2"

		boundedVars := middleware.GetBoundedVars(trackHandler.GetBoundedAlbumTracks, trackHandler.Log)
		vars := middleware.SetTripleVars(boundedVars, artistId, start, end)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := track.NewMockUseCase(ctrl)
		trackHandler.TrackUC = m

		var startUint uint64 = 0
		var endUint uint64 = 2

		testError := errors.New("testError")

		m.EXPECT().
			GetBoundedTracksByAlbumId(artistId, startUint, endUint).
			Return([]models.Track{}, testError)


		apitest.New("GetBoundedAlbumTracks-UseCaseError").
			Handler(vars).
			Method("Get").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}