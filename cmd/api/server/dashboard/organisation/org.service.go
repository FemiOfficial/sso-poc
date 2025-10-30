package organisation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	vaultRepository        *repositories.VaultRepository
}

func CreateOrganizationService(db *db.Database, redis *redis.Client, vaultEncrypt *crypto.TokenEncryption) *OrganizationService {
	return &OrganizationService{
		db:                     db,
		redis:                  redis,
		validator:              validator.New(),
		vaultEncrypt:           vaultEncrypt,
		organisationRepository: repositories.CreateOrganisationRepository(db.DB),
		userRepository:         repositories.CreateUserRepository(db.DB),
		vaultRepository:        repositories.CreateVaultRepository(db.DB),
	}
}

func (s *OrganizationService) CreateOrganization(ctx *gin.Context) (*entitities.Organization, error, *int) {
	var organization *entitities.Organization
	var statusCode int = http.StatusInternalServerError
	var createOrganizationRequest types.CreateOrganizationRequest = ctx.MustGet("request").(types.CreateOrganizationRequest)
	err := s.db.DB.Transaction(func(tx *gorm.DB) error {
		err := s.validateOrganizationParams(tx, createOrganizationRequest.Email, createOrganizationRequest.Domain)
		if err != nil {
			statusCode = http.StatusBadRequest
			return err
		}

		organization, err := s.organisationRepository.FindOneByFilter(repositories.OrganisationFilter{
			Email: createOrganizationRequest.Email,
		}, tx)
		if errors.Is(err, gorm.ErrRecordNotFound) != true {
			statusCode = http.StatusBadRequest
			return errors.New("organization already exists with this email")
		}

		organization, err = s.organisationRepository.Create(&createOrganizationRequest, tx)
		if err != nil {
			return err
		}

		user := &entitities.User{
			Email:          createOrganizationRequest.Email,
			EmailVerified:  false,
			MfaEnabled:     false,
			OrganizationID: organization.ID,
		}
		err = s.userRepository.Create(user, tx)
		if err != nil {
			return err
		}

		if err := s.cacheVerificationCode(ctx, tx, nil, user.OrganizationID, user.ID, user.Email); err != nil {
			return err
		}

		return nil
	})
	if err == nil {
		statusCode = http.StatusOK
	}
	return organization, err, &statusCode
}

func (s *OrganizationService) VerifyOrganizationEmail(ctx *gin.Context) (error, *string, *int) {

	var statusCode = http.StatusInternalServerError
	var message string = "email verified successfully"
	var verifyEmailRequest types.VerifyEmailRequest = ctx.MustGet("request").(types.VerifyEmailRequest)

	value, err := s.redis.Get(ctx, fmt.Sprintf("email_verification_token:%s", verifyEmailRequest.Otp)).Result()
	if err != nil {
		message = "invalid otp token, please try again"
		statusCode = http.StatusBadRequest
		return err, &message, &statusCode
	}

	var cacheValue map[string]string
	if err := json.Unmarshal([]byte(value), &cacheValue); err != nil {
		message = "something went wrong with otp validation"
		statusCode = http.StatusInternalServerError
		return err, &message, &statusCode
	}

	userId := cacheValue["user_id"]

	err = s.db.DB.Transaction(func(tx *gorm.DB) error {
		user, err := s.userRepository.FindOneByFilter(repositories.UserFilter{
			ID: userId,
		}, tx)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = http.StatusBadRequest
			return errors.New("something went wrong, invalid user")
		}

		user.EmailVerified = true
		user.EmailVerifiedAt = time.Now()
		err = tx.Save(user).Error
		if err != nil {
			return err
		}

		err = s.saveNewPassword(tx, user, verifyEmailRequest.NewPassword, verifyEmailRequest.OldPassword)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		message = err.Error()
		statusCode = http.StatusInternalServerError
		return err, &message, &statusCode
	}

	statusCode = http.StatusOK
	if verifyEmailRequest.NewPassword != "" && verifyEmailRequest.OldPassword != "" {
		message = "password updated successfully"
	}
	return nil, &message, &statusCode
}

