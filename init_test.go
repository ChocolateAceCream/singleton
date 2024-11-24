package singleton

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewSingleton_LocalRedis(t *testing.T) {
	// Redis client options for your local Redis server
	redisOptions := &redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
	}

	// Create the singleton instance with the Redis client
	s, err := NewSingleton(WithRedisClient(redisOptions))
	if err != nil {
		t.Fatalf("failed to create singleton: %v", err)
	}

	// Check that the singleton and Redis client are initialized
	assert.NotNil(t, s, "expected singleton instance to be non-nil")
	assert.NotNil(t, s.Redis, "expected Redis client to be initialized")

	// Test if Redis client is functional
	ctx := context.Background()
	err = s.Redis.Set(ctx, "test_key", "test_value", time.Minute).Err()
	if err != nil {
		t.Fatalf("failed to set key in Redis: %v", err)
	}
	log.Println("Successfully set test_key in Redis")

	val, err := s.Redis.Get(ctx, "test_key").Result()
	if err != nil {
		t.Fatalf("failed to get key from Redis: %v", err)
	}
	assert.Equal(t, "test_value", val, "expected Redis value to match")

	// Clean up
	err = s.Redis.Del(ctx, "test_key").Err()
	if err != nil {
		t.Fatalf("failed to delete key from Redis: %v", err)
	}
	log.Println("Successfully deleted test_key from Redis")
}
