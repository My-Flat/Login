package controller

import (
	"net/http"

	"my-flat-login/internal/service" // Replace with your actual module path

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Login(ctx *gin.Context) {
	// 1. Get ID token from the request
	idToken := ctx.GetHeader("Authorization") // Assuming the token is in the Authorization header
	if idToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID token is required"})
		return
	}

	// 2. Call the Login service
	user, tokenString, err := c.authService.Login(ctx, idToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Set the JWT in the response header
	ctx.Header("Authorization", "Bearer "+tokenString)

	// 4. Return the user data (optional)
	ctx.JSON(http.StatusOK, gin.H{"user": user, "token": tokenString})
}
