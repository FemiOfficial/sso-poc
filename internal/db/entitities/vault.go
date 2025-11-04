package entitities

type VaultOwnerType string
type VaultKey string

const (
	OrganizationVerificationSecret VaultKey = "organization_verification_secret"
	UserVerificationSecret VaultKey = "user_verification_secret"
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
	Object    string    `gorm:"type:text;not null" json:"object"` // encrypted json string
	Key       VaultKey  `gorm:"type:varchar(255);not null" json:"key"`
	OwnerID   string    `gorm:"not null;type:varchar(255)" json:"owner_id"`
	OwnerType VaultOwnerType `gorm:"not null;default:organization" json:"owner_type"`
}

func (Vault) TableName() string {
	return "vaults"
}
