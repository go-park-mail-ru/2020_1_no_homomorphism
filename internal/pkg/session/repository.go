package session

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Repository interface {
	Create(login string, expire time.Duration) (uuid.UUID, error)
	Delete(sessionID uuid.UUID) error
	GetLoginBySessionID(sessionID uuid.UUID) (string, error)
}
