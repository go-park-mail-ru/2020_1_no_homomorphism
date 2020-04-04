package usecase

import (
	"errors"
	"fmt"
	"io"
	"no_homomorphism/internal/pkg/models"
	users "no_homomorphism/internal/pkg/user"
	"os"
	"path/filepath"
)

type UserUseCase struct {
	Repository users.Repository
	AvatarDir  string
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
	fileName := user.Id
	filePath := filepath.Join(os.Getenv("MUSIC_PROJ_DIR"), uc.AvatarDir, fileName+"."+fileType)

	newFile, err := os.Create(filePath)
	if err != nil {
		return "", errors.New("failed to create file")
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return "", errors.New("error while writing to file")
	}

	err = uc.Repository.UpdateAvatar(user, filepath.Join(uc.AvatarDir, fileName+"."+fileType))
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (uc *UserUseCase) GetUserByLogin(user string) (models.User, error) {
	return uc.Repository.GetUserByLogin(user)
}

func (uc *UserUseCase) GetProfileByLogin(login string) (models.Profile, error) {
	user, err := uc.Repository.GetUserByLogin(login)
	if err != nil {
		return models.Profile{}, err
	}
	return uc.GetProfileByUser(user), nil
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

func (uc *UserUseCase) GetProfileByUser(user models.User) models.Profile {
	return models.Profile{
		Name:  user.Name,
		Login: user.Login,
		Sex:   user.Sex,
		Image: user.Image,
		Email: user.Email,
	}
}

func (uc *UserUseCase) GetUserStat(id string) (models.UserStat, error) {
	return uc.Repository.GetUserStat(id)
}

func (uc *UserUseCase) CheckUserPassword(userPassword string, inputPassword string) error {
	return uc.Repository.CheckUserPassword(userPassword, inputPassword)
}
