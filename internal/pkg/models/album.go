package models

type Album struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	ArtistId string `json:"artist_id"`
}

type UserAlbums struct {
	Count  int               `json:"count"`
	Albums []AlbumWithArtist `json:"albums"`
}

type AlbumWithArtist struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Artist Artist `json:"artist"`
}

type AlbumTracks struct {
	Album  Album   `json:"album"`
	Count  int     `json:"count"`
	Tracks []Track `json:"tracks"`
}
