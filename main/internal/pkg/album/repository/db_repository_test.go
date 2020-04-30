package repository

import (
	"database/sql"
	"errors"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository DbAlbumRepository
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

	s.repository = NewDbAlbumRepository(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestAlbumById() {
	album := models.Album{
		Id:       "2434234",
		Name:     "test-name",
		Release:  "12-01-1999",
		Image:    "img-test",
		ArtistId: "3487919",
	}

	loc := time.Local
	testTime := time.Date(1999, 1, 12, 0, 0, 0, 0, loc)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums" WHERE (id = $1)`)).
		WithArgs(album.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "release", "image", "artist_id"}).
			AddRow(album.Id, album.Name, testTime, album.Image, album.ArtistId))

	res, err := s.repository.GetAlbumById(album.Id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(album, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WithArgs(album.Id).WillReturnError(dbError)

	_, err = s.repository.GetAlbumById(album.Id)

	require.Error(s.T(), err)
}

func (s *Suite) TestGetUserAlbums() {
	album := []models.Album{{
		Id:         "123",
		Name:       "test-name",
		Image:      "img-test",
		Release:    "12-01-1999",
		ArtistName: "artist_name",
		ArtistId:   "1234123",
	},
	}

	loc := time.Local
	testTime := time.Date(1999, 1, 12, 0, 0, 0, 0, loc)

	s.mock.ExpectQuery("SELECT album_id as id, album_name as name, album_image as image, artist_name, artist_id FROM user_albums WHERE user_id = ?").
		WithArgs(album[0].Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "release", "image", "artist_id", "artist_name"}).
			AddRow(album[0].Id, album[0].Name, testTime, album[0].Image, album[0].ArtistId, album[0].ArtistName))

	res, err := s.repository.GetUserAlbums(album[0].Id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(album, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WithArgs(album[0].Id).WillReturnError(dbError)

	_, err = s.repository.GetUserAlbums(album[0].Id)

	require.Error(s.T(), err)
}

func (s *Suite) TestSearch() {
	album := []models.AlbumSearch{
		{
			AlbumID:    "23423",
			AlbumName:  "testName",
			ArtistID:   "124252",
			ArtistName: "TestArtist",
			Image:      "default.png",
		},
	}

	text := "artis"
	count := 5

	s.mock.ExpectQuery("SELECT").
		WithArgs("%" + text + "%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "artist_id", "artist_name", "image"}).
			AddRow(album[0].AlbumID, album[0].AlbumName, album[0].ArtistID, album[0].ArtistName, album[0].Image))

	res, err := s.repository.Search(text, uint(count))

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(album, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WillReturnError(dbError)

	_, err = s.repository.Search(text, uint(count))

	require.Error(s.T(), err)
}
