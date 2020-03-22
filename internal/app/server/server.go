package server

import (
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"no_homomorphism/internal/pkg/middleware"
	playlistDelivery "no_homomorphism/internal/pkg/playlist/delivery"
	playlistRepo "no_homomorphism/internal/pkg/playlist/repository"
	playlistUC "no_homomorphism/internal/pkg/playlist/usecase"
	sessionDelivery "no_homomorphism/internal/pkg/session/delivery"
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
	"time"
)

func InitNewHandler(mainLogger *logger.MainLogger, db *gorm.DB, redis redis.Conn) (
	*userDelivery.Handler,
	*trackDelivery.TrackHandler,
	*playlistDelivery.PlaylistHandler,
	*middleware.Middleware) {
	sesRep := sessionRepo.NewRedisSessionManager(redis)
	trackRep := trackRepo.NewDbTrackRepo(db)
	playlistRep := playlistRepo.NewDbPlaylistRepository(db)
	dbRep := userRepo.NewDbUserRepository(db, "/static/img/avatar/default.png") //todo add to config

	PlaylistUC := playlistUC.PlaylistUseCase{
		PlRepository:    playlistRep,
		TrackRepository: trackRep,
	}

	SessionUC := sessionUC.SessionUseCase{
		Repository: sesRep,
	}

	SessionDelivery := sessionDelivery.SessionDelivery{
		UseCase:    &SessionUC,
		ExpireTime: 24 * 31 * time.Hour,
	}
	UserUC := userUC.UserUseCase{
		Repository: dbRep,
		AvatarDir:  "/static/img/avatar/",
	}
	TrackUC := trackUC.TrackUseCase{
		Repository: trackRep,
	}

	playlistHandler := &playlistDelivery.PlaylistHandler{
		PlaylistUC: &PlaylistUC,
		Log:        mainLogger,
	}

	h := &userDelivery.Handler{
		SessionDelivery: &SessionDelivery,
		UserUC:          &UserUC,
		Log:             mainLogger,
	}

	trackHandler := &trackDelivery.TrackHandler{
		TrackUC: &TrackUC,
		Log:     mainLogger,
	}

	m := middleware.NewMiddleware(&SessionDelivery, &UserUC, &TrackUC, &PlaylistUC)

	return h, trackHandler, playlistHandler, m
}

func StartNew() {
	connStr := "user=postgres password=postgres dbname=music_app" //TODO получать из конфига

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to start db: " + err.Error())
	}
	defer db.Close()

	db.DB().SetMaxOpenConns(10)

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

	redisAddr := flag.String("addr", "redis://user:@localhost:6379/0", "redis addr")

	redisConn, err := redis.DialURL(*redisAddr)
	if err != nil {
		log.Fatalf("cant connect to redis")
	}

	user, track, playlist, m := InitNewHandler(customLogger, db, redisConn)

	r.HandleFunc("/profile/settings", user.Update).Methods("PUT")
	r.HandleFunc("/profile/me", user.SelfProfile).Methods("GET")
	r.HandleFunc("/profile/playlists", playlist.GetUserPlaylists).Methods("GET")
	r.HandleFunc("/playlists/{id:[0-9]+}", playlist.GetPlaylistTracks).Methods("GET")
	r.HandleFunc("/profiles/{profile}", user.Profile)
	r.HandleFunc("/image", user.UpdateAvatar).Methods("POST")
	r.HandleFunc("/user", user.CheckAuth)
	r.HandleFunc("/signup", user.Create).Methods("POST")
	r.HandleFunc("/login", user.Login).Methods("POST")
	r.HandleFunc("/logout", user.Logout).Methods("DELETE")
	r.HandleFunc("/track/{id:[0-9]+}", track.GetTrack).Methods("GET")
	authHandler := m.CheckAuthMiddleware(r)
	fmt.Println("Starts server at 8081")

	accessMiddleware := middleware.AccessLogMiddleware(authHandler, user.Log)
	panicMiddleware := middleware.PanicMiddleware(accessMiddleware, user.Log)

	err = http.ListenAndServe(":8081", c.Handler(panicMiddleware))
	if err != nil {
		log.Println(err)
		return
	}
}
