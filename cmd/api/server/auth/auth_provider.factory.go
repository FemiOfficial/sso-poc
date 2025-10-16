package auth

import (
	"encoding/json"
	"fmt"
	"sso-poc/internal/crypto"
	"sso-poc/internal/db/entitities"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/linkedin"
)

func CreateProvider(appIdentityProvider *entitities.AppIdentityProvider, vaultEncrypt *crypto.TokenEncryption, callbackURL string) (goth.Provider, error) {
	scopes := appIdentityProvider.App.Scopes
	if scopes == "" {
		scopes = getDefaultScopes(appIdentityProvider.IdentityProvider.Name)
	}

	decryptedKeys, err := vaultEncrypt.Decrypt(appIdentityProvider.Vault.Object)

	var keys map[string]string
	err = json.Unmarshal([]byte(decryptedKeys), &keys)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	switch appIdentityProvider.IdentityProvider.Name {
	case "google":
		return google.New(keys["key"], keys["secret"], callbackURL, scopes), nil
	case "github":
		return github.New(keys["key"], keys["secret"], callbackURL, scopes), nil
	case "facebook":
		return facebook.New(keys["key"], keys["secret"], callbackURL, scopes), nil
	case "linkedin":
		return linkedin.New(keys["key"], keys["secret"], callbackURL, scopes), nil
	default:
		return nil, fmt.Errorf("invalid provider: %s", appIdentityProvider.IdentityProvider.Name)
	}
}

func getDefaultScopes(providerName string) string {
	defaults := map[string]string{
		"google":   "email profile",
		"github":   "user:email, user:read, user:email, user:read.email, read:org",
		"facebook": "email,public_profile",
		"linkedin": "email,profile,openid",
	}

	if scope, ok := defaults[providerName]; ok {
		return scope
	}
	return "email profile"
}
