package main

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/fileserver/proto/delivery"
	"github.com/2020_1_no_homomorphism/fileserver/proto/filetransfer"
	"google.golang.org/grpc"
	"net"

	"log"
	"net/http"
)

func main() {
	lis, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}
	server := grpc.NewServer()

	filetransfer.RegisterUploadServiceServer(server, delivery.NewFileTransferDelivery())

	fmt.Println("starting grpc server at :8084")
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()

	fmt.Println("Starts server at 8082")
	//err := http.ListenAndServeTLS(":8082","fullchain.pem","privkey.pem", http.FileServer(http.Dir("./resources")))
	err = http.ListenAndServe(":8082", http.FileServer(http.Dir("./resources")))
	if err != nil {
		log.Println(err)
		return
	}
}
