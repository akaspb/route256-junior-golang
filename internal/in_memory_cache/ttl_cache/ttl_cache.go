package ttl_cache

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/in_memory_cache"
	"time"
)

type CacheData[D any] struct {
	validUntil time.Time
	data       D
}

func (c *CacheData[D]) Value() D {
	return c.data
}

func (c *CacheData[D]) Validate(now time.Time) bool {
	return c.validUntil.Before(now)
}

type TTLCacheFactory[K comparable, D any] struct {
	validTime time.Duration
}

func NewTTLCacheFactory[K comparable, D any](validTime time.Duration) *TTLCacheFactory[K, D] {
	return &TTLCacheFactory[K, D]{validTime: validTime}
}

func (f *TTLCacheFactory[K, D]) Create(
	data D, now time.Time,
) in_memory_cache.Cache[time.Time, D] {
	return &CacheData[D]{
		validUntil: now.Add(f.validTime),
		data:       data,
	}
}

func NewTTLCache[K comparable, D any](validTime time.Duration) *in_memory_cache.InMemoryCache[time.Time, D] {
	cacheFactory := NewTTLCacheFactory[K, D](validTime)
	return in_memory_cache.NewInMemoryCache[time.Time, D](cacheFactory)
}
