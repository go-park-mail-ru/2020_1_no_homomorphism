package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	uuid "github.com/satori/go.uuid"
	. "no_homomorphism/handlers"
	"no_homomorphism/models"
	"sync"
)

func main() {
	r := mux.NewRouter()
	api := MyHandler{Sessions: make(map[uuid.UUID]uuid.UUID, 10),
		UsersStorage: &models.UsersStorage{
			Users: map[string]*models.User{
				"test": {
					Id:       uuid.NewV4(),
					Login:    "test",
					Password: "123",
				},
				"test2": {
					Id:       uuid.NewV4(),
					Login:    "test2",
					Password: "456",
				},
				"test3": {
					Id:       uuid.NewV4(),
					Login:    "test3",
					Password: "789",
				},
			},
			Mutex: sync.RWMutex{},
		},
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
	r.HandleFunc("/profile/settings", api.SettingsHandler).Methods("PUT")
	r.HandleFunc("/profiles/{profile}", api.GetProfileHandler)
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
