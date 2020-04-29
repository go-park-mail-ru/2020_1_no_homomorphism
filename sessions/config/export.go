package config

import (
	"os"

	"github.com/spf13/viper"
)

var ConfigFields = struct {
	RedisAddr string
	TcpPort string
	ExpireTime string
} {
	RedisAddr : "redis.addr",
	TcpPort: "tcp.port",
	ExpireTime: "cookie.expired",
}

func ExportConfig() error {
	viper.AddConfigPath(os.Getenv("SESSIONS_CONFIG_PATH"))
	viper.SetConfigName(os.Getenv("SESSIONS_CONFIG_NAME"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

