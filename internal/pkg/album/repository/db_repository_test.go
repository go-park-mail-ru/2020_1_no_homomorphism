package repository

import (
	"database/sql"
	"errors"
	"fmt"
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

	repository *DbAlbumRepository
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
	var id, artistId uint64
	artistId = 123
	id = 12345

	album := models.Album{
		Id:       fmt.Sprint(id),
		Name:     "test-name",
		Image:    "img-test",
		ArtistId: fmt.Sprint(artistId),
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums" WHERE (id = $1)`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image", "artist_id"}).
			AddRow(id, album.Name, album.Image, artistId))

	res, err := s.repository.GetAlbumById(id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(album, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WithArgs(id).WillReturnError(dbError)

	_, err = s.repository.GetAlbumById(id)

	require.Error(s.T(), err)
}

func (s *Suite) TestGetUserAlbums() {
	var id, artistId uint64
	artistId = 123
	id = 12345

	album := []models.AlbumWithArtist{{
		Id:    fmt.Sprint(id),
		Name:  "test-name",
		Image: "img-test",
		Artist: models.Artist{
			Id:    fmt.Sprint(artistId),
			Name:  "curName",
			Image: "keklol123",
			Genre: "rock",
		},
	},
	}

	s.mock.ExpectQuery("SELECT album_id as id, album_name as name, album_image as image, artist_id as artists_id, artist_id, artist_name, artist_genre, artist_image FROM user_albums WHERE user_id = ?").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image", "artist_id","artists_id", "artist_name", "artist_genre", "artist_image"}).
			AddRow(id, album[0].Name, album[0].Image, artistId, artistId, album[0].Artist.Name, album[0].Artist.Genre, album[0].Artist.Image))

	res, err := s.repository.GetUserAlbums(id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(album, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WithArgs(id).WillReturnError(dbError)

	_, err = s.repository.GetUserAlbums(id)

	require.Error(s.T(), err)
}
