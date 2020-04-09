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
	csrfLib "no_homomorphism/internal/pkg/csrf"
	m "no_homomorphism/internal/pkg/middleware"
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

func InitHandler(mainLogger *logger.MainLogger, db *gorm.DB, redis *redis.Pool, csrfToken csrfLib.CryptToken) (
	userDelivery.UserHandler,
	trackDelivery.TrackHandler,
	playlistDelivery.PlaylistHandler,
	albumDelivery.AlbumHandler,
	artistDelivery.ArtistHandler,
	m.AuthMidleware,
	m.CsrfMiddleware,
) {

	sesRep := sessionRepo.NewRedisSessionManager(redis)
	trackRep := trackRepo.NewDbTrackRepo(db)
	playlistRep := playlistRepo.NewDbPlaylistRepository(db)
	albumRep := albumRepo.NewDbAlbumRepository(db)
	artistRep := artistRepo.NewDbArtistRepository(db)
	dbRep := userRepo.NewDbUserRepository(db, constants.AvatarDefault, constants.AvatarDir)

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

	auth := m.NewAuthMiddleware(&SessionDelivery, &UserUC)
	csrf := m.NewCsrfMiddleware(csrfToken)

	return userHandler, trackHandler, playlistHandler, albumHandler, artistHandler, auth, csrf
}

func InitRouter(customLogger *logger.MainLogger, db *gorm.DB, redisConn *redis.Pool, csrfToken csrfLib.CryptToken) http.Handler {
	user, track, playlist, album, artist, auth, csrf := InitHandler(customLogger, db, redisConn, csrfToken)

	r := mux.NewRouter().PathPrefix(constants.ApiPrefix).Subrouter()

	r.Handle("/users/albums", auth.AuthMiddleware(album.GetUserAlbums)).Methods("GET")
	r.HandleFunc("/albums/{id:[0-9]+}", album.GetFullAlbum).Methods("GET")
	r.Handle("/artists/{id:[0-9]+}/albums/{start:[0-9]+}/{end:[0-9]+}", m.GetBoundedVars(album.GetBoundedAlbumsByArtistId, user.Log)).Methods("GET")

	r.HandleFunc("/artists/{id:[0-9]+}", artist.GetFullArtistInfo).Methods("GET")
	r.HandleFunc("/artists/{id:[0-9]+}/stat", artist.GetArtistStat).Methods("GET")
	r.HandleFunc("/artists/{start:[0-9]+}/{end:[0-9]+}", artist.GetBoundedArtists).Methods("GET")

	r.Handle("/users/playlists", auth.AuthMiddleware(playlist.GetUserPlaylists)).Methods("GET")
	r.Handle("/playlists/{id:[0-9]+}", auth.AuthMiddleware(playlist.GetFullPlaylistById)).Methods("GET")
	r.Handle("/playlists/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", auth.AuthMiddleware(m.GetBoundedVars(playlist.GetBoundedPlaylistTracks, user.Log))).Methods("GET")

	r.HandleFunc("/tracks/{id:[0-9]+}", track.GetTrack).Methods("GET")
	r.Handle("/albums/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", m.GetBoundedVars(track.GetBoundedAlbumTracks, user.Log)).Methods("GET")
	r.Handle("/artists/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", m.GetBoundedVars(track.GetBoundedArtistTracks, user.Log)).Methods("GET")

	r.Handle("/users", auth.AuthMiddleware(user.CheckAuth))
	r.Handle("/users/me", auth.AuthMiddleware(user.SelfProfile)).Methods("GET")
	r.Handle("/users/login", auth.AuthMiddleware(user.Login)).Methods("POST")
	r.Handle("/users/logout", auth.AuthMiddleware(user.Logout)).Methods("DELETE")
	r.Handle("/users/images", auth.AuthMiddleware(csrf.CSRFCheckMiddleware(user.UpdateAvatar))).Methods("POST")
	r.Handle("/users/signup", auth.AuthMiddleware(user.Create)).Methods("POST")
	r.Handle("/users/settings", auth.AuthMiddleware(csrf.CSRFCheckMiddleware(user.Update))).Methods("PUT")
	r.Handle("/users/profiles/{profile}", auth.AuthMiddleware(user.Profile)).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}/stat", user.GetUserStat).Methods("GET")

	accessMiddleware := m.AccessLogMiddleware(r, user.Log)
	panicMiddleware := m.PanicMiddleware(accessMiddleware, user.Log)

	return panicMiddleware
}

func StartNew() {
	db, err := gorm.Open("postgres", constants.DbConn)
	if err != nil {
		log.Fatalf("Failed to start db: %v", err)
	}

	defer db.Close()

	db.DB().SetMaxOpenConns(constants.DbMaxConnN)

	if err := db.DB().Ping(); err != nil {
		log.Fatalf("Failed to ping db: %v", err)
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
				log.Fatalf("failed to init redis pool: %v", err)
			}
			return conn, err
		},
	}
	defer redisConn.Close()

	csrfToken, err := csrfLib.NewAesCryptHashToken(constants.CsrfSecret)
	if err != nil {
		log.Fatalf("failed to init csrf token: %v", err)
	}

	routes := InitRouter(customLogger, db, redisConn, csrfToken)

	fmt.Println("Starts server at 8081")
	err = http.ListenAndServe(":8081", c.Handler(routes))
	if err != nil {
		log.Println(err)
		return
	}
}
