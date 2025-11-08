package ports

import (
	"context"
	"time"
)

// CacheRepository is an interface for interacting with cache-related business logic
type CacheRepository interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error // Set stores the value in the cache
	Get(ctx context.Context, key string) ([]byte, error)                        // Get retrieves the value from the cache
	Delete(ctx context.Context, key string) error                               // Delete removes the value from the cache
	DeleteByPrefix(ctx context.Context, prefix string) error                    // DeleteByPrefix removes the value from the cache with the given prefix
}
