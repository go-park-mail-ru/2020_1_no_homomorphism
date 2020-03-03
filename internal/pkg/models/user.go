package models

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Sex      string `json:"sex"`
	Password string `json:"password"`
	Email    string `json:"email"`
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

type UserSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
