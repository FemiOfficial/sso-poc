package auth

import (
	"net/http"
	"os"
	"fmt"
	"sso-poc/internal/crypto"
	"sso-poc/internal/db"
	"sso-poc/internal/db/entitities"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
)

type AuthLib struct {
	db           *db.Database
	redis        *redis.Client
	sessionStore *sessions.CookieStore
	vaultEncrypt *crypto.TokenEncryption
}

func CreateAuthLib(db *db.Database, redis *redis.Client, vaultEncrypt *crypto.TokenEncryption) *AuthLib {
	environment := os.Getenv("ENVIRONMENT")

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = environment == "production"
	return &AuthLib{db: db, redis: redis, vaultEncrypt: vaultEncrypt}
}

func (lib *AuthLib) InitiateAuthSession(context *gin.Context, app *entitities.App, providers []string) (*string, error, int, gin.H) {
	sessionId := uuid.New().String()
	authRequest := &entitities.AuthRequest{
		SessionID: sessionId,
		AppID:     app.ID,
		Providers: providers,
		State:     entitities.AuthRequestState{Status: "initiated"},
	}

	err := lib.redis.Set(context, sessionId, authRequest, 60*time.Minute).Err()
	if err != nil {
		message := "Something went wrong while initiating auth session"
		return &message, err, http.StatusInternalServerError, nil
	}
	message := "Auth session initiated successfully"
	data := gin.H{
		"sessionId": sessionId,
		"authRequest": authRequest,
		"link": fmt.Sprintf("%s/auth/%s?session_d=%s", os.Getenv("APP_URL"), providers[0], sessionId),
	}	
	return &message, nil, http.StatusOK, data
}
