package main

import (
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
	"log"
	"net"
	"no_homomorphism/configs/proto/session"
	"no_homomorphism/internal/pkg/constants"
	"no_homomorphism/sessions/internal/delivery"
	"no_homomorphism/sessions/internal/repository"
	"no_homomorphism/sessions/internal/usecase"
)

func main() {
	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	redisAddr := flag.String("addr", constants.RedisAddr, "redis addr")
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

	session.RegisterAuthCheckerServer(server, delivery.NewSessionDelivery(&sessionUseCase, constants.CookieExpireTime))

	fmt.Println("starting server at :8083")
	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
