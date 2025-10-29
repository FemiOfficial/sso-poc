package repositories

import (
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserFilter struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	OrganizationID string `json:"organization_id"`
}

func CreateUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entitities.User, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}

	return tx.Create(user).Error
}

func (r *UserRepository) FindOneByFilter(filter UserFilter, tx *gorm.DB) (*entitities.User, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.User{}).Preload("Organization")

	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}

	if filter.OrganizationID != "" {
		query = query.Where("organization_id = ?", filter.OrganizationID)
	}

	user := &entitities.User{}
	return user, query.First(user).Error
}

func (r *UserRepository) FindAllByFilter(filter UserFilter, tx *gorm.DB) ([]*entitities.User, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.User{}).Preload("Organization")

	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}

	if filter.OrganizationID != "" {
		query = query.Where("organization_id = ?", filter.OrganizationID)
	}

	users := []*entitities.User{}
	return users, query.Find(&users).Error
}
