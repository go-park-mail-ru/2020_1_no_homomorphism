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

// func NewUserRepository(mutex *sync.Mutex) *MemUserRepository {
// 	return &MemUserRepository{
// 		Users: make(map[string]*models.User),
// 		Count: 0,
// 	}
// }

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

func (ur *MemUserRepository) Update(user *models.User) error {
	ur.Users[user.Login] = user
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

// --------------------------------------------------
/*
func NewUsersStorage() *UsersStorage {
	return &UsersStorage{
		Users: make(map[string]*User),
		Mutex: &sync.Mutex{},
	}
}

type UsersStorage struct {
	Users map[string]*User
	Mutex *sync.Mutex
}

func (us *UsersStorage) AddUser(input *User) (uuid.UUID, error) {

	if input == nil {
		return uuid.UUID{0}, errors.New("nil input")
	}
	input.Id = uuid.NewV4()
	us.Mutex.Lock()
	defer us.Mutex.Unlock()
	_, e := us.Users[input.Login]
	if e {
		return uuid.UUID{0}, errors.New("user with current login is already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	input.Password = string(hash)
	us.Users[input.Login] = input
	return input.Id, nil
}

func (us *UsersStorage) GetProfileByLogin(login string) (*Profile, error) {
	us.Mutex.Lock()
	defer us.Mutex.Unlock()
	_, e := us.Users[login]
	if !e {
		return nil, errors.New("no user with that name")
	}
	profile := &Profile{
		Login: us.Users[login].Login,
		Sex:   us.Users[login].Sex,
		Name:  us.Users[login].Name,
		Image: "",
		Email: us.Users[login].Email,
	}
	return profile, nil
}

func (us *UsersStorage) GetIdByLogin(login string) uuid.UUID {
	us.Mutex.Lock()
	defer us.Mutex.Unlock()
	return us.Users[login].Id
}

func (us *UsersStorage) GetUserById(id uuid.UUID) (*User, error) {
	us.Mutex.Lock()
	defer us.Mutex.Unlock()
	for _, user := range us.Users {
		if user.Id == id {
			return user, nil
		}
	}
	return nil, errors.New("user with this id does not exists: " + id.String())
}

func (us *UsersStorage) GetUserPassword(login string) (string, error) {
	us.Mutex.Lock()
	defer us.Mutex.Unlock()
	if user, ok := us.Users[login]; !ok {
		return "", errors.New("user with this login does not exists: " + login)
	} else {
		return user.Password, nil
	}
}

func (us *UsersStorage) GetFullUserInfo(login string) (User, error) {
	us.Mutex.Lock()
	defer us.Mutex.Unlock()
	if user, ok := us.Users[login]; !ok {
		return User{}, errors.New("user with this login does not exists: " + login)
	} else {
		return *user, nil
	}
}

func (us *UsersStorage) EditUser(user *User, newUserData *UserSettings) error {
	if user == nil || newUserData == nil {
		return errors.New("input data is nil")
	}
	us.Mutex.Lock()
	defer us.Mutex.Unlock()
	hash, err := bcrypt.GenerateFromPassword([]byte(newUserData.NewPassword), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	if newUserData.Email != "" {
		user.Email = newUserData.Email
	}
	if newUserData.NewPassword != "" {
		user.Password = string(hash)
	}
	if newUserData.Sex != "" {
		user.Password = string(hash)
	}
	if newUserData.Name != "" {
		user.Name = newUserData.Name
	}

	fmt.Println("this is user:", user)

	return nil
} */