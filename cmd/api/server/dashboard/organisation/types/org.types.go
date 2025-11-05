package organisationTypes

import "time"

type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Domain      string `json:"domain" binding:"required,min=3,max=255"`
	Logo        string `json:"logo" binding:"omitempty,min=3,max=255"`
	Description string `json:"description" binding:"required,min=3,max=255"`
	Location    string `json:"location" binding:"required,min=3,max=255"`
	Industry    string `json:"industry" binding:"required,min=3,max=255"`
	Size        int    `json:"size" binding:"required,min=1,max=1000000"`
	Email       string `json:"email" binding:"required,email"`
}

type LoginOrganizationRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type VerifyEmailRequest struct {
	Otp         string `json:"otp" binding:"required,min=6,max=6"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=255"`
	OldPassword string `json:"old_password" binding:"omitempty,min=8,max=255"`
}

type LoginOrganizationResponseData struct {
	UserId             string    `json:"user_id"`
	Token              string    `json:"token"`
	RefreshToken       string    `json:"refresh_token"`
	ExpiresIn          int       `json:"expires_in"`
	Email              string    `json:"email"`
	EmailVerified      bool      `json:"email_verified"`
	MfaEnabled         bool      `json:"mfa_enabled"`
	EmailVerifiedAt    time.Time `json:"email_verified_at"`
	OrganizationId     string    `json:"organization_id"`
	OrganizationName   string    `json:"organization_name"`
	OrganizationDomain string    `json:"organization_domain"`
	OrganizationLogo   string    `json:"organization_logo"`
}

type ResendEmailVerificationOtpRequest struct {
	Email string `json:"email" binding:"required"`
}
