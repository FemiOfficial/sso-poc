package auth

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"sso-poc/cmd/lib/auth/factories"
	"sso-poc/internal/crypto"
	"sso-poc/internal/db"
	"sso-poc/internal/db/entitities"
	"sso-poc/internal/db/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/redis/go-redis/v9"
)

type AuthLib struct {
	db                            *db.Database
	redis                         *redis.Client
	sessionStore                  *sessions.CookieStore
	vaultEncrypt                  *crypto.TokenEncryption
	authRequestRepository         *repositories.AuthRequestRepository
	appIdentityProviderRepository *repositories.AppIdentityProviderRepository
}

func CreateAuthLib(db *db.Database, redis *redis.Client, vaultEncrypt *crypto.TokenEncryption, authRequestRepository *repositories.AuthRequestRepository) *AuthLib {
	environment := os.Getenv("ENVIRONMENT")

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = environment == "production"
	appIdentityProviderRepository := repositories.CreateAppIdentityProviderRepository(db.DB)
	authRequestRepository = repositories.CreateAuthRequestRepository(db.DB)

	return &AuthLib{
		db:                            db,
		redis:                         redis,
		vaultEncrypt:                  vaultEncrypt,
		authRequestRepository:         authRequestRepository,
		appIdentityProviderRepository: appIdentityProviderRepository,
	}
}

func (lib *AuthLib) InitiateAuthSession(context *gin.Context, app *entitities.App, providers []string) (*string, error, int, gin.H) {
	sessionId := uuid.New().String()

	if len(providers) == 0 || providers[0] == "" {
		appIdentityProviders, err := lib.appIdentityProviderRepository.FindAllByFilter(repositories.AppIdentityProviderFilter{AppID: app.ID}, nil)
		if err != nil {
			message := "Something went wrong while getting app identity providers"
			return &message, err, http.StatusInternalServerError, nil
		}

		for _, appIdentityProvider := range appIdentityProviders {
			providers = append(providers, appIdentityProvider.IdentityProvider.Name)
		}

		if len(providers) > 0 && providers[0] == "" {
			providers = providers[1:]
		}
	}

	authRequest := &entitities.AuthRequest{
		SessionID:   sessionId,
		AppID:       app.ID,
		ProviderIDs: providers,
		State:       entitities.AuthRequestState{Status: "initiated"},
	}

	err := lib.authRequestRepository.Create(authRequest, nil)
	if err != nil {
		message := "Something went wrong while creating auth request"
		return &message, err, http.StatusInternalServerError, nil
	}

	message := "Auth request created successfully"
	data := gin.H{
		"sessionId":   sessionId,
		"authRequest": authRequest,
		"link":        fmt.Sprintf("%s/auth/%s?session_id=%s", os.Getenv("APP_URL"), providers[0], sessionId),
	}
	return &message, nil, http.StatusOK, data
}

func (lib *AuthLib) LoginUser(context *gin.Context, app *entitities.App, provider string, sessionId string) (*string, error, int, gin.H) {
	authRequest := &entitities.AuthRequest{}

	if authRequest.State.Status != "initiated" {
		message := "auth request is not initiated"
		return &message, nil, http.StatusBadRequest, nil
	}

	if !slices.Contains(authRequest.ProviderIDs, provider) {
		message := "provider is not valid for this session"
		return &message, nil, http.StatusBadRequest, nil
	}

	callbackURL := fmt.Sprintf("%s/auth/%s/%s/callback", os.Getenv("APP_URL"), provider, sessionId)

	appIdentityProvider, err := lib.appIdentityProviderRepository.FindOneByFilter(repositories.AppIdentityProviderFilter{AppID: authRequest.AppID, Provider: provider}, nil)

	providerInstance, err := factories.CreateProvider(appIdentityProvider, lib.vaultEncrypt, callbackURL)
	if err != nil {
		return nil, err, http.StatusInternalServerError, nil
	}

	message := "auth request found successfully"
	providerInstance.SetName(provider)

	goth.UseProviders(providerInstance)

	gothic.BeginAuthHandler(context.Writer, context.Request)
	return &message, nil, http.StatusOK, nil
}
