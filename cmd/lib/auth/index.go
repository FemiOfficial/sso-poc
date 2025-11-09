package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"sso-poc/cmd/lib/auth/factories"
	"sso-poc/cmd/lib/auth/types"
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
	"gorm.io/gorm"
)

type AuthLib struct {
	db                             *db.Database
	redis                          *redis.Client
	sessionStore                   *sessions.CookieStore
	vaultEncrypt                   *crypto.TokenEncryption
	authRequestRepository          *repositories.AuthRequestRepository
	appIdentityProviderRepository  *repositories.AppIdentityProviderRepository
	authIdentityProviderRepository *repositories.AuthIdentityProviderRepository
}

func (lib *AuthLib) GetDB() *db.Database {
	return lib.db
}

func CreateAuthLib(
	db *db.Database,
	redis *redis.Client,
	vaultEncrypt *crypto.TokenEncryption,
	authRequestRepository *repositories.AuthRequestRepository,
	appIdentityProviderRepository *repositories.AppIdentityProviderRepository,
	authIdentityProviderRepository *repositories.AuthIdentityProviderRepository,
) *AuthLib {
	environment := os.Getenv("ENVIRONMENT")

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = environment == "production"
	appIdentityProviderRepository = repositories.CreateAppIdentityProviderRepository(db.DB)
	authRequestRepository = repositories.CreateAuthRequestRepository(db.DB)
	authIdentityProviderRepository = repositories.CreateAuthIdentityProviderRepository(db.DB)

	return &AuthLib{
		db:                             db,
		redis:                          redis,
		vaultEncrypt:                   vaultEncrypt,
		authRequestRepository:          authRequestRepository,
		appIdentityProviderRepository:  appIdentityProviderRepository,
		authIdentityProviderRepository: authIdentityProviderRepository,
	}
}

func (lib *AuthLib) InitiateAuthSession(context *gin.Context,
	app *entitities.App, providers []string) (*string, error, int, gin.H) {
	sessionId := uuid.New().String()
	var statusCode int = http.StatusInternalServerError
	var message string = "something went wrong while intiating session"

	authRequest := &entitities.AuthRequest{
		SessionID: sessionId,
		AppID:     app.ID,
		State:     entitities.AuthRequestState{Status: "initiated"},
	}

	err := lib.GetDB().DB.Transaction(func(tx *gorm.DB) error {
		authIdentityProviders := []entitities.AuthIdentityProvider{}
		providersQuery := repositories.AppIdentityProviderFilter{AppID: app.ID}

		if len(providers) == 0 || providers[0] != "" {
			providersQuery.ProviderIDs = providers
		}

		appIdentityProviders, err := lib.appIdentityProviderRepository.FindAllByFilter(providersQuery, tx)
		if err != nil {
			return err
		}

		err = lib.authRequestRepository.Create(authRequest, tx)
		if err != nil {
			return err
		}

		for _, appIdentityProvider := range appIdentityProviders {
			authIdentityProvider := &entitities.AuthIdentityProvider{
				AuthRequestID:      authRequest.ID,
				IdentityProviderID: appIdentityProvider.IdentityProvider.ID,
			}

			// TODO: make this better not multiple creates
			err := lib.authIdentityProviderRepository.Create(authIdentityProvider, tx)
			authIdentityProviders = append(authIdentityProviders, *authIdentityProvider)
			if err != nil {
				return err
			}
		}

		authRequest.AuthIdentityProviders = authIdentityProviders
		err = tx.Save(authRequest).Error
		if err != nil {
			return err
		}

		return err
	})

	if err != nil {
		return &message, err, statusCode, nil
	}

	message = "Auth request created successfully"
	data := gin.H{
		"sessionId":   sessionId,
		"authRequest": authRequest,
		"link":        fmt.Sprintf("%s/auth/%s", os.Getenv("APP_URL"), sessionId),
	}
	return &message, nil, http.StatusOK, data
}

func (lib *AuthLib) ResolveSession(sessionId string) (*types.ResolveSessionResponse, error) {
	authRequest, err := lib.authRequestRepository.FindByFilter(repositories.AuthRequestFilter{SessionID: sessionId}, nil)
	if err != nil {
		return nil, err
	}

	if authRequest.State.Status != "initiated" {
		message := "invalid auth request"
		return nil, errors.New(message)
	}

	providerResponses := make([]types.SessionProviders, len(authRequest.AuthIdentityProviders))
	for i, provider := range authRequest.AuthIdentityProviders {
		providerResponses[i] = types.SessionProviders{
			ID:               provider.IdentityProvider.ID,
			Name:             provider.IdentityProvider.Name,
			Logo:             provider.IdentityProvider.Logo,
			Scopes:           provider.IdentityProvider.Scopes,
			DisplayName:      provider.IdentityProvider.DisplayName,
			CredentialFields: provider.IdentityProvider.CredentialFields,
			Status:           provider.IdentityProvider.Status,
		}
	}

	return &types.ResolveSessionResponse{
		SessionID: authRequest.SessionID,
		AppID:     authRequest.AppID,
		Status:    authRequest.State.Status,
		State: types.AuthRequestState{
			Status: authRequest.State.Status,
			Data:   authRequest.State.Data,
		},
		Providers: providerResponses,
	}, nil
}

func (lib *AuthLib) LoginUser(context *gin.Context, app *entitities.App, provider string, sessionId string) (*string, error, int, gin.H) {
	authRequest, err := lib.authRequestRepository.FindByFilter(repositories.AuthRequestFilter{SessionID: sessionId}, nil)
	if err != nil {
		message := "auth request not found"
		return &message, err, http.StatusNotFound, nil
	}

	if authRequest.State.Status != "initiated" {
		message := "auth request is not initiated"
		return &message, nil, http.StatusBadRequest, nil
	}

	providerIds := []string{}
	for _, authRequestProvider := range authRequest.AuthIdentityProviders {
		providerIds = append(providerIds, authRequestProvider.IdentityProvider.ID)
	}

	if !slices.Contains(providerIds, provider) {
		message := "provider is not valid for this session"
		return &message, nil, http.StatusBadRequest, nil
	}

	callbackURL := fmt.Sprintf("%s/auth/%s/%s/callback", os.Getenv("APP_URL"), provider, sessionId)

	appIdentityProvider, err := lib.appIdentityProviderRepository.FindOneByFilter(repositories.AppIdentityProviderFilter{
		AppID:    authRequest.AppID,
		Provider: provider,
	}, nil)
	if err != nil {
		message := "provider configuration not found"
		return &message, err, http.StatusInternalServerError, nil
	}

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
