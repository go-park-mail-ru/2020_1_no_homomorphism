package models

type SearchResult struct {
	Artists []ArtistSearch `json:"artists"`
	Albums  []AlbumSearch  `json:"albums"`
	Tracks  []TrackSearch  `json:"tracks"`
}
