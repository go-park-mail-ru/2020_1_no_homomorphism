package session

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
)

type UseCase interface {
	Create(user *models.User) (*http.Cookie, error)
	Delete(sessionID uuid.UUID) error
	GetLoginBySessionID(sessionID uuid.UUID) (string, error)
}
