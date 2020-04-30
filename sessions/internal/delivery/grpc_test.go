package delivery

import (
	"context"
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

		m := session.NewMockUseCase(ctrl)

		in := &session.Session{Login: "testLogin"}

		delivery := NewSessionDelivery(m, 23525)

		uID := uuid.NewV4()

		m.
			EXPECT().
			Create(in.Login, delivery.ExpireTime).
			Return(uID, nil)

		sessID := &session.SessionID{ID: uID.String()}

		sessionID, err := delivery.Create(context.TODO(), in)
		assert.NoError(t, err)
		assert.Equal(t, sessID, sessionID)
	})
	t.Run("Create-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)

		in := &session.Session{Login: "testLogin"}

		delivery := NewSessionDelivery(m, 23525)

		m.
			EXPECT().
			Create(in.Login, delivery.ExpireTime).
			Return(uuid.UUID{}, testError)

		_, err := delivery.Create(context.TODO(), in)
		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	testError := errors.New("something go wrong")

	t.Run("Delete-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)

		delivery := NewSessionDelivery(m, 23525)

		uID := uuid.NewV4()
		sessID := &session.SessionID{ID: uID.String()}

		m.
			EXPECT().
			Delete(uuid.FromStringOrNil(sessID.ID)).
			Return(nil)

		_, err := delivery.Delete(context.TODO(), sessID)
		assert.NoError(t, err)
	})
	t.Run("Delete-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)

		delivery := NewSessionDelivery(m, 23525)

		uID := uuid.NewV4()
		sessID := &session.SessionID{ID: uID.String()}

		m.
			EXPECT().
			Delete(uuid.FromStringOrNil(sessID.ID)).
			Return(testError)

		_, err := delivery.Delete(context.TODO(), sessID)
		assert.Error(t, err)
	})
}

func TestCheck(t *testing.T) {
	testError := errors.New("something go wrong")

	t.Run("Check-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)

		delivery := NewSessionDelivery(m, 23525)

		uID := uuid.NewV4()
		sessID := &session.SessionID{ID: uID.String()}
		out := &session.Session{Login: "testLogin"}

		m.
			EXPECT().
			Check(uuid.FromStringOrNil(sessID.ID)).
			Return(out.Login, nil)

		res, err := delivery.Check(context.TODO(), sessID)
		assert.NoError(t, err)
		assert.Equal(t, out, res)
	})
	t.Run("Check-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)

		delivery := NewSessionDelivery(m, 23525)
		uID := uuid.NewV4()
		sessID := &session.SessionID{ID: uID.String()}

		m.
			EXPECT().
			Check(uuid.FromStringOrNil(sessID.ID)).
			Return("", testError)

		_, err := delivery.Check(context.TODO(), sessID)
		assert.Error(t, err)
	})
}
