package repository

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/internal/pkg/models"
)

type MemUserRepository struct {
	Users map[string]*models.User
	Count uint
}

func (ur *MemUserRepository) Create(user *models.User) error {
	user.Id = ur.Count
	ur.Count++
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return nil
	}
	user.Password = string(hash)
	ur.Users[user.Login] = user
	return nil
}

func (ur *MemUserRepository) Update(user *models.User, input *models.UserSettings) error {
	if input.NewPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
		if err != nil {
			return err
		}
		user.Password = string(hash)
	}
	user.Email = input.Email
	return nil
}

func (ur *MemUserRepository) UpdateAvatar(user *models.User, filePath string) {
	user.Image = filePath
}

func (ur *MemUserRepository) GetUserByLogin(login string) (*models.User, error) {
	user, ok := ur.Users[login]
	if !ok {
		return nil, errors.New("user with this login does not exists")
	}
	return user, nil
}

func (ur *MemUserRepository) PrintUserList() {
	fmt.Println("[USERS LIST]")
	for _, r := range ur.Users {
		fmt.Println(r)
	}
}
