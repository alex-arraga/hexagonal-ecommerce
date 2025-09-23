package cachekeys

import "fmt"

// GenerateCacheKey generates a cache key based on the input parameters
func generateCacheKey(prefix string, params any) string {
	return fmt.Sprintf("%s:%v", prefix, params)
}
