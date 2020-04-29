package main

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/fileserver/config"
	"github.com/2020_1_no_homomorphism/fileserver/proto/delivery"
	"github.com/2020_1_no_homomorphism/fileserver/proto/filetransfer"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to export env vars: %v", err)
	}

	if err := config.ExportConfig(); err != nil {
		log.Fatalf("can't export config %v", err)
	}


	lis, err := net.Listen("tcp", viper.GetString(config.ConfigFields.Port))
	if err != nil {
		log.Fatalln("cant listen port", err)
	}
	server := grpc.NewServer()

	filetransfer.RegisterUploadServiceServer(server, delivery.NewFileTransferDelivery())

	fmt.Println("starting grpc server at ", viper.GetString(config.ConfigFields.Port))
	err = server.Serve(lis)
	if err != nil {
		fmt.Println("failed to serve grpc")
	}

	fmt.Println("Starts server at ", viper.GetString(config.ConfigFields.PortTLS))
	//err := http.ListenAndServeTLS(viper.GetString(config.ConfigFields.PortTLS),"fullchain.pem","privkey.pem", http.FileServer(http.Dir(viper.GetString(config.ConfigFields.Dir)))))
	err = http.ListenAndServe(viper.GetString(config.ConfigFields.Port), http.FileServer(http.Dir(viper.GetString(config.ConfigFields.Dir))))
	if err != nil {
		log.Println(err)
		return
	}
}
