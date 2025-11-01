package app

import (
	"fmt"
	"sso-poc/internal/utils"

	"github.com/gin-gonic/gin"
)

type AppController struct {
	appService *AppService
}

func CreateAppController(appService *AppService) *AppController {
	return &AppController{appService: appService}
}

func (c *AppController) CreateApp(ctx *gin.Context) {
	appId, err, statusCode := c.appService.CreateApp(ctx)
	message := "App created successfully"
	if err != nil {
		fmt.Println("Error creating app: ", err, statusCode)
		message = err.Error()
		ctx.JSON(*statusCode, utils.GenericApiResponse(*statusCode, &message, nil))
		return
	}
	ctx.JSON(*statusCode, utils.GenericApiResponse(*statusCode, &message, appId))
}
