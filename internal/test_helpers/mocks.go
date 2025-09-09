package testhelpers

import (
	"context"
	"go-ecommerce/internal/core/ports"
	"time"
)

type MockRedis struct {
	store map[string]string
}

func NewMockRedis() ports.CacheRepository {
	return &MockRedis{store: make(map[string]string)}
}

func (m *MockRedis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	m.store[key] = string(value)
	return nil
}

func (m *MockRedis) Get(ctx context.Context, key string) ([]byte, error) {
	v, ok := m.store[key]
	if !ok {
		return nil, nil
	}
	return []byte(v), nil
}

func (m *MockRedis) Delete(ctx context.Context, key string) error {
	return nil
}

func (m *MockRedis) DeleteByPrefix(ctx context.Context, prefix string) error {
	return nil
}

func (m *MockRedis) Close() error {
	return nil
}
