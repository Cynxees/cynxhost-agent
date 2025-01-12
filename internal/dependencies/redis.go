package dependencies

import (
	"strconv"

	"github.com/go-redis/redis"
)

func NewRedisClient(config *Config) *redis.Client {

	addr := config.Database.Redis.Host + ":" + strconv.Itoa(config.Database.Redis.Port)

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}
