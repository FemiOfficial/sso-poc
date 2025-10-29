package entitities

type IdentityProvider struct {
	BaseEntity
	Name string `gorm:"not null"`
	Status string `gorm:"not null;enum:active,inactive"`
}

func (IdentityProvider) TableName() string {
	return "identity_providers"
}
