package user

import (
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"io"
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
	GetProfileByLogin(login string) (models.User, error)
	GetOutputUserData(user models.User) models.User
	CheckUserPassword(userPassword string, InputPassword string) error
	GetUserStat(id string) (models.UserStat, error)
}
