package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateRequestBody[T any]() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var request T
		err := ctx.ShouldBind(&request); 

		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errrors := make(map[string]string)
				for _, fieldErr := range validationErrors {
					field := fieldErr.Field()
					tag := fieldErr.Tag()
					fieldKey := strings.ToLower(field[:1]) + field[1:]
					switch tag {
					case "required":
						errrors[fieldKey] = "This field is required"
					case "min":
						errrors[fieldKey] = "This field must be at least " + fieldErr.Param() + " characters long"
					case "max":
						errrors[fieldKey] = "This field must be at most " + fieldErr.Param() + " characters long"
					}
				}
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "invalid request body", "data": errrors})
				ctx.Abort()
				return
			}
			fmt.Println("Error: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "invalid request body, please check the request body and try again", "data": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set("request", request)
		ctx.Next()
	}
}
