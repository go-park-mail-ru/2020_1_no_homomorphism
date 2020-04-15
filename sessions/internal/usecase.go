package session

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type UseCase interface {
	Create(login string, expires time.Duration) (uuid.UUID, error)
	Delete(sessionID uuid.UUID) error
	Check(sessionID uuid.UUID) (string, error)
}
