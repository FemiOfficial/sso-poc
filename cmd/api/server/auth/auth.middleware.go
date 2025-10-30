package auth

import (
	"net/http"
	"sso-poc/internal/db"
	"sso-poc/internal/db/entitities"
	"sso-poc/internal/utils"
	"github.com/gin-gonic/gin"
)

func ClientAuthMiddleware(db *db.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		clientId := ctx.GetHeader("Client-Id")
		clientSecret := ctx.GetHeader("Client-Secret")

		var message string = "Missing client credentials"

		if clientId == "" || clientSecret == "" {
			ctx.JSON(http.StatusUnauthorized, utils.GenericApiResponse(http.StatusUnauthorized, &message, nil))
			ctx.Abort()
			return
		}

		app := &entitities.App{}
		db.DB.Where("client_id = ?", clientId).First(app)

		if app.ID == "" || app.ClientSecret != clientSecret {
			message = "Invalid client credentials"
			ctx.JSON(http.StatusUnauthorized, utils.GenericApiResponse(http.StatusUnauthorized, &message, nil))
			ctx.Abort()
			return
		}

		// Store app in context for use in handlers
		ctx.Set("app", app)
		ctx.Next()
	}
}
