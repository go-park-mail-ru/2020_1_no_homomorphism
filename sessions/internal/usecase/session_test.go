package usecase

import (
	"errors"
	session "github.com/2020_1_no_homomorphism/no_homo_sessions/internal"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	testError := errors.New("something go wrong")

	t.Run("Create-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockRepository(ctrl)

		testLogin := "testLogin"
		var expires uint64 = 237498234

		m.
			EXPECT().
			Create(gomock.Any(), testLogin, expires).
			Return(nil)

		useCase := SessionUseCase{
			Repository: m,
		}

		uuid, err := useCase.Create(testLogin, expires)
		assert.NoError(t, err)
		assert.NotNil(t, uuid)
	})

	t.Run("Create-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockRepository(ctrl)

		testLogin := "testLogin"
		var expires uint64 = 237498234

		m.
			EXPECT().
			Create(gomock.Any(), testLogin, expires).
			Return(testError)

		useCase := SessionUseCase{
			Repository: m,
		}

		_, err := useCase.Create(testLogin, expires)
		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	testError := errors.New("something go wrong")

	t.Run("Delete-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockRepository(ctrl)

		id := uuid.NewV4()

		m.
			EXPECT().
			GetLoginBySessionID("sessions:" + id.String()).
			Return("", nil)
		m.
			EXPECT().
			Delete("sessions:" + id.String()).
			Return(nil)

		useCase := SessionUseCase{
			Repository: m,
		}

		err := useCase.Delete(id)
		assert.NoError(t, err)
	})

	t.Run("Delete-GetLoginError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockRepository(ctrl)

		id := uuid.NewV4()

		m.
			EXPECT().
			GetLoginBySessionID("sessions:" + id.String()).
			Return("", testError)

		useCase := SessionUseCase{
			Repository: m,
		}

		err := useCase.Delete(id)
		assert.Error(t, err)
	})
	t.Run("Delete-DeleteError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockRepository(ctrl)

		id := uuid.NewV4()

		m.
			EXPECT().
			GetLoginBySessionID("sessions:" + id.String()).
			Return("", nil)
		m.
			EXPECT().
			Delete("sessions:" + id.String()).
			Return(testError)

		useCase := SessionUseCase{
			Repository: m,
		}

		err := useCase.Delete(id)
		assert.Error(t, err)
	})
}
