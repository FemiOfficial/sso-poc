package repositories

import (
	entitities "sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

type AppIdentityProviderRepository struct {
	db *gorm.DB
}

type AppIdentityProviderFilter struct {
	ID        string   `json:"id"`
	IDs       []string `json:"ids"`
	Status    string   `json:"status"` // active, inactive
	Provider  string   `json:"provider"`
	AppID     string   `json:"app_id"`
	IsDefault bool     `json:"is_default"`
}

func CreateAppIdentityProviderRepository(db *gorm.DB) *AppIdentityProviderRepository {
	return &AppIdentityProviderRepository{db: db}
}

func (r *AppIdentityProviderRepository) CreateMany(appIdentityProviders []*entitities.AppIdentityProvider, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}

	return tx.CreateInBatches(appIdentityProviders, 100).Error
}

func (r *AppIdentityProviderRepository) FindOneByFilter(filter AppIdentityProviderFilter, tx *gorm.DB) (*entitities.AppIdentityProvider, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.AppIdentityProvider{}).Preload("IdentityProvider")

	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.AppID != "" {
		query = query.Where("app_id = ?", filter.AppID)
	}

	if filter.Provider != "" {
		query = query.Where("identity_provider.name = ?", filter.Provider)
	}

	appIdentityProvider := &entitities.AppIdentityProvider{}
	return appIdentityProvider, query.First(appIdentityProvider).Error
}

func (r *AppIdentityProviderRepository) FindAllByFilter(filter AppIdentityProviderFilter, tx *gorm.DB) ([]*entitities.AppIdentityProvider, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.AppIdentityProvider{}).Preload("IdentityProvider")

	if len(filter.IDs) > 0 {
		query = query.Where("id IN (?)", filter.IDs)
	}

	if filter.AppID != "" {
		query = query.Where("app_id = ?", filter.AppID)
	}

	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.IsDefault {
		query = query.Where("is_default = ?", filter.IsDefault)
	}

	appIdentityProviders := []*entitities.AppIdentityProvider{}
	return appIdentityProviders, query.Find(&appIdentityProviders).Error
}
