package singleton

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGSQLOptions struct {
	Source          string
	MaxConns        int32
	MaxConnIdleTime time.Duration
}

func WithPGSQL(options PGSQLOptions) Option {
	return func(s *Singleton) (err error) {
		config, err := pgxpool.ParseConfig(options.Source)
		if err != nil {
			return
		}
		config.MaxConns = options.MaxConns
		config.MaxConnIdleTime = options.MaxConnIdleTime
		var pgPool *pgxpool.Pool
		for retries := 0; retries < 5; retries++ {
			pgPool, err = pgxpool.NewWithConfig(context.TODO(), config)
			if err == nil {
				s.PGPool = pgPool
				return
			}
			time.Sleep(time.Duration(retries) * time.Second)
		}
		return
	}
}
