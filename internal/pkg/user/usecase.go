package user

import (
	"io"
	"no_homomorphism/internal/pkg/models"
)

type SameUserExists int8

const (
	NO SameUserExists = iota
	LOGIN
	EMAIL
	FULL
)

type UseCase interface {
	Create(user models.User) (SameUserExists, error)
	Update(user models.User, input models.UserSettings) (SameUserExists, error)
	Login(input models.UserSignIn) (models.User, error)
	UpdateAvatar(user models.User, file io.Reader, fileType string) (string, error)
	GetUserByLogin(login string) (models.User, error)
	GetProfileByLogin(login string) (models.Profile, error)
	GetProfileByUser(user models.User) models.Profile
	CheckUserPassword(userPassword string, InputPassword string) error
}
