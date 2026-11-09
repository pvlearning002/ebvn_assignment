package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func InitMockRedis(t *testing.T) *redis.Client {
	// Create a mock Redis client using miniredis
	mockRedis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start mock Redis server: %v", err)
	}

	// Create a Redis client that connects to the mock Redis server
	rdb := redis.NewClient(&redis.Options{
		Addr:     mockRedis.Addr(),
		Password: "", // No password for the mock server
		DB:       0,  // Default DB
	})

	return rdb
}
