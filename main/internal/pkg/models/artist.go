package models

type Artist struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Genre string `json:"genre"`
}

type ArtistStat struct {
	ArtistId    string `json:"artist_id"`
	Tracks      uint64 `json:"tracks"`
	Albums      uint64 `json:"albums"`
	Subscribers uint64 `json:"subscribers"`
}

type ArtistSearch struct {
	ArtistID string `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
}
