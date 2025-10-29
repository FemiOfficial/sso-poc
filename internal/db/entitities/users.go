package entitities

import "time"

type User struct {
	BaseEntity
	FirstName       string       `gorm:"not null"`
	LastName        string       `gorm:"not null"`
	Email           string       `gorm:"not null; unique"`
	Password        string       `gorm:"type:varchar(255); null"`
	MfaEnabled      bool         `gorm:"default:false"`
	EmailVerified   bool         `gorm:"default:false"`
	EmailVerifiedAt time.Time    `gorm:"type:timestamp; null; default:null"`
	OrganizationID  string       `gorm:"not null"`
	Organization    Organization `gorm:"foreignKey:OrganizationID;references:ID"`
}

func (User) TableName() string {
	return "users"
}
