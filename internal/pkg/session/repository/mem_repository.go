package repository

import (
	"errors"

	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
)

type SessionRepository struct {
	sessions map[uuid.UUID]uuid.UUID
}

func InitRepository() *SessionRepository {
	return &SessionRepository{
		sessions: make(map[uuid.UUID]uuid.UUID),
	}
}

func (sr *SessionRepository) Create(user models.User) uuid.UUID {
	newUUID := uuid.NewV4()
	sr.sessions[newUUID] = user.Id
	return newUUID
}

func (sr *SessionRepository) Delete(sessionID uuid.UUID) error {
	_, ok := sr.sessions[sessionID]
	if !ok {
		return errors.New("can't delete session because it does not exists : " + sessionID.String())
	}
	delete(sr.sessions, sessionID)
	return nil
}

