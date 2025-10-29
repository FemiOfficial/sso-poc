package organisation

import (
	"fmt"
	"net/http"
	"sso-poc/internal/utils"
	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	organizationService *OrganizationService
}

func CreateOrganizationController(organizationService *OrganizationService) *OrganizationController {
	return &OrganizationController{organizationService: organizationService}
}

func (c *OrganizationController) CreateOrganization(ctx *gin.Context) {
	organization, err := c.organizationService.CreateOrganization(ctx)
	if err != nil {
		fmt.Println("Error creating organization: ", err)
		ctx.JSON(http.StatusInternalServerError, utils.GenericApiResponse(http.StatusInternalServerError, nil, err.Error()))
		return
	}
	message := "Organization created successfully"
	ctx.JSON(http.StatusOK, utils.GenericApiResponse(http.StatusOK, &message, organization))
}

func (c *OrganizationController) VerifyOrganizationEmail(ctx *gin.Context) {
	err, message := c.organizationService.VerifyOrganizationEmail(ctx)
	if err != nil {
		fmt.Println("Error verifying organization email: ", err)
		ctx.JSON(http.StatusInternalServerError, utils.GenericApiResponse(http.StatusInternalServerError, nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.GenericApiResponse(http.StatusOK, &message, nil))
}

func (c *OrganizationController) LoginOrganization(ctx *gin.Context) {
	err, response := c.organizationService.LoginOrganization(ctx)
	if err != nil {
		fmt.Println("Error logging in organization user: ", err)
		ctx.JSON(http.StatusInternalServerError, utils.GenericApiResponse(http.StatusInternalServerError, nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.GenericApiResponse(http.StatusOK, &response.Message, response.Data))
}	