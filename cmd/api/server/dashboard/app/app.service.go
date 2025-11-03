package app

import (
	"net/http"
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
}

func CreateAppService(db *db.Database, redis *redis.Client, vaultEncrypt *crypto.TokenEncryption) *AppService {
	return &AppService{
		db:                            db,
		redis:                         redis,
		vaultEncrypt:                  vaultEncrypt,
		appRepository:                 repositories.CreateAppRepository(db.DB),
		appIdentityProviderRepository: repositories.CreateAppIdentityProviderRepository(db.DB),
		identityProviderRepository:    repositories.CreateIdentityProviderRepository(db.DB),
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

func (s *AppService) GetApp(ctx *gin.Context) (*entitities.App, error, *int) {
	var app *entitities.App
	var statusCode int = http.StatusInternalServerError

	var appId string = ctx.Param("app_id")

	app, err := s.appRepository.FindOneByFilter(repositories.AppFilter{ID: appId}, nil)
	if err != nil {
		return nil, err, &statusCode
	}

	statusCode = http.StatusOK
	return app, nil, &statusCode
}
