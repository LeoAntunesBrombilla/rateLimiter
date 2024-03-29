package redisRepository

import (
	"github.com/go-redis/redis/v8"
	"os"
)

func Config() *redis.Client {
	redisAddr, exists := os.LookupEnv("REDIS_ADDR")
	if !exists {
		redisAddr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return client
}
