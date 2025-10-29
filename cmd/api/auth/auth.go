package auth

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	key    = "key"
	isProd = false
	maxAge = 86400 * 30
)

type Auth struct {
	Store sessions.Store
}

func NewAuth() *Auth {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")
	if callbackURL == "" {
		callbackURL = "http://localhost:8080/auth/google/callback"
	}

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	store.Options.SameSite = http.SameSiteLaxMode

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, callbackURL, "email", "profile"),
	)

	return &Auth{
		Store: store,
	}
}

func (a *Auth) BeginAuth(w http.ResponseWriter, r *http.Request) {
	gothic.Store = a.Store
	gothic.BeginAuthHandler(w, r)
}

func (a *Auth) CompleteAuth(w http.ResponseWriter, r *http.Request) (goth.User, error) {
	gothic.Store = a.Store
	return gothic.CompleteUserAuth(w, r)
}

func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) error {
	gothic.Store = a.Store
	return gothic.Logout(w, r)
}
