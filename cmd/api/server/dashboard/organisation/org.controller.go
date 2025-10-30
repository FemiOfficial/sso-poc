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
	organization, err, statusCode := c.organizationService.CreateOrganization(ctx)
	message := "Organization created successfully"

	if err != nil {
		fmt.Println("Error creating organization: ", err, statusCode)
		message = err.Error()
		ctx.JSON(*statusCode, utils.GenericApiResponse(*statusCode, &message, nil))
		return
	}
	ctx.JSON(http.StatusOK, utils.GenericApiResponse(http.StatusOK, &message, organization))
}

func (c *OrganizationController) VerifyOrganizationEmail(ctx *gin.Context) {
	err, message, statusCode := c.organizationService.VerifyOrganizationEmail(ctx)
	if err != nil {
		fmt.Println("Error verifying organization email: ", err, statusCode)
		ctx.JSON(*statusCode, utils.GenericApiResponse(*statusCode, message, nil))
		return
	}
	ctx.JSON(*statusCode, utils.GenericApiResponse(*statusCode, message, nil))
}

func (c *OrganizationController) LoginOrganization(ctx *gin.Context) {
	err, response, statusCode := c.organizationService.LoginOrganization(ctx)
	message := "login successful"
	if err != nil {
		fmt.Println("Error logging in organization user: ", err, statusCode)
		message = err.Error()
		ctx.JSON(*statusCode, utils.GenericApiResponse(*statusCode, &message, nil))
		return
	}
	ctx.JSON(*statusCode, utils.GenericApiResponse(*statusCode, &message, response))
}

func (c *OrganizationController) ResendEmailVerificationOtp(ctx *gin.Context) {
	err, response, statusCode := c.organizationService.ResendEmailVerificationOtp(ctx)
	message := "otp resent successfully"
	if err != nil {
		fmt.Println("Error resending verification otp: ", err, statusCode)
		message = err.Error()
		ctx.JSON(*statusCode, utils.GenericApiResponse(*statusCode, &message, nil))
		return
	}
	ctx.JSON(*statusCode, utils.GenericApiResponse(*statusCode, &message, response))
}
