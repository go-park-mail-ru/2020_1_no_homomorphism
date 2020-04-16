package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"strconv"
)

type Artists struct {
	Id    uint64 `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	Image string `gorm:"column:image"`
	Genre string `gorm:"column:genre"`
}

type DbArtistRepository struct {
	db *gorm.DB
}

func NewDbArtistRepository(database *gorm.DB) DbArtistRepository {
	return DbArtistRepository{
		db: database,
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
