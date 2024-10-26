package in_memory_cache

import (
	"sync"
)

type InMemoryCache[K, V comparable, D any] struct {
	cacheFactory CacheFactory[V, D]
	memory       sync.Map
}

func NewInMemoryCache[K, V comparable, D any](cacheFactory CacheFactory[V, D]) *InMemoryCache[K, V, D] {
	return &InMemoryCache[K, V, D]{
		cacheFactory: cacheFactory,
	}
}

func (c *InMemoryCache[K, V, D]) Get(key K, validateArg V) (res D, ok bool) {
	cacheAny, ok := c.memory.Load(key)
	if !ok {
		return
	}

	cache := cacheAny.(Cache[V, D])
	res = cache.Value()
	ok = cache.Validate(validateArg)
	if !ok {
		c.memory.Delete(key)
		return
	}

	return
}

func (c *InMemoryCache[K, V, D]) Set(key K, data D, validateValue V) {
	cache := c.cacheFactory.Create(data, validateValue)
	c.memory.Store(key, cache)
	return
}

func (c *InMemoryCache[K, V, D]) Delete(key K) (ok bool) {
	_, ok = c.memory.LoadAndDelete(key)
	return
}
