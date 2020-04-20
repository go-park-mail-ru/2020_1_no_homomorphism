package models

type Playlist struct {
	Id     string `json:"id"`
	Name   string `json:"name,omitempty"`
	Image  string `json:"image,omitempty"`
	UserId string `json:"-"`
}

type PlaylistTracks struct {
	PlaylistID string `json:"playlist_id"`
	TrackID    string `json:"track_id"`
	Index      string `json:"index,omitempty"`
	Image      string `json:"image"`
}
