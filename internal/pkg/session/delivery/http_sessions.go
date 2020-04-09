package delivery

import (
	"fmt"
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
		Path:     "/",
		Expires:  time.Now().Add(uc.ExpireTime),
	}, nil
}

func (uc *SessionDelivery) Delete(sessionID string) error {
	sid, err := uuid.FromString(sessionID)
	if err != nil {
		return fmt.Errorf("can't parse uuid from string")
	}
	return uc.UseCase.Delete(sid)
}

func (uc *SessionDelivery) GetLoginBySessionID(sessionID string) (string, error) {
	sid, err := uuid.FromString(sessionID)
	if err != nil {
		return "", fmt.Errorf("can't parse uuid from string")
	}
	return uc.UseCase.GetLoginBySessionID(sid)
}
