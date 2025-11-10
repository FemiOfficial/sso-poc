package repositories

import (
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

type IdentityProviderRepository struct {
	db *gorm.DB
}

type IdentityProviderFilter struct {
	Status    string   `json:"status"`
	IDs       []string `json:"ids"`
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Names     []string `json:"names"`
	Scopes    []string `json:"scopes"`
	IsDefault bool     `json:"is_default"`
}

func CreateIdentityProviderRepository(db *gorm.DB) *IdentityProviderRepository {
	return &IdentityProviderRepository{db: db}
}

func (r *IdentityProviderRepository) CreateMany(identityProviders []*entitities.IdentityProvider, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}

	return tx.CreateInBatches(identityProviders, 100).Error
}

func (r *IdentityProviderRepository) FindAllByFilter(filter IdentityProviderFilter, tx *gorm.DB)  ([]*entitities.IdentityProvider, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.IdentityProvider{})

	if len(filter.IDs) > 0 {
		query = query.Where("id IN (?)", filter.IDs)
	}

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}

	if filter.IsDefault {
		query = query.Where("is_default = ?", filter.IsDefault)
	}

	if filter.Names != nil {
		query = query.Where("name IN (?)", filter.Names)
	}

	if len(filter.Scopes) > 0 {
		query = query.Where("scopes @> ?", filter.Scopes)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	} else {
		query = query.Where("status = ?", "active")
	}

	identityProviders := []*entitities.IdentityProvider{}
	return identityProviders, query.Find(&identityProviders).Error
}

func (r *IdentityProviderRepository) FindOneByFilter(filter IdentityProviderFilter, tx *gorm.DB) (*entitities.IdentityProvider, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.IdentityProvider{})
	
	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}
	

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	identityProvider := &entitities.IdentityProvider{}
	return identityProvider, query.First(identityProvider).Error
}