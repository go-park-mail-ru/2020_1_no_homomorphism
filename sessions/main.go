package main

import (
	"flag"
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_sessions/config"
	session "github.com/2020_1_no_homomorphism/no_homo_sessions/internal"
	"github.com/2020_1_no_homomorphism/no_homo_sessions/internal/delivery"
	"github.com/2020_1_no_homomorphism/no_homo_sessions/internal/repository"
	"github.com/2020_1_no_homomorphism/no_homo_sessions/internal/usecase"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to export env vars: %v", err)
	}

	if err := config.ExportConfig(); err != nil {
		log.Fatalf("can't export config %v", err)
	}

	tcpPort := viper.GetString(config.ConfigFields.TcpPort)
	lis, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatalln("cant listen port", err)
	}
	server := grpc.NewServer()

	redisAddr := flag.String("addr", viper.GetString(config.ConfigFields.RedisAddr), "redis addr")

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

	session.RegisterAuthCheckerServer(server, delivery.NewSessionDelivery(&sessionUseCase, viper.GetUint64(config.ConfigFields.ExpireTime)))

	fmt.Printf("starting server at %s\n", tcpPort)
	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
