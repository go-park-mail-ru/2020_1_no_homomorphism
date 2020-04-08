package constants

import (
	"time"

	"github.com/rs/cors"
)

var CorsOptions = cors.Options{
	AllowedOrigins:   []string{"http://89.208.199.170:3000", "http://195.19.37.246:10982", "http://89.208.199.170:3001", "http://localhost:3000"},
	AllowCredentials: true,
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	AllowedHeaders:   []string{"Content-Type", "X-Content-Type-Options", "Csrf-Token"},
	Debug:            false,
}

var (
	DbConn = "host=localhost port=5432 user=postgres password=postgres dbname=music_app"
	DbMaxConnN = 10
)

var LogFile = "logfile.log"

var RedisAddr = "redis://user:@localhost:6379/0"

var CsrfSecret = "qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB"

var (
	FileServerRoot = "resources"
	FileServerAddr = "http://localhost:8082/"
	AvatarDefault  = "http://localhost:8082/avatar/default.jpg"
	AvatarDir      = "/avatar"
	AvatarTypes    = map[string]string{"image/jpeg": "jpg", "image/png": "png", "image/gif": "gif"}
)

var CookieExpireTime = 24 * 31 * time.Hour
