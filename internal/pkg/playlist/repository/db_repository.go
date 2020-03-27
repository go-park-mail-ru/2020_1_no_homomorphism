package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"no_homomorphism/internal/pkg/models"
)

type Playlists struct {
	Id     uint64 `sql:"AUTO_INCREMENT" gorm:"column:id"`
	Name   string `gorm:"column:name"`
	Image  string `gorm:"column:image"`
	UserId uint64 `gorm:"column:user_id"`
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
		Id:    fmt.Sprint(pl.Id),
		Name:  pl.Name,
		Image: pl.Image,
	}
}

func (pr *DbPlaylistRepository) GetUserPlaylists(uId string) ([]models.Playlist, error) {
	var dbPlaylists []Playlists

	db := pr.db.Where("user_ID = ?", uId).Find(&dbPlaylists)
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

	db := pr.db.Where("id = ?", pId).First(&dbPlaylists)
	err := db.Error
	if err != nil {
		return models.Playlist{}, err
	}
	return toModel(dbPlaylists), nil
}
