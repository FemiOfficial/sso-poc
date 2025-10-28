package organisation

import (
	"encoding/json"
	"fmt"
	"sso-poc/internal/crypto"
	"sso-poc/internal/db"
	"sso-poc/internal/db/entitities"
	"sso-poc/internal/utils"

	"github.com/gin-gonic/gin"

	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"

	"sso-poc/cmd/api/server/dashboard/organisation/types"

	"github.com/redis/go-redis/v9"
)

type OrganizationService struct {
	db           *db.Database
	redis        *redis.Client
	validator    *validator.Validate
	vaultEncrypt *crypto.TokenEncryption
}

func CreateOrganizationService(db *db.Database, redis *redis.Client, vaultEncrypt *crypto.TokenEncryption) *OrganizationService {
	return &OrganizationService{db: db, redis: redis, validator: validator.New(), vaultEncrypt: vaultEncrypt}
}

func (s *OrganizationService) CreateOrganization(ctx *gin.Context) (*entitities.Organization, error) {
	var organization *entitities.Organization
	var createOrganizationRequest types.CreateOrganizationRequest = ctx.MustGet("request").(types.CreateOrganizationRequest)

	err := s.db.DB.Transaction(func(tx *gorm.DB) error {
		organization = &entitities.Organization{
			Name:          createOrganizationRequest.Name,
			Domain:        createOrganizationRequest.Domain,
			Logo:          createOrganizationRequest.Logo,
			Description:   createOrganizationRequest.Description,
			Location:      createOrganizationRequest.Location,
			Industry:      createOrganizationRequest.Industry,
			Size:          createOrganizationRequest.Size,
			EmailVerified: false,
		}

		err := tx.Create(organization).Error
		if err != nil {
			return err
		}

		secret, err := utils.GenerateRandomString(32)
		if err != nil {
			return err
		}

		otp, err := totp.GenerateCode(secret, time.Now())
		if err != nil {
			return err
		}

		cacheKey := fmt.Sprintf("email_verification_token:%s", otp)

		cacheValue := map[string]string{
			"secret":          secret,
			"organization_id": organization.ID,
		}
		cacheValueJSON, err := json.Marshal(cacheValue)
		if err != nil {
			return err
		}

		err = s.redis.Set(ctx, cacheKey, cacheValueJSON, 1*time.Hour).Err()
		if err != nil {
			return err
		}

		err = s.saveOrgVerificationSecret(tx, secret, organization.ID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return organization, nil

	// send verification email to org
	// save token to redis (1hour
}

// func (s *OrganizationService) VerifyOrganizationEmail(organizationId string) error {
// 	// organization := &entitities.Organization{}
// 	// err := s.db.DB.Where("id = ?", organizationId).First(organization).Error
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	organization.EmailVerified = true
// 	return nil
// }

func (s *OrganizationService) LoginOrganization(loginOrganizationRequest types.LoginOrganizationRequest) (*entitities.Organization, error) {
	organization := &entitities.Organization{}
	err := s.db.DB.Where("email = ?", loginOrganizationRequest.Email).First(organization).Error
	if err != nil {
		return nil, err
	}
	return organization, nil
}

func (s *OrganizationService) saveOrgVerificationSecret(tx *gorm.DB, secret string, organizationId string) error {
	vaultObject := map[string]string{
		"secret": secret,
	}

	vaultObjectJSON, err := json.Marshal(vaultObject)
	if err != nil {
		return err
	}
	vaultObjectEncrypted, err := s.vaultEncrypt.Encrypt(string(vaultObjectJSON))
	if err != nil {
		return err
	}

	vault := &entitities.Vault{
		Key:       entitities.OrganizationVerificationSecret,
		OwnerID:   organizationId,
		OwnerType: "organization",
		Object:    vaultObjectEncrypted,
	}
	err = tx.Create(vault).Error
	if err != nil {
		return err
	}
	return nil
}
