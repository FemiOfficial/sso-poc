package misc

import (
	"fmt"
	"net/http"
	"sso-poc/internal/utils"

	"github.com/gin-gonic/gin"
)

type MiscController struct {
	miscService *MiscService
}

func CreateMiscController(miscService *MiscService) *MiscController {
	return &MiscController{miscService: miscService}
}

func (c *MiscController) GetIdentityProviders(ctx *gin.Context) {
	response, err := c.miscService.GetIdentityProviders(ctx)
	message := "Identity providers fetched successfully"
	if err != nil {
		fmt.Println("Error fetching identity providers: ", err)
		message = err.Error()
		ctx.JSON(http.StatusInternalServerError, utils.GenericApiResponse(http.StatusInternalServerError, &message, nil))
		return
	}
	ctx.JSON(http.StatusOK, utils.GenericApiResponse(http.StatusOK, &message, response))
}
