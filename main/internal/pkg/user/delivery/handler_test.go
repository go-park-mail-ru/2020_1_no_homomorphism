package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/csrf"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/middleware"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/session"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/user"
	"github.com/2020_1_no_homomorphism/no_homo_main/pkg/logger"
	"os"
	"testing"
)

var userHandlers UserHandler

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
	userHandlers.Log = logger.NewLogger(os.Stdout)
}

func TestLogin(t *testing.T) {
	t.Run("Login-OK", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)
		s := session.NewMockDelivery(ctrl)

		testInput := models.UserSignIn{
			Login:    testUser.Login,
			Password: testUser.Password,
		}

		m.EXPECT().
			Login(testInput).
			Return(testUser, nil)

		s.EXPECT().
			Create(testUser).
			Return(http.Cookie{
				Name:  "session_id",
				Value: "testValue",
			}, nil)

		userHandlers.UserUC = m
		userHandlers.SessionDelivery = s

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, false, models.User{}, "")

		apitest.New("Login-OK").
			Handler(middlewareMock).
			Method("Post").
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	//test auth
	t.Run("Login-auth", func(t *testing.T) {

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, true, testUser, "")

		apitest.New("Login-auth").
			Handler(middlewareMock).
			Method("Post").
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	//test on login error
	t.Run("Login-UseCaseError", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, false, models.User{}, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		testInput := models.UserSignIn{
			Login:    testUser.Login,
			Password: testUser.Password,
		}
		testError := errors.New("test error")

		m.EXPECT().
			Login(testInput).
			Return(models.User{}, testError)

		userHandlers.UserUC = m

		apitest.New("Login-UseCaseError").
			Handler(middlewareMock).
			Method("Post").
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Login-SessionError", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, false, models.User{}, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)
		s := session.NewMockDelivery(ctrl)

		testInput := models.UserSignIn{
			Login:    testUser.Login,
			Password: testUser.Password,
		}

		testError := errors.New("test error")

		m.EXPECT().
			Login(testInput).
			Return(testUser, nil)

		s.EXPECT().
			Create(testUser).
			Return(http.Cookie{}, testError)

		userHandlers.UserUC = m
		userHandlers.SessionDelivery = s

		apitest.New("Login-SessionError").
			Handler(middlewareMock).
			Method("Post").
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("Login-FailedToParseJSON", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, false, models.User{}, "")

		apitest.New("Login-FailedToParseJSON").
			Handler(middlewareMock).
			Method("Post").
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "sdmfasd%so23}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestCreate(t *testing.T) {
	t.Run("Create-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)
		s := session.NewMockDelivery(ctrl)

		inputData := testUser
		inputData.Id = ""
		inputData.Image = ""

		m.EXPECT().
			Create(inputData).
			Return(user.NO, nil)

		s.EXPECT().
			Create(inputData).
			Return(http.Cookie{
				Name:  "session_id",
				Value: "testValue",
			}, nil)

		userHandlers.UserUC = m
		userHandlers.SessionDelivery = s

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false, models.User{}, "")

		apitest.New("Create-OK").
			Handler(middlewareMock).
			Method("Post").
			URL("/signup").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s"}`,
				testUser.Login,
				testUser.Password,
				testUser.Email,
				testUser.Sex,
				testUser.Name,
			)).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})
	t.Run("Create-ErrorJSON", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false, models.User{}, "")

		apitest.New("Create-ErrorJSON").
			Handler(middlewareMock).
			Method("Post").
			URL("/signup").
			Body(`{"loginsdf":sdf%sSkd}`).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Create-auth", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, true, testUser, "")

		apitest.New("Create-auth").
			Handler(middlewareMock).
			Method("Post").
			URL("/signup").
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})
	t.Run("Create-UseCaseError", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false, models.User{}, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		testError := errors.New("test error")

		inputData := testUser
		inputData.Id = ""
		inputData.Image = ""

		m.EXPECT().
			Create(inputData).
			Return(user.NO, testError)

		userHandlers.UserUC = m

		apitest.New("Create-UseCaseError").
			Handler(middlewareMock).
			Method("Post").
			URL("/signup").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s"}`,
				testUser.Login,
				testUser.Password,
				testUser.Email,
				testUser.Sex,
				testUser.Name,
			)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Create-SessionError", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false, models.User{}, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)
		s := session.NewMockDelivery(ctrl)

		inputData := testUser
		inputData.Id = ""
		inputData.Image = ""

		testError := errors.New("test error")

		m.EXPECT().
			Create(inputData).
			Return(user.NO, nil)

		s.EXPECT().
			Create(inputData).
			Return(http.Cookie{}, testError)

		userHandlers.UserUC = m
		userHandlers.SessionDelivery = s

		apitest.New("Create-SessionError").
			Handler(middlewareMock).
			Method("Post").
			URL("/signup").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s"}`,
				testUser.Login,
				testUser.Password,
				testUser.Email,
				testUser.Sex,
				testUser.Name,
			)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("Create-EmailExists", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false, models.User{}, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		inputData := testUser
		inputData.Id = ""
		inputData.Image = ""

		m.EXPECT().
			Create(inputData).
			Return(user.EMAIL, nil)

		userHandlers.UserUC = m

		apitest.New("Create-EmailExists").
			Handler(middlewareMock).
			Method("Post").
			URL("/signup").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s"}`,
				testUser.Login,
				testUser.Password,
				testUser.Email,
				testUser.Sex,
				testUser.Name,
			)).
			Expect(t).
			Body(`{"login_exists":false, "email_exists":true}`).
			Status(http.StatusConflict).
			End()
	})

	t.Run("Create-LoginExists", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false, models.User{}, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		inputData := testUser
		inputData.Id = ""
		inputData.Image = ""

		m.EXPECT().
			Create(inputData).
			Return(user.LOGIN, nil)

		userHandlers.UserUC = m

		apitest.New("Create-LoginExists").
			Handler(middlewareMock).
			Method("Post").
			URL("/signup").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s"}`,
				testUser.Login,
				testUser.Password,
				testUser.Email,
				testUser.Sex,
				testUser.Name,
			)).
			Expect(t).
			Body(`{"login_exists":true, "email_exists":false}`).
			Status(http.StatusConflict).
			End()
	})

	t.Run("Create-ExistsFull", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false, models.User{}, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		inputData := testUser
		inputData.Id = ""
		inputData.Image = ""

		m.EXPECT().
			Create(inputData).
			Return(user.FULL, nil)

		userHandlers.UserUC = m

		apitest.New("Create-ExistsFull").
			Handler(middlewareMock).
			Method("Post").
			URL("/signup").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s"}`,
				testUser.Login,
				testUser.Password,
				testUser.Email,
				testUser.Sex,
				testUser.Name,
			)).
			Expect(t).
			Body(`{"login_exists":true, "email_exists":true}`).
			Status(http.StatusConflict).
			End()
	})
}

