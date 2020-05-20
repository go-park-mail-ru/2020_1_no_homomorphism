package repository

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"strconv"
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
	pl1 := models.Playlist{
		Id:     "342354",
		Name:   "name",
		Image:  "custom/img",
		UserId: "24123",
	}
	pl2 := models.Playlist{
		Id:     "423516514",
		Name:   "my_second_playlist",
		Image:  "custom/img/2",
		UserId: "24123",
	}
	pls := []models.Playlist{pl1, pl2}

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "playlists" WHERE (user_ID = $1)`)).
		WithArgs(pl1.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image", "user_id"}).
			AddRow(pl1.Id, pl1.Name, pl1.Image, pl1.UserId).AddRow(pl2.Id, pl2.Name, pl2.Image, pl2.UserId))

	res, err := s.repository.GetUserPlaylists(pl1.UserId)

	require.NoError(s.T(), err)

	for i, elem := range res {
		require.Nil(s.T(), deep.Equal(pls[i], elem))
	}

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(pl1.UserId).WillReturnError(dbError)

	_, err = s.repository.GetUserPlaylists(pl1.UserId)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, dbError)
}

func (s *Suite) TestGetPlaylistById() {
	pl1 := models.Playlist{
		Id:     "5234523",
		Name:   "name",
		Image:  "custom/img",
		UserId: "4123123",
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "playlists" WHERE (id = $1) ORDER BY "playlists"."id" ASC LIMIT 1`)).
		WithArgs(pl1.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image", "user_id"}).
			AddRow(pl1.Id, pl1.Name, pl1.Image, pl1.UserId))

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

func (s *Suite) TestGetUserPlaylistsIdByTrack() {
	pl1 := models.Playlist{
		Id:     "5234523",
		Name:   "name",
		Image:  "custom/img",
		UserId: "4123123",
	}

	trID := "234234"

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT playlist_ID as id FROM playlist_tracks as pt join playlists as p on p.id = pt.playlist_id`)).
		WithArgs(pl1.UserId, trID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image", "user_id"}).
			AddRow(pl1.Id, pl1.Name, pl1.Image, pl1.UserId))

	res, err := s.repository.GetUserPlaylistsIdByTrack(pl1.UserId, trID)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal([]string{pl1.Id}, res))

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(pl1.UserId, trID).WillReturnError(dbError)

	_, err = s.repository.GetUserPlaylistsIdByTrack(pl1.UserId, trID)

	require.Error(s.T(), err)
}

func (s *Suite) TestDeleteTrackFromPlaylist() {
	var plID int64 = 53452
	pl1 := models.Playlist{
		Id:     strconv.Itoa(int(plID)),
		Name:   "name",
		Image:  "custom/img",
		UserId: "4123123",
	}

	var trID int64 = 235235
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM`)).
		WithArgs(plID, trID).WillReturnResult(driver.RowsAffected(1))
	s.mock.ExpectCommit()
	err := s.repository.DeleteTrackFromPlaylist(pl1.Id, strconv.FormatInt(trID, 10))

	require.NoError(s.T(), err)

	//test on db error
	dbError := errors.New("db_error")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE`)).
		WithArgs(plID, trID).WillReturnError(dbError)
	s.mock.ExpectRollback()

	err = s.repository.DeleteTrackFromPlaylist(pl1.Id, strconv.FormatInt(trID, 10))

	require.Error(s.T(), err)
}

func (s *Suite) TestAddTrackToPlaylist() {
	var plID int64 = 53452
	var tID int64 = 23423425
	var index int64 = 0

	pl := models.PlaylistTracks{
		PlaylistID: strconv.FormatInt(plID, 10),
		TrackID:    strconv.FormatInt(tID, 10),
		Index:      strconv.FormatInt(index, 10),
		Image:      "default",
	}

	query := `INSERT INTO "playlist_tracks" ("playlist_id","track_id","index","image") VALUES ($1,$2,$3,$4) RETURNING "playlist_tracks"."playlist_id"`

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(plID, tID, index, pl.Image).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(plID))

	s.mock.ExpectCommit()

	err := s.repository.AddTrackToPlaylist(pl)

	require.NoError(s.T(), err)

	//test on db error
	dbError := errors.New("db_error")

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(plID, tID, index, pl.Image).WillReturnError(dbError)

	s.mock.ExpectRollback()

	err = s.repository.AddTrackToPlaylist(pl)

	require.Error(s.T(), err)
}
