package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"no_homomorphism/internal/pkg/models"
	"testing"
)

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository *DbTrackRepository
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

	s.repository = NewDbTrackRepo(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetTrackByID() {
	name := "test-name"
	artist := "artist-name"
	link := "test_link"
	var id, duration uint
	duration = 123
	id = 12345

	s.mock.ExpectQuery("SELECT track_id, track_name, artist_name, duration, link FROM full_track_info WHERE track_id = ?").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"track_id", "track_name", "artist_name", "duration", "link"}).
			AddRow(id, name, artist, duration, link))

	res, err := s.repository.GetTrackById(id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&models.Track{
		Id:       fmt.Sprint(id),
		Name:     name,
		Artist:   artist,
		Duration: duration,
		Link:     link,
	}, res))
}
