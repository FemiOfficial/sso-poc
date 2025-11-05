package entitities

type AppIdentityProvider struct {
	BaseEntity
	AppID              string           `gorm:"not null" json:"app_id"`
	App                *App              `gorm:"foreignKey:AppID" json:"app"`
	IdentityProviderID string           `gorm:"not null" json:"identity_provider_id"`
	IdentityProvider   *IdentityProvider `gorm:"foreignKey:IdentityProviderID" json:"identity_provider"`
	Status             string           `gorm:"not null;enum:active,inactive;default:active" json:"status"`
	Scopes             StringArray      `gorm:"type:text[];null;default:null" json:"scopes"`
	VaultId            string           `gorm:"type:varchar(255);null;default:null" json:"vault_id"`
	Vault              *Vault            `gorm:"foreignKey:VaultId" json:"vault"`
}

func (AppIdentityProvider) TableName() string {
	return "app_identity_providers"
}
