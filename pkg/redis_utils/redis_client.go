package redis_utils

import (
	"context"
	"fmt"
	exec "os/exec"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var dockerInspect = func() ([]byte, error) {
	cmd := exec.Command(
		"docker",
		"inspect",
		"-f",
		"{{.State.Running}}",
		"redis",
	)
	return cmd.Output()
}

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

func CheckRedisConnection(
	ctx context.Context,
	client *redis.Client,
) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	return nil
}

func IsRedisContainerRunning() bool {
	output, err := dockerInspect()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(output)) == "true"
}
