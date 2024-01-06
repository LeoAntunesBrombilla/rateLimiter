package redisRepository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) Set(key string, value interface{}) error {
	ctx := context.Background()

	expiration := time.Duration(0)
	fmt.Println("Saving")
	fmt.Println(key)
	fmt.Println(value)
	fmt.Println(expiration)
	fmt.Println("DONE")
	err := r.client.Set(ctx, key, value, expiration).Err()

	if err != nil {
		log.Printf("Error setting key in Redis: %v", err)
		return err
	}

	return nil
}

func (r *RedisRepository) Get(key string) (interface{}, error) {
	ctx := context.Background()

	result, err := r.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}
