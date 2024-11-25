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
	s := &Singleton{}
	// Create the singleton instance with the Redis client
	err := s.AddPlugin(WithRedisClient(redisOptions))
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

type Config struct {
	Version string `mapstructure:"version" yaml:"version"`
	Level   string `mapstructure:"level" yaml:"level"`
}

func TestNewSingleton_WithViper(t *testing.T) {
	// Define the configuration file details
	options := ViperOptions{
		Path:     "./",      // Current directory
		FileName: "config",  // Name of the config file without extension
		FileType: "yaml",    // File type
		EnvName:  "release", // Environment section to read
		Target:   &Config{}, // Target to unmarshal into
	}

	s := &Singleton{}

	// Add the Viper plugin to the Singleton instance
	err := s.AddPlugin(WithViper(options))
	if err != nil {
		t.Fatalf("failed to create singleton with Viper: %v", err)
	}

	// Assert that the Singleton instance and Viper client are initialized
	assert.NotNil(t, s, "expected singleton instance to be non-nil")
	assert.NotNil(t, s.Viper, "expected Viper instance to be initialized")

	// Assert that the configuration was read and unmarshaled correctly
	config := options.Target.(*Config) // Cast the target back to the Config type
	assert.Equal(t, "1.0.0", config.Version, "expected Release.Version to match")
	assert.Equal(t, "release", config.Level, "expected Debug.Level to match")

	log.Println("Successfully read configuration using Viper:")
	log.Printf("Release Version: %s", config.Version)
	log.Printf("Debug Level: %s", config.Level)
}

func TestNewSingleton_WithPGSQL(t *testing.T) {
	// Define the PostgreSQL options
	options := PGSQLOptions{
		Source:          "postgresql://nuodi:123qwe@localhost:5555/iot_backend?sslmode=disable",
		MaxConns:        10,
		MaxConnIdleTime: time.Hour,
	}

	s := &Singleton{}

	// Add the PostgreSQL plugin to the Singleton instance
	err := s.AddPlugin(WithPGSQL(options))
	if err != nil {
		t.Fatalf("failed to create singleton with PGSQL: %v", err)
	}

	// Assert that the Singleton instance and PostgreSQL connection pool are initialized
	assert.NotNil(t, s, "expected singleton instance to be non-nil")
	assert.NotNil(t, s.PGPool, "expected PostgreSQL pool to be initialized")

	// Check the PostgreSQL connection pool configuration
	config := s.PGPool.Config()
	assert.Equal(t, options.MaxConns, config.MaxConns, "expected MaxConns to match")
	assert.Equal(t, options.MaxConnIdleTime, config.MaxConnIdleTime, "expected MaxConnIdleTime to match")

	// Test if the PostgreSQL connection pool is functional
	ctx := context.Background()
	conn, err := s.PGPool.Acquire(ctx)
	if err != nil {
		t.Fatalf("failed to acquire connection from PGPool: %v", err)
	}
	defer conn.Release()

	// Verify that the connection is valid
	err = conn.Conn().Ping(ctx)
	if err != nil {
		t.Fatalf("failed to ping PostgreSQL database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
}
