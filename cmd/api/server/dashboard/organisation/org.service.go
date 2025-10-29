package organisation

import (
	"encoding/json"
	"errors"
	"fmt"
	"sso-poc/internal/crypto"
	"sso-poc/internal/db"
	"sso-poc/internal/db/entitities"
	"sso-poc/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"

	"sso-poc/cmd/api/server/dashboard/organisation/types"
	"sso-poc/internal/db/repositories"

	"github.com/redis/go-redis/v9"
)

type OrganizationService struct {
	db                     *db.Database
	redis                  *redis.Client
	validator              *validator.Validate
	vaultEncrypt           *crypto.TokenEncryption
	organisationRepository *repositories.OrganisationRepository
	userRepository         *repositories.UserRepository
}

func CreateOrganizationService(db *db.Database, redis *redis.Client, vaultEncrypt *crypto.TokenEncryption) *OrganizationService {
	return &OrganizationService{
		db:                     db,
		redis:                  redis,
		validator:              validator.New(),
		vaultEncrypt:           vaultEncrypt,
		organisationRepository: repositories.CreateOrganisationRepository(db.DB),
		userRepository:         repositories.CreateUserRepository(db.DB),
	}
}

func (s *OrganizationService) CreateOrganization(ctx *gin.Context) (*entitities.Organization, error) {
	var organization *entitities.Organization
	var createOrganizationRequest types.CreateOrganizationRequest = ctx.MustGet("request").(types.CreateOrganizationRequest)
	err := s.db.DB.Transaction(func(tx *gorm.DB) error {
		err := s.validateOrganizationParams(tx, createOrganizationRequest.Email, createOrganizationRequest.Domain); if err != nil {
			return err
		}

		err = s.organisationRepository.Create(&createOrganizationRequest, tx); if err != nil {
			return err
		}

		user := &entitities.User{
			Email:          createOrganizationRequest.Email,
			EmailVerified:  false,
			MfaEnabled:     false,
			OrganizationID: organization.ID,
		}
		err = s.userRepository.Create(user, tx); if err != nil {
			return err
		}	

		secret, err := utils.GenerateRandomString(32); if err != nil {
			return err
		}

		otp, err := totp.GenerateCode(secret, time.Now()); if err != nil {
			return err
		}

		cacheValue := map[string]string{
			"secret":          secret,
			"organization_id": organization.ID,
			"user_id":         user.ID,
			"email":           user.Email,
		}
		cacheValueJSON, err := json.Marshal(cacheValue); if err != nil {
			return err
		}

		err = s.redis.Set(ctx, fmt.Sprintf("email_verification_token:%s", otp), cacheValueJSON, 1*time.Hour).Err(); if err != nil {
			return err
		}

		err = s.saveUserVerificationSecret(tx, secret, user.ID); if err != nil {
			return err
		}

		return nil
	})

	return organization, err
}

func (s *OrganizationService) VerifyOrganizationEmail(ctx *gin.Context) (error, string) {

	var verifyEmailRequest types.VerifyEmailRequest = ctx.MustGet("request").(types.VerifyEmailRequest)

	cacheValue := map[string]string{}

	err := s.redis.Get(ctx, fmt.Sprintf("email_verification_token:%s", verifyEmailRequest.Otp)).Scan(&cacheValue)
	if err != nil {
		return err, ""
	}

	secret := cacheValue["secret"]
	userId := cacheValue["user_id"]

	isValid := totp.Validate(verifyEmailRequest.Otp, secret)
	if !isValid {
		return errors.New("invalid otp token"), ""
	}

	err = s.db.DB.Transaction(func(tx *gorm.DB) error {
		user := &entitities.User{}
		err := tx.Where("id = ?", userId).First(user).Error
		if err != nil {
			return err
		}

		user.EmailVerified = true
		user.EmailVerifiedAt = time.Now()
		err = tx.Save(user).Error
		if err != nil {
			return err
		}

		err = s.saveNewPassword(tx, user, verifyEmailRequest.NewPassword, verifyEmailRequest.OldPassword)
		return err
	})
	if err != nil {
		return err, ""
	}
	return nil, "password updated successfully"
}