func (s *OrganizationService) LoginOrganization(ctx *gin.Context) (error, *types.LoginOrganizationResponseData, *int) {
	var loginOrganizationRequest types.LoginOrganizationRequest = ctx.MustGet("request").(types.LoginOrganizationRequest)

	var statusCode int = http.StatusInternalServerError
	var message string = "login successful"

	user, err := s.userRepository.FindOneByFilter(repositories.UserFilter{
		Email: loginOrganizationRequest.Email,
	}, s.db.DB)
	if err != nil {
		message = err.Error()
		statusCode = http.StatusInternalServerError
		return err, nil, &statusCode
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		message = "user not found"
		statusCode = http.StatusBadRequest
		return errors.New(message), nil, &statusCode
	}

	if user.EmailVerified == false {
		message = "email not verified"
		statusCode = http.StatusBadRequest
		return errors.New("email not verified"), nil, &statusCode
	}

	if user.Password == "" {
		message = "please choose a password to login"
		statusCode = http.StatusBadRequest
		return errors.New(message), nil, &statusCode
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginOrganizationRequest.Password))
	if err != nil {
		message = "invalid password"
		statusCode = http.StatusBadRequest
		return errors.New(message), nil, &statusCode
	}

	token, refreshToken, jwtExpiration, err := utils.GenerateJWT(user.ID, user.OrganizationID, user.Email)
	if err != nil {
		message = err.Error()
		statusCode = http.StatusInternalServerError
		return err, nil, &statusCode
	}

	statusCode = http.StatusOK
	return nil,
		&types.LoginOrganizationResponseData{
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
		&statusCode
}

func (s *OrganizationService) ResendEmailVerificationOtp(ctx *gin.Context) (error, *string, *int) {
	var statusCode = http.StatusInternalServerError
	var message string = "email verification otp sent successfully"
	var resendEmailVerificationOtpRequest types.ResendEmailVerificationOtpRequest = ctx.MustGet("request").(types.ResendEmailVerificationOtpRequest)

	err := s.db.DB.Transaction(func(tx *gorm.DB) error {
		user, err := s.userRepository.FindOneByFilter(repositories.UserFilter{
			Email: resendEmailVerificationOtpRequest.Email,
		}, tx)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = http.StatusBadRequest
			return errors.New("user not found")
		}

		if err != nil {
			return err
		}

		var secretValue string
		vault, err := s.vaultRepository.FindOneByFilter(repositories.VaultFilter{
			OwnerID:   user.ID,
			OwnerType: "user",
			Key:       string(entitities.UserVerificationSecret),
		}, tx)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			secretValue = ""
		} else if err != nil {
			return err
		} else {
			vaultObject, err := s.vaultEncrypt.Decrypt(vault.Object)
			if err != nil {
				return err
			}
			var obj map[string]string
			if err := json.Unmarshal([]byte(vaultObject), &obj); err != nil {
				return err
			}
			secretValue = obj["secret"]
		}

		if err := s.cacheVerificationCode(ctx, tx, &secretValue, user.OrganizationID, user.ID, user.Email); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		message = err.Error()
		return err, &message, &statusCode
	}

	statusCode = http.StatusOK
	return nil, &message, &statusCode
}

func (s *OrganizationService) cacheVerificationCode(ctx *gin.Context, tx *gorm.DB, secret *string, organisationId string, userId string, userEmail string) error {
	var secretValue string

	if secret == nil || *secret == "" {
		generated, _ := utils.GenerateRandomString(32)
		secretValue = generated

		if err := s.saveUserVerificationSecret(tx, secretValue, userId); err != nil {
			return err
		}

	} else {
		secretValue = *secret
	}

	otp, err := totp.GenerateCode(secretValue, time.Now())
	fmt.Printf("otp value: %s", otp)
	if err != nil {
		return err
	}

	cacheValue := map[string]string{
		"secret":          secretValue,
		"organization_id": organisationId,
		"user_id":         userId,
		"email":           userEmail,
	}
	cacheValueJSON, err := json.Marshal(cacheValue)
	if err != nil {
		return err
	}

	if err := s.redis.Set(ctx, fmt.Sprintf("email_verification_token:%s", otp), cacheValueJSON, 1*time.Hour).Err(); err != nil {
		return err
	}

	return nil
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
	err = s.vaultRepository.Create(vault, tx)
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
	_, err := s.organisationRepository.FindOneByFilter(repositories.OrganisationFilter{
		Email: email,
	}, tx)
	if err != nil && err != gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func (s *OrganizationService) assertUserEmailDoesNotExists(tx *gorm.DB, email string) bool {
	_, err := s.userRepository.FindOneByFilter(repositories.UserFilter{
		Email: email,
	}, tx)
	if err != nil && err != gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func (s *OrganizationService) assertDomainDoesNotExists(tx *gorm.DB, domain string) bool {
	_, err := s.organisationRepository.FindOneByFilter(repositories.OrganisationFilter{
		Domain: domain,
	}, tx)
	if err != nil && err != gorm.ErrRecordNotFound {
		return true
	}
	return false
}
