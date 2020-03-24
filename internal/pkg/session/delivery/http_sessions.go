package delivery

import (
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/session"
)

type SessionDelivery struct {
	UseCase    session.UseCase
	ExpireTime time.Duration
}

func (uc *SessionDelivery) Create(user models.User) (http.Cookie, error) {
	sid, err := uc.UseCase.Create(user, uc.ExpireTime)
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     "session_id",
		Value:    sid.String(),
		HttpOnly: true,
		Expires:  time.Now().Add(uc.ExpireTime),
	}, nil
}

func (uc *SessionDelivery) Delete(sessionID uuid.UUID) error {
	return uc.UseCase.Delete(sessionID)
}

func (uc *SessionDelivery) GetLoginBySessionID(sessionID uuid.UUID) (string, error) {
	return uc.UseCase.GetLoginBySessionID(sessionID)
}
