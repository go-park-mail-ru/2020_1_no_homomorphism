package models

import (
	"errors"
	"sync"

	uuid "github.com/satori/go.uuid"
)

func NewUsersStorage() *UsersStorage {
	return &UsersStorage{
		Users: make(map[string]*User),
		Mutex: sync.RWMutex{},
	}
}

type UsersStorage struct {
	Users map[string]*User
	Mutex sync.RWMutex
}

type User struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Login    string    `json:"login"`
	Sex      string    `json:"sex"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
}

type UserInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserSettings struct {
	NewPassword string `json:"newPassword"`
	Name     string    `json:"name"`
	Login    string    `json:"login"`
	Sex      string    `json:"sex"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
}

func (us *UsersStorage) AddUser(input *UserInput) (uuid.UUID, error) {
	us.Mutex.Lock()
	user := &User{
		Id:       uuid.NewV4(),
		Login:    input.Login,
		Password: input.Password,
	}
	us.Users[user.Login] = user
	us.Mutex.Unlock()
	return user.Id, nil
}
func (us *UsersStorage) GetIdByUsername(username string) uuid.UUID{
	return us.Users[username].Id
}
func (us *UsersStorage) GetUserById(id uuid.UUID) (*User, error) {
	for _, user := range us.Users {
		if user.Id == id {
			return user, nil
		}
	}
	return nil, errors.New("user with this id does not exists: " + id.String())
}

func (us *UsersStorage) EditUser(user *User, newUserData *UserSettings) {
	newUser := &User{
		Id: user.Id,
		Name: newUserData.Name,
		Login: newUserData.Login,
		Password: newUserData.NewPassword,
		Email: newUserData.Email,
		Sex: newUserData.Sex,
	}
	delete(us.Users, user.Login)
	us.Users[newUserData.Login] = newUser
}
