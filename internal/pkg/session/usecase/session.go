package usecase

import (
	"sync"

	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/session"
	"no_homomorphism/internal/pkg/session/repository"
)


type sessionUseCase struct {
	repository session.Repository
	mutex *sync.Mutex
}

func NewSessionUseCase(mutex *sync.Mutex) session.UseCase {
	return &sessionUseCase{
		repository: repository.NewSessionRepository(mutex),
		mutex: mutex,
	}
}
