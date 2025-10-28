package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"sso-poc/cmd/api/server/auth"
	"sso-poc/internal/cache"
	"sso-poc/internal/db"
	"sso-poc/internal/crypto"
	"sso-poc/cmd/api/server/dashboard/organisation"
)

type Server struct {
	port int
	db             *db.Database
	authController *auth.AuthController
	organizationController *organisation.OrganizationController
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := db.InitializeDB()
	redis := cache.CreateRedisClient()
	vaultEncrypt, err := crypto.NewTokenEncryption()
	if err != nil {
		panic(fmt.Sprintf("Failed to create vault helper: %v", err))
	}
	NewServer := &Server{
		port: port,
		// auth:           auth.NewAuth(),
		db:             db,
		authController: auth.CreateAuthController(auth.CreateAuthService(db, redis, vaultEncrypt)),
		organizationController: organisation.CreateOrganizationController(organisation.CreateOrganizationService(db)),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
