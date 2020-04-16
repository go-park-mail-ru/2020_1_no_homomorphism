package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	session "github.com/2020_1_no_homomorphism/no_homo_sessions/internal"
	"github.com/2020_1_no_homomorphism/no_homo_sessions/internal/delivery"
	"github.com/2020_1_no_homomorphism/no_homo_sessions/internal/repository"
	"github.com/2020_1_no_homomorphism/no_homo_sessions/internal/usecase"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	redisAddr := flag.String("addr", redisAddr, "redis addr")

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

	sessionUseCase := usecase.SessionUseCase{
		Repository: repository.NewRedisSessionManager(redisConn),
	}

	session.RegisterAuthCheckerServer(server, delivery.NewSessionDelivery(&sessionUseCase, expireTime))

	fmt.Printf("starting server at %s\n", tcpPort)
	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
