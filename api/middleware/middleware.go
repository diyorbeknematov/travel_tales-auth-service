package middleware

import (
	"auth-service/api/handler/token"
	"auth-service/logs"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		valid, err := token.ValidateToken(authHeader)
		if err != nil || !valid {
			c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Errorf("invalid token: %s", err))
			return
		}

		claims, err := token.ExtractClaimsAccess(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims.ExpiresAt < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		// Claimsdan ma'lumotlarni kontekstga qo'shish
		c.Set("user_id", claims.UserId)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logs.Logger.Info("Request received",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)

		c.Next()

		logs.Logger.Info("Response sent",
			slog.Int("status", c.Writer.Status()),
		)
	}
}
