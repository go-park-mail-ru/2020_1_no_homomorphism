package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/internal/pkg/models"
)

type User struct {
	Id       uint64 `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Login    string
	Password []byte
	Name     string
	Email    string
	Sex      string
	Image    string
}

type DbUserRepository struct {
	db           *gorm.DB
	defaultImage string
}

func NewTestDbUserRepository(database *gorm.DB, defaultImage string) *DbUserRepository {
	return &DbUserRepository{
		db:           database,
		defaultImage: defaultImage,
	}
}

func (ur *DbUserRepository) prepareDbUser(user *models.User) (*User, error) {
	ok := IsModelFieldsNotEmpty(user)
	if !ok {
		return nil, errors.New("some input fields are empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
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

func toModel(user *User) *models.User {
	return &models.User{
		Id:       string(user.Id),
		Login:    user.Login,
		Password: string(user.Password),
		Name:     user.Name,
		Email:    user.Email,
		Sex:      user.Sex,
		Image:    user.Image,
	}
}

func (ur *DbUserRepository) Create(user *models.User) error {
	dbUser, err := ur.prepareDbUser(user)
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

func (ur *DbUserRepository) Update(user *models.User, input *models.UserSettings) error {
	var dbUser User
	db := ur.db.Raw("SELECT * FROM users WHERE login=?", user.Login).Scan(&dbUser)
	err := db.Error
	if err != nil {
		return err
	}
	if input.NewPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
		if err != nil {
			return err
		}
		dbUser.Password = hash
	}
	dbUser.Name = input.Name
	dbUser.Email = input.Email
	db = ur.db.Save(&dbUser)
	//if db.RowsAffected < 1 {
	//	return errors.New("user is not modified")
	//}
	err = db.Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *DbUserRepository) UpdateAvatar(user *models.User, filePath string) error {
	var dbUser User
	db := ur.db.Raw("SELECT * FROM users WHERE login=?", user.Login).Scan(&dbUser)
	err := db.Error
	if err != nil {
		return err
	}
	dbUser.Image = filePath

	db = ur.db.Save(&dbUser)
	err = db.Error
	if err != nil {
		return err //todo error wrapper
	}
	return nil
}

func (ur *DbUserRepository) GetUserByLogin(login string) (*models.User, error) {
	var results User
	db := ur.db.Raw("SELECT * FROM users WHERE login=?", login).Scan(&results)
	err := db.Error
	if err != nil {
		return nil, err
	}
	user := toModel(&results)
	return user, nil
}

func (ur *DbUserRepository) CheckIfExists(login string, email string) (bool, error) {
	var results User
	db := ur.db.Raw("SELECT * FROM users WHERE login=? or email=?", login, email).Scan(&results)
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
