package repositories

import (
	"errors"
	"sso-poc/internal/crypto"
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

type VaultRepository struct {
	db           *gorm.DB
	vaultEncrypt *crypto.TokenEncryption
}

type VaultFilter struct {
	ID        string `json:"id"`
	OwnerID   string `json:"owner_id"`
	OwnerType string `json:"owner_type"`
	Key       string `json:"key"`
}

func CreateVaultRepository(db *gorm.DB, vaultEncrypt *crypto.TokenEncryption) *VaultRepository {
	return &VaultRepository{db: db, vaultEncrypt: vaultEncrypt}
}

func (r *VaultRepository) Create(vault *entitities.Vault, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}

	return tx.Create(vault).Error
}

func (r *VaultRepository) GetDecryptedObject(vaultId string, tx *gorm.DB) (*string, error) {
	vault, err := r.FindOneByFilter(VaultFilter{ID: vaultId}, tx)

	if vault == nil {
		return nil, errors.New("vault not found")
	}

	if vault.Object == "" {
		return nil, errors.New("vault object is empty")
	}

	decryptedObject, err := r.vaultEncrypt.Decrypt(vault.Object)
	if err != nil {
		return nil, err
	}

	return &decryptedObject, nil
}

func (r *VaultRepository) EncryptObjectAndSave(object string, vault *entitities.Vault, tx *gorm.DB) error {

	if tx == nil {
		tx = r.db
	}

	encryptedObject, err := r.vaultEncrypt.Encrypt(object)
	if err != nil {
		return err
	}

	vault.Object = encryptedObject
	return tx.Save(vault).Error
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
