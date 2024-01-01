package redisRepository

import "github.com/go-redis/redis/v8"

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) Set(key string, value interface{}) error {
	panic("Implement me!")
}

func (r *RedisRepository) Get(key string) (error, value interface{}) {
	panic("Implement me!")
}
