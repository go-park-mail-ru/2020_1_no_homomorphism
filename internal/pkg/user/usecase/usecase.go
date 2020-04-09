package usecase

import (
	"fmt"
	"io"

	"no_homomorphism/internal/pkg/models"
	users "no_homomorphism/internal/pkg/user"
)

type UserUseCase struct {
	Repository users.Repository
}

func (uc *UserUseCase) Create(user models.User) (users.SameUserExists, error) {
	loginExists, emailExists, err := uc.Repository.CheckIfExists(user.Login, user.Email)
	if err != nil {
		return users.FULL, err
	}
	if loginExists && emailExists {
		return users.FULL, nil
	}
	if loginExists {
		return users.LOGIN, nil
	}
	if emailExists {
		return users.EMAIL, nil
	}
	return users.NO, uc.Repository.Create(user)
}

func (uc *UserUseCase) Update(user models.User, input models.UserSettings) (users.SameUserExists, error) {
	if user.Email != input.Email {
		_, emailExists, err := uc.Repository.CheckIfExists("", input.Email)
		if err != nil {
			return users.FULL, fmt.Errorf("failed to check email existing: %v", err)
		}
		if emailExists {
			return users.EMAIL, nil
		}
	}
	return users.NO, uc.Repository.Update(user, input)
}

func (uc *UserUseCase) UpdateAvatar(user models.User, file io.Reader, fileType string) (string, error) {

	return uc.Repository.UpdateAvatar(user, file, fileType)
}

func (uc *UserUseCase) GetUserByLogin(user string) (models.User, error) {
	return uc.Repository.GetUserByLogin(user)
}

func (uc *UserUseCase) GetProfileByLogin(login string) (models.User, error) {
	user, err := uc.Repository.GetUserByLogin(login)
	if err != nil {
		return models.User{}, err
	}
	return uc.GetOutputUserData(user), nil
}

func (uc *UserUseCase) Login(input models.UserSignIn) (models.User, error) {
	user, err := uc.GetUserByLogin(input.Login)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %v", err)
	}
	err = uc.CheckUserPassword(user.Password, input.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("wrong password: %v", err)
	}
	return user, nil
}

func (uc *UserUseCase) GetOutputUserData(user models.User) models.User {
	return models.User{
		Id:       user.Id,
		Name:     user.Name,
		Login:    user.Login,
		Sex:      user.Sex,
		Image:    user.Image,
		Email:    user.Email,
	}
}

func (uc *UserUseCase) GetUserStat(id string) (models.UserStat, error) {
	return uc.Repository.GetUserStat(id)
}

func (uc *UserUseCase) CheckUserPassword(userPassword string, inputPassword string) error {
	return uc.Repository.CheckUserPassword(userPassword, inputPassword)
}
