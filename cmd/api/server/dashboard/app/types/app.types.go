package appTypes

import "sso-poc/internal/db/entitities"


type IdentityProviders struct {
	ID string `json:"id"`
	Scopes []string `json:"scopes" binding:"omitempty,min=1,max=2"`
}
type CreateAppRequest struct {
	Name string `json:"name" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"required,min=3,max=255"`
	RedirectURI string `json:"redirect_uri" binding:"required,url"`
	MfaEnabled bool `json:"mfa_enabled" binding:"omitempty,boolean"`
	Scopes []entitities.AppScope `json:"scopes" binding:"omitempty,min=1,max=2"`
	IdentityProviders []IdentityProviders `json:"identity_providers" binding:"omitempty,min=1,max=2"`
}

type GetAppIdentityProviderResponse struct {
	IdentityProviderID string `json:"identity_provider_id"`
	VaultId string `json:"vault_id"`
	Status string `json:"status"`
	Scopes []string `json:"scopes"`
}

type GetAppResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	RedirectURI string `json:"redirect_uri"`
	MfaEnabled bool `json:"mfa_enabled"`
	ClientID string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackURI string `json:"callback_uri"`
	Live bool `json:"live"`
	Scopes []string`json:"scopes"`
	IdentityProviders []GetAppIdentityProviderResponse `json:"identity_providers"`
}

type AppIdentityProviderCredentials struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
type UpdateAppIdentityProviderRequest struct {
	ID string `json:"id" binding:"required"`
	Scopes []string `json:"scopes" binding:"omitempty,min=1,max=2"`
	Credentials []AppIdentityProviderCredentials `json:"credentials" binding:"omitempty,min=1,max=2"`
}
