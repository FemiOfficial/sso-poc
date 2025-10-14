package entitities


// this will be changed  to hashicorp vault
type Vault struct {
	BaseEntity
	EncryptedAccessTokens  string           `gorm:"not null"`
	EncryptedRefreshTokens string           `gorm:"not null"`
	IdentityProviderID     string           `gorm:"not null"`
	IdentityProvider       IdentityProvider `gorm:"foreignKey:IdentityProviderID"`
	OrganizationID         string           `gorm:"not null"`
	Organization           Organization     `gorm:"foreignKey:OrganizationID"`
	AppID                  string           `gorm:"not null"`
	App                    App              `gorm:"foreignKey:AppID"`
	Scopes                 string           `gorm:"type:text[]"`
	UserID                 string           `gorm:"not null"`
	User                   User             `gorm:"foreignKey:UserID"`
}

func (Vault) TableName() string {
	return "vaults"
}
