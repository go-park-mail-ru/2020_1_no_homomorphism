package session

import (
	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
)

type Repository interface {
	Create(user *models.User) (uuid.UUID, error)
	Delete(sessionID uuid.UUID) error
	GetUserBySessionID(sessionID uuid.UUID) (*models.User, error)
}