package models

type User struct {
	Id       uint   `json:"id"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Sex      string `json:"sex"`
	Image    string `json:"image"`
	Email    string `json:"email"`
}

type UserSettings struct {
	NewPassword string `json:"newPassword"`
	User
}

