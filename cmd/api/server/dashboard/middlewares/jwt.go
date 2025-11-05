package middlewares

import (
	"net/http"
	"strings"

	"sso-poc/internal/db"
	"sso-poc/internal/utils"

	"github.com/gin-gonic/gin"
)

func JwtMiddleware(db *db.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		var message string = "Unauthorized request"

		if token == "" {
			ctx.JSON(http.StatusUnauthorized, utils.GenericApiResponse(http.StatusUnauthorized, &message, nil))
			ctx.Abort()
			return
		}

		token = strings.Replace(token, "Bearer ", "", 1)

		claims, err := utils.VerifyJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.GenericApiResponse(http.StatusUnauthorized, &message, nil))
			ctx.Abort()
			return
		}

		ctx.Set("user", claims)
		ctx.Next()
	}
}
