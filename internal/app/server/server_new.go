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
	"no_homomorphism/internal/pkg/track/delivery"
	repository3 "no_homomorphism/internal/pkg/track/repository"
	usecase3 "no_homomorphism/internal/pkg/track/usecase"
	. "no_homomorphism/internal/pkg/user/delivery"
	repository2 "no_homomorphism/internal/pkg/user/repository"
	usecase2 "no_homomorphism/internal/pkg/user/usecase"
)

func InitNewHandler() *Handler {
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
				Image:    "/static/img/avatar/default.png",
			},
			"test2": &models.User{
				Id:       1,
				Login:    "test2",
				Name:     "User2",
				Email:    "user2@da.tu",
				Password: "$2a$04$r/rWIhO8ptZAxheWs9cXmeG8fKhICfA5Gko3Qr61ae0.71CwjyODC",
				Image:    "/static/img/avatar/default.png",
			},
			"test3": &models.User{
				Id:       2,
				Login:    "test3",
				Name:     "User3",
				Email:    "user3@da.tu",
				Password: "$2a$04$8G8SC41DvtOYD04qVizzbek.uL9zEI5zlQ3q2Cg.DYekuzMWFsoLa",
				Image:    "/static/img/avatar/default.png",
			},
		},
		Count: 3,
	}

	SessionUC := usecase.SessionUseCase{
		Repository: &sesRep,
	}
	UserUC := usecase2.UserUseCase{
		Repository: &userRep,
		AvatarDir:  "/static/img/avatar/",
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
		AllowedOrigins:   []string{"http://89.208.199.170:3000", "http://195.19.37.246:10982", "http://89.208.199.170:3001","http://localhost:3000",},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := InitNewHandler()
	trackHandler := &delivery.TrackHandler{
		TrackUC: &usecase3.TrackUseCase{
			Repository: repository3.NewTestRepo(),
		},
	}

	fmt.Printf("Starts server at 8081\n")
	r.HandleFunc("/signup", handler.Create).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/logout", handler.Logout).Methods("DELETE")
	r.HandleFunc("/profile/settings", handler.Update).Methods("PUT")
	r.HandleFunc("/profiles/{profile}", handler.Profile)
	r.HandleFunc("/profile/me", handler.SelfProfile).Methods("GET")
	r.HandleFunc("/image", handler.UpdateAvatar).Methods("POST")
	r.HandleFunc("/track/{id:[0-9]+}", trackHandler.GetTrack).Methods("GET")
	r.HandleFunc("/debug", handler.Debug)
	r.HandleFunc("/user", handler.CheckAuth)

	if err := http.ListenAndServe(":8081", c.Handler(r)); err != nil {
		log.Println(err)
		return
	}

}
