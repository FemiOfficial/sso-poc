package entitities

type IdentityProvider struct {
	BaseEntity
	Name string `gorm:"not null"`
}

func (IdentityProvider) TableName() string {
	return "identity_providers"
}
