package user

import (
	"io"

	"no_homomorphism/internal/pkg/models"
)

type Repository interface {
	Create(user models.User) error
	Update(user models.User, input models.UserSettings) error
	UpdateAvatar(user models.User, file io.Reader, fileType string) (string, error)
	GetUserByLogin(login string) (models.User, error)
	CheckIfExists(login string, email string) (loginExists bool, emailExists bool, err error)
	CheckUserPassword(userPassword string, inputPassword string) error
	GetUserStat(id string) (models.UserStat, error)
}
