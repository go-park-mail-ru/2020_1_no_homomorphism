package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
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
	Image    string `sql:"null"`
}

type DbUserRepository struct {
	db *gorm.DB
}

func NewTestDbUserRepository(database *gorm.DB) *DbUserRepository {
	return &DbUserRepository{
		db: database,
	}
}

func toDbUserWithHash(user *models.User, hash []byte) *User {
	return &User{
		Login:    user.Login,
		Password: hash,
		Name:     user.Name,
		Email:    user.Email,
		Sex:      user.Sex,
	}
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
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	dbUser := toDbUserWithHash(user, hash)
	db := ur.db.Create(&dbUser)
	err = db.Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *DbUserRepository) Update(user *models.User, input *models.UserSettings) error {
	dbUser := User{}
	db := ur.db.Find(&dbUser)
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
	err = db.Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *DbUserRepository) UpdateAvatar(user *models.User, filePath string) {
	user.Image = filePath
}

func (ur *DbUserRepository) GetUserByLogin(login string) (*models.User, error) {
	dbUser := User{}
	db := ur.db.Find(dbUser)
	logrus.Info(dbUser)
	err := db.Error
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}
	user := toModel(&dbUser)
	logrus.Info(user)
	return user, nil //todo errors
}

func (ur *DbUserRepository) PrintUserList() {
	//fmt.Println("[USERS LIST]")
	//for _, r := range ur.Users {
	//	fmt.Println(r)
	//}
}
