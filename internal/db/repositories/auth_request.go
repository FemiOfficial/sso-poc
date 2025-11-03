package repositories

import (
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

type AuthRequestRepository struct {
	db *gorm.DB
}

type AuthRequestFilter struct {
	SessionID string `json:"session_id"`
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

func (r *AuthRequestRepository) FindByFilter(filter AuthRequestFilter, tx *gorm.DB) (*entitities.AuthRequest, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.AuthRequest{})

	if filter.SessionID != "" {
		query = query.Where("session_id = ?", filter.SessionID)
	}

	authRequest := &entitities.AuthRequest{}
	return authRequest, query.First(authRequest).Preload("App").Error
}