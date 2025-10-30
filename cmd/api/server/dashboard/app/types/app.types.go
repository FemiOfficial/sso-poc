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
