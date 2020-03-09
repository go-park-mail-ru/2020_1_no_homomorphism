package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"no_homomorphism/pkg/logger"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

func InitNewHandler(mainLogger *logger.MainLogger) (*userDelivery.Handler, *trackDelivery.TrackHandler, *middleware.Middleware) {
	sesRep := sessionRepo.NewSessionRepository()
	userRep := userRepo.NewTestMemUserRepository()
	trackRep := trackRepo.NewTestTrackRepo()

	SessionUC := sessionUC.SessionUseCase{
		Repository: sesRep,
	}
	UserUC := userUC.UserUseCase{
		Repository: userRep,
		AvatarDir:  "/static/img/avatar/",
	}
	TrackUC := trackUC.TrackUseCase{
		Repository: trackRep,
	}

	h := &userDelivery.Handler{
		SessionUC: &SessionUC,
		UserUC:    &UserUC,
		Log:       mainLogger,
	}

	trackHandler := &trackDelivery.TrackHandler{
		TrackUC: &TrackUC,
		Log:     mainLogger,
	}
	m := middleware.NewMiddleware(&SessionUC, &UserUC, &TrackUC)

	return h, trackHandler, m
}

func StartNew() {

	r := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://89.208.199.170:3000", "http://195.19.37.246:10982", "http://89.208.199.170:3001", "http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		Debug:            false,
	})

	var customLogger *logger.MainLogger

	filename := "logfile.log"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logrus.Error("Failed to open logfile:", err)
		customLogger = logger.NewLogger(os.Stdout)
	} else {
		customLogger = logger.NewLogger(f)
	}
	defer f.Close()

	handler, trackHandler, m := InitNewHandler(customLogger)

	r.HandleFunc("/profile/settings", handler.Update).Methods("PUT")
	r.HandleFunc("/profile/me", handler.SelfProfile).Methods("GET")
	r.HandleFunc("/profiles/{profile}", handler.Profile)
	r.HandleFunc("/image", handler.UpdateAvatar).Methods("POST")
	r.HandleFunc("/user", handler.CheckAuth)
	r.HandleFunc("/signup", handler.Create).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/logout", handler.Logout).Methods("DELETE")
	r.HandleFunc("/track/{id:[0-9]+}", trackHandler.GetTrack).Methods("GET")
	authHandler := m.CheckAuthMiddleware(r)
	fmt.Println("Starts server at 8081")

	accessMiddleware := middleware.AccessLogMiddleware(authHandler, handler.Log)
	panicMiddleware := middleware.PanicMiddleware(accessMiddleware, handler.Log)

	err = http.ListenAndServe(":8081", c.Handler(panicMiddleware))
	if err != nil {
		log.Println(err)
		return
	}
}
