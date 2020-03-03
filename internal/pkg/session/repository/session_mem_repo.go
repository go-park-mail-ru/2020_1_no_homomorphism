package repository

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
)

type SessionRepository struct {
	Sessions map[uuid.UUID]*models.User
}

func (sr *SessionRepository) Create(user *models.User) uuid.UUID {
	newUUID := uuid.NewV4()
	sr.Sessions[newUUID] = user
	return newUUID
}

func (sr *SessionRepository) Delete(sessionID uuid.UUID) {
	delete(sr.Sessions, sessionID)
}

func (sr *SessionRepository) GetUserBySessionID(sessionID uuid.UUID) (*models.User, error) {
	user, ok := sr.Sessions[sessionID]
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
