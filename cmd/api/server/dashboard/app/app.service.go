package app

import (
	"sso-poc/internal/db"
	"sso-poc/internal/crypto"
	"github.com/redis/go-redis/v9"
)

type AppService struct {
	db *db.Database
	redis *redis.Client
	vaultEncrypt *crypto.TokenEncryption
}

func CreateAppService(db *db.Database, redis *redis.Client, vaultEncrypt *crypto.TokenEncryption) *AppService {
	return &AppService{
		db: db,
		redis: redis,
		vaultEncrypt: vaultEncrypt,
	}
}