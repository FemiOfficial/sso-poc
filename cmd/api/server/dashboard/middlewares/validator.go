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
		err := ctx.ShouldBind(&request)

		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errrors := formatValidationErrors(validationErrors)
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

func ValidateRequestQuery[T any]() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request T
		err := ctx.ShouldBindQuery(&request)
		if err != nil {

			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errrors := formatValidationErrors(validationErrors)
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "invalid request query", "data": errrors})
				ctx.Abort()
				return
			}

			fmt.Println("Error: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "failed", 
				"message": "invalid request query, please check the request query and try again", 
				"data": err.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Set("request", request)
		ctx.Next()
	}
}

func formatValidationErrors(validationErrors validator.ValidationErrors) map[string]string {
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
		case "email":
			errrors[fieldKey] = "This field must be a valid email address"
		case "url":
			errrors[fieldKey] = "This field must be a valid URL"
		case "boolean":
			errrors[fieldKey] = "This field must be a valid boolean"
		case "integer":
			errrors[fieldKey] = "This field must be a valid integer"
		case "float":
			errrors[fieldKey] = "This field must be a valid float"
		case "array":
			errrors[fieldKey] = "This field must be a valid array"
		case "object":
			errrors[fieldKey] = "This field must be a valid object"
		case "enum":
			errrors[fieldKey] = "This field must be a valid enum"
		case "uuid":
			errrors[fieldKey] = "This field must be a valid UUID"
		case "date":
			errrors[fieldKey] = "This field must be a valid date"
		case "time":
			errrors[fieldKey] = "This field must be a valid time"
		}
	}
	return errrors
}
