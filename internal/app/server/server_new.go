package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/session/repository"
	"no_homomorphism/internal/pkg/session/usecase"
	. "no_homomorphism/internal/pkg/user/delivery"
	repository2 "no_homomorphism/internal/pkg/user/repository"
	usecase2 "no_homomorphism/internal/pkg/user/usecase"
)

func InitNewStorages() *Handler {
	sesRep := repository.SessionRepository{
		Sessions: make(map[uuid.UUID]*models.User),
	}
	userRep := repository2.MemUserRepository{
		Users: map[string]*models.User{
			"test": &models.User{
				Id:       0,
				Login:    "test",
				Name:     "Rita",
				Email:    "rita@margarita.tyt",
				Password: "$2a$04$0GzSltexrV9gQjFwv5BYuebu7/F13cX.NOupseJQUwqHWDucyBBgO",
			},
			"test2": &models.User{
				Id:       1,
				Login:    "test2",
				Name:     "User2",
				Email:    "user2@da.tu",
				Password: "$2a$04$r/rWIhO8ptZAxheWs9cXmeG8fKhICfA5Gko3Qr61ae0.71CwjyODC",
			},
			"test3": &models.User{
				Id:       2,
				Login:    "test3",
				Name:     "User3",
				Email:    "user3@da.tu",
				Password: "$2a$04$8G8SC41DvtOYD04qVizzbek.uL9zEI5zlQ3q2Cg.DYekuzMWFsoLa",
			},
		},
		Count: 3,
	}

	SessionUC := usecase.SessionUseCase{
		Repository: &sesRep,
	}
	UserUC := usecase2.UserUseCase{
		Repository: &userRep,
		AvatarDir: "/static/img/",
	}
	h := &Handler{
		SessionUC: &SessionUC,
		UserUC:    &UserUC,
	}

	return h
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
	r.HandleFunc("/login", api.Login).Methods("POST")
	r.HandleFunc("/logout", api.Logout).Methods("DELETE")
	r.HandleFunc("/profile/settings", api.Update).Methods("PUT")
	r.HandleFunc("/profiles/{profile}", api.Profile)
	r.HandleFunc("/profile/me", api.SelfProfile).Methods("GET")
	r.HandleFunc("/image", api.UpdateAvatar).Methods("POST")
	r.HandleFunc("/debug", api.Debug)

	if err := http.ListenAndServe(":8081", c.Handler(r)); err != nil {
		log.Println(err)
		return
	}

}
