package api

import (
	"auth-service/api/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "auth-service/api/handler/docs"
)

// @title Auth Service API
// @version 1.0
// @description This is a sample server for Auth Service.
// @host localhost:8081
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /api/v1
// @schemes http
func NewRouter(handle *handler.Handler) *gin.Engine {
	router := gin.Default()

	// Swagger endpointini sozlash
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", handle.RegisterHandler)
		auth.POST("/login", handle.LoginHandler)
		auth.POST("/reset-password", handle.ResetPasswordHandler)
		auth.POST("/reset-password/new-password", handle.UpdatePasswordHandler)
		auth.POST("/refresh", handle.RefreshToken)
		auth.POST("/logout", handle.LogoutUserHandler)
	}

	return router
}
