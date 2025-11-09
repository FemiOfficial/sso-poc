package utils

import (
	"crypto/rand"
	"encoding/base32"

	"github.com/gin-gonic/gin"
)

func GenericErrorMessages() map[int]string {
	return map[int]string{
		500: "Something went wrong please try again",
		401: "Unauthorized request",
		404: "Not Found",
	}
}

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

func ConvertToSet(values []string) []string {

	tracker := make(map[string]bool)
	scopeSet := []string{}
	for _, value := range values {
		if tracker[value] {
			continue
		}
		tracker[value] = true
		scopeSet = append(scopeSet, value)
	}

	return scopeSet
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

func Contains(base []string, subset []string) bool {
	if len(subset) == 0 {
		return true
	}

	seen := make(map[string]struct{}, len(base))
	for _, n := range base {
		seen[n] = struct{}{}
	}

	for _, c := range subset {
		if _, ok := seen[c]; !ok {
			return false
		}
	}
	return true
}