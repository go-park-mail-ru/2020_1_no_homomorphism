package user

import (
	"no_homomorphism/internal/pkg/models"
)

type UseCase interface {
	Create(user *models.User)  error
	Update(user *models.User) error
}