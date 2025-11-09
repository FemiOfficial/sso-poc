package repositories

import (
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

type AuthIdentityProviderRepository struct {
	db *gorm.DB
}

func CreateAuthIdentityProviderRepository(db *gorm.DB) *AuthIdentityProviderRepository {
	return &AuthIdentityProviderRepository{db: db}
}

func (r *AuthIdentityProviderRepository) Create(authIdentityProvider *entitities.AuthIdentityProvider, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(authIdentityProvider).Error
}
