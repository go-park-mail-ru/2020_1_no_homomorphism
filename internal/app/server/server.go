package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"no_homomorphism/internal/pkg/middleware"
	sessionRepo "no_homomorphism/internal/pkg/session/repository"
	sessionUC "no_homomorphism/internal/pkg/session/usecase"
	trackDelivery "no_homomorphism/internal/pkg/track/delivery"
	trackRepo "no_homomorphism/internal/pkg/track/repository"
	trackUC "no_homomorphism/internal/pkg/track/usecase"
	userDelivery "no_homomorphism/internal/pkg/user/delivery"
	userRepo "no_homomorphism/internal/pkg/user/repository"
	userUC "no_homomorphism/internal/pkg/user/usecase"
)

func InitNewHandler() *userDelivery.Handler {
	sesRep := sessionRepo.NewSessionRepository()
	userRep := userRepo.NewTestMemUserRepository()

	SessionUC := sessionUC.SessionUseCase{
		Repository: sesRep,
	}
	UserUC := userUC.UserUseCase{
		Repository: userRep,
		AvatarDir:  "/static/img/avatar/",
	}
	h := &userDelivery.Handler{
		SessionUC: &SessionUC,
		UserUC:    &UserUC,
	}

	return h
}

func StartNew() {

	r := mux.NewRouter()

	c := middleware.InitCors()

	handler := InitNewHandler()
	trackHandler := &trackDelivery.TrackHandler{
		TrackUC: &trackUC.TrackUseCase{
			Repository: trackRepo.NewTestTrackRepo(),
		},
	}

	fmt.Println("Starts server at 8081")
	r.HandleFunc("/signup", handler.Create).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/logout", handler.Logout).Methods("DELETE")
	r.HandleFunc("/profile/settings", handler.Update).Methods("PUT")
	r.HandleFunc("/profiles/{profile}", handler.Profile)
	r.HandleFunc("/profile/me", handler.SelfProfile).Methods("GET")
	r.HandleFunc("/image", handler.UpdateAvatar).Methods("POST")
	r.HandleFunc("/track/{id:[0-9]+}", trackHandler.GetTrack).Methods("GET")
	// r.HandleFunc("/debug", handler.Debug)
	r.HandleFunc("/user", handler.CheckAuth)
	err := http.ListenAndServe(":8081", c.Handler(r))
	if err != nil {
		log.Println(err)
		return
	}
}
