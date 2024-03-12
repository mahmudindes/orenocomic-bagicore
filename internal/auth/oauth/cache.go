package oauth

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	bagicore "github.com/mahmudindes/orenocomic-bagicore"
	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

type cTokenCacheStore interface {
	GetToken(ctx context.Context, id string) (*accessToken, error)
	SetToken(ctx context.Context, id string, token *accessToken) error
}

func newCTokenCache(rdb Redis) cTokenCacheStore {
	if rdb != nil && !reflect.ValueOf(rdb).IsNil() {
		return cTokenCacheRedis{rdb}
	}
	memory := &cTokenCacheMemory{tokens: make(map[string]accessToken)}
	go memory.BackgroundPurger()
	return memory
}

type cTokenCacheMemory struct {
	tokens map[string]accessToken
	mu     sync.Mutex
}

const cTokenCachePurge = 1 * time.Hour

func (ctcm *cTokenCacheMemory) GetToken(ctx context.Context, id string) (*accessToken, error) {
	ctcm.mu.Lock()
	defer ctcm.mu.Unlock()

	val, ok := ctcm.tokens[id]
	if !ok || time.Until(val.Expiration)/2 <= 0 {
		delete(ctcm.tokens, id)
		return nil, nil
	}
	return &val, nil
}

func (ctcm *cTokenCacheMemory) SetToken(ctx context.Context, id string, token *accessToken) error {
	ctcm.mu.Lock()
	defer ctcm.mu.Unlock()

	if token != nil {
		ctcm.tokens[id] = *token
	}
	return nil
}

func (ctcm *cTokenCacheMemory) BackgroundPurger() {
	for {
		for key, val := range ctcm.tokens {
			if time.Until(val.Expiration)/2 <= 0 {
				ctcm.mu.Lock()
				delete(ctcm.tokens, key)
				ctcm.mu.Unlock()
			}
		}
		time.Sleep(cTokenCachePurge)
	}
}

type cTokenCacheRedis struct {
	rdb Redis
}

var cTokenCacheRedisPrefix = bagicore.ID + ":oauth:token:"

func (ctcr cTokenCacheRedis) GetToken(ctx context.Context, id string) (*accessToken, error) {
	var token *accessToken
	if err := ctcr.rdb.GobGet(ctx, cTokenCacheRedisPrefix+id, token); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			return nil, nil
		}
		return nil, model.CacheError(err)
	}
	return token, nil
}

func (ctcr cTokenCacheRedis) SetToken(ctx context.Context, id string, token *accessToken) error {
	expiry := time.Until(token.Expiration) / 2
	if expiry <= 0 {
		return nil
	}
	if err := ctcr.rdb.GobSet(ctx, cTokenCacheRedisPrefix+id, token, expiry); err != nil {
		return model.CacheError(err)
	}
	return nil
}
