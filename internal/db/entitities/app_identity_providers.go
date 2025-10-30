package entitities

type AppIdentityProvider struct {
	BaseEntity
	AppID              string           `gorm:"not null"`
	App                App              `gorm:"foreignKey:AppID"`
	IdentityProviderID string           `gorm:"not null"`
	IdentityProvider   IdentityProvider `gorm:"foreignKey:IdentityProviderID"`
	Status             string           `gorm:"not null;enum:active,inactive;default:active"`
	IsDefault          bool             `gorm:"not null;default:false"`
	Scopes             []string         `gorm:"type:text[];not null"`
	VaultId            *string          `gorm:"type:varchar(255);null"`
	Vault              Vault            `gorm:"Id:VaultId,Type:vault"`
}

func (AppIdentityProvider) TableName() string {
	return "app_identity_providers"
}
