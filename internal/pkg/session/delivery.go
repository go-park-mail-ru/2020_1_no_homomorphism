package session

import (
	"net/http"
	"no_homomorphism/internal/pkg/models"
)

type Delivery interface {
	Create(user models.User) (http.Cookie, error)
	Delete(sessionID string) error
	GetLoginBySessionID(sessionID string) (string, error)
}
