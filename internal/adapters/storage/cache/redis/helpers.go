package redis

import "fmt"

// GenerateCacheKey generates a cache key based on the input parameters
func GenerateCacheKey(prefix string, params any) string {
	return fmt.Sprintf("%s:%v", prefix, params)
}

// GenerateCacheParams generates a cache params based on the input parameters
func GenerateCacheKeyParams(params ...any) string {
	var str string

	for i, param := range params {
		str += fmt.Sprintf("%v", param)

		last := len(params) - 1
		if i != last {
			str += "-"
		}
	}

	return str
}
