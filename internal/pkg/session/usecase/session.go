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
}

func (uc *SessionUseCase) Create(user *models.User) (*http.Cookie, error) {
	sid, err := uc.Repository.Create(user)
	if err != nil {
		for i := 0; i < 5; i++ {
			sid, err = uc.Repository.Create(user)
			if err == nil {
				break
			}
		}
		return nil, errors.New("could not create new session and cookie")
	}
	return &http.Cookie{
		Name:     "session_id",
		Value:    sid.String(),
		HttpOnly: true,
		Expires:  time.Now().Add(24 * 30 * time.Hour),
	}, nil
}

func (uc *SessionUseCase) Delete(sessionID uuid.UUID) error {
	_, err := uc.Repository.GetUserBySessionID(sessionID)
	if err != nil {
		return errors.New("can't delete session because it does not exists : " + sessionID.String())
	}
	uc.Repository.Delete(sessionID)
	return nil
}

func (uc *SessionUseCase) GetUserBySessionID(sessionID uuid.UUID) (*models.User, error) {
	return uc.Repository.GetUserBySessionID(sessionID)
}

// func (uc *SessionUseCase) PrintSessionList() {
// 	uc.Repository.PrintSessionList()
// }
