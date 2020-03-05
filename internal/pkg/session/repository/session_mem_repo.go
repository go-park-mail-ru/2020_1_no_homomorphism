package repository

import (
	"errors"
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
)

type SessionRepository struct {
	Sessions map[uuid.UUID]*models.User
	mutex    *sync.RWMutex
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		Sessions: make(map[uuid.UUID]*models.User),
		mutex:    &sync.RWMutex{},
	}
}

func (sr *SessionRepository) Create(user *models.User) (uuid.UUID, error) {
	newUUID := uuid.NewV4()
	if _, ok := sr.Sessions[newUUID]; ok {
		return uuid.Nil, errors.New("can't create session with this uuid")
	}
	sr.mutex.Lock()
	sr.Sessions[newUUID] = user
	sr.mutex.Unlock()
	return newUUID, nil
}

func (sr *SessionRepository) Delete(sessionID uuid.UUID) {
	sr.mutex.Lock()
	delete(sr.Sessions, sessionID)
	sr.mutex.Unlock()
}

func (sr *SessionRepository) GetUserBySessionID(sessionID uuid.UUID) (*models.User, error) {
	sr.mutex.Lock()
	user, ok := sr.Sessions[sessionID]
	sr.mutex.Unlock()
	if !ok {
		return nil, errors.New("session does not exists")
	}
	return user, nil
}

func (sr *SessionRepository) PrintSessionList() {
	fmt.Println("[SESSIONSS LIST]")
	for _, r := range sr.Sessions {
		fmt.Println(r)
	}
}
