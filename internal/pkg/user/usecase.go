package user

import (
	"mime/multipart"

	"no_homomorphism/internal/pkg/models"
)

type UseCase interface {
	Create(user *models.User) error
	Update(user *models.User, input *models.UserSettings) error
	UpdateAvatar(user *models.User, file multipart.File) error
	GetUserByLogin(login string) (*models.User, error)
	PrintUserList()
	GetProfileByLogin(login string) (*models.Profile, error)
	GetProfileByUser(user *models.User) *models.Profile
}