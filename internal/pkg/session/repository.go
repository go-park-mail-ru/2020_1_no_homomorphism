package session

import (
	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
)

type Repository interface {
	Create(user *models.User) uuid.UUID
	Delete(sessionID uuid.UUID)
	GetUserBySessionID(sessionID uuid.UUID) (*models.User, error)
}
