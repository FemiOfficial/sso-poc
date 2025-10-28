package utils

import (
	"crypto/rand"
	"encoding/base32"
)

func GenerateRandomString(length int) (string, error) {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(buf), nil
}