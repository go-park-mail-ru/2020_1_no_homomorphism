package config

import (
	"os"

	"github.com/spf13/viper"
)

var ConfigFields = struct {
	GRPC    string
	PortTLS string
	Dir     string
}{
	GRPC:    "grpc",
	PortTLS: "port_tls",
	Dir:     "dir",
}

func ExportConfig() error {
	viper.AddConfigPath(os.Getenv("FILESERVER_CONFIG_PATH"))
	viper.SetConfigName(os.Getenv("FILESERVER_CONFIG_NAME"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
