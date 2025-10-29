package repositories

import (
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

type VaultRepository struct {
	db *gorm.DB
}

type VaultFilter struct {
	ID        string `json:"id"`
	OwnerID   string `json:"owner_id"`
	OwnerType string `json:"owner_type"`
	Key       string `json:"key"`
}

func CreateVaultRepository(db *gorm.DB) *VaultRepository {
	return &VaultRepository{db: db}
}

func (r *VaultRepository) Create(vault *entitities.Vault, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}

	return tx.Create(vault).Error
}

func (r *VaultRepository) FindOneByFilter(filter VaultFilter, tx *gorm.DB) (*entitities.Vault, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.Vault{})

	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.OwnerID != "" {
		query = query.Where("owner_id = ?", filter.OwnerID)
	}

	if filter.OwnerType != "" {
		query = query.Where("owner_type = ?", filter.OwnerType)
	}

	if filter.Key != "" {
		query = query.Where("key = ?", filter.Key)
	}

	vault := &entitities.Vault{}
	return vault, query.First(vault).Error
}

func (r *VaultRepository) FindAllByFilter(filter VaultFilter, tx *gorm.DB) ([]*entitities.Vault, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.Vault{})

	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.OwnerID != "" {
		query = query.Where("owner_id = ?", filter.OwnerID)
	}

	if filter.OwnerType != "" {
		query = query.Where("owner_type = ?", filter.OwnerType)
	}

	if filter.Key != "" {
		query = query.Where("key = ?", filter.Key)
	}

	vaults := []*entitities.Vault{}
	return vaults, query.Find(&vaults).Error
}
