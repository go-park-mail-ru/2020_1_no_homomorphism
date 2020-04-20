package models

type Track struct { //TODO add ArtistId field
	Id       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Duration uint   `json:"duration"`
	Image    string `json:"image"`
	ArtistID string `json:"artist_id"`
	Link     string `json:"link"`
}
