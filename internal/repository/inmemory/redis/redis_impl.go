package redis

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisImpl struct {
	Redis *redis.Client
}

func New(redisClient *redis.Client) *RedisImpl {
	return &RedisImpl{
		Redis: redisClient,
	}
}

func (r *RedisImpl) GetRedis(key string) (string, error) {

	val, err := r.Redis.Get(key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisImpl) SetRedis(key string, value string, ttl time.Duration) error {
	
	err := r.Redis.Set(key, value, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}