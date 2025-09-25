package cachettl

import "time"

const (
	User     = 40 * time.Minute
	Product  = 30 * time.Minute
	Category = 10 * time.Minute
	Order    = 20 * time.Minute
	Cart     = 0 // always is in cache
)
