package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type (
	Redis struct {
		client redis.UniversalClient
	}

	Config struct {
		Enable bool   `conf:"enable"`
		URL    string `conf:"url"`
	}
)

func New(ctx context.Context, cfg Config) (*Redis, error) {
	if !cfg.Enable {
		return nil, nil
	}

	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{client: client}, nil
}
