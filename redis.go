package singleton

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func WithRedisClient(options *redis.Options) Option {
	return func(s *Singleton) (err error) {
		var client *redis.Client
		for retries := 0; retries < 5; retries++ {
			client = redis.NewClient(options)
			if _, err := client.Ping(context.Background()).Result(); err == nil {
				s.Redis = client
				return nil
			}
			time.Sleep(time.Duration(retries) * time.Second)
		}
		err = fmt.Errorf("failed to connect to Redis after 5 retries")
		return
	}
}
