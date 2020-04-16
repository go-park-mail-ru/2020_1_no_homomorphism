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
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"regexp"
	"testing"
)

type Suite struct {
	suite.Suite
	DB     *gorm.DB
	mock   sqlmock.Sqlmock
	tracks []models.Track

	repository DbTrackRepository
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

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
	tr3 := models.Track{
		Id:       "2452345",
		Name:     "test23452345",
		Artist:   "artist-name",
		Duration: 345,
		Link:     "test_link3",
	}
	s.tracks = []models.Track{tr1, tr2, tr3}

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
	testTrack := s.tracks[0]

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "full_track_info" WHERE (track_id = $1)`)).
		WithArgs(testTrack.Id).
		WillReturnRows(sqlmock.NewRows([]string{"track_id", "track_name", "artist_name", "duration", "link"}).
			AddRow(testTrack.Id, testTrack.Name, testTrack.Artist, testTrack.Duration, testTrack.Link))

	res, err := s.repository.GetTrackById(testTrack.Id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(models.Track{
		Id:       fmt.Sprint(testTrack.Id),
		Name:     testTrack.Name,
		Artist:   testTrack.Artist,
		Duration: testTrack.Duration,
		Link:     testTrack.Link,
	}, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WithArgs(testTrack.Id).WillReturnError(dbError)

	_, err = s.repository.GetTrackById(testTrack.Id)

	require.Error(s.T(), err)
}

func (s *Suite) TestGetBoundedPlaylistTracks() {
	plId := "4123"

	tr1 := s.tracks[0]
	tr2 := s.tracks[1]
	tr3 := s.tracks[2]

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tracks_in_playlist" WHERE (playlist_id = $1) ORDER BY "index" LIMIT 3 OFFSET 0`)).
		WithArgs(plId).
		WillReturnRows(sqlmock.NewRows([]string{"track_id", "track_name", "artist_name", "duration", "link"}).
			AddRow(tr1.Id, tr1.Name, tr1.Artist, tr1.Duration, tr1.Link).
			AddRow(tr2.Id, tr2.Name, tr2.Artist, tr2.Duration, tr2.Link).
			AddRow(tr3.Id, tr3.Name, tr3.Artist, tr3.Duration, tr3.Link))

	res, err := s.repository.GetBoundedTracksByPlaylistId(plId, 0, 3)

	require.NoError(s.T(), err)
	for i, elem := range res {
		require.Nil(s.T(), deep.Equal(s.tracks[i], elem))
	}

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WithArgs(plId).WillReturnError(dbError)

	_, err = s.repository.GetBoundedTracksByPlaylistId(plId, 0, 3)

	require.Error(s.T(), err)
}

func (s *Suite) TestGetAlbumTracks() {
	aId := "4123"

	tr1 := s.tracks[0]
	tr2 := s.tracks[1]

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tracks_in_album" WHERE (album_id = $1) ORDER BY "index" LIMIT 2 OFFSET 0`)).
		WithArgs(aId).
		WillReturnRows(sqlmock.NewRows([]string{"track_id", "track_name", "artist_name", "duration", "link"}).
			AddRow(tr1.Id, tr1.Name, tr1.Artist, tr1.Duration, tr1.Link).AddRow(tr2.Id, tr2.Name, tr2.Artist, tr2.Duration, tr2.Link))

	res, err := s.repository.GetBoundedTracksByAlbumId(aId, 0, 2)

	require.NoError(s.T(), err)
	for i, elem := range res {
		require.Nil(s.T(), deep.Equal(s.tracks[i], elem))
	}

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WithArgs(aId).WillReturnError(dbError)

	_, err = s.repository.GetBoundedTracksByAlbumId(aId, 0, 2)

	require.Error(s.T(), err)
}
