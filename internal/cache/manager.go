package cache

import "time"

var GlobalCache *Cache

func InitGlobalCache(ttl time.Duration) {
	GlobalCache = New(ttl)
}
