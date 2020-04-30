package models

type Track struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Duration uint   `json:"duration"`
	Image    string `json:"image"`
	ArtistID string `json:"artist_id"`
	Link     string `json:"link"`
}

type TrackSearch struct {
	TrackID    string `json:"id"`
	TrackName  string `json:"name"`
	ArtistName string `json:"artist"`
	ArtistID   string `json:"artist_id"`
	Image      string `json:"image"`
}
