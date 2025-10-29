package repositories

import (
	"sso-poc/cmd/api/server/dashboard/organisation/types"
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

type OrganisationRepository struct {
	db *gorm.DB
}

func CreateOrganisationRepository(db *gorm.DB) *OrganisationRepository {
	return &OrganisationRepository{db: db}
}

func (r *OrganisationRepository) Create(organisation *types.CreateOrganizationRequest, tx *gorm.DB) error {
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
	return tx.Create(organization).Error
}

// func (r *OrganisationRepository) GetOrganisationByEmail(email string) (*entitities.Organization, error) {
// 	organisation := &entitities.Organization{}
// 	return organisation, r.db.Where("email = ?", email).First(organisation).Error
// }

// func (r *OrganisationRepository) GetOrganisationByDomain(domain string) (*entitities.Organization, error) {
// 	organisation := &entitities.Organization{}
// 	return organisation, r.db.Where("domain = ?", domain).First(organisation).Error
// }
