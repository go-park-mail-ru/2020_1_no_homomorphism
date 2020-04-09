package session

import (
	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
	"time"
)

type UseCase interface {
	Create(user models.User, expires time.Duration) (uuid.UUID, error)
	Delete(sessionID uuid.UUID) error
	GetLoginBySessionID(sessionID uuid.UUID) (string, error)
}
