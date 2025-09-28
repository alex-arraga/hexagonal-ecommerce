package cachekeys

func Category(id uint64) string {
	return generateCacheKey("category", id)
}

func AllCategories() string {
	return "categories:all"
}
