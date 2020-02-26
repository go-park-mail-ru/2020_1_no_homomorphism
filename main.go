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
	trackStorage := models.NewTrackStorage()
	userStorage:= models.NewUsersStorage()

	api := &MyHandler{
		Sessions:     make(map[uuid.UUID]uuid.UUID, 10),
		UsersStorage: userStorage,
		TrackStorage: trackStorage,
		Mutex:        &sync.Mutex{},
		AvatarDir:    "/home/ubuntu/2020_1_no_homomorphism/static/img/avatar/",
	}

	user1 := models.User{
		Id:       uuid.NewV4(),
		Login:    "test",
		Password: "123",
	}

	user2 := models.User{
		Id:       uuid.NewV4(),
		Login:    "test2",
		Password: "456",
	}

	user3 := models.User{
		Id:       uuid.NewV4(),
		Login:    "test3",
		Password: "789",
	}

	api.UsersStorage.Users["test"] = &user1
	api.UsersStorage.Users["test2"] = &user2
	api.UsersStorage.Users["test3"] = &user3

	fmt.Printf("Starts server at 8080\n")
	r.HandleFunc("/", api.MainHandler)
	r.HandleFunc("/login", api.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", api.LogoutHandler).Methods("DELETE")
	r.HandleFunc("/signup", api.SignUpHandler).Methods("POST")
	r.HandleFunc("/profile/settings", api.SettingsHandler).Methods("PUT")
	r.HandleFunc("/profiles/{profile}", api.GetProfileHandler)
	r.HandleFunc("/profile/me", api.GetProfileByCookieHandler).Methods("GET")
	r.HandleFunc("/image", api.PostImageHandler).Methods("POST")
	//r.HandleFunc("/image", api.GetImageURLHandler).Methods("GET")
	r.HandleFunc("/image", api.GetUserImageHandler).Methods("GET")
	r.HandleFunc("/track/{id:[0-9]+}", api.GetTrackHandler).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
