package cache

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func CreateRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
		DB:   0,
	})
}
