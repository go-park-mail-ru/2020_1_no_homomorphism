package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
	. "no_homomorphism/internal/pkg/user/delivery"
)

func InitNewStorages() *Handler {
	api := NewUserHandler()
	user1 := models.User{
		Id:       uuid.NewV4(),
		Login:    "test",
		Name:     "Rita",
		Email:    "rita@margarita.tyt",
		Password: "$2a$04$0GzSltexrV9gQjFwv5BYuebu7/F13cX.NOupseJQUwqHWDucyBBgO",
	}

	user2 := models.User{
		Id:       uuid.NewV4(),
		Login:    "test2",
		Name:     "User2",
		Email:    "user2@da.tu",
		Password: "$2a$04$r/rWIhO8ptZAxheWs9cXmeG8fKhICfA5Gko3Qr61ae0.71CwjyODC",
	}

	user3 := models.User{
		Id:       uuid.NewV4(),
		Login:    "test3",
		Name:     "User3",
		Email:    "user3@da.tu",
		Password: "$2a$04$8G8SC41DvtOYD04qVizzbek.uL9zEI5zlQ3q2Cg.DYekuzMWFsoLa",
	}

	api.UserRepo.Users["test"] = &user1
	api.UserRepo.Users["test2"] = &user2
	api.UserRepo.Users["test3"] = &user3
	return api
}

func StartNew() {

	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://89.208.199.170:3000", "http://195.19.37.246:10982", "http://89.208.199.170:3001", "http://194.186.188.240"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	api := InitNewStorages()

	fmt.Printf("Starts server at 8081\n")
	r.HandleFunc("/signup", api.Create).Methods("POST")
	r.HandleFunc("/profile/settings", api.Update).Methods("PUT")
	err := http.ListenAndServe(":8081", c.Handler(r))
	if err != nil {
		fmt.Println(err)
		return
	}
}
