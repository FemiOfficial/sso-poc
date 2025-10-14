package auth

import (
	"fmt"
	"net/http"
	"os"
	"sso-poc/internal/db/entitities"
	"time"

	"sso-poc/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	db    *db.Database
	redis *redis.Client
}

func CreateAuthService(db *db.Database, redis *redis.Client) *AuthService {
	return &AuthService{db: db, redis: redis}
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
		State:     entitities.AuthRequestStateInitiated,
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

// func (s *AuthService) GetAuthRequest(ctx *gin.Context) {
// 	sessionId := ctx.Query("session_id")
// 	authRequest := &entitities.AuthRequest{}
// 	s.redis.Get(ctx, sessionId).Scan(authRequest)
// 	ctx.JSON(http.StatusOK, authRequest)
// }

func (s *AuthService) GetAuthProfileData(ctx *gin.Context) {
	sessionId := ctx.Query("session_id")

	

}