func (s *OrganizationService) LoginOrganization(ctx *gin.Context) (error, *types.LoginOrganizationResponse) {
	var loginOrganizationRequest types.LoginOrganizationRequest = ctx.MustGet("request").(types.LoginOrganizationRequest)

	user := &entitities.User{}
	err := s.db.DB.Joins("Organization").Where("email = ?", loginOrganizationRequest.Email).First(user).Error
	if err != nil {
		return err, nil
	}

	if user.EmailVerified == false {
		return errors.New("email not verified"), nil
	}

	if user.Password == "" {
		return errors.New("please choose a password to login"), nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginOrganizationRequest.Password))
	if err != nil {
		return errors.New("invalid password"), nil
	}

	token, refreshToken, jwtExpiration, err := utils.GenerateJWT(user.ID, user.OrganizationID, user.Email)
	if err != nil {
		return err, nil
	}

	return nil, &types.LoginOrganizationResponse{
		Message: "login successful",
		Data: types.LoginOrganizationResponseData{
			UserId:             user.ID,
			Token:              *token,
			RefreshToken:       *refreshToken,
			ExpiresIn:          *jwtExpiration,
			Email:              user.Email,
			EmailVerified:      user.EmailVerified,
			MfaEnabled:         user.MfaEnabled,
			EmailVerifiedAt:    user.EmailVerifiedAt,
			OrganizationId:     user.OrganizationID,
			OrganizationName:   user.Organization.Name,
			OrganizationDomain: user.Organization.Domain,
			OrganizationLogo:   user.Organization.Logo,
		},
	}
}

func (s *OrganizationService) saveNewPassword(tx *gorm.DB, user *entitities.User, newPassword string, oldPassword string) error {
	if user.Password != "" && oldPassword == "" {
		return errors.New("old password is required")
	}

	if user.Password == "" && newPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
		if err != nil {
			return errors.New("invalid old password")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to save new password")
		}
		user.Password = string(hashedPassword)

	}

	err := tx.Save(user).Error
	if err != nil {
		return errors.New("failed to save new password")
	}

	return nil
}

func (s *OrganizationService) saveUserVerificationSecret(tx *gorm.DB, secret string, userId string) error {
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
		Key:       entitities.UserVerificationSecret,
		OwnerID:   userId,
		OwnerType: "user",
		Object:    vaultObjectEncrypted,
	}
	err = tx.Create(vault).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *OrganizationService) validateOrganizationParams(tx *gorm.DB, email string, domain string) error {
	doesOrganizationExist := s.assertOrganisationEmailDoesNotExists(tx, email)
	if doesOrganizationExist {
		return errors.New("organization already exists with this email")
	}

	doesDomainExist := s.assertDomainDoesNotExists(tx, domain)
	if doesDomainExist {
		return errors.New("domain already exists with this domain")
	}

	doesUserExist := s.assertUserEmailDoesNotExists(tx, email)
	if doesUserExist {
		return errors.New("user already exists with this email")
	}
	return nil
}

func (s *OrganizationService) assertOrganisationEmailDoesNotExists(tx *gorm.DB, email string) bool {
	doesOrganizationExist := tx.Where("email = ?", email).First(&entitities.Organization{}).Error
	if doesOrganizationExist == nil {
		return false
	}
	return true
}

func (s *OrganizationService) assertUserEmailDoesNotExists(tx *gorm.DB, email string) bool {
	doesUserExist := tx.Where("email = ?", email).First(&entitities.User{}).Error
	if doesUserExist == nil {
		return false
	}
	return true
}

func (s *OrganizationService) assertDomainDoesNotExists(tx *gorm.DB, domain string) bool {
	doesDomainExist := tx.Where("domain = ?", domain).First(&entitities.Organization{}).Error
	if doesDomainExist == nil {
		return false
	}
	return true
}
