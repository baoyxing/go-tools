package v1

import (
	genericsCache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/clock"
	v1 "github.com/baoyxing/go-tools/cache/v1"
	"time"
)

type ClockCache[K comparable, V any] struct {
	cache *genericsCache.Cache[K, V]
}

func (c ClockCache[K, V]) Get(key K) (value V, ok bool) {
	return c.cache.Get(key)
}

func (c ClockCache[K, V]) Set(key K, val V) {
	c.cache.Set(key, val)
}

func (c ClockCache[K, V]) SetWithExp(key K, val V, exp time.Duration) {
	c.cache.Set(key, val, genericsCache.WithExpiration(exp))
}

func (c ClockCache[K, V]) Delete(key K) {
	c.cache.Delete(key)
}

func (c ClockCache[K, V]) Keys() []K {
	return c.cache.Keys()
}
func (c ClockCache[K, V]) Contains(key K) bool {
	return c.cache.Contains(key)
}
func (c ClockCache[K, V]) Empty() {
	keys := c.Keys()
	for _, key := range keys {
		c.Delete(key)
	}
}
func (c ClockCache[K, V]) Len() int {
	return len(c.cache.Keys())
}

func NewClockCache[K comparable, V any](capacity int) v1.Cache[K, V] {
	opts := genericsCache.AsClock[K, V](clock.WithCapacity(capacity))
	c := genericsCache.New[K, V](opts)
	return &ClockCache[K, V]{
		c,
	}
}
