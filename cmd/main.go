package main

import (
	"auth-service/api"
	"auth-service/api/handler"
	"auth-service/cmd/server"
	"auth-service/config"
	"auth-service/logs"
	"auth-service/storage/postgres"
	"auth-service/storage/redis"
	"log"
	"log/slog"
	"sync"
)

func main() {
	logs.InitLogger()

	logs.Logger.Info("Starting the server ...")
	db, err := postgres.ConnectDB()
	if err != nil {
		logs.Logger.Error("Error connection th postgres", slog.String("error", err.Error()))
		log.Fatal(err)
	}
	defer db.Close()

	redisClient := redis.NewRedisClient()

	cfg := config.Load()
	handle := handler.NewHandler(postgres.NewUserRepo(db), logs.Logger, redisClient)
	router := api.NewRouter(handle)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		logs.Logger.Info("server is running", "PORT", cfg.HTTP_PORT)
		err := router.Run(cfg.HTTP_PORT)
		if err != nil {
			logs.Logger.Error("Faild server is running", "error", err.Error())
			log.Fatal(err)
		}
	}()

	server.RunServer(postgres.NewUserRepo(db), redis.NewRedisClient())

	wg.Wait()
}
