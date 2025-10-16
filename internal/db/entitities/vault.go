package entitities

type VaultOwner struct {
	Id   string `gorm:"not null;type:varchar(255)"`
	Type string `gorm:"not null;enum:organization,user,app,identity_provider"`
}

// this will be changed  to hashicorp vault
type Vault struct {
	BaseEntity
	Object      string `gorm:"type:text;not null"` // encrypted json string
	Owner              VaultOwner        `gorm:"Id:Id,Type:Type"`
}

func (Vault) TableName() string {
	return "vaults"
}
