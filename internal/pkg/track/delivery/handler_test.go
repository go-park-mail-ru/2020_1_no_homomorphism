package delivery

import (
	"encoding/json"
	"errors"
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
