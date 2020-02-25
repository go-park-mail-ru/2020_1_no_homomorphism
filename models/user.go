package models

import (
	"errors"
	"sync"

	uuid "github.com/satori/go.uuid"
)

func NewUsersStorage() *UsersStorage {
	return &UsersStorage{
		Users:  make(map[string]*User),
		Mutex:  sync.RWMutex{},
	}
}

type UsersStorage struct {
	Users  map[string]*User
	Mutex  sync.RWMutex
}

type User struct {
	Id        uuid.UUID   `json:"id"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"avatar_url"`
}

type UserInput struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func (us *UsersStorage) AddUser(input *UserInput) (uuid.UUID, error) {
	us.Mutex.Lock()
	user := &User{
		Id:       uuid.NewV4(),
		Nickname: input.Nickname,
		Password: input.Password,
	}
	us.Users[user.Nickname] = user
	us.Mutex.Unlock()
	return user.Id, nil
}
func (us *UsersStorage) GetByUsername(username string) uuid.UUID {
	return us.Users[username].Id
}
func (us *UsersStorage) GetById(id uuid.UUID) (*User, error) {
	for _, user := range us.Users {
		if user.Id == id {
			return user, nil
		}
	}
	return nil, errors.New("user with this id does not exists: " + id.String())
}

func (us *UsersStorage) EditUser(user *User, newUserData *User){

	newUserData.Id = user.Id
	delete(us.Users, user.Nickname)
	us.Users[newUserData.Nickname] = newUserData
}
