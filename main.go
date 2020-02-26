package main

import (
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	. "no_homomorphism/handlers"
	"no_homomorphism/models"
	"sync"
)

func main() {
	r := mux.NewRouter()
	mu := &sync.Mutex{}
	api := &MyHandler{
		Sessions:     make(map[uuid.UUID]uuid.UUID, 10),
		UsersStorage: models.NewUsersStorage(mu),
		TrackStorage: models.NewTrackStorage(mu),
		Mutex:        mu,
	}

	fmt.Printf("Starts server at 8080\n")
	r.HandleFunc("/", api.MainHandler)
	r.HandleFunc("/login", api.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", api.LogoutHandler).Methods("DELETE")
	r.HandleFunc("/signup", api.SignUpHandler).Methods("POST")
	r.HandleFunc("/image", api.PostImageHandler).Methods("POST")
	r.HandleFunc("/image", api.GetUserImageHandler).Methods("GET")
	r.HandleFunc("/profile/settings", api.SettingsHandler).Methods("PUT")
	r.HandleFunc("/track/{id:[0-9]+}", api.GetTrackHandler).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
