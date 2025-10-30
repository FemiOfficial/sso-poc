package repositories

import (
	"sso-poc/cmd/api/server/dashboard/organisation/types"
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

type OrganisationRepository struct {
	db *gorm.DB
}

type OrganisationFilter struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Domain string `json:"domain"`
	Name string `json:"name"`
}

func CreateOrganisationRepository(db *gorm.DB) *OrganisationRepository {
	return &OrganisationRepository{db: db}
}

func (r *OrganisationRepository) Create(organisation *organisationTypes.CreateOrganizationRequest, tx *gorm.DB) (*entitities.Organization, error) {
	organization := &entitities.Organization{
		Name:        organisation.Name,
		Domain:      organisation.Domain,
		Logo:        organisation.Logo,
		Description: organisation.Description,
		Location:    organisation.Location,
		Industry:    organisation.Industry,
		Size:        organisation.Size,
		Email:       organisation.Email,
	}

	if tx == nil {
		tx = r.db
	}
	return organization, tx.Create(organization).Error
}

func (r *OrganisationRepository) FindOneByFilter(filter OrganisationFilter, tx *gorm.DB) (*entitities.Organization, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.Organization{})

	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}

	if filter.Domain != "" {
		query = query.Where("domain = ?", filter.Domain)
	}

	organisation := &entitities.Organization{}
	return organisation, query.First(organisation).Error
}

func (r *OrganisationRepository) FindAllByFilter(filter OrganisationFilter, tx *gorm.DB) ([]*entitities.Organization, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.Organization{})

	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}

	if filter.Domain != "" {
		query = query.Where("domain = ?", filter.Domain)
	}

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}

	organisations := []*entitities.Organization{}
	return organisations, query.Find(&organisations).Error
}
