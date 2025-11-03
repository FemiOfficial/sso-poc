package auth

import (
	"os"
	"sso-poc/internal/crypto"
	"sso-poc/internal/db"
	"sso-poc/internal/db/entitities"

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

func (a *AuthLib) InitiateAuthSession(app *entitities.App, providers []string) {
	environment := os.Getenv("ENVIRONMENT")

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

}
