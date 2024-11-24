package singleton

import (
	"github.com/redis/go-redis/v9"
)

type Singleton struct {
	Redis *redis.Client
}

type Option func(*Singleton) error

func NewSingleton(opts ...Option) (s *Singleton, err error) {
	s = &Singleton{}
	for _, opt := range opts {
		err = opt(s)
		if err != nil {
			return
		}
	}
	return
}
