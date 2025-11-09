package entitities

type AuthIdentityProvider struct {
	BaseEntity
	AuthRequestID string `gorm:"not null" json:"auth_request_id"`
	AuthRequest *AuthRequest `gorm:"foreignKey:AuthRequestID" json:"auth_request"`
	IdentityProviderID string `gorm:"not null" json:"identity_provider_id"`
	IdentityProvider *IdentityProvider `gorm:"foreignKey:IdentityProviderID" json:"identity_provider"`
}

func (AuthIdentityProvider) TableName() string {
	return "auth_identity_providers"
}