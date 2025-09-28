package cachekeys

func Product(id string) string {
	return generateCacheKey("product", id)
}

func AllProducts() string {
	return "products:all"
}
