package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			for _, err := range ctx.Errors {
				fmt.Println("Error: ", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}
