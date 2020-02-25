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
func (us *UsersStorage) GetByUsername(username string) (uuid.UUID) {
	return us.Users[username].Id
}

func ValidateEmail(email string) bool {

	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|" +
		"}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\." +
		"[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return pattern.MatchString(email)
}

func ValidateNickname(nickname string) bool {
	return true
}