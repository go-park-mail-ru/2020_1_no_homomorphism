package repository

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/jinzhu/gorm"
	"strconv"
)

type Tracks struct {
	Id       uint64 `gorm:"column:track_id"`
	Name     string `gorm:"column:track_name"`
	Artist   string `gorm:"column:artist_name"`
	ArtistID uint64 `gorm:"column:artist_id"`
	Duration uint   `gorm:"column:duration"`
	Image    string `gorm:"column:track_image"`
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

func toSearchModel(dbTrack Tracks) models.TrackSearch {
	return models.TrackSearch{
		TrackID:    strconv.FormatUint(dbTrack.Id, 10),
		TrackName:  dbTrack.Name,
		ArtistName: dbTrack.Artist,
		ArtistID:   strconv.FormatUint(dbTrack.ArtistID, 10),
		Image:      dbTrack.Image,
	}
}

func toModel(dbTrack Tracks) models.Track {
	return models.Track{
		Id:       strconv.FormatUint(dbTrack.Id, 10),
		Name:     dbTrack.Name,
		Artist:   dbTrack.Artist,
		ArtistID: strconv.FormatUint(dbTrack.ArtistID, 10),
		Duration: dbTrack.Duration,
		Image:    dbTrack.Image,
		Link:     dbTrack.Link,
	}
}

func (tr *DbTrackRepository) GetTrackById(id string) (models.Track, error) {
	var track Tracks

	db := tr.db.
		Table("full_track_info").
		Where("track_id = ?", id).
		Find(&track)

	err := db.Error
	if err != nil {
		return models.Track{}, fmt.Errorf("query error: %v", err)
	}
	return toModel(track), nil
}

func (tr *DbTrackRepository) GetBoundedTracksByArtistId(id string, start, end uint64) ([]models.Track, error) {
	var tracks []Tracks
	limit := end - start

	db := tr.db.
		Table("full_track_info").
		Where("artist_id = ?", id).
		Order("track_name").
		Limit(limit).
		Offset(start).
		Find(&tracks)

	err := db.Error
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}

	modTracks := make([]models.Track, len(tracks))
	for i, elem := range tracks {
		modTracks[i] = toModel(elem)
	}
	return modTracks, nil
}

func (tr *DbTrackRepository) GetBoundedTracksByPlaylistId(plId string, start, end uint64) ([]models.Track, error) {
	var tracks []Tracks
	limit := end - start

	db := tr.db.
		Table("tracks_in_playlist").
		Where("playlist_id = ?", plId).
		Order("index").
		Limit(limit).
		Offset(start).
		Find(&tracks)

	err := db.Error
	if err != nil {
		return nil, fmt.Errorf("failed to select query: %e", err)
	}

	modTracks := make([]models.Track, len(tracks))
	for i, elem := range tracks {
		modTracks[i] = toModel(elem)
	}
	return modTracks, nil
}

func (tr *DbTrackRepository) GetBoundedTracksByAlbumId(aId string, start, end uint64) ([]models.Track, error) {
	var tracks []Tracks
	limit := end - start

	db := tr.db.
		Table("tracks_in_album").
		Where("album_id = ?", aId).
		Order("index").
		Limit(limit).
		Offset(start).
		Find(&tracks)

	err := db.Error
	if err != nil {
		return nil, fmt.Errorf("failed to select query: %e", err)
	}

	modTracks := make([]models.Track, len(tracks))
	for i, elem := range tracks {
		modTracks[i] = toModel(elem)
	}
	return modTracks, nil
}

func (tr *DbTrackRepository) Search(text string, count uint) ([]models.TrackSearch, error) {
	var tracks []Tracks

	db := tr.db.
		Table("full_track_info").
		Where("track_name ILIKE ?", "%"+text+"%").
		Limit(count).
		Find(&tracks)

	if err := db.Error; err != nil {
		return nil, fmt.Errorf("failed to search tracks: %v", err)
	}

	trackSearch := make([]models.TrackSearch, len(tracks))
	for i, elem := range tracks {
		trackSearch[i] = toSearchModel(elem)
	}
	return trackSearch, nil
}
