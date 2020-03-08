package models

type Track struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Duration uint   `json:"duration"`
	Image    string `json:"image"`
	Link     string `json:"link"`
}
