package appTypes

type CreateAppRequest struct {
	Name string `json:"name" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"required,min=3,max=255"`
	RedirectURI string `json:"redirect_uri" binding:"required,url"`
	Scopes string `json:"scopes" binding:"required,min=3,max=255"`
	MfaEnabled bool `json:"mfa_enabled" binding:"required,boolean"`
}
