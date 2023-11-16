package v1

import (
	genericsCache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/mru"
	v1 "github.com/baoyxing/go-tools/cache/v1"
	"time"
)

type MruCache[K comparable, V any] struct {
	cache *genericsCache.Cache[K, V]
}

func (c MruCache[K, V]) Get(key K) (value V, ok bool) {
	return c.cache.Get(key)
}

func (c MruCache[K, V]) Set(key K, val V) {
	c.cache.Set(key, val)
}

func (c MruCache[K, V]) SetWithExp(key K, val V, exp time.Duration) {
	c.cache.Set(key, val, genericsCache.WithExpiration(exp))
}

func (c MruCache[K, V]) Delete(key K) {
	c.cache.Delete(key)
}

func (c MruCache[K, V]) Keys() []K {
	return c.cache.Keys()
}
func (c MruCache[K, V]) Contains(key K) bool {
	return c.cache.Contains(key)
}
func (c MruCache[K, V]) Empty() {
	keys := c.Keys()
	for _, key := range keys {
		c.Delete(key)
	}
}
func (c MruCache[K, V]) Len() int {
	return len(c.cache.Keys())
}
func NewMruCache[K comparable, V any](capacity int) v1.Cache[K, V] {
	opts := genericsCache.AsMRU[K, V](mru.WithCapacity(capacity))
	c := genericsCache.New[K, V](opts)
	return &MruCache[K, V]{
		c,
	}
}
