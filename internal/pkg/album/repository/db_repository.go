package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"no_homomorphism/internal/pkg/models"
)

type Albums struct {
	AlbumId    uint64 `gorm:"column:id"`
	AlbumName  string `gorm:"column:name"`
	AlbumImage string `gorm:"column:image"`
	ArtistId   uint64 `gorm:"column:artist_id"`
}

type Artists struct {
	Id          uint64 `gorm:"column:artists_id"`
	ArtistName  string `gorm:"column:artist_name"`
	ArtistImage string `gorm:"column:artist_image"`
	ArtistGenre string `gorm:"column:artist_genre"`
}
type AlbumWithArtist struct {
	Albums
	Artists
}

type DbAlbumRepository struct {
	db *gorm.DB
}

func NewDbAlbumRepository(database *gorm.DB) *DbAlbumRepository {
	return &DbAlbumRepository{
		db: database,
	}
}

func toModel(album Albums) models.Album {
	return models.Album{
		Id:       fmt.Sprint(album.AlbumId),
		Name:     album.AlbumName,
		Image:    album.AlbumImage,
		ArtistId: fmt.Sprint(album.ArtistId),
	}
}

func toFullModel(al AlbumWithArtist) models.AlbumWithArtist {
	return models.AlbumWithArtist{
		Id:    fmt.Sprint(al.AlbumId),
		Name:  al.AlbumName,
		Image: al.AlbumImage,
		Artist: models.Artist{
			Id:    fmt.Sprint(al.Id),
			Name:  al.ArtistName,
			Image: al.ArtistImage,
			Genre: al.ArtistGenre,
		},
	}
}

func (ar *DbAlbumRepository) GetUserAlbums(uId uint64) ([]models.AlbumWithArtist, error) {
	var dbAlbum []AlbumWithArtist

	db := ar.db.Raw("SELECT album_id as id, album_name as name, album_image as image, artist_id as artists_id, artist_id, artist_name, artist_genre, artist_image FROM user_albums WHERE user_id = ?", uId).
		Scan(&dbAlbum)
	err := db.Error
	if err != nil {
		return nil, err
	}

	albums := make([]models.AlbumWithArtist, len(dbAlbum))
	for i, elem := range dbAlbum {
		albums[i] = toFullModel(elem)
	}
	return albums, nil
}

func (ar *DbAlbumRepository) GetAlbumById(aId uint64) (models.Album, error) {
	var dbAlbum Albums

	db := ar.db.Where("id = ?", aId).Find(&dbAlbum)
	err := db.Error
	if err != nil {
		return models.Album{}, err
	}
	return toModel(dbAlbum), nil
}
