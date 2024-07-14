package handler

import (
	"auth-service/api/token"
	"auth-service/generated/user"
	"auth-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param Register body user.RegisterRequest true "User Registration"
// @Success 201 {object} user.RegisterResponse
// @Failure 400 {object} models.Errors
// @Router /auth/register [post]
func (h *Handler) RegisterHandler(ctx *gin.Context) {
	var signUp user.RegisterRequest

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

	resp, err := h.UserRepo.CreateUser(&signUp)
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
// @Param Login body user.LoginRequest true "User Login"
// @Success 200 {object} models.Token
// @Failure 400 {object} models.Errors
// @Failure 404 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Router /auth/login [post]
func (h *Handler) LoginHandler(ctx *gin.Context) {
	var signIn user.LoginRequest

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

	h.Logger.Info("user login successfully")
	newToken := models.Token{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}
	ctx.JSON(http.StatusOK, newToken)
}
