package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"
)

func (rdb Redis) GobSet(ctx context.Context, key string, v any, exp time.Duration) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return err
	}
	return rdb.Set(ctx, key, buf, exp)
}

func (rdb Redis) GobGet(ctx context.Context, key string, v any) error {
	buf, err := rdb.GetBytes(ctx, key)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf)).Decode(v)
}
