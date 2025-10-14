package entitities

type App struct {
	BaseEntity
	Name           string       `gorm:"not null"`
	ClientID       string       `gorm:"not null"`
	ClientSecret   string       `gorm:"not null"`
	RedirectURI    string       `gorm:"not null"`
	Scopes         string       `gorm:"not null"`
	MfaEnabled     bool         `gorm:"not null"`
	OrganizationID string       `gorm:"not null"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
}

func (App) TableName() string {
	return "apps"
}
