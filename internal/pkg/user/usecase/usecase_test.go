package usecase

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/user"
	"testing"
)

var testUser = models.User{
	Id:       "1234",
	Password: "76453647fvd",
	Name:     "TestName",
	Login:    "nnnagibator",
	Sex:      "Man",
	Image:    "/static/avatar/default.png",
	Email:    "klsJDLKfj@mail.ru",
}

func TestCreate(t *testing.T) {
	testError := errors.New("something go wrong")

	t.Run("Create-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)

		m.
			EXPECT().
			CheckIfExists(gomock.Eq(testUser.Login), gomock.Eq(testUser.Email)).
			Return(false, false, nil)
		m.
			EXPECT().
			Create(gomock.Eq(testUser)).
			Return(nil)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		exists, err := useCase.Create(testUser)
		assert.NoError(t, err)
		assert.Equal(t, exists, user.NO)
	})

	t.Run("Create-UserLoginExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)
		m.
			EXPECT().
			CheckIfExists(gomock.Eq(testUser.Login), gomock.Eq(testUser.Email)).
			Return(true, false, nil)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}
		exists, err := useCase.Create(testUser)

		assert.NoError(t, err)
		assert.Equal(t, exists, user.LOGIN)
	})

	t.Run("Create-UserEmailExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)
		m.
			EXPECT().
			CheckIfExists(gomock.Eq(testUser.Login), gomock.Eq(testUser.Email)).
			Return(false, true, nil)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}
		exists, err := useCase.Create(testUser)

		assert.NoError(t, err)
		assert.Equal(t, exists, user.EMAIL)
	})

	t.Run("Create-UserFullExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)
		m.
			EXPECT().
			CheckIfExists(gomock.Eq(testUser.Login), gomock.Eq(testUser.Email)).
			Return(true, true, nil)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}
		exists, err := useCase.Create(testUser)

		assert.NoError(t, err)
		assert.Equal(t, exists, user.FULL)
	})

	t.Run("Create-CheckUserError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)

		m.
			EXPECT().
			CheckIfExists(gomock.Eq(testUser.Login), gomock.Eq(testUser.Email)).
			Return(true, true, testError)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}
		_, err := useCase.Create(testUser)

		assert.Error(t, err)
		assert.Equal(t, err, testError)
	})

	t.Run("Create-OnCreateError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)

		m.
			EXPECT().
			CheckIfExists(gomock.Eq(testUser.Login), gomock.Eq(testUser.Email)).
			Return(false, false, nil)
		m.
			EXPECT().
			Create(gomock.Eq(testUser)).
			Return(testError)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}
		_, err := useCase.Create(testUser)

		assert.Error(t, err)
		assert.Equal(t, err, testError)
	})
}

func TestUpdate(t *testing.T) {
	testInput := models.UserSettings{
		NewPassword: "01238401ksjdf20934",
		User:        testUser,
	}

	testError := errors.New("something go wrong")

	t.Run("Update-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)

		m.
			EXPECT().
			Update(testUser, testInput).
			Return(nil)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		exists, err := useCase.Update(testUser, testInput)
		assert.NoError(t, err)
		assert.Equal(t, exists, user.NO)
	})

	t.Run("Update-WithoutNewPass", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)

		testInput2 := testInput
		testInput2.NewPassword = ""

		m.
			EXPECT().
			Update(testUser, testInput2).
			Return(nil)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		exists, err := useCase.Update(testUser, testInput2)
		assert.NoError(t, err)
		assert.Equal(t, exists, user.NO)
	})

	t.Run("Update-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)

		m.
			EXPECT().
			Update(testUser, testInput).
			Return(testError)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		_, err := useCase.Update(testUser, testInput)
		assert.Error(t, err)
		assert.Equal(t, err, testError)
	})

	t.Run("Update-Exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)

		testInput.Email = "newemail@mail.ru"

		m.
			EXPECT().
			CheckIfExists(gomock.Eq(""), gomock.Eq(testInput.Email)).
			Return(false, true, nil)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		exists, err := useCase.Update(testUser, testInput)
		assert.NoError(t, err)
		assert.Equal(t, exists, user.EMAIL)
	})
	t.Run("Update-ExistsError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)

		m.
			EXPECT().
			CheckIfExists(gomock.Eq(""), gomock.Eq(testInput.Email)).
			Return(false, true, testError)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		_, err := useCase.Update(testUser, testInput)
		assert.Error(t, err)
		assert.Equal(t, err, fmt.Errorf("failed to check email existing: %v", testError))
	})
}

func TestGetProfileByLogin(t *testing.T) {
	testProfile := models.Profile{
		Name:  testUser.Name,
		Login: testUser.Login,
		Sex:   testUser.Sex,
		Image: testUser.Image,
		Email: testUser.Email,
	}

	t.Run("GetProfileByLogin-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)
		m.
			EXPECT().
			GetUserByLogin(testUser.Login).
			Return(testUser, nil)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		userData, err := useCase.GetProfileByLogin(testUser.Login)
		assert.NoError(t, err)
		assert.Equal(t, userData, testProfile)
	})

	t.Run("GetProfileByLogin-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		testError := errors.New("testError")
		m := user.NewMockRepository(ctrl)
		m.
			EXPECT().
			GetUserByLogin(testUser.Login).
			Return(testUser, testError)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		_, err := useCase.GetProfileByLogin(testUser.Login)
		assert.Error(t, err)
		assert.Equal(t, err, testError)
	})
}

func TestGetUserByLogin(t *testing.T) {
	t.Run("GetProfileByLogin-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)
		m.
			EXPECT().
			GetUserByLogin(testUser.Login).
			Return(testUser, nil)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		userData, err := useCase.GetUserByLogin(testUser.Login)
		assert.NoError(t, err)
		assert.Equal(t, userData, testUser)
	})

	t.Run("GetProfileByLogin-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		testError := errors.New("testError")
		m := user.NewMockRepository(ctrl)
		m.
			EXPECT().
			GetUserByLogin(testUser.Login).
			Return(testUser, testError)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		_, err := useCase.GetUserByLogin(testUser.Login)
		assert.Error(t, err)
		assert.Equal(t, err, testError)
	})
}

func TestLogin(t *testing.T) {
	testError := errors.New("some test error")

	testInput := models.UserSignIn{
		Login:    testUser.Login,
		Password: testUser.Password,
	}

	t.Run("GetLogin-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)
		m.
			EXPECT().
			GetUserByLogin(testInput.Login).
			Return(testUser, nil)
		m.
			EXPECT().
			CheckUserPassword(testUser.Password, testInput.Password).
			Return(nil)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		userData, err := useCase.Login(testInput)
		assert.NoError(t, err)
		assert.Equal(t, userData, testUser)
	})

	t.Run("GetLogin-NoSuchUser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)
		m.
			EXPECT().
			GetUserByLogin(testInput.Login).
			Return(models.User{}, testError)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		_, err := useCase.Login(testInput)
		assert.Error(t, err)
	})

	t.Run("GetLogin-WrongPassword", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepository(ctrl)
		m.
			EXPECT().
			GetUserByLogin(testInput.Login).
			Return(testUser, nil)
		m.
			EXPECT().
			CheckUserPassword(testUser.Password, testInput.Password).
			Return(testError)

		useCase := UserUseCase{
			Repository: m,
			AvatarDir:  "/test",
		}

		_, err := useCase.Login(testInput)
		assert.Error(t, err)
	})
}
