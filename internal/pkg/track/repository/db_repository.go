package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"no_homomorphism/internal/pkg/models"
)

type DbTrack struct {
	Id       string `gorm:"column:track_id"`
	Name     string `gorm:"column:track_name"`
	Artist   string `gorm:"column:artist_name"`
	Duration uint   `gorm:"column:duration"`
	Link     string `gorm:"column:link"`
}

type DbTrackRepository struct {
	db *gorm.DB
}

func NewDbTrackRepo(db *gorm.DB) *DbTrackRepository {
	return &DbTrackRepository{
		db: db,
	}
}

func toModel(dbTrack DbTrack) *models.Track {
	return &models.Track{
		Id:       dbTrack.Id,
		Name:     dbTrack.Name,
		Artist:   dbTrack.Artist,
		Duration: dbTrack.Duration,
		Image:    "", //todo подумать над фото трека
		Link:     dbTrack.Link,
	}
}

func (tr *DbTrackRepository) GetTrackById(id uint) (*models.Track, error) {
	var track DbTrack
	db := tr.db.Raw("SELECT track_id,  track_name, artist_name, duration, link FROM full_track_info WHERE track_id = ?", id).Scan(&track)
	err := db.Error
	if err != nil {
		return nil, errors.New("query error: " + err.Error())
	}
	return toModel(track), nil
}

func (tr *DbTrackRepository) GetArtistTracks(artistId uint) ([]*models.Track, error) {
	var track []DbTrack
	db := tr.db.Raw("SELECT track_id,  track_name, artist_name, duration, link FROM full_track_info WHERE artist_id = ?", artistId).Scan(&track)
	err := db.Error
	if err != nil {
		return nil, errors.New("query error: " + err.Error())
	}
	var tracks []*models.Track

	rowsNum := db.RowsAffected
	var i int64
	for i = 0; i < rowsNum; i++ {
		tracks = append(tracks, toModel(track[i]))
	}
	return tracks, nil
}
