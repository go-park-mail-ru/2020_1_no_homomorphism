package main

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/fileserver/proto/delivery"
	"github.com/2020_1_no_homomorphism/fileserver/proto/filetransfer"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	lis, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatalln("cant listet port", err)
	}
	server := grpc.NewServer()

	filetransfer.RegisterUploadServiceServer(server, delivery.NewFileTransferDelivery())

	fmt.Println("starting grpc server at :8084")
	err = server.Serve(lis)
	if err != nil {
		fmt.Println("failed to serve grpc")
	}

	fmt.Println("Starts server at 8082")
	//err := http.ListenAndServeTLS(":8082","fullchain.pem","privkey.pem", http.FileServer(http.Dir("./resources")))
	err = http.ListenAndServe(":8082", http.FileServer(http.Dir("./resources")))
	if err != nil {
		log.Println(err)
		return
	}
}
