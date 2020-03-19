package usecase

import (
	"errors"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/session"
)

type SessionUseCase struct {
	Repository session.Repository
	ExpireTime time.Duration
}

func (uc *SessionUseCase) Create(user *models.User) (*http.Cookie, error) {
	sid, err := uc.Repository.Create(user.Login, uc.ExpireTime)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:     "session_id",
		Value:    sid.String(),
		HttpOnly: true,
		Expires:  time.Now().Add(uc.ExpireTime),
	}, nil
}

func (uc *SessionUseCase) Delete(sessionID uuid.UUID) error {
	_, err := uc.Repository.GetLoginBySessionID(sessionID)
	if err != nil {
		return errors.New("can't find session: " + sessionID.String() + " error:" + err.Error())
	}
	err = uc.Repository.Delete(sessionID)
	if err != nil {
		return errors.New("can't delete session: " + sessionID.String() + " error:" + err.Error())

	}
	return nil
}

func (uc *SessionUseCase) GetLoginBySessionID(sessionID uuid.UUID) (string, error) {
	return uc.Repository.GetLoginBySessionID(sessionID)
}

// func (uc *SessionUseCase) PrintSessionList() {
// 	uc.Repository.PrintSessionList()
// }
