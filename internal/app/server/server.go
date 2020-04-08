package server

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/kabukky/httpscerts"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	albumDelivery "no_homomorphism/internal/pkg/album/delivery"
	albumRepo "no_homomorphism/internal/pkg/album/repository"
	albumUC "no_homomorphism/internal/pkg/album/usecase"
	artistDelivery "no_homomorphism/internal/pkg/artist/delivery"
	artistRepo "no_homomorphism/internal/pkg/artist/repository"
	artistUC "no_homomorphism/internal/pkg/artist/usecase"
	"no_homomorphism/internal/pkg/constants"
	"no_homomorphism/internal/pkg/csrf"
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
)

func InitHandler(mainLogger *logger.MainLogger, db *gorm.DB, redis *redis.Pool, csrfToken csrf.CryptToken) (
	userDelivery.UserHandler,
	trackDelivery.TrackHandler,
	playlistDelivery.PlaylistHandler,
	albumDelivery.AlbumHandler,
	artistDelivery.ArtistHandler,
	middleware.MiddlewareManager) {

	sesRep := sessionRepo.NewRedisSessionManager(redis)
	trackRep := trackRepo.NewDbTrackRepo(db)
	playlistRep := playlistRepo.NewDbPlaylistRepository(db)
	albumRep := albumRepo.NewDbAlbumRepository(db)
	artistRep := artistRepo.NewDbArtistRepository(db)
	dbRep := userRepo.NewDbUserRepository(db, constants.AvatarDefault,  constants.AvatarDir) // todo add to config

	ArtistUC := artistUC.ArtistUseCase{
		ArtistRepository: &artistRep,
	}

	AlbumUC := albumUC.AlbumUseCase{
		AlbumRepository: &albumRep,
	}

	PlaylistUC := playlistUC.PlaylistUseCase{
		PlRepository: &playlistRep,
	}

	SessionUC := sessionUC.SessionUseCase{
		Repository: &sesRep,
	}

	SessionDelivery := sessionDelivery.SessionDelivery{
		UseCase:    &SessionUC,
		ExpireTime: constants.CookieExpireTime,
	}
	UserUC := userUC.UserUseCase{
		Repository: &dbRep,
	}
	TrackUC := trackUC.TrackUseCase{
		Repository: &trackRep,
	}

	playlistHandler := playlistDelivery.PlaylistHandler{
		PlaylistUC: &PlaylistUC,
		TrackUC:    &TrackUC,
		Log:        mainLogger,
	}

	artistHandler := artistDelivery.ArtistHandler{
		ArtistUC: &ArtistUC,
		TrackUC:  &TrackUC,
		Log:      mainLogger,
	}

	userHandler := userDelivery.UserHandler{
		SessionDelivery: &SessionDelivery,
		UserUC:          &UserUC,
		Log:             mainLogger,
		ImgTypes:        constants.AvatarTypes,
		CSRF:            csrfToken,
	}

	trackHandler := trackDelivery.TrackHandler{
		TrackUC: &TrackUC,
		Log:     mainLogger,
	}

	albumHandler := albumDelivery.AlbumHandler{
		AlbumUC: &AlbumUC,
		TrackUC: &TrackUC,
		Log:     mainLogger,
	}

	m := middleware.NewMiddlewareManager(&SessionDelivery, &UserUC, &TrackUC, &PlaylistUC, csrfToken)

	return userHandler, trackHandler, playlistHandler, albumHandler, artistHandler, m
}

