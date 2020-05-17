package models

type Playlist struct {
	Id      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Image   string `json:"image,omitempty"`
	UserId  string `json:"user_id"`
	Private bool   `json:"private,omitempty"`
}

type PlaylistTracks struct {
	PlaylistID string `json:"playlist_id"`
	TrackID    string `json:"track_id"`
	Index      string `json:"index,omitempty"`
	Image      string `json:"image"`
}

type PlaylistsID struct {
	IDs []string `json:"playlists"`
}

type PlaylistTracksArray struct {
	Id     string  `json:"id"`
	Tracks []Track `json:"tracks"`
}
