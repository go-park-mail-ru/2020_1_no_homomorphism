package session

import (
	uuid "github.com/satori/go.uuid"
)

type UseCase interface {
	Create(login string, expires uint64) (uuid.UUID, error)
	Delete(sessionID uuid.UUID) error
	Check(sessionID uuid.UUID) (string, error)
}
