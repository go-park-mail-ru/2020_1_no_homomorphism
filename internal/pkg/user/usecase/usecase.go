package usecase

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/user"
)

var mutex = &sync.Mutex{}

type UserUseCase struct {
	Repository user.Repository
}

func (uc *UserUseCase) Create(user *models.User) error {
	_, ok := uc.GetUserByLogin(user.Login)
	if ok == nil {
		return errors.New("user with this login is already exists")
	}
	return uc.Repository.Create(user)
}
func (uc *UserUseCase) Update(user *models.UserSettings) error {
	oldUser, err := uc.GetUserByLogin(user.Login)
	if err != nil {
		return errors.New("user with this login does not exists")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(oldUser.Password), []byte(user.Password)); err != nil {
		return errors.New("old password is wrong")
	}
	if user.NewPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.NewPassword), bcrypt.MinCost)
		if err != nil {
			return err
		}
		user.Password = string(hash)
	}
	return uc.Repository.Update(&user.User)
}

func (uc *UserUseCase) GetUserByLogin(user string) (*models.User, error) {
	return uc.Repository.GetUserByLogin(user)
}

func (uc *UserUseCase) PrintUserList() {
	uc.Repository.PrintUserList()
}
