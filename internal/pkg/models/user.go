package models

import (
	uuid "github.com/satori/go.uuid"
)


type User struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Login    string    `json:"login"`
	Sex      string    `json:"sex"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
}

type Profile struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	Sex   string `json:"sex"`
	Image string `json:"image"`
	Email string `json:"email"`
}

type UserSettings struct {
	NewPassword string `json:"newPassword"`
	User
}

