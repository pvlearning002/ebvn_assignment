package repository

import (
	"context"

	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
)

type URLStorage interface {
	StoreURL(ctx context.Context, key, url string) error
	GetURL(ctx context.Context, key string) (string, error)
}

type urlStorage struct {
	redisUtils utils.RedisClientUtils
}

func NewURLStorage(redisUtils utils.RedisClientUtils) URLStorage {
	return &urlStorage{
		redisUtils: redisUtils,
	}
}

func (self *urlStorage) StoreURL(ctx context.Context, key, url string) error {
	return self.redisUtils.StoreStringValue(ctx, key, url)
}

func (self *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	return self.redisUtils.GetStringValue(ctx, key)
}
