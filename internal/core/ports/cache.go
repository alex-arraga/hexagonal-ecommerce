package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"
	"time"
)

// CacheRepository is an interface for interacting with cache-related business logic
type CacheRepository interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error // Set stores the value in the cache
	Get(ctx context.Context, key string) ([]byte, error)                        // Get retrieves the value from the cache
	Delete(ctx context.Context, key string) error                               // Delete removes the value from the cache
	DeleteByPrefix(ctx context.Context, prefix string) error                    // DeleteByPrefix removes the value from the cache with the given prefix
	Close() error                                                               // Close closes the connection to the cache server
}

// UserCacheRepository is an interface for interacting with cache-related user business logic
type UserCacheRepository interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	SetUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
	DeleteAllUsers(ctx context.Context) error
}
