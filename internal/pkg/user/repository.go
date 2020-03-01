package user

import (
	"no_homomorphism/internal/pkg/models"
)

type Repository interface {
	Create(user models.User) (*models.User, error)
	Update(user models.UserSettings) (*models.User, error)
}
