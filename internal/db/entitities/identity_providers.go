package entitities

type IdentityProvider struct {
	BaseEntity
	Name string `gorm:"not null;unique" json:"name"`
	DisplayName string `gorm:"not null" json:"display_name"`
	Logo string `gorm:"null" json:"logo"`
	Scopes StringArray `gorm:"type:text[];null;default:null" json:"scopes"`
	Status string `gorm:"not null;enum:active,inactive" json:"status"`
	IsDefault bool `gorm:"not null;default:false" json:"is_default"`
	CredentialFields StringArray `gorm:"type:text[];null;default:null" json:"credentials"`
}

func (IdentityProvider) TableName() string {
	return "identity_providers"
}
