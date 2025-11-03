package repositories

import (
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

type AuthRequestRepository struct {
	db *gorm.DB
}

func CreateAuthRequestRepository(db *gorm.DB) *AuthRequestRepository {
	return &AuthRequestRepository{db: db}
}

func (r *AuthRequestRepository) Create(authRequest *entitities.AuthRequest, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}

	return tx.Create(authRequest).Error
}