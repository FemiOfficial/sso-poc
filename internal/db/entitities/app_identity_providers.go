package entitities

type AppIdentityProvider struct {
	BaseEntity
	AppID              string           `gorm:"not null"`
	App                App              `gorm:"foreignKey:AppID"`
	IdentityProviderID string           `gorm:"not null"`
	IdentityProvider   IdentityProvider `gorm:"foreignKey:IdentityProviderID"`
	Status             string           `gorm:"not null;enum:active,inactive;default:active"`
	Scopes             StringArray      `gorm:"type:text[];null;default:null"`
	VaultId            string           `gorm:"type:varchar(255);null;default:null"`
	Vault              Vault            `gorm:"foreignKey:VaultId"`
}

func (AppIdentityProvider) TableName() string {
	return "app_identity_providers"
}
