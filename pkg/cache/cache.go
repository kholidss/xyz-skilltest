package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Closure func(bytes []byte) error

const (
	cacheNil string = `redis: nil`
)

type Cacher interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val interface{}, duration time.Duration) error
	Delete(ctx context.Context, key ...string) error
}

type cache struct {
	redis           *redis.Client
	retentionSecond time.Duration
}

func NewCache(redis *redis.Client) Cacher {
	return &cache{
		redis: redis,
	}
}

func (c *cache) Set(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	cmd := c.redis.Set(ctx, key, val, exp)
	return cmd.Err()
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	cmd := c.redis.Get(ctx, key)
	b, e := cmd.Bytes()

	if e != nil {
		if e.Error() == cacheNil {
			return b, nil
		}
	}

	return b, e
}

func (c *cache) Delete(ctx context.Context, key ...string) error {
	cmd := c.redis.Del(ctx, key...)
	return cmd.Err()
}
