package server

import (
	"flag"
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/constants"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/csrf/repository"
	"github.com/2020_1_no_homomorphism/no_homo_main/pkg/logger"
	"github.com/2020_1_no_homomorphism/no_homo_main/proto/session"
	"github.com/kabukky/httpscerts"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"

	albumDelivery "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/album/delivery"
	albumRepo "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/album/repository"
	albumUC "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/album/usecase"
	artistDelivery "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist/delivery"
	artistRepo "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist/repository"
	artistUC "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist/usecase"
	csrfLib "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/csrf/usecase"
	m "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/middleware"
	playlistDelivery "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist/delivery"
	playlistRepo "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist/repository"
	playlistUC "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist/usecase"
	trackDelivery "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track/delivery"
	trackRepo "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track/repository"
	trackUC "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track/usecase"
	userDelivery "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/user/delivery"
	userRepo "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/user/repository"
	userUC "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/user/usecase"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

func InitHandler(mainLogger *logger.MainLogger, db *gorm.DB, csrfToken csrfLib.CryptToken, sessManager session.AuthCheckerClient) (
	userDelivery.UserHandler,
	trackDelivery.TrackHandler,
	playlistDelivery.PlaylistHandler,
	albumDelivery.AlbumHandler,
	artistDelivery.ArtistHandler,
	m.AuthMidleware,
	m.CsrfMiddleware,
) {

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
		SessionDelivery: sessManager,
		UserUC:          &UserUC,
		Log:             mainLogger,
		ImgTypes:        constants.AvatarTypes,
		CSRF:            &csrfToken,
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

	auth := m.NewAuthMiddleware(sessManager, &UserUC, mainLogger)
	csrf := m.NewCsrfMiddleware(&csrfToken)

	return userHandler, trackHandler, playlistHandler, albumHandler, artistHandler, auth, csrf
}

func InitRouter(customLogger *logger.MainLogger, db *gorm.DB, csrfToken csrfLib.CryptToken, sessManager session.AuthCheckerClient) http.Handler {
	user, track, playlist, album, artist, auth, csrf := InitHandler(customLogger, db, csrfToken, sessManager)

	r := mux.NewRouter().PathPrefix(constants.ApiPrefix).Subrouter()

	r.Handle("/users/albums", auth.AuthMiddleware(album.GetUserAlbums, false)).Methods("GET")
	r.HandleFunc("/albums/{id:[0-9]+}", album.GetFullAlbum).Methods("GET")
	r.Handle("/artists/{id:[0-9]+}/albums/{start:[0-9]+}/{end:[0-9]+}", m.GetBoundedVars(album.GetBoundedAlbumsByArtistId, user.Log)).Methods("GET")

	r.HandleFunc("/artists/{id:[0-9]+}", artist.GetFullArtistInfo).Methods("GET")
	r.HandleFunc("/artists/{id:[0-9]+}/stat", artist.GetArtistStat).Methods("GET")
	r.HandleFunc("/artists/{start:[0-9]+}/{end:[0-9]+}", artist.GetBoundedArtists).Methods("GET")

	r.Handle("/users/playlists", auth.AuthMiddleware(playlist.GetUserPlaylists, false)).Methods("GET")
	r.Handle("/playlists/{id:[0-9]+}", auth.AuthMiddleware(playlist.GetFullPlaylistById, false)).Methods("GET")
	r.Handle("/playlists/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", auth.AuthMiddleware(m.GetBoundedVars(playlist.GetBoundedPlaylistTracks, user.Log), false)).Methods("GET")

	r.HandleFunc("/tracks/{id:[0-9]+}", track.GetTrack).Methods("GET")
	r.Handle("/albums/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", m.GetBoundedVars(track.GetBoundedAlbumTracks, user.Log)).Methods("GET")
	r.Handle("/artists/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", m.GetBoundedVars(track.GetBoundedArtistTracks, user.Log)).Methods("GET")

	r.Handle("/users", auth.AuthMiddleware(user.CheckAuth, false))
	r.Handle("/users/token", auth.AuthMiddleware(user.GetCSRF, false)).Methods("GET")
	r.Handle("/users/me", auth.AuthMiddleware(user.SelfProfile, false)).Methods("GET")
	r.Handle("/users/login", auth.AuthMiddleware(user.Login, true)).Methods("POST")
	r.Handle("/users/logout", auth.AuthMiddleware(user.Logout, false)).Methods("DELETE")
	r.Handle("/users/images", auth.AuthMiddleware(csrf.CSRFCheckMiddleware(user.UpdateAvatar), false)).Methods("POST")
	r.Handle("/users/signup", auth.AuthMiddleware(user.Create, true)).Methods("POST")
	r.Handle("/users/settings", auth.AuthMiddleware(csrf.CSRFCheckMiddleware(user.Update), false)).Methods("PUT")
	r.Handle("/users/profiles/{profile}", auth.AuthMiddleware(user.Profile, false)).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}/stat", user.GetUserStat).Methods("GET")

	accessMiddleware := m.AccessLogMiddleware(r, user.Log)
	panicMiddleware := m.PanicMiddleware(accessMiddleware, user.Log)

	return panicMiddleware
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

	csrfRepo := repository.NewRedisTokenManager(redisConn)

	csrfToken, err := csrfLib.NewAesCryptHashToken(constants.CsrfSecret, constants.CsrfDuration, &csrfRepo)
	if err != nil {
		log.Fatalf("failed to init csrf token: %v", err)
	}

	grcpConn, err := grpc.Dial(
		"127.0.0.1:8083",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	sessManager := session.NewAuthCheckerClient(grcpConn)

	routes := InitRouter(customLogger, db, csrfToken, sessManager)

	fmt.Println("Starts server at 8081")
	//err = http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", c.Handler(m.HeadersHandler(routes)))
	err = http.ListenAndServe(":8081", c.Handler(m.HeadersHandler(routes)))
	if err != nil {
		log.Println(err)
		return
	}
}
