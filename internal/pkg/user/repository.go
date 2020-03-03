package user

import (
	"no_homomorphism/internal/pkg/models"
)

type Repository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	GetUserByLogin(login string) (*models.User, error)
	PrintUserList()
}
