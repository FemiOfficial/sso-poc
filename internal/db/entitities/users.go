package entitities

type User struct {
	BaseEntity
	FirstName      string       `gorm:"not null"`
	LastName       string       `gorm:"not null"`
	Email          string       `gorm:"not null"`
	Verified       bool         `gorm:"not null"`
	VaultReference string       `gorm:"not null"`
	OrganizationID string       `gorm:"not null"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
}

func (User) TableName() string {
	return "users"
}
