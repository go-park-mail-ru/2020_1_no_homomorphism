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
	TrackID    string `json:"track_id"`
	TrackName  string `json:"track_name"`
	ArtistName string `json:"artist_name"`
	ArtistID   string `json:"artist_id"`
	Image      string `json:"image"`
}
