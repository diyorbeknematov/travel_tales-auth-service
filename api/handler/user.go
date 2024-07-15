package handler

import (
	"auth-service/api/handler/token"
	"auth-service/models"
	"auth-service/pkg"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param Register body models.Register true "User Registration"
// @Success 201 {object} models.Register
// @Failure 400 {object} models.Errors
// @Router /auth/register [post]
func (h *Handler) RegisterHandler(ctx *gin.Context) {
	var signUp models.Register

	if err := ctx.ShouldBindJSON(&signUp); err != nil {
		h.Logger.Error("Error bind json")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(signUp.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Error("Error generating hashed password", "error", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	signUp.Password = string(hashedPass)

	resp, err := h.UserRepo.CreateUser(signUp)
	if err != nil {
		h.Logger.Error("Error register user", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// @Summary Login a user
// @Description Login a user with email and password
// @Tags Auth
// @Accept json
// Produce json
// @Param Login body models.UserLogin true "User Login"
// @Success 200 {object} models.Token
// @Failure 400 {object} models.Errors
// @Failure 404 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Router /auth/login [post]
func (h *Handler) LoginHandler(ctx *gin.Context) {
	var signIn models.UserLogin

	if err := ctx.ShouldBindJSON(&signIn); err != nil {
		h.Logger.Error("Error bind json")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.UserRepo.GetUserByEmail(signIn.Email)
	if err != nil {
		h.Logger.Error("Error getting user by email", "error", err.Error())
		ctx.JSON(http.StatusNotFound, gin.H{
			"Error": err.Error(),
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(signIn.Password), []byte(user.Password)); err != nil {
		h.Logger.Error("Invalid password", "error", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	accessToken, err := token.GenerateAccessJWT(user)
	if err != nil {
		h.Logger.Error("Error generating access token:", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := token.GenerateRefreshJWT(user)
	if err != nil {
		h.Logger.Error("Error generating refresh token:", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	err = h.UserRepo.SaveRefreshToken(user.Username, refreshToken, time.Now().Add(7*24*time.Hour))
	if err != nil {
		h.Logger.Error("error in save refresh token", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInsufficientStorage, gin.H{"error": "Faild in save refresh token"})
		return
	}
	h.Logger.Info("user login successfully")
	newToken := models.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	ctx.JSON(http.StatusOK, newToken)
}

// @Summary Logout user
// @Description Logout the authenticated user
// @Tags Auth
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Param Authorization header string true "Logout User"
// @Success 200 {object} models.Success
// @Failure 404 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Router /auth/logout [post]
func (h *Handler) LogoutUserHandler(ctx *gin.Context) {
	accessToken := ctx.GetHeader("Authorization")

	accessClaims, err := token.ExtractClaims(accessToken)
	if err != nil {
		h.Logger.Error("eror tokenni extract qilishda", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Refresh tokenni bekor qilish
	err = h.UserRepo.InvalidateRefreshToken(accessClaims.Username)
	if err != nil {
		h.Logger.Error("Error invalidate token", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.RedisClient.BlacklistToken(accessToken, time.Duration(accessClaims.ExpiresAt))
	ctx.JSON(http.StatusOK, models.Success{
		Message: "user logouted successfully",
	})
}

// @Summary Reset Password
// @Description Userni parolini qayta tiklash
// Tags Auth
// Accept json
// @Security ApiKeyAuth
// Produce json
// Param ResetPassword body models.ResetPassword true "Reset password"
// Success 200 {object} models.Success
// Failure 404 {object} models.Errors
// Failure 400 {object} models.Errors
// Failure 500 {object} models.Errors
// Router /auth/reset-password [post]
func (h *Handler) ResetPasswordHandler(ctx *gin.Context) {
	var email models.ResetPassword

	if err := ctx.ShouldBindJSON(&email); err != nil {
		h.Logger.Error("Error bind json")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newtoken, err := token.GenerateAccessJWT(&models.UserLogin{
		ID:       ctx.GetString("user-id"),
		Username: ctx.GetString("username"),
		Email:    ctx.GetString("email"),
	})
	if err != nil {
		h.Logger.Error("Error generated token", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	resetLink := pkg.CreateResetLink(ctx.Request.URL.Path, newtoken)

	// Email content
	subject := "Password Reset"
	body := "Click on the link to reset your password: " + resetLink

	// Emailni yuborish
	err = pkg.SendEmail(email.Email, subject, body)
	if err != nil {
		h.Logger.Error("Error in send email reset password link", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInsufficientStorage, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Success{
		Message: "Send email succussfully",
	})
}

// @Summary Update Parol
// @Description Parolni emailga yuborilgan linkda yangilash
// @Tags Auth
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Param UpdatePassword body models.UpdatePassword true "Reset Password"
// @Success 200 {object} models.Success
// @Failure 400 {object} models.Errors
// @Failure 500 {object} models.Errors
// Router /auth/reset-password/new-password [post]
func (h *Handler) UpdatePassword(ctx *gin.Context) {
	var pass models.UpdatePassword
	accesstoken := ctx.Query("token")
	claims, err := token.ExtractClaims(accesstoken)
	
	if err != nil {
		h.Logger.Error("Error tokenni tekshirishda", slog.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, models.Errors{
			Message: "tokenni tekshirishda xatolik",
		})
		return
	}
	pass.ID = claims.Id

	if err := ctx.ShouldBindJSON(&pass); err != nil {
		h.Logger.Error("Error bind json")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	resp, err := h.UserRepo.UpdatePassword(pass)
	if err != nil {
		h.Logger.Error("Error in reset password", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, models.Errors{
			Message: "Error in reset password",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Refresh access token
// @Description Refresh the access token using the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Refresh token"
// @Success 200 {object} models.Token
// @Failure 400 {object} models.Errors
// @Failure 401 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Security ApiKeyAuth
// @Router /auth/refresh_token [get]
func (h *Handler) RefreshToken(c *gin.Context) {
	h.Logger.Info("Handling RefreshToken request")

	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	claims, err := token.ExtractClaims(refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	if claims.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		return
	}

	newAccessToken, err := token.GenerateAccessJWT(&models.UserLogin{
		ID:       claims.UserId,
		Username: claims.Username,
		Email:    claims.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	c.JSON(http.StatusOK, models.Token{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
	})
}
