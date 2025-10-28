package types

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required,min=3,max=255"`
	Domain string `json:"domain" binding:"required,min=3,max=255"`
	Logo string `json:"logo" binding:"omitempty,min=3,max=255"`
	Description string `json:"description" binding:"required,min=3,max=255"`
	Location string `json:"location" binding:"required,min=3,max=255"`
	Industry string `json:"industry" binding:"required,min=3,max=255"`
	Size int `json:"size" binding:"required,min=1,max=1000000"`
	MfaEnabled bool `json:"mfa_enabled" binding:"omitempty,boolean"`
}

type LoginOrganizationRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}