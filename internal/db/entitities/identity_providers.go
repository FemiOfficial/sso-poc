package entitities


type IdentityProvider struct {
	BaseEntity
	Name string `gorm:"not null;unique" json:"name"`
	DisplayName string `gorm:"not null" json:"display_name"`
	Scopes StringArray `gorm:"type:text[];null;default:null" json:"scopes"`
	Status string `gorm:"not null;enum:active,inactive" json:"status"`
	IsDefault bool `gorm:"not null;default:false" json:"is_default"`
}

func (IdentityProvider) TableName() string {
	return "identity_providers"
}
