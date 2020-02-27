package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	uuid "github.com/satori/go.uuid"
	"net/http"
	. "no_homomorphism/handlers"
	"no_homomorphism/models"
	"sync"
)

func InitStorages() *MyHandler {
	trackStorage := models.NewTrackStorage()
	userStorage := models.NewUsersStorage()

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
		Name: "Rita",
		Email: "rita@margarita.tyt",
		Password: "$2a$04$0GzSltexrV9gQjFwv5BYuebu7/F13cX.NOupseJQUwqHWDucyBBgO",
	}

	user2 := models.User{
		Id:       uuid.NewV4(),
		Login:    "test2",
		Name:   "User2",
		Email: "user2@da.tu",
		Password: "$2a$04$r/rWIhO8ptZAxheWs9cXmeG8fKhICfA5Gko3Qr61ae0.71CwjyODC",
	}

	user3 := models.User{
		Id:       uuid.NewV4(),
		Login:    "test3",
		Name:   "User3",
		Email: "user3@da.tu",
		Password: "$2a$04$8G8SC41DvtOYD04qVizzbek.uL9zEI5zlQ3q2Cg.DYekuzMWFsoLa",
	}

	api.UsersStorage.Users["test"] = &user1
	api.UsersStorage.Users["test2"] = &user2
	api.UsersStorage.Users["test3"] = &user3
	return api
}

func main() {
	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://89.208.199.170:3000", "http://195.19.37.246:10982", "http://89.208.199.170:3001"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	api := InitStorages()

	fmt.Printf("Starts server at 8081\n")
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
	r.HandleFunc("/debug", api.Debug)
	//handler := c.Handler(r)
	err := http.ListenAndServe(":8081", c.Handler(r))
	if err != nil {
		fmt.Println(err)
		return
	}
}
