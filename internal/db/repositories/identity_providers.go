package repositories

import (
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

type IdentityProviderRepository struct {
	db *gorm.DB
}

type IdentityProviderFilter struct {
	Status string   `json:"status"`
	IDs    []string `json:"ids"`
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

func CreateIdentityProviderRepository(db *gorm.DB) *IdentityProviderRepository {
	return &IdentityProviderRepository{db: db}
}

func (r *IdentityProviderRepository) FindAllByFilter(filter IdentityProviderFilter, tx *gorm.DB) ([]*entitities.IdentityProvider, error) {
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
