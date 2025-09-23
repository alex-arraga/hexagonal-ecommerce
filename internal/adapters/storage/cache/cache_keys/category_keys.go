package cachekeys

func Product(id string) string {
	return generateCacheKey("product:%s", id)
}

func AllProducts() string {
	return "products:all"
}
