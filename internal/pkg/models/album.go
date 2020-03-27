package models

type Album struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ArtistName string `json:"artist_name"`
	ArtistId   string `json:"artist_id"`
}
