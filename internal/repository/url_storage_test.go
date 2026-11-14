package repository

import (
	"context"
	"testing"

	config "github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/redis_utils"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestURLStorage_ShortenURL(t *testing.T) {
	t.Parallel()
	key_test := "testKey"
	url_test := "https://www.google.com/"
	testCases := []struct {
		name           string
		setupMockRedis func(t *testing.T) *goredis.Client
		expectedError  error
		verifyFunc     func(ctx context.Context, r *goredis.Client)
	}{
		{
			name: "normal case",
			setupMockRedis: func(t *testing.T) *goredis.Client {
				return redis_utils.InitMockRedis(t)
			},
			expectedError: nil,
			verifyFunc: func(ctx context.Context, r *goredis.Client) {
				result, err := r.Get(ctx, key_test).Result()
				assert.NoError(t, err)
				assert.Equal(t, url_test, result)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRedis := tc.setupMockRedis(t)
			defer mockRedis.Close()
			testRepo := NewURLStorage(mockRedis, &config.Config{RedisExpiredTime: 60})
			err := testRepo.StoreURL(ctx, key_test, url_test)
			assert.ErrorIs(t, tc.expectedError, err)
			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mockRedis)
			}
		})
	}
}
