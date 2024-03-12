package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
	"github.com/mahmudindes/orenocomic-bagicore/internal/utila"
)

func (rdb Redis) Set(ctx context.Context, key string, val any, exp time.Duration) error {
	return rdb.client.Set(ctx, key, val, exp).Err()
}

func (rdb Redis) get(ctx context.Context, key string) (*redis.StringCmd, error) {
	cmd := rdb.client.Get(ctx, key)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil, model.NotFoundError(cmd.Err())
		}
		return nil, cmd.Err()
	}
	return cmd, nil
}

func (rdb Redis) GetBytes(ctx context.Context, key string) ([]byte, error) {
	cmd, err := rdb.get(ctx, key)
	if err != nil {
		return nil, err
	}
	return cmd.Bytes()
}

func (rdb Redis) GetInt(ctx context.Context, key string) (int, error) {
	cmd, err := rdb.get(ctx, key)
	if err != nil {
		return 0, err
	}
	return cmd.Int()
}

func (rdb Redis) GetUint(ctx context.Context, key string) (uint, error) {
	cmd, err := rdb.get(ctx, key)
	if err != nil {
		return 0, err
	}
	return utila.Atou(cmd.Val())
}

func (rdb Redis) Delete(ctx context.Context, key string) error {
	return rdb.client.Del(ctx, key).Err()
}

func (rdb Redis) Expire(ctx context.Context, key string, exp time.Duration) error {
	return rdb.client.Expire(ctx, key, exp).Err()
}

func (rdb Redis) Increment(ctx context.Context, key string, val int64) error {
	if val == 1 {
		return rdb.client.Incr(ctx, key).Err()
	}
	return rdb.client.IncrBy(ctx, key, val).Err()
}

func (rdb Redis) Decrement(ctx context.Context, key string, val int64) error {
	if val == 1 {
		return rdb.client.Decr(ctx, key).Err()
	}
	return rdb.client.DecrBy(ctx, key, val).Err()
}
