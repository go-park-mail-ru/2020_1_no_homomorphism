package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/2020_1_no_homomorphism/no_homo_main/config"
	albumDelivery "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/album/delivery"
	albumRepo "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/album/repository"
	albumUC "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/album/usecase"
	artistDelivery "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist/delivery"
	artistRepo "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist/repository"
	artistUC "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist/usecase"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/csrf/repository"
	csrfLib "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/csrf/usecase"
	m "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/middleware"
	playlistDelivery "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist/delivery"
	playlistRepo "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist/repository"
	playlistUC "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist/usecase"
	searchDelivery "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/search/delivery"
	searchUC "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/search/usecase"
	trackDelivery "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track/delivery"
	trackRepo "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track/repository"
	trackUC "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track/usecase"
	userDelivery "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/user/delivery"
	userRepo "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/user/repository"
	userUC "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/user/usecase"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/2020_1_no_homomorphism/no_homo_main/proto/filetransfer"
	"github.com/2020_1_no_homomorphism/no_homo_main/proto/session"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func getInterceptor(mainLogger *logger.MainLogger) func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		mainLogger.Tracef("call=%v req=%#v reply=%#v time=%v err=%v",
			method, req, reply, time.Since(start), err)
		return err
	}
}

func InitHandler(mainLogger *logger.MainLogger, db *gorm.DB, csrfToken csrfLib.CryptToken, sessManager session.AuthCheckerClient, fileserver filetransfer.UploadServiceClient) (
	userDelivery.UserHandler,
	trackDelivery.TrackHandler,
	playlistDelivery.PlaylistHandler,
	albumDelivery.AlbumHandler,
	artistDelivery.ArtistHandler,
	searchDelivery.SearchHandler,
	m.AuthMidleware,
	m.CsrfMiddleware,
) {

	trackRep := trackRepo.NewDbTrackRepo(db)
	playlistRep := playlistRepo.NewDbPlaylistRepository(db)
	albumRep := albumRepo.NewDbAlbumRepository(db)
	artistRep := artistRepo.NewDbArtistRepository(db)
	dbRep := userRepo.NewDbUserRepository(db, viper.GetString(config.ConfigFields.AvatarDefault))

	ArtistUC := artistUC.ArtistUseCase{
		ArtistRepository: &artistRep,
	}

	AlbumUC := albumUC.AlbumUseCase{
		AlbumRepository: &albumRep,
	}

	PlaylistUC := playlistUC.PlaylistUseCase{
		PlRepository: &playlistRep,
		FileService:  fileserver,
		AvatarDir:    viper.GetString(config.ConfigFields.PlaylistAvatarDir),
	}

	UserUC := userUC.UserUseCase{
		Repository:  &dbRep,
		FileService: fileserver,
		AvatarDir:   viper.GetString(config.ConfigFields.AvatarDir),
	}
	TrackUC := trackUC.TrackUseCase{
		Repository: &trackRep,
	}

	playlistHandler := playlistDelivery.PlaylistHandler{
		PlaylistUC: &PlaylistUC,
		TrackUC:    &TrackUC,
		Log:        mainLogger,
		ImgTypes:   viper.GetStringMapString(config.ConfigFields.AvatarTypes),
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
		ImgTypes:        viper.GetStringMapString(config.ConfigFields.AvatarTypes),
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

	searchHandler := searchDelivery.SearchHandler{
		SearchUC: searchUC.SearchUseCase{
			ArtistRepo: &artistRep,
			AlbumRepo:  &albumRep,
			TrackRepo:  &trackRep,
		},
		Log: mainLogger,
	}

	auth := m.NewAuthMiddleware(sessManager, &UserUC, mainLogger)
	csrf := m.NewCsrfMiddleware(&csrfToken)

	return userHandler, trackHandler, playlistHandler, albumHandler, artistHandler, searchHandler, auth, csrf
}

func InitRouter(customLogger *logger.MainLogger, db *gorm.DB, csrfToken csrfLib.CryptToken, sessManager session.AuthCheckerClient, fileserver filetransfer.UploadServiceClient) http.Handler {
	user, track, playlist, album, artist, search, auth, csrf := InitHandler(customLogger, db, csrfToken, sessManager, fileserver)

	r := mux.NewRouter().PathPrefix(viper.GetString(config.ConfigFields.ApiPrefix)).Subrouter()

	r.Handle("/users/albums", auth.Auth(album.GetUserAlbums, false)).Methods("GET")
	r.Handle("/albums/{id:[0-9]+}", auth.Auth(album.GetFullAlbum, true)).Methods("GET")
	r.Handle("/albums/{id:[0-9]+}/rating", auth.Auth(album.RateAlbum, false)).Methods("POST")
	r.Handle("/artists/{id:[0-9]+}/albums/{start:[0-9]+}/{end:[0-9]+}", m.BoundedVars(album.GetBoundedAlbumsByArtistId, user.Log)).Methods("GET")

	r.Handle("/users/artists", auth.Auth(artist.SubscriptionList, false)).Methods("GET")
	r.Handle("/artists/{id:[0-9]+}", auth.Auth(artist.GetFullArtistInfo, true)).Methods("GET")
	r.HandleFunc("/artists/{id:[0-9]+}/stat", artist.GetArtistStat).Methods("GET")
	r.HandleFunc("/artists/{start:[0-9]+}/{end:[0-9]+}", artist.GetBoundedArtists).Methods("GET")
	r.Handle("/artists/{id:[0-9]+}/subscription", auth.Auth(artist.Subscribe, false)).Methods("POST") //todo csrf
	r.HandleFunc("/artists/top", artist.GetTopArtists).Methods("GET")

	r.Handle("/users/playlists", auth.Auth(playlist.GetUserPlaylists, false)).Methods("GET")
	r.Handle("/playlists/{id:[0-9]+}", auth.Auth(playlist.GetFullPlaylistById, true)).Methods("GET")
	r.Handle("/playlists/tracks/{id:[0-9]+}", auth.Auth(playlist.GetPlaylistsIDByTrack, false)).Methods("GET")
	r.Handle("/playlists/tracks", auth.Auth(csrf.CSRFCheck(playlist.AddTrackToPlaylist), false)).Methods("POST")
	r.Handle("/playlists/new/{name}", auth.Auth(csrf.CSRFCheck(playlist.CreatePlaylist), false)).Methods("POST")
	r.Handle("/playlists/{id:[0-9]+}", auth.Auth(csrf.CSRFCheck(playlist.DeletePlaylist), false)).Methods("DELETE")
	r.Handle("/playlists/{playlist:[0-9]+}/tracks/{track:[0-9]+}", auth.Auth(playlist.DeleteTrackFromPlaylist, false)).Methods("DELETE")
	r.Handle("/playlists/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", auth.Auth(m.BoundedVars(playlist.GetBoundedPlaylistTracks, user.Log), true)).Methods("GET")
	r.Handle("/playlists/{id:[0-9]+}/privacy", auth.Auth(playlist.ChangePrivacy, false)).Methods("POST")    //todo csrf
	r.Handle("/playlists/shared/{id:[0-9]+}", auth.Auth(playlist.AddSharedPlaylist, false)).Methods("POST") //todo csrf
	r.Handle("/playlists/{id:[0-9]+}/image", auth.Auth(csrf.CSRFCheck(playlist.UpdatePlaylistAvatar), false)).Methods("POST")
	r.Handle("/playlists/{id:[0-9]+}/update/{name}", auth.Auth(csrf.CSRFCheck(playlist.Update), false)).Methods("POST")

	r.Handle("/users/tracks", auth.Auth(track.GetUserTracks, false)).Methods("GET")
	r.HandleFunc("/tracks/{id:[0-9]+}", track.GetTrack).Methods("GET")
	r.Handle("/tracks/{id:[0-9]+}/rating", auth.Auth(track.RateTrack, false)).Methods("POST")
	r.Handle("/albums/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", auth.Auth(m.BoundedVars(track.GetBoundedAlbumTracks, user.Log), true)).Methods("GET")
	r.Handle("/artists/{id:[0-9]+}/tracks/{start:[0-9]+}/{end:[0-9]+}", auth.Auth(m.BoundedVars(track.GetBoundedArtistTracks, user.Log), true)).Methods("GET")
	r.Handle("/albums/newest", auth.Auth(album.GetNewestReleases, false)).Methods("GET")
	r.HandleFunc("/albums/worldnews", album.GetWorldNews).Methods("GET")
	r.Handle("/users", auth.Auth(user.CheckAuth, false))
	r.HandleFunc("/users/{id:[0-9]+}/stat", user.GetUserStat).Methods("GET")
	r.Handle("/users/login", auth.Auth(user.Login, true)).Methods("POST")
	r.Handle("/users/signup", auth.Auth(user.Create, true)).Methods("POST")
	r.Handle("/users/token", auth.Auth(user.GetCSRF, false)).Methods("GET")
	r.Handle("/users/me", auth.Auth(user.SelfProfile, false)).Methods("GET")
	r.Handle("/users/logout", auth.Auth(user.Logout, false)).Methods("DELETE") //todo убрать глаголы
	r.Handle("/users/profiles/{profile}", auth.Auth(user.Profile, false)).Methods("GET")
	r.Handle("/users/settings", auth.Auth(csrf.CSRFCheck(user.Update), false)).Methods("PUT")
	r.Handle("/users/images", auth.Auth(csrf.CSRFCheck(user.UpdateAvatar), false)).Methods("POST")


	r.HandleFunc("/media/{text}/{count:[0-9]+}", search.Search).Methods("GET")

	r.Handle("/metrics", promhttp.Handler())

	accessMiddleware := m.AccessLogMiddleware(r, user.Log)
	panicMiddleware := m.PanicMiddleware(accessMiddleware, user.Log)

	return panicMiddleware
}

func StartNew() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to export env vars: %v", err)
	}

	if err := config.ExportConfig(); err != nil {
		log.Fatalf("Failed to export config: %v", err)
	}

	db, err := gorm.Open("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		log.Fatalf("Failed to start db: %v", err)
	}

	defer db.Close()

	db.DB().SetMaxOpenConns(viper.GetInt(config.ConfigFields.DBMaxConnNum))

	if err := db.DB().Ping(); err != nil {
		log.Fatalf("Failed to ping db: %v", err)
	}

	c := cors.New(config.CorsInit())

	var customLogger *logger.MainLogger
	f, err := os.OpenFile(viper.GetString(config.ConfigFields.LogFile), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logrus.Error("Failed to open logfile:", err)
		customLogger = logger.NewLogger(os.Stdout)
	} else {
		customLogger = logger.NewLogger(f)
	}
	customLogger.SetLevel(logrus.TraceLevel)
	defer f.Close()

	redisAddr := flag.String("addr", viper.GetString(config.ConfigFields.RedisAddr), "redis addr")
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

	csrfToken, err := csrfLib.NewAesCryptHashToken(os.Getenv("CSRF_SECRET"), viper.GetInt64(config.ConfigFields.CsrfDuration), &csrfRepo)
	if err != nil {
		log.Fatalf("failed to init csrf token: %v", err)
	}

	grpcSessionsConn, err := grpc.Dial(
		viper.GetString(config.ConfigFields.GRPCsessions),
		grpc.WithUnaryInterceptor(getInterceptor(customLogger)),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grpcSessionsConn.Close()

	sessManager := session.NewAuthCheckerClient(grpcSessionsConn)

	grpcFileserverConn, err := grpc.Dial(
		viper.GetString(config.ConfigFields.GRPCfs),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grpcFileserverConn.Close()

	fileserver := filetransfer.NewUploadServiceClient(grpcFileserverConn)

	routes := InitRouter(customLogger, db, csrfToken, sessManager, fileserver)

	fmt.Println("Starts server at ", viper.GetString(config.ConfigFields.MainAddr))
	//err = http.ListenAndServeTLS(viper.GetString(config.ConfigFields.MainAddr), viper.GetString(config.ConfigFields.SSLfullchain), viper.GetString(config.ConfigFields.SSLkey), c.Handler(m.HeadersHandler(routes)))
	err = http.ListenAndServe(viper.GetString(config.ConfigFields.MainAddr), c.Handler(m.HeadersHandler(routes)))
	if err != nil {
		log.Println(err)
		return
	}
}
