package cache

import (
	"../lib/cache"
)

// Singleton
var instantiated *cache.Cache = nil

func CachePtr() *cache.Cache {
	return instantiated
}

func Init(t int64) {
	instantiated = &cache.Cache{
		Cachetime: t,
		Name: "node",
	}
}