package organisation

import (
	"sso-poc/internal/db"
	"sso-poc/internal/db/entitities"
	"time"

	"sso-poc/cmd/api/server/dashboard/organisation/types"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type OrganizationService struct {
	db        *db.Database
	validator *validator.Validate
}

func CreateOrganizationService(db *db.Database) *OrganizationService {
	return &OrganizationService{db: db, validator: validator.New()}
}

func (s *OrganizationService) CreateOrganization(createOrganizationRequest types.CreateOrganizationRequest) (*entitities.Organization, error) {

	// send verification email to org
	// save token to redis (1hour)

	token := crypto.GenerateToken()

	err := s.redis.Set(ctx, token, organizationId, 1*time.Hour).Err()
	if err != nil {
		return nil, err
	}
	return organization, nil
}

func (s *OrganizationService) VerifyOrganizationEmail(organizationId string) error {
	// organization := &entitities.Organization{}
	// err := s.db.DB.Where("id = ?", organizationId).First(organization).Error
	// if err != nil {
	// 	return err
	// }
	organization.EmailVerified = true
	return nil
}

func (s *OrganizationService) LoginOrganization(loginOrganizationRequest types.LoginOrganizationRequest) (*entitities.Organization, error) {
	organization := &entitities.Organization{}
	err := s.db.DB.Where("email = ?", loginOrganizationRequest.Email).First(organization).Error
	if err != nil {
		return nil, err
	}
	return organization, nil
}
