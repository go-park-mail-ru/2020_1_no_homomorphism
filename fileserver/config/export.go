package config

import (
	"os"

	"github.com/spf13/viper"
)
var ConfigFields = struct{
	Port string
	PortTLS string
	Dir string
} {
	Port: "port",
	PortTLS: "port_tls",
	Dir: "dir",
}

func ExportConfig() error {
	viper.AddConfigPath(os.Getenv("FILESERVER_CONFIG_PATH"))
	viper.SetConfigName(os.Getenv("FILESERVER_CONFIG_NAME"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}