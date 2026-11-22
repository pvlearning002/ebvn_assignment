package utils

import (
	"context"
	"errors"
	"fmt"
	exec "os/exec"
	"strings"
	"time"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClientUtils interface {
	NewRedisClient(cfg config.ConfigRedis) (*redis.Client, error)
	DockerInspect() ([]byte, error)
	IsRedisConnect(ctx context.Context) (bool, error)
	IsRedisContainerRunning() bool
	StoreStringValue(ctx context.Context, key, stringValue string) error
	GetStringValue(ctx context.Context, key string) (string, error)
}

type redisClientUtils struct {
	client          *redis.Client
	cfg             config.Config
	cfgRedis        config.ConfigRedis
	cfgError        config.ConfigError
	dockerInspectFn func() ([]byte, error)
}

func NewRedisClientUtils(
	cfg config.Config,
	cfgRedis config.ConfigRedis,
	cfgError config.ConfigError) RedisClientUtils {
	client := redis.NewClient(&redis.Options{
		Addr:     cfgRedis.Address,
		Password: cfgRedis.Password,
		DB:       cfgRedis.DB,
	})
	return &redisClientUtils{
		client:   client,
		cfg:      cfg,
		cfgRedis: cfgRedis,
		cfgError: cfgError,
	}
}

// NewRedisClient creates a new Redis client using the provided environment prefix for configuration.
func (self *redisClientUtils) NewRedisClient(cfg config.ConfigRedis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return rdb, nil
}

func (self *redisClientUtils) DockerInspect() ([]byte, error) {
	if self.dockerInspectFn != nil {
		return self.dockerInspectFn()
	}
	cmd := exec.Command(
		"docker",
		"inspect",
		"-f",
		"{{.State.Running}}",
		"redis",
	)
	return cmd.Output()

}

func (self *redisClientUtils) IsRedisConnect(
	ctx context.Context,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := self.client.Ping(ctx).Err(); err != nil {
		return false, fmt.Errorf("redis connection failed: %w", err)
	}

	return true, nil
}

func (self *redisClientUtils) IsRedisContainerRunning() bool {
	output, err := self.DockerInspect()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

func (self *redisClientUtils) StoreStringValue(ctx context.Context, key, stringValue string) error {
	if isConnected, err := self.IsRedisConnect(ctx); err == nil && isConnected {
		expTime := time.Duration(self.cfg.RedisExpiredTime) * time.Second
		return self.client.Set(ctx, key, stringValue, expTime).Err()
	}
	return errors.New(self.cfgError.ErrorRedisConnection)
}

func (self *redisClientUtils) GetStringValue(ctx context.Context, key string) (string, error) {
	if isConnected, err := self.IsRedisConnect(ctx); err == nil && isConnected {
		return self.client.Get(ctx, key).Result()
	}
	return self.cfg.StringEmpty, errors.New(self.cfgError.ErrorRedisConnection)

}
