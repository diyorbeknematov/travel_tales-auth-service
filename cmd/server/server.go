package server

import (
	"auth-service/config"
	"auth-service/generated/user"
	"auth-service/logs"
	"auth-service/service"
	"auth-service/storage/postgres"
	"auth-service/storage/redis"
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunServer(userRepo *postgres.UserRepo, redisClient *redis.RedisClient) {
	logs.InitLogger()
	cfg := config.Load()
	listener, err := net.Listen("tcp", cfg.GRPC_PORT)
	if err != nil {
		logs.Logger.Error("Error create to new listener", "error", err.Error())
		log.Fatal(err)
	}

	s := grpc.NewServer()
	srv := service.UserService{
		UserRepo: userRepo,
		RedisClient: redisClient,
		Logger: logs.Logger,
	}

	user.RegisterAuthServiceServer(s, &srv)

	logs.Logger.Info("server is running ", "PORT", cfg.GRPC_PORT)

	log.Printf("server is running on %v...", listener.Addr())
	if err := s.Serve(listener); err != nil {
		logs.Logger.Error("Faild server is running", "error", err.Error())
		log.Fatal(err)
	}
}
