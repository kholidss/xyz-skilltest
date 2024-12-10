package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockRedisLock struct {
	mock.Mock
}

func (m *MockRedisLock) Obtain(ctx context.Context, key string, ttlMilliSecond int) error {
	args := m.Called(ctx, key, ttlMilliSecond)
	return args.Error(0)
}

func (m *MockRedisLock) ReleaseLock(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockRedisLock) TTLCheck(ctx context.Context) (time.Duration, error) {
	args := m.Called(ctx)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockRedisLock) Refresh(ctx context.Context, ttlMilliSecond int) error {
	args := m.Called(ctx, ttlMilliSecond)
	return args.Error(0)
}
