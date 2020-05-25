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

type LikedAlbums struct {
	ArtistID string `gorm:"column:album_id"`
	UserID   string `gorm:"column:user_id"`
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
		albums[i].IsLiked = true
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

func (ar *DbAlbumRepository) RateAlbum(aID, uID string) error {
	liked := LikedAlbums{}

	db := ar.db.Raw("select * from liked_albums where user_id = ? and album_id = ?", uID, aID).Scan(&liked)
	switch db.Error {
	case gorm.ErrRecordNotFound:
		db := ar.db.Exec("insert into liked_albums(user_id, album_id) values (?, ?)", uID, aID)
		if err := db.Error; err != nil {
			return err
		}
	case nil:
		db := ar.db.Exec("delete from liked_albums where user_id = ? and album_id = ?", uID, aID)
		if err := db.Error; err != nil {
			return err
		}
	default:
		return db.Error
	}
	return nil
}

func (ar *DbAlbumRepository) CheckLike(aID, uID string) bool {
	var liked LikedAlbums

	db := ar.db.
		Table("liked_albums").
		Where("album_id = ? and user_id = ?", aID, uID).
		Find(&liked)

	if err := db.Error; err != nil {
		return false
	}

	return true
}

type NewestReleases struct {
	Albums
	artist_image string `gorm:"artist_image" json:"artist_image"`
}

func toNewestReleasesModel(r NewestReleases) models.NewestReleases{
   return models.NewestReleases{
	   Album:       toModel(r.Albums),
	   ArtistImage: r.artist_image,
   }
}

func (ar *DbAlbumRepository) GetNewestReleases(uID string, begin int, end int) ([]models.NewestReleases, error) {
	var newestReleases []NewestReleases
	dif := end - begin
	db := ar.db.Raw("SELECT albums.*, sub_artists.image as artist_image " +
		"FROM sub_artists " +
		"JOIN albums on sub_artists.artist_id = albums.artist_id " +
		"WHERE user_id = ? " +
		"ORDER BY release " +
		"LIMIT ? " +
		"OFFSET ? ", uID, dif, begin).Scan(&newestReleases)

	if err := db.Error; err != nil {
		return nil, err
	}
	newestReleasesModel := make([]models.NewestReleases, len(newestReleases))
	for i, r := range newestReleases {
		newestReleasesModel[i] = toNewestReleasesModel(r)
	}
	return newestReleasesModel, nil
}