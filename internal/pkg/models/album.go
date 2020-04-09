package models

type Album struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	Release    string `json:"release"`
	ArtistName string `json:"artist_name"`
	ArtistId   string `json:"artist_id"`
}
