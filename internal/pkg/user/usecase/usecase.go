package usecase

import (
	"sync"

	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/user"
	"no_homomorphism/internal/pkg/user/repository"
)

var mutex = &sync.Mutex{}

type userUseCase struct{
	repo *repository.MemUserRepository
	Mutex *sync.Mutex
}

func newUserUseCase() user.UseCase{
	return &userUseCase{
		repo: repository.NewUserRepository(mutex),
		Mutex: mutex,
	}
}

func (uc *userUseCase) Create(user *models.User) error {
	return uc.repo.Create(user)
}
func (uc *userUseCase) Update(user *models.User) error {
	return uc.Update(user)
}

// ----------------------------------------------------------

