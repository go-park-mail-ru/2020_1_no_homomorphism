package repository

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type Albums struct {
	Id         uint64    `gorm:"column:id"`
	Name       string    `gorm:"column:name"`
	Image      string    `gorm:"column:image"`
	Release    time.Time `gorm:"column:release"`
	ArtistName string    `gorm:"column:artist_name"`
	ArtistId   uint64    `gorm:"column:artist_id"`
}

type DbAlbumRepository struct {
	db *gorm.DB
}

func NewDbAlbumRepository(database *gorm.DB) DbAlbumRepository {
	return DbAlbumRepository{
		db: database,
	}
}

func toSearchModel(albums Albums) models.AlbumSearch {
	return models.AlbumSearch{
		AlbumID:    strconv.FormatUint(albums.Id, 10),
		AlbumName:  albums.Name,
		ArtistID:   strconv.FormatUint(albums.ArtistId, 10),
		ArtistName: albums.ArtistName,
		Image:      albums.Image,
	}
}

func toModel(album Albums) models.Album {
	return models.Album{
		Id:         strconv.FormatUint(album.Id, 10),
		Name:       album.Name,
		Image:      album.Image,
		Release:    album.Release.Format("02-01-2006"),
		ArtistName: album.ArtistName,
		ArtistId:   strconv.FormatUint(album.ArtistId, 10),
	}
}

func (ar *DbAlbumRepository) GetUserAlbums(id string) ([]models.Album, error) {
	var dbAlbum []Albums

	sqlQuery := "SELECT album_id as id, album_name as name, album_image as image, artist_name, artist_id FROM user_albums WHERE user_id = ?"

	db := ar.db.
		Raw(sqlQuery, id).
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

	db := ar.db.
		Where("id = ?", id).
		Find(&dbAlbum)

	err := db.Error
	if err != nil {
		return models.Album{}, err
	}
	return toModel(dbAlbum), nil
}

func (ar *DbAlbumRepository) GetBoundedAlbumsByArtistId(id string, start, end uint64) ([]models.Album, error) {
	var dbAlbum []Albums
	limit := end - start

	db := ar.db.
		Where("artist_id = ?", id).
		Order("release").
		Limit(limit).
		Offset(start).
		Find(&dbAlbum)

	err := db.Error
	if err != nil {
		return []models.Album{}, err
	}

	albumsArray := make([]models.Album, len(dbAlbum))
	for i, elem := range dbAlbum {
		albumsArray[i] = toModel(elem)
	}
	return albumsArray, nil
}

func (ar *DbAlbumRepository) Search(text string, count uint) ([]models.AlbumSearch, error) {
	var albums []Albums

	db := ar.db.
		Table("albums").
		Where("name ILIKE ?", "%"+text+"%").
		Limit(count).
		Find(&albums)

	if err := db.Error; err != nil {
		return nil, fmt.Errorf("failed to search albums: %v", err)
	}

	albumSearch := make([]models.AlbumSearch, len(albums))
	for i, elem := range albums {
		albumSearch[i] = toSearchModel(elem)
	}
	return albumSearch, nil
}
