package factories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sso-poc/internal/crypto"
	"sso-poc/internal/db/entitities"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/linkedin"

	// "github.com/markbates/goth/providers/microsoft"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/okta"
	"github.com/markbates/goth/providers/twitter"

	// "github.com/markbates/goth/providers/openid"
	"github.com/markbates/goth/providers/auth0"
)

func CreateProvider(appIdentityProvider *entitities.AppIdentityProvider, vaultEncrypt *crypto.TokenEncryption, callbackURL string) (goth.Provider, error) {
	provider := appIdentityProvider.IdentityProvider
	decryptedKeys, err := vaultEncrypt.Decrypt(appIdentityProvider.Vault.Object)

	var keys map[string]string
	err = json.Unmarshal([]byte(decryptedKeys), &keys)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}

	switch appIdentityProvider.IdentityProvider.Name {
	case "google":
		return google.New(keys["key"], keys["secret"], callbackURL, strings.Join(provider.Scopes, ",")), nil
	case "github":
		return github.New(keys["key"], keys["secret"], callbackURL, strings.Join(provider.Scopes, ",")), nil
	case "facebook":
		return facebook.New(keys["key"], keys["secret"], callbackURL, strings.Join(provider.Scopes, ",")), nil
	case "linkedin":
		return linkedin.New(keys["key"], keys["secret"], callbackURL, strings.Join(provider.Scopes, ",")), nil
	case "twitter":
		return twitter.New(keys["key"], keys["secret"], callbackURL), nil
	case "apple":
		return apple.New(keys["key"], keys["secret"], callbackURL, httpClient, strings.Join(provider.Scopes, ",")), nil
	case "amazon":
		return amazon.New(keys["key"], keys["secret"], callbackURL), nil
	case "okta":
		return okta.New(keys["id"], keys["secret"], keys["domain"], callbackURL, strings.Join(provider.Scopes, ",")), nil
	// case "openid":
	// 	return openid.New(keys["key"], keys["secret"], keys["discovery_url"], callbackURL, keys["discovery_url"]), nil
	case "auth0":
		return auth0.New(keys["domain"], callbackURL, keys["domain"], strings.Join(provider.Scopes, ",")), nil
	// case "microsoft":
	// 	return microsoft.New(keys["online_key"], keys["online_secret"], callbackURL), nil
	default:
		return nil, fmt.Errorf("invalid provider: %s", appIdentityProvider.IdentityProvider.Name)
	}
}