func TestSelfProfile(t *testing.T) {
	t.Run("SelfProfile-OK", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.SelfProfile, true, testUser, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		profile := models.User{
			Id:    testUser.Id,
			Name:  testUser.Name,
			Login: testUser.Login,
			Sex:   testUser.Sex,
			Image: testUser.Image,
			Email: testUser.Email,
		}

		m.EXPECT().
			GetOutputUserData(testUser).
			Return(profile)

		userHandlers.UserUC = m

		apitest.New("SelfProfile-OK").
			Handler(middlewareMock).
			Method("Get").
			URL("/profile/me").
			Expect(t).
			Body(fmt.Sprintf(`{"id":"%s", "name":"%s", "login":"%s", "sex":"%s", "image":"%s", "email":"%s"}`,
				profile.Id,
				profile.Name,
				profile.Login,
				profile.Sex,
				profile.Image,
				profile.Email,
			)).
			Status(http.StatusOK).
			End()
	})
}

func TestGetProfile(t *testing.T) {
	t.Run("GetProfile-OK", func(t *testing.T) {
		middlewareMock := middleware.SetMuxVars(userHandlers.Profile, "profile", testUser.Login)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		profile := models.User{
			Id:    testUser.Id,
			Name:  testUser.Name,
			Login: testUser.Login,
			Sex:   testUser.Sex,
			Image: testUser.Image,
			Email: testUser.Email,
		}

		m.EXPECT().
			GetProfileByLogin(testUser.Login).
			Return(profile, nil)

		userHandlers.UserUC = m

		apitest.New("GetProfile-OK").
			Handler(middlewareMock).
			Method("Get").
			URL("/profile/keklol").
			Expect(t).
			Body(fmt.Sprintf(`{"id":"%s", "name":"%s", "login":"%s", "sex":"%s", "image":"%s", "email":"%s"}`,
				profile.Id,
				profile.Name,
				profile.Login,
				profile.Sex,
				profile.Image,
				profile.Email,
			)).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetProfile-UseCaseError", func(t *testing.T) {
		middlewareMock := middleware.SetMuxVars(userHandlers.Profile, "profile", testUser.Login)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		testError := errors.New("testError")

		m.EXPECT().
			GetProfileByLogin(testUser.Login).
			Return(models.User{}, testError)

		userHandlers.UserUC = m

		apitest.New("GetProfile-UseCaseError").
			Handler(middlewareMock).
			Method("Get").
			URL("/profile/keklol").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetProfile-NoMuxVars", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		apitest.New("GetProfile-NoMuxVars").
			Handler(http.HandlerFunc(userHandlers.Profile)).
			Method("Get").
			URL("/profile/keklol").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestLogout(t *testing.T) {
	t.Run("Logout-OK", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Logout, true, testUser, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := session.NewMockDelivery(ctrl)

		cookieValue := "89273894cjawiue983nc29384c2n23cu9"

		s.EXPECT().
			Delete(cookieValue).
			Return(nil)

		userHandlers.SessionDelivery = s

		apitest.New("Logout-OK").
			Handler(middlewareMock).
			Method("Delete").
			URL("/logout").
			Cookie("session_id", cookieValue).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Logout-NotAuth", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Logout, false, models.User{}, "")

		apitest.New("Logout-NotAuth").
			Handler(middlewareMock).
			Method("Delete").
			URL("/logout").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Logout-SessionDeleteError", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Logout, true, testUser, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := session.NewMockDelivery(ctrl)

		cookieValue := "89273894cjawiue983nc29384c2n23cu9"
		testError := errors.New("test error")

		s.EXPECT().
			Delete(cookieValue).
			Return(testError)

		userHandlers.SessionDelivery = s

		apitest.New("Logout-SessionDeleteError").
			Handler(middlewareMock).
			Method("Delete").
			URL("/logout").
			Cookie("session_id", cookieValue).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update-OK", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Update, true, testUser, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		inputData := models.UserSettings{
			NewPassword: "23hf9823d9123d",
			User:        testUser,
		}
		inputData.Id = ""
		inputData.Image = ""

		m.EXPECT().
			Update(testUser, inputData).
			Return(user.NO, nil)

		userHandlers.UserUC = m

		apitest.New("Update-OK").
			Handler(middlewareMock).
			Method("Put").
			URL("/users/settings").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s", "new_password":"%s"}`,
				inputData.Login,
				inputData.Password,
				inputData.Email,
				inputData.Sex,
				inputData.Name,
				inputData.NewPassword,
			)).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Update-BadInput", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Update, true, testUser, "")

		apitest.New("Update-BadInput").
			Handler(middlewareMock).
			Method("Put").
			URL("/users/settings").
			Body(`{askd:"oaskdoepkwr23r:,kapod}"`).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Update-EmailExists", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Update, true, testUser, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		inputData := models.UserSettings{
			NewPassword: "23hf9823d9123d",
			User:        testUser,
		}
		inputData.Id = ""
		inputData.Email = "exists@email.true"
		inputData.Image = ""

		m.EXPECT().
			Update(testUser, inputData).
			Return(user.EMAIL, nil)

		userHandlers.UserUC = m

		apitest.New("Update-EmailExists").
			Handler(middlewareMock).
			Method("Put").
			URL("/users/settings").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s", "new_password":"%s"}`,
				inputData.Login,
				inputData.Password,
				inputData.Email,
				inputData.Sex,
				inputData.Name,
				inputData.NewPassword,
			)).
			Expect(t).
			Status(http.StatusConflict).
			End()
	})

	t.Run("Update-UseCaseError", func(t *testing.T) {
		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Update, true, testUser, "")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		inputData := models.UserSettings{
			NewPassword: "23hf9823d9123d",
			User:        testUser,
		}
		inputData.Id = ""
		inputData.Image = ""

		testError := errors.New("test error")

		m.EXPECT().
			Update(testUser, inputData).
			Return(user.NO, testError)

		userHandlers.UserUC = m

		apitest.New("Update-UseCaseError").
			Handler(middlewareMock).
			Method("Put").
			URL("/users/settings").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s", "new_password":"%s"}`,
				inputData.Login,
				inputData.Password,
				inputData.Email,
				inputData.Sex,
				inputData.Name,
				inputData.NewPassword,
			)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Update-NoCsrf", func(t *testing.T) {
		apitest.New("Update-NoCsrf").
			Handler(http.HandlerFunc(userHandlers.Update)).
			Method("Put").
			URL("/users/settings").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})
}

