package auth

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"sso-poc/internal/db/entitities"
	"time"

	"sso-poc/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"

	"sso-poc/cmd/api/server/auth/types"
	"sso-poc/internal/crypto"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	db           *db.Database
	redis        *redis.Client
	sessionStore *sessions.CookieStore
	vaultEncrypt *crypto.TokenEncryption
}

func CreateAuthService(db *db.Database, redis *redis.Client, vaultEncrypt *crypto.TokenEncryption) *AuthService {

	environment := os.Getenv("ENVIRONMENT")

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = environment == "production"
	return &AuthService{db: db, redis: redis, sessionStore: store, vaultEncrypt: vaultEncrypt}
}

func (s *AuthService) InitiateAuthSession(ctx *gin.Context) {

	clientId := ctx.GetHeader("Client-Id")
	clientSecret := ctx.GetHeader("Client-Secret")

	dbConnection := s.db.DB

	providers := []string{}

	if ctx.Query("provider") != "" {
		dbConnection.Where("name = ?", ctx.Query("provider")).First(entitities.IdentityProvider{})
		isValidProvider := &entitities.IdentityProvider{}
		dbConnection.Where("name = ?", ctx.Query("provider")).First(isValidProvider)
		if isValidProvider.ID != "" {
			providers = append(providers, isValidProvider.ID)
		}
	}

	app := &entitities.App{}
	dbConnection.Where("client_id = ?", clientId).First(app)

	if app.ClientSecret != clientSecret {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid client secret"})
		return
	}

	if len(providers) > 0 {
		appIdentityProviders := []entitities.AppIdentityProvider{}
		dbConnection.Where("app_id = ?", app.ID).Find(&appIdentityProviders)
		if len(appIdentityProviders) == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid provider"})
			return
		}

		for _, appIdentityProvider := range appIdentityProviders {
			providers = append(providers, appIdentityProvider.IdentityProviderID)
		}
	}

	sessionId := uuid.New().String()
	authRequest := &entitities.AuthRequest{
		SessionID: sessionId,
		AppID:     app.ID,
		Providers: providers,
		State:     entitities.AuthRequestState{Status: "initiated"},
	}

	// cache this in redis
	err := s.redis.Set(ctx, sessionId, authRequest, 3*time.Minute).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong while initiating auth session"})
		return
	}

	link := fmt.Sprintf("%s/auth/%s?session_id=%s", os.Getenv("APP_URL"), providers[0], sessionId)

	ctx.JSON(http.StatusOK, gin.H{"session_id": sessionId, "auth_request": authRequest, "link": link})
}

func (s *AuthService) LoginUser(ctx *gin.Context) {
	// create new auth with goth
	// begin auth session or complete it if that is the case
	// begin auth will calls .BeginAuth from goth api

	var loginUserRequest types.LoginUserRequest
	if err := ctx.ShouldBindJSON(&loginUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider := loginUserRequest.Provider
	sessionId := loginUserRequest.SessionID

	app := ctx.MustGet("app").(*entitities.App)

	authRequest := &entitities.AuthRequest{}
	err := s.redis.Get(ctx, sessionId).Scan(authRequest)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, session not found"})
		return
	}

	if authRequest.State.Status != "initiated" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, session is not yet initiated"})
		return
	}

	if !slices.Contains(authRequest.Providers, provider) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid provider, provider is not valid for this session"})
		return
	}

	callbackURL := fmt.Sprintf("%s/auth/%s/%s/callback", os.Getenv("APP_URL"), provider, sessionId)

	appIdentityProvider := &entitities.AppIdentityProvider{}
	s.db.DB.Joins("IdentityProvider").Where("app_id = ?", app.ID).Where("identity_provider.name = ?", provider).First(appIdentityProvider)
	providerInstance, err := CreateProvider(appIdentityProvider, s.vaultEncrypt, callbackURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	providerInstance.SetName(provider)

	goth.UseProviders(providerInstance)

	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (s *AuthService) Callback(ctx *gin.Context) {
	provider := ctx.Param("provider")
	sessionId := ctx.Query("session_id")
	
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

func (s *AuthService) GetAuthProfileData(ctx *gin.Context) {
	sessionId := ctx.Query("session_id")
	dbConnection := s.db.DB

	// Get app from context (set by middleware)
	app := ctx.MustGet("app").(*entitities.App)

	authRequest := &entitities.AuthRequest{}
	err := s.redis.Get(ctx, sessionId).Scan(authRequest)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, session not found"})
		return
	}

	if authRequest.AppID != app.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, session is not for this app"})
		return
	}

	if authRequest.State.Status != "auth_completed" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, session is not yet completed"})
		return
	}

	user := &entitities.User{}

	email := authRequest.State.Data["email"].Email

	if email == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, email not found"})
		return
	}

	dbConnection.Where("email = ?", authRequest.State.Data["email"].Email).First(user)
	if user.ID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
