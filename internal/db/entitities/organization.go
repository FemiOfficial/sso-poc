package entitities

type Organization struct {
	BaseEntity
	Name string `gorm:"not null"`
	MfaEnabled bool `gorm:"not null"`
}

func (Organization) TableName() string {
	return "organizations"
}
