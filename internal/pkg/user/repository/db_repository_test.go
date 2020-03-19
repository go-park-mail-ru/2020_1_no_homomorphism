package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/internal/pkg/models"
	"testing"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	user       *models.User
	repository *DbUserRepository
	bdError    error
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)
	s.user = &models.User{
		Id:       "1",
		Password: "12345678",
		Name:     "Vasya",
		Login:    "pupkin",
		Sex:      "male",
		Image:    "/image/test",
		Email:    "test@email.test",
	}

	s.bdError = errors.New("some bd error")

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	s.DB.LogMode(false)

	s.repository = NewDbUserRepository(s.DB, "/image/test")
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) getMockSelectAll(user *models.User, hash []byte) {
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "name", "sex", "image", "email"}).
			AddRow(user.Id, user.Login, hash, user.Name, user.Sex, user.Image, user.Email))
}

func (s *Suite) TestGetUserByLogin() {
	user := s.user

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	require.NoError(s.T(), err)

	s.getMockSelectAll(user, hash)

	getUser, err := s.repository.GetUserByLogin(user.Login)

	require.NoError(s.T(), err)
	user.Password = string(hash)
	require.Nil(s.T(), deep.Equal(user, getUser))

	//test on bd error
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login).WillReturnError(s.bdError)

	_, err = s.repository.GetUserByLogin(user.Login)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestUpdate() {
	user := s.user

	userSettings := &models.UserSettings{
		NewPassword: "",
		User: models.User{
			Name:  "lol",
			Email: "newemail@mail.ru",
		},
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	require.Nil(s.T(), err)

	s.getMockSelectAll(user, hash)

	var id int64
	id = 1
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").WithArgs(user.Login, hash, userSettings.Name, userSettings.Email, user.Sex, user.Image, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err = s.repository.Update(user, userSettings, []byte{})

	require.NoError(s.T(), err)

	//test on hash update
	s.getMockSelectAll(user, hash)

	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").WithArgs(user.Login, hash, userSettings.Name, userSettings.Email, user.Sex, user.Image, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err = s.repository.Update(user, userSettings, hash)

	require.NoError(s.T(), err)

	//test on first query bd error
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login).WillReturnError(s.bdError)

	err = s.repository.Update(user, userSettings, hash)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)

	//test on second query bd error
	s.getMockSelectAll(user, hash)

	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").WithArgs(user.Login, hash, userSettings.Name, userSettings.Email, user.Sex, user.Image, id).
		WillReturnError(s.bdError)
	s.mock.ExpectRollback()

	err = s.repository.Update(user, userSettings, hash)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestCreate() {
	user := s.user
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	require.NoError(s.T(), err)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").WithArgs(user.Login, hash, user.Name, user.Email, user.Sex, user.Image).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.Id))
	s.mock.ExpectCommit()

	err = s.repository.Create(user, hash)
	require.NoError(s.T(), err)

	//test on empty input
	user.Email = ""

	err = s.repository.Create(user, hash)

	require.Equal(s.T(), err, errors.New("some input fields are empty"))

	//test on bd error
	user.Email = "mail@mai.ru"

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").WithArgs(user.Login, hash, user.Name, user.Email, user.Sex, user.Image).
		WillReturnError(s.bdError)
	s.mock.ExpectRollback()

	err = s.repository.Create(user, hash)
	require.Equal(s.T(), err, s.bdError)
}
func (s *Suite) TestCheckIfExists() {
	user := s.user

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery("SELECT").WithArgs(user.Login, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "name", "sex", "image", "email"}).
			AddRow(user.Id, user.Login, hash, user.Name, user.Sex, user.Image, user.Email))

	exists, err := s.repository.CheckIfExists(user.Login, user.Email)
	require.NoError(s.T(), err)
	require.Equal(s.T(), exists, true)

	//test on not exists
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login, user.Email).WillReturnError(gorm.ErrRecordNotFound)

	exists, err = s.repository.CheckIfExists(user.Login, user.Email)
	require.NoError(s.T(), err)
	require.Equal(s.T(), exists, false)

	//test on other error
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login, user.Email).WillReturnError(s.bdError)

	exists, err = s.repository.CheckIfExists(user.Login, user.Email)
	require.Error(s.T(), err)
	require.Equal(s.T(), exists, true)
}

func (s *Suite) TestUpdateAvatar() {
	user := s.user

	filePath := "new/user/filepath"

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	require.Nil(s.T(), err)

	s.getMockSelectAll(user, hash)

	var id int64
	id = 1
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").WithArgs(user.Login, hash, user.Name, user.Email, user.Sex, filePath, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err = s.repository.UpdateAvatar(user, filePath)

	require.NoError(s.T(), err)

	//test on first query bd error
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login).WillReturnError(s.bdError)

	err = s.repository.UpdateAvatar(user, filePath)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)

	//test on second query bd error
	s.getMockSelectAll(user, hash)

	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").WithArgs(user.Login, hash, user.Name, user.Email, user.Sex, filePath, id).WillReturnError(s.bdError)
	s.mock.ExpectRollback()

	err = s.repository.UpdateAvatar(user, filePath)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)

}
