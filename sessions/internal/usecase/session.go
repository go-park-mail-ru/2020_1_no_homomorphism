package usecase

import (
	"errors"
	session "github.com/2020_1_no_homomorphism/no_homo_sessions/internal"
	uuid "github.com/satori/go.uuid"
)

type SessionUseCase struct {
	Repository session.Repository
}

func addPrefix(id uuid.UUID) string {
	return "sessions:" + id.String()
}

func (uc *SessionUseCase) Create(login string, expires uint64) (uuid.UUID, error) {
	id := uuid.NewV4()
	sId := addPrefix(id)
	err := uc.Repository.Create(sId, login, expires)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (uc *SessionUseCase) Delete(sessionID uuid.UUID) error {
	sId := addPrefix(sessionID)
	_, err := uc.Repository.GetLoginBySessionID(sId)
	if err != nil {
		return errors.New("can't find session: " + sessionID.String() + " error:" + err.Error())
	}
	err = uc.Repository.Delete(sId)
	if err != nil {
		return errors.New("can't delete session: " + sessionID.String() + " error:" + err.Error())
	}
	return nil
}

func (uc *SessionUseCase) Check(sessionID uuid.UUID) (string, error) {
	return uc.Repository.GetLoginBySessionID(addPrefix(sessionID))
}
