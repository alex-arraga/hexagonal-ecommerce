package cachekeys

func Order(id string) string {
	return generateCacheKey("order", id)
}

func AllOrders() string {
	return "orders:all"
}
