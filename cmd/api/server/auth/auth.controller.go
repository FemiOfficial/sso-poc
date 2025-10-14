package auth

import "github.com/gin-gonic/gin"

type AuthController struct {
	authService *AuthService
}

func CreateAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) InitiateAuthSession(ctx *gin.Context) {
	c.authService.InitiateAuthSession(ctx)
}
