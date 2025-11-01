package miscTypes

import "sso-poc/internal/db/entitities"

type GetIDPRequest struct {
	Status string `json:"status" binding:"omitempty,oneof=active inactive"`
	IDs    []string `json:"ids" binding:"omitempty,min=1"`
	Name   string `json:"name" binding:"omitempty,min=1"`
	Scopes []string `json:"scopes" binding:"omitempty,min=1"`
}

type GetIDPResponse struct {
	IdentityProviders []entitities.IdentityProvider `json:"identity_providers"`
}