func InitRouter(customLogger *logger.MainLogger, db *gorm.DB, redisConn *redis.Pool, csrfToken csrf.CryptToken) http.Handler {
	user, track, playlist, album, artist, m := InitHandler(customLogger, db, redisConn, csrfToken)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/users/settings", user.Update).Methods("PUT")
	r.HandleFunc("/api/v1/users/me", user.SelfProfile).Methods("GET")
	r.HandleFunc("/api/v1/users/playlists", playlist.GetUserPlaylists).Methods("GET")
	r.HandleFunc("/api/v1/users/albums", album.GetUserAlbums).Methods("GET")
	r.HandleFunc("/api/v1/playlists/{id:[0-9]+}", playlist.GetFullPlaylistById).Methods("GET")
	r.HandleFunc("/api/v1/albums/{id:[0-9]+}", album.GetFullAlbum).Methods("GET")
	r.HandleFunc("/api/v1/users/profiles/{profile}", user.Profile)
	r.HandleFunc("/api/v1/users/images", user.UpdateAvatar).Methods("POST")
	r.HandleFunc("/api/v1/users", user.CheckAuth)
	r.HandleFunc("/api/v1/users/signup", user.Create).Methods("POST")
	r.HandleFunc("/api/v1/users/login", user.Login).Methods("POST")
	r.HandleFunc("/api/v1/users/logout", user.Logout).Methods("DELETE")
	r.HandleFunc("/api/v1/tracks/{id:[0-9]+}", track.GetTrack).Methods("GET")
	r.Handle("/api/v1/artists/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", middleware.GetBoundedVars(artist.GetBoundedArtistTracks, user.Log)).Methods("GET")
	r.Handle("/api/v1/artists/{id:[0-9]+}/albums/{start:[0-9]+}/{end:[0-9]+}", middleware.GetBoundedVars(album.GetBoundedAlbumsByArtistId, user.Log)).Methods("GET")
	r.HandleFunc("/api/v1/artists/{start:[0-9]+}/{end:[0-9]+}", artist.GetBoundedArtists).Methods("GET")
	r.HandleFunc("/api/v1/artists/{id:[0-9]+}", artist.GetFullArtistInfo).Methods("GET")
	r.HandleFunc("/api/v1/artists/{id:[0-9]+}/stat", artist.GetArtistStat).Methods("GET")
	r.HandleFunc("/api/v1/users/{id:[0-9]+}/stat", user.GetUserStat).Methods("GET")
	r.Handle("/api/v1/albums/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", middleware.GetBoundedVars(album.GetBoundedAlbumTracks, user.Log)).Methods("GET")
	r.Handle("/api/v1/playlists/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", middleware.GetBoundedVars(playlist.GetBoundedPlaylistTracks, user.Log)).Methods("GET")

	csrfMiddleware := m.CSRFCheckMiddleware(r)
	authHandler := m.CheckAuthMiddleware(csrfMiddleware)

	accessMiddleware := middleware.AccessLogMiddleware(authHandler, user.Log)
	panicMiddleware := middleware.PanicMiddleware(accessMiddleware, user.Log)
	return panicMiddleware
}

func StartNew() {

	db, err := gorm.Open("postgres", constants.DbConn)
	if err != nil {
		log.Fatal("Failed to start db: " + err.Error())
	}

	defer db.Close()

	db.DB().SetMaxOpenConns(constants.DbMaxConnN)

	if err := db.DB().Ping(); err != nil {
		log.Fatal("Failed to ping db: " + err.Error())
	}

	c := cors.New(constants.CorsOptions)

	var customLogger *logger.MainLogger
	f, err := os.OpenFile(constants.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logrus.Error("Failed to open logfile:", err)
		customLogger = logger.NewLogger(os.Stdout)
	} else {
		customLogger = logger.NewLogger(f)
	}
	defer f.Close()

	redisAddr := flag.String("addr", constants.RedisAddr, "redis addr")
	redisConn := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(*redisAddr)
			if err != nil {
				log.Fatal("fail init redis pool: ", err)
			}
			return conn, err
		},
	}
	defer redisConn.Close()

	csrfToken, err := csrf.NewAesCryptHashToken(constants.CsrfSecret)
	if err != nil {
		log.Fatal("fail init csrf token")
	}

	middlewares := InitRouter(customLogger, db, redisConn, csrfToken)
	//generateSSL()

	fmt.Println("Starts server at 8081")
	//err = http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", c.Handler(panicMiddleware))
	err = http.ListenAndServe(":8081", c.Handler(middlewares))
	if err != nil {
		log.Println(err)
		return
	}
}

func generateSSL() {
	// Проверяем, доступен ли cert файл.
	err := httpscerts.Check("cert.pem", "key.pem")
	// Если он недоступен, то генерируем новый.
	if err != nil {
		err = httpscerts.Generate("cert.pem", "key.pem", "https://127.0.0.1:8081")
		//err = httpscerts.Generate("cert.pem", "key.pem", "http://89.208.199.170:8001")
		if err != nil {
			logrus.Fatal("failed to generate https cert")
		}
	}
}
