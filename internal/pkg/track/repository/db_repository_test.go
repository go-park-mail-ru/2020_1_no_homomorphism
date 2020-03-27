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
	"testing"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository DbTrackRepository
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
	id := "12345"
	var duration uint
	duration = 123

	s.mock.ExpectQuery("SELECT track_id, track_name, artist_name, duration, link FROM full_track_info WHERE track_id = ?").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"track_id", "track_name", "artist_name", "duration", "link"}).
			AddRow(id, name, artist, duration, link))

	res, err := s.repository.GetTrackById(id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(models.Track{
		Id:       fmt.Sprint(id),
		Name:     name,
		Artist:   artist,
		Duration: duration,
		Link:     link,
	}, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WithArgs(id).WillReturnError(dbError)

	_, err = s.repository.GetTrackById(id)

	require.Error(s.T(), err)
}

func (s *Suite) TestGetPlaylistTracks() {
	plId := "4123"

	tr1 := models.Track{
		Id:       "12345",
		Name:     "test-name",
		Artist:   "artist-name",
		Duration: 123,
		Link:     "test_link1",
	}
	tr2 := models.Track{
		Id:       "4532",
		Name:     "test-namqweqwee",
		Artist:   "artist-name",
		Duration: 5235,
		Link:     "test_link2",
	}

	trs := []models.Track{tr1, tr2}

	s.mock.ExpectQuery("SELECT track_id, track_name, artist_name, duration, link FROM tracks_in_playlist WHERE playlist_id = ?").
		WithArgs(plId).
		WillReturnRows(sqlmock.NewRows([]string{"track_id", "track_name", "artist_name", "duration", "link"}).
			AddRow(tr1.Id, tr1.Name, tr1.Artist, tr1.Duration, tr1.Link).AddRow(tr2.Id, tr2.Name, tr2.Artist, tr2.Duration, tr2.Link))

	res, err := s.repository.GetPlaylistTracks(plId)

	require.NoError(s.T(), err)
	for i, elem := range res {
		require.Nil(s.T(), deep.Equal(trs[i], elem))
	}

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WithArgs(plId).WillReturnError(dbError)

	_, err = s.repository.GetPlaylistTracks(plId)

	require.Error(s.T(), err)
}

func (s *Suite) TestGetAlbumTracks() {
	aId := "4123"

	tr1 := models.Track{
		Id:       "12345",
		Name:     "test-name",
		Artist:   "artist-name",
		Duration: 123,
		Link:     "test_link1",
	}
	tr2 := models.Track{
		Id:       "4532",
		Name:     "test-namqweqwee",
		Artist:   "artist-name",
		Duration: 5235,
		Link:     "test_link2",
	}

	trs := []models.Track{tr1, tr2}

	s.mock.ExpectQuery("SELECT track_id, track_name, artist_name, duration, link FROM tracks_in_album WHERE album_id = ?").
		WithArgs(aId).
		WillReturnRows(sqlmock.NewRows([]string{"track_id", "track_name", "artist_name", "duration", "link"}).
			AddRow(tr1.Id, tr1.Name, tr1.Artist, tr1.Duration, tr1.Link).AddRow(tr2.Id, tr2.Name, tr2.Artist, tr2.Duration, tr2.Link))

	res, err := s.repository.GetTracksByAlbumId(aId)

	require.NoError(s.T(), err)
	for i, elem := range res {
		require.Nil(s.T(), deep.Equal(trs[i], elem))
	}

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WithArgs(aId).WillReturnError(dbError)

	_, err = s.repository.GetTracksByAlbumId(aId)

	require.Error(s.T(), err)
}
