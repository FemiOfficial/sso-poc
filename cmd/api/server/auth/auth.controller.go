package auth

import (
	"fmt"
	"sso-poc/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *AuthService
}

func CreateAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) InitiateAuthSession(ctx *gin.Context) {
	message, err, statusCode, data := c.authService.InitiateAuthSession(ctx)
	if err != nil {
		fmt.Println("Error initiating auth session: ", err, statusCode)
		ctx.JSON(statusCode, utils.GenericApiResponse(statusCode, message, nil))
		return
	}
	ctx.JSON(statusCode, utils.GenericApiResponse(statusCode, message, data))
}

// func(c * AuthController) LoginUser(ctx *gin.Context) {
// 	c.authService.LoginUser(ctx)
// }

// func (c * AuthController) Callback(ctx *gin.Context) {
// 	c.authService.Callback(ctx)
// }

// func (c * AuthController) GetAuthProfileData(ctx *gin.Context) {
// 	c.authService.GetAuthProfileData(ctx)
// }
