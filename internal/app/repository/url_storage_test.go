package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type redisClientUtilsMock struct {
	storeStringValue func(context.Context, string, string) error
	getStringValue   func(context.Context, string) (string, error)
}

func (m *redisClientUtilsMock) NewRedisClient(config.ConfigRedis) (*redis.Client, error) {
	return nil, nil
}
func (m *redisClientUtilsMock) DockerInspect() ([]byte, error) { return nil, nil }
func (m *redisClientUtilsMock) IsRedisConnect(context.Context) (bool, error) {
	return false, nil
}
func (m *redisClientUtilsMock) IsRedisContainerRunning() bool { return false }
func (m *redisClientUtilsMock) StoreStringValue(ctx context.Context, key, value string) error {
	return m.storeStringValue(ctx, key, value)
}
func (m *redisClientUtilsMock) GetStringValue(ctx context.Context, key string) (string, error) {
	return m.getStringValue(ctx, key)
}

var _ utils.RedisClientUtils = (*redisClientUtilsMock)(nil)

func TestNewURLStorage(t *testing.T) {
	repo := NewURLStorage(&redisClientUtilsMock{})
	require.NotNil(t, repo)
	assert.IsType(t, &urlStorage{}, repo)
}

func TestURLStorage_StoreURL(t *testing.T) {
	ctx := context.Background()
	const key = "test-key"
	const url = "https://example.com/path"
	var called bool
	mockRedis := &redisClientUtilsMock{
		storeStringValue: func(gotCtx context.Context, gotKey, gotURL string) error {
			called = true
			assert.Equal(t, ctx, gotCtx)
			assert.Equal(t, key, gotKey)
			assert.Equal(t, url, gotURL)
			return nil
		},
	}

	err := NewURLStorage(mockRedis).StoreURL(ctx, key, url)
	require.NoError(t, err)
	assert.True(t, called, "StoreStringValue should be called")
}

func TestURLStorage_StoreURL_ReturnsRedisError(t *testing.T) {
	wantErr := errors.New("redis unavailable")
	mockRedis := &redisClientUtilsMock{
		storeStringValue: func(context.Context, string, string) error { return wantErr },
	}

	err := NewURLStorage(mockRedis).StoreURL(context.Background(), "key", "https://example.com")
	assert.ErrorIs(t, err, wantErr)
}

func TestURLStorage_GetURL(t *testing.T) {
	ctx := context.Background()
	const key = "test-key"
	const wantURL = "https://example.com/path"
	var called bool
	mockRedis := &redisClientUtilsMock{
		getStringValue: func(gotCtx context.Context, gotKey string) (string, error) {
			called = true
			assert.Equal(t, ctx, gotCtx)
			assert.Equal(t, key, gotKey)
			return wantURL, nil
		},
	}

	got, err := NewURLStorage(mockRedis).GetURL(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, wantURL, got)
	assert.True(t, called, "GetStringValue should be called")
}

func TestURLStorage_GetURL_ReturnsRedisError(t *testing.T) {
	wantErr := errors.New("key not found")
	mockRedis := &redisClientUtilsMock{
		getStringValue: func(context.Context, string) (string, error) { return "", wantErr },
	}

	got, err := NewURLStorage(mockRedis).GetURL(context.Background(), "missing")
	assert.Empty(t, got)
	assert.ErrorIs(t, err, wantErr)
}
