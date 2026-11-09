package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRedisClient(t *testing.T) {
	t.Run("success with custom config", func(t *testing.T) {
		prefix := "TEST"

		t.Setenv("TEST_REDIS_ADDRESS", "127.0.0.1:6380")
		t.Setenv("TEST_REDIS_PASSWORD", "test_password")
		t.Setenv("TEST_REDIS_DB", "2")

		client, err := NewRedisClient(prefix)

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		options := client.Options()

		assert.Equal(t, "127.0.0.1:6380", options.Addr)
		assert.Equal(t, "test_password", options.Password)
		assert.Equal(t, 2, options.DB)
	})

	t.Run("success with default config", func(t *testing.T) {
		client, err := NewRedisClient("DEFAULT_TEST")

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		options := client.Options()

		assert.Equal(t, "localhost:6379", options.Addr)
		assert.Equal(t, "", options.Password)
		assert.Equal(t, 0, options.DB)
	})

	t.Run("invalid DB config", func(t *testing.T) {
		prefix := "INVALID"

		t.Setenv("INVALID_REDIS_ADDRESS", "localhost:6379")
		t.Setenv("INVALID_REDIS_PASSWORD", "password")
		t.Setenv("INVALID_REDIS_DB", "invalid")

		client, err := NewRedisClient(prefix)

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}
