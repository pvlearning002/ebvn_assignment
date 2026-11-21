package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
	redis_utils "github.com/pvlearning002/ebvn_assignment/pkg/utils/utils_redis"
)

type URLStorage interface {
	StoreURL(ctx context.Context, key, url string) error
	GetURL(ctx context.Context, key string) (string, error)
}

type urlStorage struct {
	client      *redis.Client
	cfg         *config.Config
	utils_redis redis_utils.CustomRedisClient
}

func NewURLStorage(client *redis.Client, cfg *config.Config, utils_redis redis_utils.CustomRedisClient) URLStorage {
	return &urlStorage{
		client:      client,
		cfg:         cfg,
		utils_redis: utils_redis,
	}
}

func (self *urlStorage) StoreURL(ctx context.Context, key, url string) error {

	expTime := time.Duration(self.cfg.RedisExpiredTime) * time.Second
	return self.client.Set(ctx, key, url, expTime).Err()
}

func (self *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	return self.client.Get(ctx, key).Result()
}
