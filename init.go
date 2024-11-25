package singleton

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Singleton struct {
	Redis  *redis.Client
	Viper  *viper.Viper
	PGPool *pgxpool.Pool
}

type Option func(*Singleton) error

func (s *Singleton) AddPlugin(opts ...Option) (err error) {
	for _, opt := range opts {
		err = opt(s)
		if err != nil {
			return
		}
	}
	return
}
