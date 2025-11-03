package factories

import (
	"encoding/json"
	"fmt"
	"sso-poc/internal/crypto"
	"sso-poc/internal/db/entitities"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/linkedin"
)

func CreateProvider(appIdentityProvider *entitities.AppIdentityProvider, vaultEncrypt *crypto.TokenEncryption, callbackURL string) (goth.Provider, error) {
	provider := appIdentityProvider.IdentityProvider
	decryptedKeys, err := vaultEncrypt.Decrypt(appIdentityProvider.Vault.Object)

	var keys map[string]string
	err = json.Unmarshal([]byte(decryptedKeys), &keys)
	if err != nil {
		return nil, err
	}

	switch appIdentityProvider.IdentityProvider.Name {
	case "google":
		return google.New(keys["key"], keys["secret"], callbackURL, strings.Join(provider.Scopes, ",")), nil
	case "github":
		return github.New(keys["key"], keys["secret"], callbackURL, strings.Join(provider.Scopes, ",")), nil
	case "facebook":
		return facebook.New(keys["key"], keys["secret"], callbackURL, strings.Join(provider.Scopes, ",")), nil
	case "linkedin":
		return linkedin.New(keys["key"], keys["secret"], callbackURL, strings.Join(provider.Scopes, ",")), nil
	default:
		return nil, fmt.Errorf("invalid provider: %s", appIdentityProvider.IdentityProvider.Name)
	}
}
