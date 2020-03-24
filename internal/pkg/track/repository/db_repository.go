package repository

import (
	"errors"
	"fmt"
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

func NewDbTrackRepo(db *gorm.DB) DbTrackRepository {
	return DbTrackRepository{
		db: db,
	}
}

func toModel(dbTrack DbTrack) models.Track {
	return models.Track{
		Id:       dbTrack.Id,
		Name:     dbTrack.Name,
		Artist:   dbTrack.Artist,
		Duration: dbTrack.Duration,
		Image:    "", //todo подумать над фото трека
		Link:     dbTrack.Link,
	}
}

func (tr *DbTrackRepository) GetTrackById(id uint64) (models.Track, error) {
	var track DbTrack
	db := tr.db.Raw("SELECT track_id,  track_name, artist_name, duration, link FROM full_track_info WHERE track_id = ?", id).Scan(&track)
	err := db.Error
	if err != nil {
		return models.Track{}, errors.New("query error: " + err.Error())
	}
	return toModel(track), nil
}

func (tr *DbTrackRepository) GetArtistTracks(artistId uint64) ([]models.Track, error) {
	var tracks []DbTrack
	db := tr.db.Raw("SELECT track_id,  track_name, artist_name, duration, link FROM full_track_info WHERE artist_id = ?", artistId).Scan(&tracks)
	err := db.Error
	if err != nil {
		return nil, errors.New("query error: " + err.Error())
	}

	modTracks := make([]models.Track, len(tracks))
	for i, elem := range tracks {
		modTracks[i] = toModel(elem)
	}
	return modTracks, nil
}

func (tr *DbTrackRepository) GetPlaylistTracks(plId uint64) ([]models.Track, error) {
	var tracks []DbTrack
	db := tr.db.Raw("SELECT track_id, track_name, artist_name, duration, link FROM tracks_in_playlist WHERE playlist_id = ?", plId).
		Scan(&tracks)
	err := db.Error
	if err != nil {
		return nil, fmt.Errorf("failed to select query: %e", err)
	}
	modTracks := make([]models.Track, db.RowsAffected)

	for i, elem := range tracks {
		modTracks[i] = toModel(elem)
	}
	return modTracks, nil
}

func (tr *DbTrackRepository) GetAlbumTracks(aId uint64) ([]models.Track, error) {
	var tracks []DbTrack
	db := tr.db.Raw("SELECT track_id, track_name, artist_name, duration, link FROM tracks_in_album WHERE album_id = ?", aId).
		Scan(&tracks)
	err := db.Error
	if err != nil {
		return nil, fmt.Errorf("failed to select query: %e", err)
	}
	modTracks := make([]models.Track, db.RowsAffected)

	for i, elem := range tracks {
		modTracks[i] = toModel(elem)
	}
	return modTracks, nil
}
