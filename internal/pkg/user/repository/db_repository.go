package repository

import (
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/internal/pkg/models"
)

type User struct {
	Id       uint64 `gorm:"column:id"`
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

func NewDbUserRepository(database *gorm.DB, defaultImage string) DbUserRepository {
	return DbUserRepository{
		db:           database,
		defaultImage: defaultImage,
	}
}

func (ur *DbUserRepository) getUser(login string) (User, error) {
	var results User
	db := ur.db.Raw("SELECT id, login, password, name, email, sex, image FROM users WHERE login=?", login).Scan(&results)
	err := db.Error
	if err != nil {
		return User{}, err
	}
	return results, nil
}

func (ur *DbUserRepository) prepareDbUser(user models.User, hash []byte) (User, error) {
	ok := IsModelFieldsNotEmpty(user)
	if !ok {
		return User{}, errors.New("some input fields are empty")
	}
	return User{
		Login:    user.Login,
		Password: hash,
		Name:     user.Name,
		Email:    user.Email,
		Sex:      user.Sex,
		Image:    ur.defaultImage,
	}, nil
}

func ToModel(user User) models.User {
	return models.User{
		Id:       fmt.Sprint(user.Id),
		Login:    user.Login,
		Password: string(user.Password),
		Name:     user.Name,
		Email:    user.Email,
		Sex:      user.Sex,
		Image:    user.Image,
	}
}

func (ur *DbUserRepository) Create(user models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return fmt.Errorf("error while password hashing: %v", err)
	}
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

func (ur *DbUserRepository) Update(user models.User, input models.UserSettings) error {

	dbUser, err := ur.getUser(user.Login)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword(dbUser.Password, []byte(input.Password)); err != nil {
		return fmt.Errorf("old password is wrong : %v", err)
	}
	var hash []byte
	if input.NewPassword != "" {
		hash, err = bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
		if err != nil {
			return fmt.Errorf("error while password hashing: %v", err)
		}
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

func (ur *DbUserRepository) UpdateAvatar(user models.User, filePath string) error {
	dbUser, err := ur.getUser(user.Login)
	if err != nil {
		return err
	}
	dbUser.Image = os.Getenv("FILE_SERVER") + filePath

	db := ur.db.Save(&dbUser)
	err = db.Error
	if err != nil {
		return err //todo error wrapper
	}
	return nil
}

func (ur *DbUserRepository) GetUserByLogin(login string) (models.User, error) {
	dbUser, err := ur.getUser(login)
	if err != nil {
		return models.User{}, err
	}
	user := ToModel(dbUser)
	return user, nil
}

func (ur *DbUserRepository) CheckIfExists(login string, email string) (loginExists bool, emailExists bool, err error) {
	var results []User
	db := ur.db.Raw("SELECT id, login, email FROM users WHERE login=? or email=?", login, email).Scan(&results)
	err = db.Error
	if err == gorm.ErrRecordNotFound {
		return false, false, nil
	}
	if err != nil {
		return true, true, err
	}
	for _, elem := range results {
		if elem.Login == login {
			loginExists = true
		}
		if elem.Email == email {
			emailExists = true
		}
	}
	return loginExists, emailExists, nil
}

func (ur *DbUserRepository) CheckUserPassword(userPassword string, InputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(InputPassword)); err != nil {
		return errors.New("wrong password")
	}
	return nil
}

func (ur *DbUserRepository) GetUserStat(id string) (models.UserStat, error) {
	var stat models.UserStat

	db := ur.db.Table("user_stat").Where("user_id = ?", id).Find(&stat)
	err := db.Error
	if err != nil {
		return models.UserStat{}, err
	}
	return stat, nil
}

func IsModelFieldsNotEmpty(user models.User) bool {
	return len(user.Login) > 0 &&
		len(user.Password) > 0 &&
		len(user.Name) > 0 &&
		len(user.Email) > 0 &&
		len(user.Sex) > 0
}
