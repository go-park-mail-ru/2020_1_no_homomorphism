package session

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
	"no_homomorphism/internal/pkg/models"
)

type Delivery interface {
	Create(user models.User) (http.Cookie, error)
	Delete(sessionID string) error
	GetLoginBySessionID(sessionID uuid.UUID) (string, error)
}
