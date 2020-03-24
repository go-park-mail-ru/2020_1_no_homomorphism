package models

type Artist struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Genre string `json:"genre"`
}
