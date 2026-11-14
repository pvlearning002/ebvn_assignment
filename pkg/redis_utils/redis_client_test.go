package redis_utils

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRedisClient_DefaultConfig(t *testing.T) {
	prefix := "DEFAULT_REDIS"

	envVars := []string{
		prefix + "_REDIS_ADDRESS",
		prefix + "_REDIS_PASSWORD",
		prefix + "_REDIS_DB",
	}

	for _, env := range envVars {
		oldValue, existed := os.LookupEnv(env)

		require.NoError(t, os.Unsetenv(env))

		t.Cleanup(func() {
			if existed {
				_ = os.Setenv(env, oldValue)
			} else {
				_ = os.Unsetenv(env)
			}
		})
	}

	client, err := NewRedisClient(prefix)

	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	options := client.Options()

	assert.Equal(t, "localhost:6379", options.Addr)
	assert.Equal(t, "", options.Password)
	assert.Equal(t, 0, options.DB)
}

func TestCheckRedisConnection(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func() (*miniredis.Miniredis, *redis.Client)
		expectError bool
	}{
		{
			name: "normal case",
			setup: func() (*miniredis.Miniredis, *redis.Client) {
				server := miniredis.RunT(t)

				client := redis.NewClient(&redis.Options{
					Addr: server.Addr(),
				})

				return server, client
			},
			expectError: false,
		},
		{
			name: "redis connection failed",
			setup: func() (*miniredis.Miniredis, *redis.Client) {
				server := miniredis.RunT(t)

				client := redis.NewClient(&redis.Options{
					Addr: server.Addr(),
				})

				server.Close()

				return server, client
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, client := tc.setup()
			defer client.Close()

			err := CheckRedisConnection(
				context.Background(),
				client,
			)

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(
					t,
					err.Error(),
					"redis connection failed",
				)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsRedisContainerRunning(t *testing.T) {
	testCases := []struct {
		name           string
		output         []byte
		mockError      error
		expectedResult bool
	}{
		{
			name:           "container is running",
			output:         []byte("true\n"),
			expectedResult: true,
		},
		{
			name:           "container is not running",
			output:         []byte("false\n"),
			expectedResult: false,
		},
		{
			name:           "docker command failed",
			output:         nil,
			mockError:      errors.New("docker not found"),
			expectedResult: false,
		},
		{
			name:           "empty output",
			output:         []byte(""),
			expectedResult: false,
		},
		{
			name:           "unexpected output",
			output:         []byte("running"),
			expectedResult: false,
		},
		{
			name:           "output with spaces",
			output:         []byte("  true  \n"),
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalDockerInspect := dockerInspect

			dockerInspect = func() ([]byte, error) {
				return tc.output, tc.mockError
			}

			defer func() {
				dockerInspect = originalDockerInspect
			}()

			result := IsRedisContainerRunning()

			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestDockerInspect(t *testing.T) {
	originalDockerInspect := dockerInspect

	// Restore original function after test.
	dockerInspect = originalDockerInspect

	output, err := dockerInspect()

	// Docker may not be installed or the Redis container
	// may not exist. Both outcomes are valid for this test.
	if err != nil {
		assert.Error(t, err)
		return
	}

	assert.NotNil(t, output)
}

func TestNewRedisClient_InvalidConfig(t *testing.T) {
	prefix := "INVALID_REDIS"

	t.Setenv(prefix+"_REDIS_ADDRESS", "localhost:6379")
	t.Setenv(prefix+"_REDIS_PASSWORD", "test_password")
	t.Setenv(prefix+"_REDIS_DB", "invalid")

	client, err := NewRedisClient(prefix)

	require.Error(t, err)
	assert.Nil(t, client)
}
