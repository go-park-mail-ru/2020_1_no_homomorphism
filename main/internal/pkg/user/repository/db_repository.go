package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       uint64 `gorm:"column:id"`
	Login    string `gorm:"column:login"`
	Password []byte `gorm:"column:password"`
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	Sex      string `gorm:"column:sex"`
	Image    string `gorm:"column:image"`
	Theme    string `gorm:"column:theme"`
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

	db := ur.db.
		Table("users").
		Where("login = ?", login).
		Find(&results)

	err := db.Error
	if db.Error != nil {
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
		Theme:    user.Theme,
	}, nil
}

func ToModel(user User) models.User {
	return models.User{
		Id:       strconv.FormatUint(user.Id, 10),
		Login:    user.Login,
		Password: string(user.Password),
		Name:     user.Name,
		Email:    user.Email,
		Sex:      user.Sex,
		Image:    user.Image,
		Theme:    user.Theme,
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

	var hash []byte

	if input.NewPassword != "" {
		if err := bcrypt.CompareHashAndPassword(dbUser.Password, []byte(input.Password)); err != nil {
			return fmt.Errorf("old password is wrong : %v", err)
		}
		hash, err = bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
		if err != nil {
			return fmt.Errorf("error while password hashing: %v", err)
		}
		dbUser.Password = hash
	}
	dbUser.Name = input.Name
	dbUser.Email = input.Email
	if dbUser.Theme != "" {
		dbUser.Theme = input.Theme
	}

	db := ur.db.Save(&dbUser)

	return db.Error
}

func (ur *DbUserRepository) UpdateAvatar(user models.User, avatarDir string, fileName string) (string, error) {

	serverFilePath := os.Getenv("FILE_SERVER") + filepath.Join(avatarDir, fileName)

	db := ur.db.
		Model(&user).
		Update("image", serverFilePath)

	err := db.Error
	if err != nil {
		return "", fmt.Errorf("failed to update user: %v", err)
	}

	return serverFilePath, nil
}

func (ur *DbUserRepository) GetUserByLogin(login string) (models.User, error) {
	dbUser, err := ur.getUser(login)
	if err != nil {
		return models.User{}, err
	}
	return ToModel(dbUser), nil
}

func (ur *DbUserRepository) CheckIfExists(login string, email string) (loginExists bool, emailExists bool, err error) {
	var results []User
	db := ur.db.
		Raw("SELECT login, email FROM users WHERE login=? or email=?", login, email).
		Scan(&results)

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

	db := ur.db.
		Table("user_stat").
		Where("user_id = ?", id).
		Find(&stat)

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
