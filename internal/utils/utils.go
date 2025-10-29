package utils

import (
	"crypto/rand"
	"encoding/base32"

	"github.com/gin-gonic/gin"
)

func GenericApiResponse(code int, message *string, data any) gin.H {
	result := gin.H{}

	switch code {
	case 200:
		result["status"] = "success"
	default:
		result["status"] = "failed"
	}

	result["message"] = getDefaultMessage(code, message)

	if data != nil {
		result["data"] = data
	} else {
		result["data"] = nil
	}

	return result
}

func getDefaultMessage(code int, message *string) string {
	if message != nil {
		return *message
	}

	switch code {
	case 200:
		return "Success"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized request"
	case 500:
		return "Internal Server Error"
	default:
		return "Something went wrong"
	}
}
func GenerateRandomString(length int) (string, error) {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(buf), nil
}
