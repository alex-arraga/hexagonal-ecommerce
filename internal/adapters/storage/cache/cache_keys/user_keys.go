package cachekeys

func User(id string) string {
	return generateCacheKey("user", id)
}

func UserByEmail(email string) string {
	return generateCacheKey("user", email)
}

func AllUsers() string {
	return "users:all"
}
