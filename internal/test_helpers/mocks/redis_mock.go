package mocks

import (
	"context"
	"errors"
	"go-ecommerce/internal/core/ports"
	"sync"
	"time"
)

type MockRedis struct {
	mu     sync.RWMutex
	store  map[string]string
	closed bool
}

func NewMockRedis() ports.CacheRepository {
	return &MockRedis{store: make(map[string]string)}
}

func (m *MockRedis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return errors.New("redis connection is closed")
	}

	m.store[key] = string(value)
	return nil
}

func (m *MockRedis) Get(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, errors.New("redis connection is closed")
	}

	v, ok := m.store[key]
	if !ok {
		// Simula que Redis devuelve nil cuando la key no existe.
		return nil, nil
	}
	return []byte(v), nil
}

func (m *MockRedis) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return errors.New("redis connection is closed")
	}

	delete(m.store, key)
	return nil
}

func (m *MockRedis) DeleteByPrefix(ctx context.Context, prefix string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return errors.New("redis connection is closed")
	}

	for k := range m.store {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			delete(m.store, k)
		}
	}
	return nil
}

func (m *MockRedis) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return errors.New("redis already closed")
	}

	m.closed = true
	m.store = make(map[string]string) // limpia todo al cerrar
	return nil
}
