package entitities

import "time"

type User struct {
	BaseEntity
	FirstName       string       `gorm:"not null" json:"first_name"`
	LastName        string       `gorm:"not null" json:"last_name"`
	Email           string       `gorm:"not null; unique" json:"email"`
	Password        string       `gorm:"type:varchar(255); null" json:"password"`
	MfaEnabled      bool         `gorm:"default:false" json:"mfa_enabled"`
	EmailVerified   bool         `gorm:"default:false" json:"email_verified"`
	EmailVerifiedAt time.Time    `gorm:"type:timestamp; null; default:null" json:"email_verified_at"`
	OrganizationID  string       `gorm:"not null" json:"organization_id"`
	Organization    *Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"organization"`
}

func (User) TableName() string {
	return "users"
}
