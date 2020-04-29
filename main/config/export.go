package config

import (
	"os"

	"github.com/rs/cors"
	"github.com/spf13/viper"
)

var ConfigFields = struct {
	// db
	DBMaxConnNum string
	// logger
	LogFile string
	// redis
	RedisAddr string
	// csrf
	CsrfDuration string
	// fileserver
	FSRoot        string
	FSAddr        string
	AvatarDir     string
	AvatarDefault string
	AvatarTypes   string
	// api
	ApiPrefix string
	// cors
	CorsAllowedOrigins string
	CorsAllowedCreds   string
	CorsAllowedMethods string
	CorsAllowedHeaders string
	CorsDebug          string
	// cookie
	CookieExpireTime string
}{
	DBMaxConnNum:       "db.max_conn_num",
	LogFile:            "logger.file",
	RedisAddr:          "redis.addr",
	CsrfDuration:       "csrf.duration",
	FSRoot:             "fileserver.root",
	FSAddr:             "fileserver.addr",
	AvatarDefault:      "fileserver.avatar.default",
	AvatarDir:          "fileserver.avatar.dir",
	AvatarTypes:        "fileserver.avatar.types",
	ApiPrefix:          "api.prefix",
	CorsAllowedOrigins: "cors.allowed_origins",
	CorsAllowedCreds:   "cors.allowed_cred",
	CorsAllowedHeaders: "cors.allowed_headers",
	CorsAllowedMethods: "cors.allowed_methods",
	CorsDebug:          "cors.debug",
	CookieExpireTime:   "cookie.expire",
}

func CorsInit() cors.Options {
	return cors.Options{
		AllowedOrigins:   viper.GetStringSlice(ConfigFields.CorsAllowedOrigins),
		AllowCredentials: viper.GetBool(ConfigFields.CorsAllowedCreds),
		AllowedMethods:   viper.GetStringSlice(ConfigFields.CorsAllowedMethods),
		AllowedHeaders:   viper.GetStringSlice(ConfigFields.CorsAllowedHeaders),
		Debug:            viper.GetBool(ConfigFields.CorsDebug),
	}
}

func ExportConfig() error {
	viper.AddConfigPath(os.Getenv("MAIN_CONFIG_PATH"))
	viper.SetConfigName(os.Getenv("MAIN_CONFIG_NAME"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
