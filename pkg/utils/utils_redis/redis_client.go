package redis_utils

import (
	"context"
	"fmt"
	exec "os/exec"
	"strings"
	"time"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/redis/go-redis/v9"
)

type CustomRedisClient interface {
	IsRedisConnect(ctx context.Context, client *redis.Client) (bool, error)
}

type customRedisClient struct {
	cfg config.ConfigRedis
}

// NewRedisClient creates a new Redis client using the provided environment prefix for configuration.
func NewRedisClient(cfg config.ConfigRedis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return rdb, nil
}

func (self *customRedisClient) dockerInspect() ([]byte, error) {
	cmd := exec.Command(
		"docker",
		"inspect",
		"-f",
		"{{.State.Running}}",
		"redis",
	)
	return cmd.Output()

}

func (self *customRedisClient) IsRedisConnect(
	ctx context.Context,
	client *redis.Client,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return false, fmt.Errorf("redis connection failed: %w", err)
	}

	return true, nil
}

func (self *customRedisClient) IsRedisContainerRunning() bool {
	output, err := self.dockerInspect()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(output)) == "true"
}
