package user

import "net/http"

type UseCase interface {

	GetTrack(w http.ResponseWriter, r *http.Request)
	PostImage(w http.ResponseWriter, r *http.Request)
	GetUserImage(w http.ResponseWriter, r *http.Request)
	Index(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	Settings(w http.ResponseWriter, r *http.Request)
	GetProfile(w http.ResponseWriter, r *http.Request)
	GetProfileByCookie(w http.ResponseWriter, r *http.Request)
	Debug(w http.ResponseWriter, r *http.Request)
	CheckSession(w http.ResponseWriter, r *http.Request)
}