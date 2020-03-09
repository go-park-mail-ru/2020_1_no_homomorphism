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

type Profile struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	Sex   string `json:"sex"`
	Image string `json:"image"`
	Email string `json:"email"`
}

type UserSettings struct {
	NewPassword string `json:"new_password"`
	User
}

type UserSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
