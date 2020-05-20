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
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"regexp"
	"testing"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	user       models.User
	repository DbUserRepository
	bdError    error
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)
	s.user = models.User{
		Id:       "1",
		Password: "12345678",
		Name:     "Vasya",
		Login:    "pupkin",
		Sex:      "male",
		Image:    "/img/default/png",
		Email:    "test@email.test",
	}

	s.bdError = errors.New("some bd error")

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	s.DB.LogMode(false)

	s.repository = NewDbUserRepository(s.DB, "/img/default/png")
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) getMockSelectAll(user models.User, hash []byte) {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (login = $1)`)).WithArgs(user.Login).
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

	userSettings := models.UserSettings{
		NewPassword: "",
		User: models.User{
			Password: user.Password,
			Name:     "lol",
			Email:    "newemail@mail.ru",
		},
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	require.Nil(s.T(), err)

	s.getMockSelectAll(user, hash)

	var id int64 = 1
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").WithArgs(user.Login, hash, userSettings.Name, userSettings.Email, user.Sex, user.Image, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err = s.repository.Update(user, userSettings)

	require.NoError(s.T(), err)

	//test on hash update
	s.getMockSelectAll(user, hash)

	userSettings.NewPassword = "1235jei23"

	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").WithArgs(user.Login, sqlmock.AnyArg(), userSettings.Name, userSettings.Email, user.Sex, user.Image, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err = s.repository.Update(user, userSettings)
	require.NoError(s.T(), err)

	userSettings.NewPassword = ""
	//test on first query bd error
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login).WillReturnError(s.bdError)

	err = s.repository.Update(user, userSettings)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)

	//test on second query bd error
	s.getMockSelectAll(user, hash)

	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").WithArgs(user.Login, hash, userSettings.Name, userSettings.Email, user.Sex, user.Image, id).
		WillReturnError(s.bdError)
	s.mock.ExpectRollback()

	err = s.repository.Update(user, userSettings)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestCreate() {
	user := s.user
	//hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	//require.NoError(s.T(), err)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").WithArgs(user.Login, sqlmock.AnyArg(), user.Name, user.Email, user.Sex, user.Image).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.Id))
	s.mock.ExpectCommit()

	err := s.repository.Create(user)
	require.NoError(s.T(), err)

	//test on empty input
	user.Email = ""

	err = s.repository.Create(user)

	require.Equal(s.T(), err, errors.New("some input fields are empty"))

	//test on bd error
	user.Email = "mail@mai.ru"

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").WithArgs(user.Login, sqlmock.AnyArg(), user.Name, user.Email, user.Sex, user.Image).
		WillReturnError(s.bdError)
	s.mock.ExpectRollback()

	err = s.repository.Create(user)
	require.Equal(s.T(), err, s.bdError)
}
func (s *Suite) TestCheckIfExists() {
	user := s.user

	//test on both exists
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "email"}).
			AddRow(user.Id, user.Login, "email").AddRow(user.Id, "login", user.Email))

	loginExists, emailExists, err := s.repository.CheckIfExists(user.Login, user.Email)
	require.NoError(s.T(), err)
	require.Equal(s.T(), loginExists, true)
	require.Equal(s.T(), emailExists, true)

	//test on login exists
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "email"}).
			AddRow(user.Id, user.Login, "otherEmail"))

	loginExists, emailExists, err = s.repository.CheckIfExists(user.Login, user.Email)
	require.NoError(s.T(), err)
	require.Equal(s.T(), loginExists, true)
	require.Equal(s.T(), emailExists, false)

	//test on email exists
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "email"}).
			AddRow(user.Id, "otherName", user.Email))

	loginExists, emailExists, err = s.repository.CheckIfExists(user.Login, user.Email)
	require.NoError(s.T(), err)
	require.Equal(s.T(), loginExists, false)
	require.Equal(s.T(), emailExists, true)

	//test on not exists
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login, user.Email).WillReturnError(gorm.ErrRecordNotFound)

	loginExists, emailExists, err = s.repository.CheckIfExists(user.Login, user.Email)
	require.NoError(s.T(), err)
	require.Equal(s.T(), loginExists, false)
	require.Equal(s.T(), emailExists, false)

	//test on other error
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login, user.Email).WillReturnError(s.bdError)

	loginExists, emailExists, err = s.repository.CheckIfExists(user.Login, user.Email)
	require.Error(s.T(), err)
	require.Equal(s.T(), loginExists, true)
	require.Equal(s.T(), emailExists, true)
}

//
//func (s *Suite) TestUpdateAvatar() {
//	user := s.user
//
//	filePath := "new/user/filepath"
//
//	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
//	require.Nil(s.T(), err)
//
//	s.getMockSelectAll(user, hash)
//
//	var id int64
//	id = 1
//	s.mock.ExpectBegin()
//	s.mock.ExpectExec("UPDATE").WithArgs(user.Login, hash, user.Name, user.Email, user.Sex, filePath, id).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//	s.mock.ExpectCommit()
//
//	err = s.repository.UpdateAvatar(user, filePath)
//
//	require.NoError(s.T(), err)
//
//	//test on first query bd error
//	s.mock.ExpectQuery("SELECT").WithArgs(user.Login).WillReturnError(s.bdError)
//
//	err = s.repository.UpdateAvatar(user, filePath)
//
//	require.Error(s.T(), err)
//	require.Equal(s.T(), err, s.bdError)
//
//	//test on second query bd error
//	s.getMockSelectAll(user, hash)
//
//	s.mock.ExpectBegin()
//	s.mock.ExpectExec("UPDATE").WithArgs(user.Login, hash, user.Name, user.Email, user.Sex, filePath, id).WillReturnError(s.bdError)
//	s.mock.ExpectRollback()
//
//	err = s.repository.UpdateAvatar(user, filePath)
//
//	require.Error(s.T(), err)
//	require.Equal(s.T(), err, s.bdError)
//}

func (s *Suite) TestCheckUserPassword() {
	passOne := "sdofae87q3yncq823"
	hash, err := bcrypt.GenerateFromPassword([]byte(passOne), bcrypt.MinCost)
	require.NoError(s.T(), err)

	err = s.repository.CheckUserPassword(string(hash), passOne)
	require.NoError(s.T(), err)

	//test on different pass
	hash, err = bcrypt.GenerateFromPassword([]byte(passOne), bcrypt.MinCost)
	require.NoError(s.T(), err)

	passTwo := "aj43u8uqcmu9q8c9q2m3cmi"
	err = s.repository.CheckUserPassword(string(hash), passTwo)
	require.Error(s.T(), err)
}
