package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	appTypes "sso-poc/cmd/api/server/dashboard/app/types"
	"sso-poc/internal/crypto"
	"sso-poc/internal/db"
	"sso-poc/internal/db/entitities"
	"sso-poc/internal/db/repositories"
	"sso-poc/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppService struct {
	db                            *db.Database
	redis                         *redis.Client
	vaultEncrypt                  *crypto.TokenEncryption
	appRepository                 *repositories.AppRepository
	appIdentityProviderRepository *repositories.AppIdentityProviderRepository
	identityProviderRepository    *repositories.IdentityProviderRepository
	vaultRepository               *repositories.VaultRepository
}

func CreateAppService(db *db.Database, redis *redis.Client, vaultEncrypt *crypto.TokenEncryption) *AppService {
	return &AppService{
		db:                            db,
		redis:                         redis,
		vaultEncrypt:                  vaultEncrypt,
		appRepository:                 repositories.CreateAppRepository(db.DB),
		appIdentityProviderRepository: repositories.CreateAppIdentityProviderRepository(db.DB),
		identityProviderRepository:    repositories.CreateIdentityProviderRepository(db.DB),
		vaultRepository:               repositories.CreateVaultRepository(db.DB, vaultEncrypt),
	}
}

func (s *AppService) CreateApp(ctx *gin.Context) (*string, error, *int) {
	var app *entitities.App
	var statusCode int = http.StatusInternalServerError

	var createAppRequest appTypes.CreateAppRequest = ctx.MustGet("request").(appTypes.CreateAppRequest)
	var user *utils.CustomClaims = ctx.MustGet("user").(*utils.CustomClaims)
	var err error

	err = s.db.DB.Transaction(func(tx *gorm.DB) error {
		app, err = s.appRepository.Create(
			&createAppRequest,
			user.OrganizationID,
			tx,
			s.appIdentityProviderRepository,
			s.identityProviderRepository)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err, &statusCode
	}

	statusCode = http.StatusOK
	return &app.ID, nil, &statusCode
}

func (s *AppService) AddAppIdentityProvider(ctx *gin.Context) (*string, error, *int) {
	var app *entitities.App
	var statusCode int = http.StatusInternalServerError

	var user *utils.CustomClaims = ctx.MustGet("user").(*utils.CustomClaims)

	var appId string = ctx.Param("app_id")
	var request appTypes.AppIdentityProviderRequest = ctx.MustGet("request").(appTypes.AppIdentityProviderRequest)

	tx := s.db.DB.Begin()
	defer tx.Rollback()

	app, err := s.appRepository.FindOneByFilter(repositories.AppFilter{
		ID:             appId,
		OrganizationID: user.OrganizationID,
	}, tx)
	if err != nil {
		return nil, err, &statusCode
	}
	if app == nil {
		statusCode = http.StatusNotFound
		return nil, errors.New("app not found"), &statusCode
	}

	identityProvider, err := s.identityProviderRepository.FindOneByFilter(repositories.IdentityProviderFilter{
		ID: request.ID,
	}, tx)
	if err != nil {
		return nil, err, &statusCode
	}
	if identityProvider == nil {
		statusCode = http.StatusNotFound
		return nil, errors.New("identity provider not found"), &statusCode
	}

	// compare scope with identity scopes
	if !utils.Contains(identityProvider.Scopes, request.Scopes) {
		statusCode = http.StatusBadRequest
		return nil, errors.New("ensure all scopes are valid for this identity provider"), &statusCode
	}

	appIdentityProvider := &entitities.AppIdentityProvider{
		AppID:              app.ID,
		IdentityProviderID: identityProvider.ID,
		Status:             "active",
		Scopes:             request.Scopes,
	}

	err = s.appIdentityProviderRepository.Create(appIdentityProvider, tx)
	if err != nil {
		return nil, err, &statusCode
	}

	if len(request.ProviderCredentials) > 0 {
		err = s.saveCredentialsToVault(appIdentityProvider, tx, request.ProviderCredentials)
		if err != nil {
			return nil, err, &statusCode
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err, &statusCode
	}

	statusCode = http.StatusOK
	message := "app identity provider added successfully"
	return &message, nil, &statusCode
}

func (s *AppService) GetApp(ctx *gin.Context) (*entitities.App, error, *int) {
	var app *entitities.App
	var statusCode int = http.StatusInternalServerError

	var appId string = ctx.Param("app_id")

	app, err := s.appRepository.FindOneByFilter(repositories.AppFilter{ID: appId}, nil)
	if err != nil {
		return nil, err, &statusCode
	}

	// appIdentityProviders := app.IdentityProviders

	statusCode = http.StatusOK
	return app, nil, &statusCode
}

func (s *AppService) UpdateAppIdentityProvider(ctx *gin.Context) (*string, error, *int) {
	var app *entitities.App
	var statusCode int = http.StatusInternalServerError

	var appId string = ctx.Param("app_id")
	var user *utils.CustomClaims = ctx.MustGet("user").(*utils.CustomClaims)
	var request appTypes.AppIdentityProviderRequest = ctx.MustGet("request").(appTypes.AppIdentityProviderRequest)

	tx := s.db.DB.Begin()
	defer tx.Rollback()

	app, err := s.appRepository.FindOneByFilter(repositories.AppFilter{
		ID:             appId,
		OrganizationID: user.OrganizationID,
	}, tx)

	if err == gorm.ErrRecordNotFound {
		statusCode = http.StatusNotFound
		return nil, errors.New("Invalid app Id"), &statusCode
	}

	if err != nil {
		return nil, err, &statusCode
	}

	if app == nil {
		statusCode = http.StatusNotFound
		return nil, errors.New("app not found"), &statusCode
	}

	var appIdentityProvider *entitities.AppIdentityProvider
	appIdentityProvider, err = s.appIdentityProviderRepository.FindOneByFilter(repositories.AppIdentityProviderFilter{
		ID: request.ID,
	}, tx)
	if err == gorm.ErrRecordNotFound {
		statusCode = http.StatusNotFound
		return nil, errors.New("Invalid app identity provider Id"), &statusCode
	}
	if err != nil {
		return nil, err, &statusCode
	}

	if appIdentityProvider == nil {
		statusCode = http.StatusNotFound
		return nil, errors.New("app identity provider not found"), &statusCode
	}

	if len(request.ProviderCredentials) == 0 && len(request.Scopes) == 0 {
		statusCode = http.StatusBadRequest
		return nil, errors.New("at least one of provider_credentials or scopes is required"), &statusCode
	}

	if len(request.Scopes) > 0 {
		allscopes := append(appIdentityProvider.Scopes, request.Scopes...)
		appIdentityProvider.Scopes = utils.ConvertToSet(allscopes)

		err = tx.Save(appIdentityProvider).Error
		if err != nil {
			return nil, errors.New(utils.GenericErrorMessages()[500]), &statusCode
		}
	}

	if len(request.ProviderCredentials) > 0 {

		if appIdentityProvider.VaultId == "" {
			err = s.saveCredentialsToVault(appIdentityProvider, tx, request.ProviderCredentials)
			if err != nil {
				return nil, errors.New(utils.GenericErrorMessages()[500]), &statusCode
			}

			err = tx.Commit().Error
			if err != nil {
				return nil, errors.New(utils.GenericErrorMessages()[500]), &statusCode
			}

			statusCode = http.StatusOK
			message := "app identity provider credentials saved successfully"

			return &message, nil, &statusCode
		}

		vaultObjectstr, err := s.vaultRepository.GetDecryptedObject(appIdentityProvider.VaultId, tx)

		if err != nil {
			return nil, errors.New(utils.GenericErrorMessages()[500]), &statusCode
		}

		var vaultObject map[string]string
		if vaultObjectstr == nil || *vaultObjectstr == "" {
			vaultObject = make(map[string]string)
		} else {
			err = json.Unmarshal([]byte(*vaultObjectstr), &vaultObject)
			if err != nil {
				return nil, err, &statusCode
			}
		}

		providerCredentialFields := appIdentityProvider.IdentityProvider.CredentialFields
		for _, credential := range request.ProviderCredentials {
			if !slices.Contains(providerCredentialFields, credential.Key) {
				errMessage := fmt.Sprintf("%s is not a valid credential for %s", credential.Key, appIdentityProvider.IdentityProvider.Name)
				return nil, errors.New(errMessage), &statusCode
			}

			vaultObject[credential.Key] = credential.Value
		}

		encryptedObject, err := json.Marshal(vaultObject)
		if err != nil {
			return nil, errors.New(utils.GenericErrorMessages()[500]), &statusCode
		}

		err = s.vaultRepository.EncryptObjectAndSave(string(encryptedObject), appIdentityProvider.Vault, tx)
		if err != nil {
			return nil, errors.New(utils.GenericErrorMessages()[500]), &statusCode
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, errors.New(utils.GenericErrorMessages()[500]), &statusCode
	}

	statusCode = http.StatusOK
	message := "app identity provider updated successfully"
	return &message, nil, &statusCode
}

func (s *AppService) saveCredentialsToVault(appIdentityProvider *entitities.AppIdentityProvider, tx *gorm.DB, credentials []appTypes.AppIdentityProviderCredentials) error {
	vaultObject := make(map[string]string)
	for _, credential := range credentials {
		vaultObject[credential.Key] = credential.Value
	}

	plainObject, err := json.Marshal(vaultObject)
	if err != nil {
		return err
	}

	encryptedObject, err := s.vaultEncrypt.Encrypt(string(plainObject))
	if err != nil {
		return err
	}

	vault := &entitities.Vault{
		OwnerID:   appIdentityProvider.AppID,
		OwnerType: "app",
		Key:       "app_identity_provider_credentials",
		Object:    encryptedObject,
	}

	err = s.vaultRepository.Create(vault, tx)
	if err != nil {
		return err
	}

	appIdentityProvider.VaultId = vault.ID
	err = s.appIdentityProviderRepository.UpdateVaultId(appIdentityProvider.ID, vault.ID, tx)
	if err != nil {
		return err
	}

	return nil
}
