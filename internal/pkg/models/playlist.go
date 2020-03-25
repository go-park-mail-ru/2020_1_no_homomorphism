package models

type Playlist struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type UserPlaylists struct {
	Count     int        `json:"count"`
	Playlists []Playlist `json:"playlists"`
}

type PlaylistTracks struct {
	Playlist Playlist `json:"playlist"`
	Count    int      `json:"count"`
	Tracks   []Track  `json:"tracks"`
}
