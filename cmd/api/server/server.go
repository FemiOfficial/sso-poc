package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"sso-poc/cmd/api/server/auth"
	"sso-poc/cmd/api/server/dashboard/app"
	"sso-poc/cmd/api/server/dashboard/misc"
	"sso-poc/cmd/api/server/dashboard/organisation"
	"sso-poc/internal/cache"
	"sso-poc/internal/crypto"
	"sso-poc/internal/db"
	"sso-poc/internal/db/repositories"
	authLib "sso-poc/cmd/lib/auth"
)

type Server struct {
	port                   int
	db                     *db.Database
	authLib                *authLib.AuthLib
	authController         *auth.AuthController
	organizationController *organisation.OrganizationController
	miscController         *misc.MiscController
	appController          *app.AppController
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := db.InitializeDB()
	redis := cache.CreateRedisClient()
	vaultEncrypt, err := crypto.NewTokenEncryption()
	authLib := authLib.CreateAuthLib(db, redis, vaultEncrypt, repositories.CreateAuthRequestRepository(db.DB))

	if err != nil {
		panic(fmt.Sprintf("Failed to create vault helper: %v", err))
	}
	NewServer := &Server{
		port: port,
		// auth:           auth.NewAuth(),
		db:                     db,
		authController:         auth.CreateAuthController(auth.CreateAuthService(authLib)),
		organizationController: organisation.CreateOrganizationController(organisation.CreateOrganizationService(db, redis, vaultEncrypt)),
		miscController:         misc.CreateMiscController(misc.CreateMiscService(db)),
		appController:          app.CreateAppController(app.CreateAppService(db, redis, vaultEncrypt)),
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
