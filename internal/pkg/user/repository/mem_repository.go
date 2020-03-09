package repository

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/internal/pkg/models"
)

type MemUserRepository struct {
	Users map[string]*models.User
	count uint
	mutex *sync.RWMutex
}

func NewTestMemUserRepository() *MemUserRepository {
	return &MemUserRepository{
		Users: map[string]*models.User{
			"test": &models.User{
				Id:       0,
				Login:    "test",
				Name:     "Rita",
				Email:    "rita@margarita.tyt",
				Password: "$2a$04$0GzSltexrV9gQjFwv5BYuebu7/F13cX.NOupseJQUwqHWDucyBBgO",
				Image:    "/static/img/avatar/default.png",
			},
			"test2": &models.User{
				Id:       1,
				Login:    "test2",
				Name:     "User2",
				Email:    "user2@da.tu",
				Password: "$2a$04$r/rWIhO8ptZAxheWs9cXmeG8fKhICfA5Gko3Qr61ae0.71CwjyODC",
				Image:    "/static/img/avatar/default.png",
			},
			"test3": &models.User{
				Id:       2,
				Login:    "test3",
				Name:     "User3",
				Email:    "user3@da.tu",
				Password: "$2a$04$8G8SC41DvtOYD04qVizzbek.uL9zEI5zlQ3q2Cg.DYekuzMWFsoLa",
				Image:    "/static/img/avatar/default.png",
			},
		},
		count: 3,
		mutex: &sync.RWMutex{},
	}
}

func (ur *MemUserRepository) Create(user *models.User) error {
	user.Id = ur.count
	ur.count++
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return nil
	}
	user.Password = string(hash)
	ur.mutex.Lock()
	ur.Users[user.Login] = user
	ur.mutex.Unlock()
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
	ur.mutex.Lock()
	user, ok := ur.Users[login]
	ur.mutex.Unlock()
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
