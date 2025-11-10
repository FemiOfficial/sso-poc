package middlewares

import (
	"os"
	"net/http"
	"github.com/gin-gonic/gin"
	"sso-poc/internal/utils"
)

func WidgetAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientSecret:= ctx.GetHeader("Client-Id-Default")
		expeectClientSecret := os.Getenv("WIDGET_CLIENT_SECRET")

		if clientSecret != expeectClientSecret {
			message := "Unauthorized request"
			ctx.JSON(http.StatusUnauthorized, utils.GenericApiResponse(http.StatusUnauthorized, &message, nil))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}