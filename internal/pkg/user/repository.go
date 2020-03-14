package user

import (
	"no_homomorphism/internal/pkg/models"
)

type Repository interface {
	Create(user *models.User) error
	Update(user *models.User, input *models.UserSettings) error
	UpdateAvatar(user *models.User, avatarPath string) error
	GetUserByLogin(login string) (*models.User, error)
	CheckIfExists(login string, email string) (bool, error)
}
