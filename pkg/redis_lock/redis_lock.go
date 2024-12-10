package redislock

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

var (
	ErrNotObtained = errors.New("lock not obtained")
)

type Locker interface {
	Obtain(ctx context.Context, key string, ttlMilliSecond int) error
	ReleaseLock(ctx context.Context)
	TTLCheck(ctx context.Context) (time.Duration, error)
	Refresh(ctx context.Context, ttlMilliSecond int) error
}

type lock struct {
	redisLocker *redislock.Client
	locker      *redislock.Lock
}

func NewLocker(redisClient *redis.Client) Locker {
	return &lock{
		redisLocker: redislock.New(redisClient),
	}
}

func (l *lock) Obtain(ctx context.Context, key string, ttlMilliSecond int) error {
	lck, err := l.redisLocker.Obtain(ctx, key, time.Duration(ttlMilliSecond)*time.Millisecond, nil)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			return ErrNotObtained
		}
		return err
	}

	l.locker = lck
	return nil
}

func (l *lock) ReleaseLock(ctx context.Context) {
	_ = l.locker.Release(ctx)
}

func (l *lock) TTLCheck(ctx context.Context) (time.Duration, error) {
	return l.locker.TTL(ctx)
}

func (l *lock) Refresh(ctx context.Context, ttlMilliSecond int) error {
	return l.locker.Refresh(ctx, time.Duration(ttlMilliSecond)*time.Millisecond, nil)
}
