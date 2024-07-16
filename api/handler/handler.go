package handler

import (
	"auth-service/storage/postgres"
	"auth-service/storage/redis"
	"log/slog"
)

type Handler struct {
	UserRepo    *postgres.UserRepo
	RedisClient *redis.RedisClient
	Logger      *slog.Logger
}

func NewHandler(user *postgres.UserRepo, logger *slog.Logger, client *redis.RedisClient) *Handler {
	return &Handler{
		UserRepo: user,
		Logger:   logger,
		RedisClient: client,
	}
}
