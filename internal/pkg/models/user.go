package models

type User struct {
	Id       string `json:"id"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Sex      string `json:"sex"`
	Image    string `json:"image"`
	Email    string `json:"email"`
}

type UserSettings struct {
	NewPassword string `json:"new_password"`
	User
}

type UserSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserStat struct {
	UserId    string `json:"user_id"`
	Tracks    uint64 `json:"tracks"`
	Albums    uint64 `json:"albums"`
	Playlists uint64 `json:"playlists"`
	Artists   uint64 `json:"artists"`
}
