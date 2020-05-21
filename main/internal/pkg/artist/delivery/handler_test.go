package delivery

import (
	"encoding/json"
	"errors"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/middleware"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

var artistHandler ArtistHandler

var artist1 = models.Artist{
	Id:    "123124",
	Name:  "testName",
	Image: "testImage",
	Genre: "genre",
}

var artist2 = models.Artist{
	Id:    "36346",
	Name:  "testName2",
	Image: "test2Image",
	Genre: "rock",
}

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
	artistHandler.Log = logger.NewLogger(os.Stdout)
}

func TestGetFullArtistInfo(t *testing.T) {
	t.Run("GetFullArtistInfo-OK", func(t *testing.T) {
		id := "4234124"
		handler := middleware.SetMuxVars(artistHandler.GetFullArtistInfo, "id", id)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := artist.NewMockUseCase(ctrl)
		artistHandler.ArtistUC = m

		jsonData, err := json.Marshal(artist1)
		assert.NoError(t, err)

		m.EXPECT().
			GetArtistById(id, "").
			Return(artist1, nil)

		apitest.New("GetFullArtistInfo-OK").
			Handler(handler).
			Method("Get").
			URL("/api/v1/artists/").
			Expect(t).
			Body(string(jsonData)).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetFullArtistInfo-NoID", func(t *testing.T) {
		apitest.New("GetFullArtistInfo-NoID").
			Handler(http.HandlerFunc(artistHandler.GetFullArtistInfo)).
			Method("Get").
			URL("/api/v1/artists/").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetFullArtistInfo-GetArtistByIdError", func(t *testing.T) {
		id := "4234124"
		handler := middleware.SetMuxVars(artistHandler.GetFullArtistInfo, "id", id)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := artist.NewMockUseCase(ctrl)
		artistHandler.ArtistUC = m

		testError := errors.New("testError")

		m.EXPECT().
			GetArtistById(id, "").
			Return(models.Artist{}, testError)

		apitest.New("GetFullArtistInfo-GetArtistByIdError").
			Handler(handler).
			Method("Get").
			URL("/api/v1/artists/").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestGetBoundedArtists(t *testing.T) {
	t.Run("GetBoundedArtists-OK", func(t *testing.T) {
		id := "4234124"
		start := "0"
		end := "50"

		handler := middleware.SetTripleVars(artistHandler.GetBoundedArtists, id, start, end)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := artist.NewMockUseCase(ctrl)
		artistHandler.ArtistUC = m

		artistsArray := []models.Artist{artist1, artist2}

		artists := models.Artists{Artists: artistsArray}

		jsonData, err := json.Marshal(artists)
		assert.NoError(t, err)

		var startUint uint64 = 0
		var endUint uint64 = 50

		m.EXPECT().
			GetBoundedArtists(startUint, endUint).
			Return(artistsArray, nil)

		apitest.New("GetBoundedArtists-OK").
			Handler(handler).
			Method("Get").
			URL("/api/v1/artists/0/50").
			Expect(t).
			Body(string(jsonData)).
			Status(http.StatusOK).
			End()
	})
	t.Run("GetBoundedArtists-NoVars", func(t *testing.T) {
		apitest.New("GetBoundedArtists-NoVars").
			Handler(http.HandlerFunc(artistHandler.GetBoundedArtists)).
			Method("Get").
			URL("/api/v1/artists/0/50").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("GetBoundedArtists-Error", func(t *testing.T) {
		id := "4234124"
		start := "notInt"
		end := "notInt"

		handler := middleware.SetTripleVars(artistHandler.GetBoundedArtists, id, start, end)

		apitest.New("GetBoundedArtists-BadInt").
			Handler(handler).
			Method("Get").
			URL("/api/v1/artists/0/50").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("GetBoundedArtists-OK", func(t *testing.T) {
		id := "4234124"
		start := "0"
		end := "50"

		handler := middleware.SetTripleVars(artistHandler.GetBoundedArtists, id, start, end)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := artist.NewMockUseCase(ctrl)
		artistHandler.ArtistUC = m

		var startUint uint64 = 0
		var endUint uint64 = 50

		testError := errors.New("asdjhaoi")

		m.EXPECT().
			GetBoundedArtists(startUint, endUint).
			Return([]models.Artist{}, testError)

		apitest.New("GetBoundedArtists-Error").
			Handler(handler).
			Method("Get").
			URL("/api/v1/artists/0/50").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestGetArtistStat(t *testing.T) {
	t.Run("GetArtistStat-OK", func(t *testing.T) {
		id := "43423"
		handler := middleware.SetMuxVars(artistHandler.GetArtistStat, "id", id)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := artist.NewMockUseCase(ctrl)
		artistHandler.ArtistUC = m

		stat := models.ArtistStat{
			ArtistId:    "4243",
			Tracks:      4,
			Albums:      67,
			Subscribers: 2,
		}

		jsonData, err := json.Marshal(stat)
		assert.NoError(t, err)

		m.EXPECT().
			GetArtistStat(id).
			Return(stat, nil)

		apitest.New("GetArtistStat-OK").
			Handler(handler).
			Method("Get").
			URL("/api/v1/artists/stat").
			Expect(t).
			Body(string(jsonData)).
			Status(http.StatusOK).
			End()
	})
	t.Run("GetArtistStat-NoID", func(t *testing.T) {
		apitest.New("GetArtistStat-NoID").
			Handler(http.HandlerFunc(artistHandler.GetArtistStat)).
			Method("Get").
			URL("/api/v1/artists/stat").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("GetArtistStat-Error", func(t *testing.T) {
		id := "43423"
		handler := middleware.SetMuxVars(artistHandler.GetArtistStat, "id", id)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := artist.NewMockUseCase(ctrl)
		artistHandler.ArtistUC = m

		testError := errors.New("testError")

		m.EXPECT().
			GetArtistStat(id).
			Return(models.ArtistStat{}, testError)

		apitest.New("GetArtistStat-Error").
			Handler(handler).
			Method("Get").
			URL("/api/v1/artists/stat").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestSubscribe(t *testing.T) {
	t.Run("Subscribe-OK", func(t *testing.T) {
		id := "43423"

		handler := middleware.AuthMiddlewareMock(
			middleware.SetMuxVars(artistHandler.Subscribe, "id", id),
			true,
			testUser,
			"",
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := artist.NewMockUseCase(ctrl)
		artistHandler.ArtistUC = m

		m.EXPECT().
			Subscription(id, testUser.Id).
			Return(nil)

		apitest.New("Subscribe-OK").
			Handler(handler).
			Method("Get").
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	t.Run("Subscribe-NoAuth", func(t *testing.T) {
		apitest.New("Subscribe-NoAuth").
			Handler(http.HandlerFunc(artistHandler.Subscribe)).
			Method("Get").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
	t.Run("Subscribe-NoMuxVars", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(
			artistHandler.Subscribe,
			true,
			testUser,
			"",
		)

		apitest.New("Subscribe-NoMuxVars").
			Handler(handler).
			Method("Get").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("Subscribe-Error", func(t *testing.T) {
		id := "43423"

		handler := middleware.AuthMiddlewareMock(
			middleware.SetMuxVars(artistHandler.Subscribe, "id", id),
			true,
			testUser,
			"",
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := artist.NewMockUseCase(ctrl)
		artistHandler.ArtistUC = m

		testError := errors.New("test error")

		m.EXPECT().
			Subscription(id, testUser.Id).
			Return(testError)

		apitest.New("Subscribe-Error").
			Handler(handler).
			Method("Get").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}

func TestSubscriptionList(t *testing.T) {
	t.Run("SubscriptionList-OK", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(
			artistHandler.SubscriptionList,
			true,
			testUser,
			"",
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := artist.NewMockUseCase(ctrl)
		artistHandler.ArtistUC = m

		subList := []models.ArtistSearch{
			{
				ArtistID: "123",
				Name:     "test",
				Image:    "image",
			},
			{
				ArtistID: "43124",
				Name:     "test2",
				Image:    "image2",
			},
		}

		m.EXPECT().
			SubscriptionList(testUser.Id).
			Return(subList, nil)

		jsonData, err := json.Marshal(subList)
		assert.NoError(t, err)

		apitest.New("SubscriptionList-OK").
			Handler(handler).
			Method("Get").
			Expect(t).
			Body(string(jsonData)).
			Status(http.StatusOK).
			End()
	})

	t.Run("SubscriptionList-NoAuth", func(t *testing.T) {
		apitest.New("SubscriptionList-NoAuth").
			Handler(http.HandlerFunc(artistHandler.SubscriptionList)).
			Method("Get").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("SubscriptionList-Error", func(t *testing.T) {
		handler := middleware.AuthMiddlewareMock(
			artistHandler.SubscriptionList,
			true,
			testUser,
			"",
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := artist.NewMockUseCase(ctrl)
		artistHandler.ArtistUC = m

		testError := errors.New("test error")

		m.EXPECT().
			SubscriptionList(testUser.Id).
			Return(nil, testError)

		apitest.New("SubscriptionList-Error").
			Handler(handler).
			Method("Get").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}
