package cachekeys

func Order(id string) string {
	return generateCacheKey("order:%s", id)
}

func AllOrders() string {
	return "orders:all"
}
