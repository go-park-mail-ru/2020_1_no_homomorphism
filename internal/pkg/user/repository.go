package user

import (
	"no_homomorphism/internal/pkg/models"
)

type Repository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	UpdateAvatar(user *models.User, avatarPath string)
	GetUserByLogin(login string) (*models.User, error)
	PrintUserList()
}
