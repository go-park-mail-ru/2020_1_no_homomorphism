package repository

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/jinzhu/gorm"
	"strconv"
)

type Artists struct {
	Id    uint64 `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	Image string `gorm:"column:image"`
	Genre string `gorm:"column:genre"`
}

type LikedArtists struct {
	UserID   int64 `gorm:"column:user_id"`
	ArtistID int64 `gorm:"column:artist_id"`
}

type DbArtistRepository struct {
	db *gorm.DB
}

func NewDbArtistRepository(database *gorm.DB) DbArtistRepository {
	return DbArtistRepository{
		db: database,
	}
}

func toSearchModel(artist Artists) models.ArtistSearch {
	return models.ArtistSearch{
		ArtistID: strconv.FormatUint(artist.Id, 10),
		Name:     artist.Name,
		Image:    artist.Image,
	}
}

func toModel(artist Artists) models.Artist {
	return models.Artist{
		Id:    strconv.FormatUint(artist.Id, 10),
		Name:  artist.Name,
		Image: artist.Image,
		Genre: artist.Genre,
	}
}

func (ar *DbArtistRepository) GetArtist(id string) (models.Artist, error) {
	var dbArtist Artists

	db := ar.db.Where("id = ?", id).Find(&dbArtist)
	err := db.Error
	if err != nil {
		return models.Artist{}, err
	}
	return toModel(dbArtist), nil
}

func (ar *DbArtistRepository) GetBoundedArtists(start, end uint64) ([]models.Artist, error) {
	var artists []Artists
	limit := end - start

	db := ar.db.Order("name").Limit(limit).Offset(start).Find(&artists)
	err := db.Error
	if err != nil {
		return []models.Artist{}, err
	}

	modArtists := make([]models.Artist, len(artists))
	for i, elem := range artists {
		modArtists[i] = toModel(elem)
	}
	return modArtists, nil
}

func (ar *DbArtistRepository) GetArtistStat(id string) (models.ArtistStat, error) {
	var stat models.ArtistStat

	db := ar.db.Table("artist_stat").Where("artist_id = ?", id).Find(&stat)
	err := db.Error
	if err != nil {
		return models.ArtistStat{}, err
	}
	return stat, nil
}

func (ar *DbArtistRepository) Search(text string, count uint) ([]models.ArtistSearch, error) {
	var artists []Artists

	db := ar.db.
		Table("artists").
		Where("name ILIKE ?", "%"+text+"%").
		Limit(count).
		Find(&artists)

	if err := db.Error; err != nil {
		return nil, fmt.Errorf("failed to search artists: %v", err)
	}

	artistSearch := make([]models.ArtistSearch, len(artists))
	for i, elem := range artists {
		artistSearch[i] = toSearchModel(elem)
	}
	return artistSearch, nil
}

func (ar *DbArtistRepository) IsSubscribed(aID string, uID string) bool {
	var likedArtists LikedArtists

	artistID, err1 := strconv.ParseUint(aID, 10, 64)
	userID, err2 := strconv.ParseUint(uID, 10, 64)
	if err1 != nil || err2 != nil {
		return false
	}

	db := ar.db.Where("user_id = ? and artist_id = ?", userID, artistID).Find(&likedArtists)
	if err := db.Error; err != nil {
		return false
	}
	return true
}

func (ar *DbArtistRepository) Subscription(aID string, uID string) error {
	var likedArtists LikedArtists

	artistID, err1 := strconv.ParseInt(aID, 10, 64)
	userID, err2 := strconv.ParseInt(uID, 10, 64)
	if err1 != nil || err2 != nil {
		return fmt.Errorf("failed to parse artistID or userID: %v, %v", err1, err2)
	}

	likedArtists.ArtistID = artistID
	likedArtists.UserID = userID

	db := ar.db.Table("liked_artists").Where("user_id = ? and artist_id = ?", userID, artistID).Find(&likedArtists)
	switch db.Error {
	case gorm.ErrRecordNotFound:
		db := ar.db.Exec("insert into liked_artists (artist_id, user_id) values (?, ?)", artistID, userID)
		if err := db.Error; err != nil {
			return fmt.Errorf("failed to insert in liked_artists: %v", err)
		}
	case nil:
		db := ar.db.Table("liked_artists").Where("user_id = ? and artist_id = ?", userID, artistID).Delete(&likedArtists)
		if err := db.Error; err != nil {
			return fmt.Errorf("failed to delete in liked_artists: %v", err)
		}
	default:
		return fmt.Errorf("failed to check liked_artists: %v", db.Error)
	}

	return nil
}

func (ar *DbArtistRepository) SubscriptionsList(uID string) ([]models.ArtistSearch, error) {
	var artists []models.ArtistSearch

	db := ar.db.
		Table("sub_artists").
		Where("user_id = ?", uID).
		Find(&artists)

	if err := db.Error; err != nil {
		return nil, fmt.Errorf("failed to get subscribed artists: %v", err)
	}

	return artists, nil
}

func toArtistAndSubsModel(r ArtistsAndSubscribers) models.ArtistAndSubscribers {
	return models.ArtistAndSubscribers{
		Artist:      r.Artist,
		Subscribers: r.subscribers,
	}
}

type ArtistsAndSubscribers struct {
	models.Artist
	subscribers uint64 `gorm:"subscribers" json:"subscribers"`
}

func (ar *DbArtistRepository) GetTopArtist() ([]models.ArtistAndSubscribers, error) {
	var topArtists []ArtistsAndSubscribers

	db := ar.db.Raw("SELECT artists.*, count(user_id) as subscribers " +
		"FROM user_artists	" +
		"JOIN artists on user_artists.artist_id = artists.id 	" +
		"GROUP BY artists.id " +
		"ORDER BY count(user_id) DESC " +
		"LIMIT 20").Scan(&topArtists)
	if err := db.Error; err != nil {
		return nil, err
	}
	newestReleasesModel := make([]models.ArtistAndSubscribers, len(topArtists))
	for i, r := range topArtists {
		newestReleasesModel[i] = toArtistAndSubsModel(r)
	}
}
