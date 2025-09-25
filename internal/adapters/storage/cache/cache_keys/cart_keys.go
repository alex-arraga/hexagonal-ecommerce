package cachekeys

func Cart(id string) string {
	return generateCacheKey("cart:%s", id)
}
