package repository

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/jinzhu/gorm"
	"strconv"
)

type Playlists struct {
	Id     uint64 `gorm:"column:id"`
	Name   string `gorm:"column:name"`
	Image  string `gorm:"column:image"`
	UserId uint64 `gorm:"column:user_id"`
}

type TrackInPlaylist struct {
	PlaylistID uint64 `gorm:"column:playlist_id"`
	TrackID    uint64 `gorm:"column:track_id"`
	Index      uint8  `gorm:"column:index"`
	Image      string `gorm:"column:image"`
}

type DbPlaylistRepository struct {
	db *gorm.DB
}

func NewDbPlaylistRepository(database *gorm.DB) DbPlaylistRepository {
	return DbPlaylistRepository{
		db: database,
	}
}

func toModel(pl Playlists) models.Playlist {
	return models.Playlist{
		Id:     strconv.FormatUint(pl.Id, 10),
		Name:   pl.Name,
		Image:  pl.Image,
		UserId: strconv.FormatUint(pl.UserId, 10),
	}
}

func (pr *DbPlaylistRepository) GetUserPlaylists(uId string) ([]models.Playlist, error) {
	var dbPlaylists []Playlists

	db := pr.db.
		Where("user_ID = ?", uId).
		Find(&dbPlaylists)

	err := db.Error
	if err != nil {
		return nil, err
	}

	playlists := make([]models.Playlist, len(dbPlaylists))
	for i, elem := range dbPlaylists {
		playlists[i] = toModel(elem)
	}
	return playlists, nil
}

func (pr *DbPlaylistRepository) GetPlaylistById(pId string) (models.Playlist, error) {
	var dbPlaylists Playlists

	db := pr.db.
		Where("id = ?", pId).
		First(&dbPlaylists)

	err := db.Error
	if err != nil {
		return models.Playlist{}, err
	}
	return toModel(dbPlaylists), nil
}

func (pr *DbPlaylistRepository) CreatePlaylist(name string, uID string) (plID string, err error) {
	userID, err := strconv.ParseUint(uID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse uID: %v", err)
	}

	newPlaylist := Playlists{
		Name:   name,
		UserId: userID,
	}

	if err := pr.db.Create(&newPlaylist).Error; err != nil {
		return "", fmt.Errorf("failed to create playlist: %v", err)
	}

	return strconv.FormatUint(newPlaylist.Id, 10), nil
}

func (pr *DbPlaylistRepository) AddTrackToPlaylist(plTracks models.PlaylistTracks) error {
	playlistID, err := strconv.ParseUint(plTracks.PlaylistID, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse plID: %v", err)
	}
	trackID, err := strconv.ParseUint(plTracks.TrackID, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse trackID: %v", err)
	}

	newRelation := TrackInPlaylist{
		PlaylistID: playlistID,
		TrackID:    trackID,
		Image:      plTracks.Image,
	}

	if err := pr.db.Table("playlist_tracks").Create(&newRelation).Error; err != nil {
		return fmt.Errorf("failed to create playlist:track relation: %v", err)
	}

	return nil
}
