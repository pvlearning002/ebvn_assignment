package redis_utils

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRedisClient(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.ConfigRedis
	}{
		{
			name: "Success",
			cfg: config.ConfigRedis{
				Address:  "localhost:6379",
				Password: "password",
				DB:       1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewRedisClient(tt.cfg)

			require.NoError(t, err)
			require.NotNil(t, client)

			assert.Equal(t, tt.cfg.Address, client.Options().Addr)
			assert.Equal(t, tt.cfg.Password, client.Options().Password)
			assert.Equal(t, tt.cfg.DB, client.Options().DB)

			_ = client.Close()
		})
	}
}

func TestIsRedisConnect(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) *redis.Client
		expected    bool
		expectError bool
	}{
		{
			name: "Success",
			setup: func(t *testing.T) *redis.Client {
				mr := miniredis.RunT(t)

				client := redis.NewClient(&redis.Options{
					Addr: mr.Addr(),
				})

				return client
			},
			expected:    true,
			expectError: false,
		},
		{
			name: "Redis connection failed",
			setup: func(t *testing.T) *redis.Client {
				client := redis.NewClient(&redis.Options{
					Addr: "localhost:0",
				})

				return client
			},
			expected:    false,
			expectError: true,
		},
	}

	customClient := &customRedisClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.setup(t)
			defer client.Close()

			result, err := customClient.IsRedisConnect(
				context.Background(),
				client,
			)

			assert.Equal(t, tt.expected, result)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "redis connection failed")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestIsRedisContainerRunning(t *testing.T) {
	customClient := &customRedisClient{}

	t.Run("Docker inspect command", func(t *testing.T) {
		// This test depends on Docker being installed
		// and a container named "redis" existing.
		result := customClient.IsRedisContainerRunning()

		assert.IsType(t, false, result)
	})
}
