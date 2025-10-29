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

func(c * AuthController) LoginUser(ctx *gin.Context) {
	c.authService.LoginUser(ctx)
}

func (c * AuthController) Callback(ctx *gin.Context) {
	c.authService.Callback(ctx)
}

func (c * AuthController) GetAuthProfileData(ctx *gin.Context) {
	c.authService.GetAuthProfileData(ctx)
}
