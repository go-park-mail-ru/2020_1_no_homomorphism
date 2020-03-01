package user

import (
	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
)
type Repository interface {
	Create(user models.User) (*models.User, error)
	Update(user models.UserSettings) (*models.User, error)

	AddUser(input *models.User) (uuid.UUID, error)
	GetProfileByLogin(login string) (*models.Profile, error)
	GetIdByLogin(login string) uuid.UUID
	GetUserById(id uuid.UUID) (*models.User, error)
	GetUserPassword(login string) (string, error)
	GetFullUserInfo(login string) (models.User, error)
	EditUser(user *models.User, newUserData *models.UserSettings) error
}
