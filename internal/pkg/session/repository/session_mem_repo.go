package repository

import (
	"errors"
	"sync"

	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
)

type SessionRepository struct {
	sessions map[uuid.UUID]*models.User
	mutex *sync.Mutex
}

func NewSessionRepository(mutex *sync.Mutex) *SessionRepository {
	return &SessionRepository{
		sessions: make(map[uuid.UUID]*models.User),
		mutex: mutex,
	}
}

func (sr *SessionRepository) Create(user *models.User) uuid.UUID {
	newUUID := uuid.NewV4()
	sr.sessions[newUUID] = user
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


func (sr *SessionRepository) GetUserBySessionID(sessionID uuid.UUID) (*models.User, error){
	user, ok := sr.sessions[sessionID]
	if !ok {
		return nil, errors.New("session does not exists")
	}
	return user, nil
}