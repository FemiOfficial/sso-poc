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
)

type Server struct {
	port int
	db             *db.Database
	authController *auth.AuthController
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := db.InitializeDB()
	redis := cache.CreateRedisClient()
	NewServer := &Server{
		port: port,
		// auth:           auth.NewAuth(),
		db:             db,
		authController: auth.CreateAuthController(auth.CreateAuthService(db, redis)),
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
