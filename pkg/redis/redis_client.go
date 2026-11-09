package redis

import (
	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client using the provided environment prefix for configuration.
func NewRedisClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return rdb, nil
}
