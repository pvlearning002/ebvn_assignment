package utils

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRedisClientUtilsTest(t *testing.T, server *miniredis.Miniredis, cfg config.Config, cfgError config.ConfigError) *redisClientUtils {
	t.Helper()
	utils := NewRedisClientUtils(
		cfg,
		config.ConfigRedis{Address: server.Addr()},
		cfgError,
	).(*redisClientUtils)
	t.Cleanup(func() { _ = utils.client.Close() })
	return utils
}

func TestNewRedisClientUtils(t *testing.T) {
	server := miniredis.RunT(t)
	cfgRedis := config.ConfigRedis{Address: server.Addr(), Password: "secret", DB: 2}
	utils := NewRedisClientUtils(config.Config{}, cfgRedis, config.ConfigError{}).(*redisClientUtils)
	defer utils.client.Close()

	require.NotNil(t, utils.client)
	assert.Equal(t, cfgRedis.Address, utils.client.Options().Addr)
	assert.Equal(t, cfgRedis.Password, utils.client.Options().Password)
	assert.Equal(t, cfgRedis.DB, utils.client.Options().DB)
}

func TestNewRedisClient(t *testing.T) {
	cfg := config.ConfigRedis{Address: "localhost:6379", Password: "password", DB: 1}
	utils := NewRedisClientUtils(config.Config{}, config.ConfigRedis{}, config.ConfigError{})

	client, err := utils.NewRedisClient(cfg)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()
	assert.Equal(t, cfg.Address, client.Options().Addr)
	assert.Equal(t, cfg.Password, client.Options().Password)
	assert.Equal(t, cfg.DB, client.Options().DB)
}

func TestIsRedisConnect(t *testing.T) {
	server := miniredis.RunT(t)
	utils := newRedisClientUtilsTest(t, server, config.Config{}, config.ConfigError{})

	connected, err := utils.IsRedisConnect(context.Background())
	require.NoError(t, err)
	assert.True(t, connected)
}

func TestIsRedisConnect_Failure(t *testing.T) {
	server := miniredis.RunT(t)
	utils := newRedisClientUtilsTest(t, server, config.Config{}, config.ConfigError{})
	server.Close()

	connected, err := utils.IsRedisConnect(context.Background())
	assert.False(t, connected)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "redis connection failed")
}

func TestStoreStringValueAndGetStringValue(t *testing.T) {
	server := miniredis.RunT(t)
	utils := newRedisClientUtilsTest(t, server, config.Config{RedisExpiredTime: 60}, config.ConfigError{})
	ctx := context.Background()

	require.NoError(t, utils.StoreStringValue(ctx, "key", "value"))
	got, err := utils.GetStringValue(ctx, "key")
	require.NoError(t, err)
	assert.Equal(t, "value", got)
	assert.Equal(t, time.Minute, server.TTL("key"))
}

func TestStoreStringValue_ConnectionFailure(t *testing.T) {
	server := miniredis.RunT(t)
	cfgError := config.ConfigError{ErrorRedisConnection: "redis is unavailable"}
	utils := newRedisClientUtilsTest(t, server, config.Config{}, cfgError)
	server.Close()

	err := utils.StoreStringValue(context.Background(), "key", "value")
	require.Error(t, err)
	assert.EqualError(t, err, cfgError.ErrorRedisConnection)
}

func TestGetStringValue_KeyNotFound(t *testing.T) {
	server := miniredis.RunT(t)
	utils := newRedisClientUtilsTest(t, server, config.Config{}, config.ConfigError{})

	got, err := utils.GetStringValue(context.Background(), "missing")
	assert.Empty(t, got)
	assert.ErrorIs(t, err, redis.Nil)
}

func TestGetStringValue_ConnectionFailure(t *testing.T) {
	server := miniredis.RunT(t)
	cfg := config.Config{StringEmpty: "empty"}
	cfgError := config.ConfigError{ErrorRedisConnection: "redis is unavailable"}
	utils := newRedisClientUtilsTest(t, server, cfg, cfgError)
	server.Close()

	got, err := utils.GetStringValue(context.Background(), "key")
	assert.Equal(t, cfg.StringEmpty, got)
	require.Error(t, err)
	assert.EqualError(t, err, cfgError.ErrorRedisConnection)
}

func TestIsRedisContainerRunning(t *testing.T) {
	tests := []struct {
		name       string
		output     []byte
		inspectErr error
		want       bool
	}{
		{name: "running output with whitespace", output: []byte(" \ntrue\t"), want: true},
		{name: "container stopped", output: []byte("false"), want: false},
		{name: "unexpected output", output: []byte("running"), want: false},
		{name: "inspect command failed", inspectErr: errors.New("docker unavailable"), want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			utils := &redisClientUtils{
				dockerInspectFn: func() ([]byte, error) {
					return tc.output, tc.inspectErr
				},
			}
			assert.Equal(t, tc.want, utils.IsRedisContainerRunning())
		})
	}
}
