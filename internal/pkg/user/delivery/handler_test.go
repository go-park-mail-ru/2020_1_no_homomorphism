package delivery

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"net/http"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/session"
	"no_homomorphism/internal/pkg/user"
	"no_homomorphism/pkg/logger"
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

func authMiddlewareMock(next http.HandlerFunc, auth bool, user models.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "isAuth", auth)
		ctx = context.WithValue(ctx, "user", user)
		next(w, r.WithContext(ctx))
	})
}

func TestCheckAuth(t *testing.T) {
	falseAuthPreHandle := authMiddlewareMock(userHandlers.CheckAuth, true, models.User{})
	trueAuthPreHandle := authMiddlewareMock(userHandlers.CheckAuth, false, models.User{})

	apitest.New("CheckAuth-true").
		Handler(falseAuthPreHandle).
		Method("Get").
		Cookie("session_id", "randomSessionIdValueForTesting").
		URL("/user").
		Expect(t).
		Status(200).
		End()

	apitest.New("CheckAuth-false").
		Handler(trueAuthPreHandle).
		Method("Get").
		Cookie("session_id", "randomSessionIdValueForTesting").
		URL("/user").
		Expect(t).
		Status(401).
		End()
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

		middlewareMock := authMiddlewareMock(userHandlers.Login, false, models.User{})

		apitest.New("Login-OK").
			Handler(middlewareMock).
			Method("Post").
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(200).
			End()
	})
	//test auth
	t.Run("Login-auth", func(t *testing.T) {

		middlewareMock := authMiddlewareMock(userHandlers.Login, true, testUser)

		apitest.New("Login-auth").
			Handler(middlewareMock).
			Method("Post").
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(403).
			End()
	})

	//test on login error
	t.Run("Login-UseCaseError", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Login, false, models.User{})

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
			Status(400).
			End()
	})

	t.Run("Login-SessionError", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Login, false, models.User{})

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
			Status(500).
			End()
	})

	t.Run("Login-FailedToParseJSON", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Login, false, models.User{})

		apitest.New("Login-FailedToParseJSON").
			Handler(middlewareMock).
			Method("Post").
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "sdmfasd%so23}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(400).
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

		middlewareMock := authMiddlewareMock(userHandlers.Create, false, models.User{})

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
			Status(201).
			End()
	})
	//test auth
	t.Run("Create-ErrorJSON", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Create, false, models.User{})

		apitest.New("Create-ErrorJSON").
			Handler(middlewareMock).
			Method("Post").
			URL("/signup").
			Body(`{"loginsdf":sdf%sSkd}`).
			Expect(t).
			Status(400).
			End()
	})

	t.Run("Create-auth", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Create, true, testUser)

		apitest.New("Create-auth").
			Handler(middlewareMock).
			Method("Post").
			URL("/signup").
			Expect(t).
			Status(403).
			End()
	})

	//test on login error
	t.Run("Create-UseCaseError", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Create, false, models.User{})

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
			Status(400).
			End()
	})

	t.Run("Create-SessionError", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Create, false, models.User{})

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
			Status(500).
			End()
	})

	t.Run("Login-EmailExists", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Create, false, models.User{})

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
			Status(409).
			End()
	})

	t.Run("Login-LoginExists", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Create, false, models.User{})

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
			Status(409).
			End()
	})

	t.Run("Login-ExistsFull", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Create, false, models.User{})

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
			Status(409).
			End()
	})
}

func TestSelfProfile(t *testing.T) {
	t.Run("SelfProfile-OK", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.SelfProfile, true, testUser)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockUseCase(ctrl)

		profile := models.Profile{
			Name:  testUser.Name,
			Login: testUser.Login,
			Sex:   testUser.Sex,
			Image: testUser.Image,
			Email: testUser.Email,
		}

		m.EXPECT().
			GetProfileByUser(testUser).
			Return(profile)

		userHandlers.UserUC = m

		apitest.New("SelfProfile-OK").
			Handler(middlewareMock).
			Method("Get").
			URL("/profile/me").
			Expect(t).
			Body(fmt.Sprintf(`{"name":"%s", "login":"%s", "sex":"%s", "image":"%s", "email":"%s"}`,
				profile.Name,
				profile.Login,
				profile.Sex,
				profile.Image,
				profile.Email,
			)).
			Status(200).
			End()
	})

	t.Run("SelfProfile-NotAuth", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.SelfProfile, false, models.User{})

		apitest.New("SelfProfile-NotAuth").
			Handler(middlewareMock).
			Method("Get").
			URL("/profile/me").
			Expect(t).
			Status(401).
			End()
	})
}

func TestLogout(t *testing.T) {
	t.Run("Logout-OK", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Logout, true, testUser)

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
			Status(200).
			End()
	})

	t.Run("Logout-NotAuth", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Logout, false, testUser)

		apitest.New("Logout-NotAuth").
			Handler(middlewareMock).
			Method("Delete").
			URL("/logout").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Logout-SessionDeleteError", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Logout, true, testUser)

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
		middlewareMock := authMiddlewareMock(userHandlers.Update, true, testUser)

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
			URL("/profile/settings").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s", "new_password":"%s"}`,
				inputData.Login,
				inputData.Password,
				inputData.Email,
				inputData.Sex,
				inputData.Name,
				inputData.NewPassword,
			)).
			Expect(t).
			Body(`{"email_exists":false}`).
			Status(http.StatusOK).
			End()
	})

	t.Run("Update-BadInput", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Update, true, testUser)

		apitest.New("Update-BadInput").
			Handler(middlewareMock).
			Method("Put").
			URL("/profile/settings").
			Body(`{askd:"oaskdoepkwr23r:,kapod}"`).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Update-EmailExists", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Update, true, testUser)

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
			URL("/profile/settings").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s", "new_password":"%s"}`,
				inputData.Login,
				inputData.Password,
				inputData.Email,
				inputData.Sex,
				inputData.Name,
				inputData.NewPassword,
			)).
			Expect(t).
			Body(`{"email_exists":true}`).
			Status(http.StatusOK).
			End()
	})

	t.Run("Update-Unauthorized", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Update, false, models.User{})

		inputData := models.UserSettings{
			NewPassword: "23hf9823d9123d",
			User:        testUser,
		}
		inputData.Id = ""
		inputData.Email = "exists@email.true"
		inputData.Image = ""

		apitest.New("Update-Unauthorized").
			Handler(middlewareMock).
			Method("Put").
			URL("/profile/settings").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s", "email":"%s", "sex":"%s", "name":"%s", "new_password":"%s"}`,
				inputData.Login,
				inputData.Password,
				inputData.Email,
				inputData.Sex,
				inputData.Name,
				inputData.NewPassword,
			)).
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Update-UseCaseError", func(t *testing.T) {
		middlewareMock := authMiddlewareMock(userHandlers.Update, true, testUser)

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
			URL("/profile/settings").
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

}
