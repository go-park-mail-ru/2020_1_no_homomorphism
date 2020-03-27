package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"no_homomorphism/internal/pkg/models"
	"regexp"
	"testing"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository DbPlaylistRepository
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	s.DB.LogMode(false)

	s.repository = NewDbPlaylistRepository(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetUserPlaylists() {
	userId := "24123"

	pl1 := models.Playlist{
		Id:    "342354",
		Name:  "name",
		Image: "custom/img",
	}
	pl2 := models.Playlist{
		Id:    "423516514",
		Name:  "my_second_playlist",
		Image: "custom/img/2",
	}
	pls := []models.Playlist{pl1, pl2}

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "playlists" WHERE (user_ID = $1)`)).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image", "user_ID"}).
			AddRow(pl1.Id, pl1.Name, pl1.Image, userId).AddRow(pl2.Id, pl2.Name, pl2.Image, userId))

	res, err := s.repository.GetUserPlaylists(userId)

	require.NoError(s.T(), err)

	for i, elem := range res {
		require.Nil(s.T(), deep.Equal(pls[i], elem))
	}

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(userId).WillReturnError(dbError)

	_, err = s.repository.GetUserPlaylists(userId)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, dbError)
}

func (s *Suite) TestGetPlaylistById() {
	userId := "4123123"
	pl1 := models.Playlist{
		Id:    "5234523",
		Name:  "name",
		Image: "custom/img",
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "playlists" WHERE (id = $1) ORDER BY "playlists"."id" ASC LIMIT 1`)).
		WithArgs(pl1.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image", "user_ID"}).
			AddRow(pl1.Id, pl1.Name, pl1.Image, userId))

	res, err := s.repository.GetPlaylistById(pl1.Id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(pl1, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(pl1.Id).WillReturnError(dbError)

	_, err = s.repository.GetPlaylistById(pl1.Id)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, dbError)
}
