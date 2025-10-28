package organisation

import (
	"fmt"
	"net/http"

	"sso-poc/cmd/api/server/dashboard/organisation/types"

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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Something went wrong", "data": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Organization created successfully", "data": organization})
}

func (c *OrganizationController) LoginOrganization(ctx *gin.Context) {
	organization, err := c.organizationService.LoginOrganization(ctx.MustGet("request").(types.LoginOrganizationRequest))
	if err != nil {
		fmt.Println("Error logging in organization: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Something went wrong", "data": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Organization logged in successfully", "data": organization})
}	