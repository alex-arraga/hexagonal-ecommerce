package cachekeys

func User(id string) string {
	return generateCacheKey("user", id)
}

func AllUsers() string {
	return "users:all"
}
