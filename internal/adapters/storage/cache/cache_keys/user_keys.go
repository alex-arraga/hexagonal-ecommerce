package cachekeys

func User(id string) string {
	return generateCacheKey("user:%s", id)
}

func AllUsers() string {
	return "users:all"
}
