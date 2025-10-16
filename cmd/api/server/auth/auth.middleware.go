package auth

import (
	"net/http"
	"sso-poc/internal/db"
	"sso-poc/internal/db/entitities"

	"github.com/gin-gonic/gin"
)

func ClientAuthMiddleware(db *db.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		clientId := ctx.GetHeader("Client-Id")
		clientSecret := ctx.GetHeader("Client-Secret")

		if clientId == "" || clientSecret == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing client credentials"})
			ctx.Abort()
			return
		}

		app := &entitities.App{}
		db.DB.Where("client_id = ?", clientId).First(app)

		if app.ID == "" || app.ClientSecret != clientSecret {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid client credentials"})
			ctx.Abort()
			return
		}

		// Store app in context for use in handlers
		ctx.Set("app", app)
		ctx.Next()
	}
}
