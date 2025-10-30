package app

import (
	"sso-poc/internal/db"
	"sso-poc/internal/crypto"
	"github.com/redis/go-redis/v9"
	"github.com/gin-gonic/gin"
	appTypes "sso-poc/cmd/api/server/dashboard/app/types"
	"sso-poc/internal/db/repositories"
)

type AppService struct {
	db *db.Database
	redis *redis.Client
	vaultEncrypt *crypto.TokenEncryption
	appRepository *repositories.AppRepository
	appIdentityProviderRepository *repositories.AppIdentityProviderRepository
}

func CreateAppService(db *db.Database, redis *redis.Client, vaultEncrypt *crypto.TokenEncryption) *AppService {
	return &AppService{
		db: db,
		redis: redis,
		vaultEncrypt: vaultEncrypt,
		appRepository: repositories.CreateAppRepository(db.DB),
		appIdentityProviderRepository: repositories.CreateAppIdentityProviderRepository(db.DB),
	}
}

func (s *AppService) CreateApp(ctx *gin.Context, request appTypes.CreateAppRequest) (string, error) {
	
}