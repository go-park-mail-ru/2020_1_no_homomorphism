package repository

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"no_homomorphism/internal/pkg/models"
)

type User struct {
	Id       uint64 `sql:"AUTO_INCREMENT" gorm:"column:id"`
	Login    string `gorm:"column:login"`
	Password []byte `gorm:"column:password"`
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	Sex      string `gorm:"column:sex"`
	Image    string `gorm:"column:image"`
}

type DbUserRepository struct {
	db           *gorm.DB
	defaultImage string
}

func NewDbUserRepository(database *gorm.DB, defaultImage string) *DbUserRepository {
	return &DbUserRepository{
		db:           database,
		defaultImage: defaultImage,
	}
}

func (ur *DbUserRepository) getUser(login string) (*User, error) {
	var results User
	db := ur.db.Raw("SELECT id, login, password, name, email, sex, image FROM users WHERE login=?", login).Scan(&results)
	err := db.Error
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func (ur *DbUserRepository) prepareDbUser(user *models.User, hash []byte) (*User, error) {
	ok := IsModelFieldsNotEmpty(user)
	if !ok {
		return nil, errors.New("some input fields are empty")
	}
	return &User{
		Login:    user.Login,
		Password: hash,
		Name:     user.Name,
		Email:    user.Email,
		Sex:      user.Sex,
		Image:    ur.defaultImage,
	}, nil
}

func ToModel(user *User) *models.User {
	return &models.User{
		Id:       fmt.Sprint(user.Id),
		Login:    user.Login,
		Password: string(user.Password),
		Name:     user.Name,
		Email:    user.Email,
		Sex:      user.Sex,
		Image:    user.Image,
	}
}

func (ur *DbUserRepository) Create(user *models.User, hash []byte) error {
	dbUser, err := ur.prepareDbUser(user, hash)
	if err != nil {
		return err
	}
	db := ur.db.Create(&dbUser)
	err = db.Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *DbUserRepository) Update(user *models.User, input *models.UserSettings, hash []byte) error {
	dbUser, err := ur.getUser(user.Login)
	if err != nil {
		return err
	}

	if len(hash) > 0 {
		dbUser.Password = hash
	}
	dbUser.Name = input.Name
	dbUser.Email = input.Email
	db := ur.db.Save(&dbUser)
	err = db.Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *DbUserRepository) UpdateAvatar(user *models.User, filePath string) error {
	dbUser, err := ur.getUser(user.Login)
	if err != nil {
		return err
	}
	dbUser.Image = filePath

	db := ur.db.Save(&dbUser)
	err = db.Error
	if err != nil {
		return err //todo error wrapper
	}
	return nil
}

func (ur *DbUserRepository) GetUserByLogin(login string) (*models.User, error) {
	dbUser, err := ur.getUser(login)
	if err != nil {
		return nil, err
	}
	user := ToModel(dbUser)
	return user, nil
}

func (ur *DbUserRepository) CheckIfExists(login string, email string) (bool, error) {
	var results User
	db := ur.db.Raw("SELECT id FROM users WHERE login=? or email=?", login, email).Scan(&results)
	err := db.Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return true, nil
}

func IsModelFieldsNotEmpty(user *models.User) bool {
	return len(user.Login) > 0 &&
		len(user.Password) > 0 &&
		len(user.Name) > 0 &&
		len(user.Email) > 0 &&
		len(user.Sex) > 0
}
