package entitities

type AppIdentityProvider struct {
	BaseEntity
	AppID              string           `gorm:"not null"`
	App                App              `gorm:"foreignKey:AppID"`
	IdentityProviderID string           `gorm:"not null"`
	IdentityProvider   IdentityProvider `gorm:"foreignKey:IdentityProviderID"`
}

func (AppIdentityProvider) TableName() string {
	return "app_identity_providers"
}
