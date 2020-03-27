package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"no_homomorphism/internal/pkg/models"
)

type Albums struct {
	Id         uint64 `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	Image      string `gorm:"column:image"`
	ArtistName string `gorm:"column:artist_name"`
	ArtistId   uint64 `gorm:"column:artist_id"`
}

type DbAlbumRepository struct {
	db *gorm.DB
}

func NewDbAlbumRepository(database *gorm.DB) DbAlbumRepository {
	return DbAlbumRepository{
		db: database,
	}
}

func toModel(album Albums) models.Album {
	return models.Album{
		Id:         fmt.Sprint(album.Id),
		Name:       album.Name,
		Image:      album.Image,
		ArtistName: album.ArtistName,
		ArtistId:   fmt.Sprint(album.ArtistId),
	}
}

func (ar *DbAlbumRepository) GetUserAlbums(id string) ([]models.Album, error) {
	var dbAlbum []Albums

	db := ar.db.Raw("SELECT album_id as id, album_name as name, album_image as image, artist_name, artist_id FROM user_albums WHERE user_id = ?", id).
		Scan(&dbAlbum)
	err := db.Error
	if err != nil {
		return nil, err
	}

	albums := make([]models.Album, len(dbAlbum))

	for i, elem := range dbAlbum {
		albums[i] = toModel(elem)
	}
	return albums, nil
}

func (ar *DbAlbumRepository) GetAlbumById(id string) (models.Album, error) {
	var dbAlbum Albums

	db := ar.db.Where("id = ?", id).Find(&dbAlbum)
	err := db.Error
	if err != nil {
		return models.Album{}, err
	}
	return toModel(dbAlbum), nil
}
