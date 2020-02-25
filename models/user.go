package models

import (
	"regexp"
	"sync"

	uuid "github.com/satori/go.uuid"
)

func NewUsersStorage() *UsersStorage {
	return &UsersStorage{
		Users:  make(map[string]*User),
		Mutex:  sync.RWMutex{},
		//nextId: 0,
	}
}

type UsersStorage struct {
	Users  map[string]*User
	Mutex  sync.RWMutex
	nextId uuid.UUID
}

type User struct {
	Id        uuid.UUID   `json:"id"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"avatar_url"`
}

// func (us *UsersStorage) AddUser(user *User) (uint, error) {
// 	us.Mutex.Lock()
// 	user.Id = us.nextId
// 	us.Users[user.Nickname] = user
// 	us.nextId++
// 	us.Mutex.Unlock()
//
// 	return user.Id, nil
// }

func ValidateEmail(email string) bool {

	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|" +
		"}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\." +
		"[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return pattern.MatchString(email)
}

func ValidateNickname(nickname string) bool {
	return true
}