func TestGetUserStat(t *testing.T) {
	t.Run("GetUserStat-OK", func(t *testing.T) {
		middlewareMock := middleware.SetMuxVars(userHandlers.GetUserStat, "id", testUser.Id)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		stat := models.UserStat{
			UserId:    testUser.Id,
			Tracks:    12,
			Albums:    4234,
			Playlists: 5345,
			Artists:   123,
		}

		statMarshal, err := json.Marshal(stat)
		assert.NoError(t, err)

		m.EXPECT().
			GetUserStat(testUser.Id).
			Return(stat, nil)

		userHandlers.UserUC = m

		apitest.New("GetProfile-OK").
			Handler(middlewareMock).
			Method("Get").
			URL("/stat/keklol").
			Expect(t).
			Body(string(statMarshal)).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserStat-UseCaseError", func(t *testing.T) {
		middlewareMock := middleware.SetMuxVars(userHandlers.GetUserStat, "id", testUser.Id)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		testError := errors.New("testError")

		m.EXPECT().
			GetUserStat(testUser.Id).
			Return(models.UserStat{}, testError)

		userHandlers.UserUC = m

		apitest.New("GetProfile-UseCaseError").
			Handler(middlewareMock).
			Method("Get").
			URL("/stat/keklol").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetUserStat-NoMuxVars", func(t *testing.T) {
		apitest.New("GetProfile-NoMuxVars").
			Handler(http.HandlerFunc(userHandlers.GetUserStat)).
			Method("Get").
			URL("/stat/keklol").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestGetCSRF(t *testing.T) {
	t.Run("GetCSRF-OK", func(t *testing.T) {
		sessionId := "asdasdwer6545"
		handler := middleware.AuthMiddlewareMock(userHandlers.GetCSRF, true, models.User{}, sessionId)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := csrf.NewMockUseCase(ctrl)

		token := "jieimrc2mu49muc9824cu924cum9"

		m.EXPECT().
			Create(sessionId, gomock.Any()).
			Return(token, nil)

		userHandlers.CSRF = m

		apitest.New("GetCSRF-OK").
			Handler(handler).
			Method("Get").
			URL("/users/token").
			Expect(t).
			Header("Csrf-Token", token).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserStat-UseCaseError", func(t *testing.T) {
		sessionId := "asdasdwer6545"
		handler := middleware.AuthMiddlewareMock(userHandlers.GetCSRF, true, models.User{}, sessionId)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := csrf.NewMockUseCase(ctrl)

		testError := errors.New("testError")

		m.EXPECT().
			Create(sessionId, gomock.Any()).
			Return("", testError)

		userHandlers.CSRF = m

		apitest.New("GetProfile-UseCaseError").
			Handler(handler).
			Method("Get").
			URL("/users/token").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetUserStat-NoSession_id", func(t *testing.T) {
		apitest.New("GetProfile-NoSession_id").
			Handler(http.HandlerFunc(userHandlers.GetCSRF)).
			Method("Get").
			URL("/users/token").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}
