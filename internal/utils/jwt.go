package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID         string `json:"user_id"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
	Email          string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(userId string, organizationId string, email string) (*string, *string, *int, error) {

	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	var jwtIssuer = os.Getenv("JWT_ISSUER")
	var jwtExpiration, _ = strconv.Atoi(os.Getenv("JWT_HOURS"))
	var refreshTokenExpiration, _ = strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS"))

	expirationTime := time.Now().Add(time.Duration(jwtExpiration) * time.Hour) // 24-hour expiration
	refreshTokenExpirationTime := time.Now().Add(time.Duration(refreshTokenExpiration) * 24 * time.Hour)

	claims := &CustomClaims{
		UserID:         userId,
		OrganizationID: organizationId,
		Email:          email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtIssuer,
		},
	}

	refreshTokenClaims := &CustomClaims{
		UserID:         userId,
		OrganizationID: organizationId,
		Email:          email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, nil, nil, err
	}

	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return nil, nil, nil, err
	}

	return &tokenString, &refreshTokenString, &jwtExpiration, nil
}

func VerifyJWT(token string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return claims, err
}
