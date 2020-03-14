package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"no_homomorphism/internal/pkg/middleware"
	sessionRepo "no_homomorphism/internal/pkg/session/repository"
	sessionUC "no_homomorphism/internal/pkg/session/usecase"
	trackDelivery "no_homomorphism/internal/pkg/track/delivery"
	trackRepo "no_homomorphism/internal/pkg/track/repository"
	trackUC "no_homomorphism/internal/pkg/track/usecase"
	userDelivery "no_homomorphism/internal/pkg/user/delivery"
	userRepo "no_homomorphism/internal/pkg/user/repository"
	userUC "no_homomorphism/internal/pkg/user/usecase"
	"no_homomorphism/pkg/logger"
	"os"
)

func InitNewHandler(mainLogger *logger.MainLogger, db *gorm.DB) (*userDelivery.Handler, *trackDelivery.TrackHandler, *middleware.Middleware) {
	sesRep := sessionRepo.NewSessionRepository()
	//userRep := userRepo.NewTestMemUserRepository()
	trackRep := trackRepo.NewTestTrackRepo()
	dbRep := userRepo.NewTestDbUserRepository(db, "/static/img/avatar/default.png") //todo add to config

	SessionUC := sessionUC.SessionUseCase{
		Repository: sesRep,
	}
	UserUC := userUC.UserUseCase{
		Repository: dbRep,
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
	connStr := "user=postgres password=postgres dbname=music_app" //TODO получать из конфига

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to start db: " + err.Error())
	}
	defer db.Close()

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

	handler, trackHandler, m := InitNewHandler(customLogger, db)

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
