package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
)

type Suite struct {
	suite.Suite
	DB      *gorm.DB
	mock    sqlmock.Sqlmock
	artists []models.Artist

	repository DbArtistRepository
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	tr1 := models.Artist{
		Id:    "12345",
		Name:  "test-name",
		Image: "default.png",
		Genre: "rock",
	}
	tr2 := models.Artist{
		Id:    "4532",
		Name:  "test-namqweqwee",
		Image: "default.png",
		Genre: "rap",
	}
	tr3 := models.Artist{
		Id:    "2452345",
		Name:  "test23452345",
		Image: "default.png",
		Genre: "classic",
	}
	s.artists = []models.Artist{tr1, tr2, tr3}

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	s.DB.LogMode(false)

	s.repository = NewDbArtistRepository(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetArtist() {
	testArtist := s.artists[0]

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artists" WHERE (id = $1)`)).
		WithArgs(testArtist.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image", "genre"}).
			AddRow(testArtist.Id, testArtist.Name, testArtist.Image, testArtist.Genre))

	res, err := s.repository.GetArtist(testArtist.Id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(models.Artist{
		Id:    fmt.Sprint(testArtist.Id),
		Name:  testArtist.Name,
		Image: testArtist.Image,
		Genre: testArtist.Genre,
	}, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WithArgs(testArtist.Id).WillReturnError(dbError)

	_, err = s.repository.GetArtist(testArtist.Id)

	require.Error(s.T(), err)
}

func (s *Suite) TestGetBoundedArtists() {
	testArtist := s.artists[0]
	testArtist2 := s.artists[1]

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artists" ORDER BY "name" LIMIT 5 OFFSET 0`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image", "genre"}).
			AddRow(testArtist.Id, testArtist.Name, testArtist.Image, testArtist.Genre).
			AddRow(testArtist2.Id, testArtist2.Name, testArtist2.Image, testArtist2.Genre))

	res, err := s.repository.GetBoundedArtists(0, 5)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(s.artists[0:2], res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WillReturnError(dbError)

	_, err = s.repository.GetBoundedArtists(0, 5)

	require.Error(s.T(), err)
}

func (s *Suite) TestSearch() {

	testArtist := []models.ArtistSearch{
		{
			ArtistID: "234234",
			Name:     "dfdsf",
			Image:    "default.png",
		},
	}

	s.mock.ExpectQuery(`SELECT * `).
		WithArgs("%" + testArtist[0].Name + "%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image"}).
			AddRow(testArtist[0].ArtistID, testArtist[0].Name, testArtist[0].Image))

	res, err := s.repository.Search(testArtist[0].Name, 5)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(testArtist, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WillReturnError(dbError)

	_, err = s.repository.Search(testArtist[0].Name, 5)

	require.Error(s.T(), err)
}

func (s *Suite) TestGetArtistStat() {

	stat := models.ArtistStat{
		ArtistId:    "14124",
		Tracks:      52134,
		Albums:      532,
		Subscribers: 42324513,
	}

	query := `SELECT * FROM "artist_stat" WHERE (artist_id = $1)`

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(stat.ArtistId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "tracks", "albums", "subscribers"}).
			AddRow(stat.ArtistId, stat.Tracks, stat.Albums, stat.Subscribers))

	res, err := s.repository.GetArtistStat(stat.ArtistId)

	stat.ArtistId = ""

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(stat, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery("SELECT").
		WillReturnError(dbError)

	_, err = s.repository.GetArtistStat(stat.ArtistId)

	require.Error(s.T(), err)
}

func (s *Suite) TestIsSubscribed() {
	query := `SELECT`
	var aID uint64 = 1231412
	var uID uint64 = 5135

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(uID, aID).
		WillReturnRows(sqlmock.NewRows([]string{"artist_id", "user_id"}).
			AddRow(aID, uID))

	res := s.repository.IsSubscribed(fmt.Sprint(aID), fmt.Sprint(uID))

	require.Equal(s.T(), true, res)

	//test on false
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(uID, aID).
		WillReturnError(gorm.ErrRecordNotFound)

	res = s.repository.IsSubscribed(fmt.Sprint(aID), fmt.Sprint(uID))

	require.Equal(s.T(), false, res)
}

func (s *Suite) TestSubscriptionsList() {
	query := `SELECT * FROM "sub_artists"  WHERE (user_id = $1)`
	uID := "5135"

	testArtist := []models.ArtistSearch{
		{
			ArtistID: "234234",
			Name:     "dfdsf",
			Image:    "default.png",
		},
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(uID).
		WillReturnRows(sqlmock.NewRows([]string{"artist_id", "name", "image"}).
			AddRow(testArtist[0].ArtistID, testArtist[0].Name, testArtist[0].Image))

	res, err := s.repository.SubscriptionsList(fmt.Sprint(uID))

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(testArtist, res))

	//test on false
	dbError := errors.New("test error")

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(uID).
		WillReturnError(dbError)

	_, err = s.repository.SubscriptionsList(fmt.Sprint(uID))

	require.Error(s.T(), err)
}
