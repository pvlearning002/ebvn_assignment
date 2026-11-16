package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

type URLStorage interface {
	StoreURL(ctx context.Context, key, url string) error
	GetURL(ctx context.Context, key string) (string, error)
}

type urlStorage struct {
	client *redis.Client
	cfg    *config.Config
}

func NewURLStorage(client *redis.Client, cfg *config.Config) URLStorage {
	return &urlStorage{
		client: client,
		cfg:    cfg,
	}
}

func (self *urlStorage) StoreURL(ctx context.Context, key, url string) error {
	expTime := time.Duration(self.cfg.RedisExpiredTime) * time.Second
	return self.client.Set(ctx, key, url, expTime).Err()
}

func (self *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	return self.client.Get(ctx, key).Result()
}
