package user

import (
	"no_homomorphism/internal/pkg/models"
)

type UseCase interface {
	Create(user *models.User) error
	Update(user *models.UserSettings) error
	GetUserByLogin(login string) (*models.User, error)
	PrintUserList()
	GetProfileByLogin(login string) (*models.Profile, error)
}
