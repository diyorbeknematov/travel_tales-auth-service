package handler

import (
	"auth-service/storage/postgres"
	"log/slog"
)

type Handler struct {
	UserRepo *postgres.UserRepo
	Logger *slog.Logger
}

func NewHandler(user *postgres.UserRepo, logger *slog.Logger) *Handler {
	return &Handler{
		UserRepo: user,
		Logger: logger,
	}
}