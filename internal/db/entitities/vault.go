package entitities

type VaultOwnerType string
type VaultKey string

const (
	EmailVerificationToken VaultKey = "email_verification_token"
	PasswordResetToken     VaultKey = "password_reset_token"
)

const (
	VaultKeyOrganization     VaultKey = "organization"
	VaultKeyUser             VaultKey = "user"
	VaultKeyApp              VaultKey = "app"
	VaultKeyIdentityProvider VaultKey = "identity_provider"
)

// this will be changed  to hashicorp vault
type Vault struct {
	BaseEntity
	Object    string    `gorm:"type:text;not null"` // encrypted json string
	Key       VaultKey  `gorm:"type:varchar(255);not null"`
	OwnerID   string    `gorm:"not null;type:varchar(255)"`
	OwnerType VaultOwnerType `gorm:"not null;default:organization"`
}

func (Vault) TableName() string {
	return "vaults"
}
