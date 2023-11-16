package v1

import (
	genericsCache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/fifo"
	v1 "github.com/baoyxing/go-tools/cache/v1"
	"time"
)

type FifoCache[K comparable, V any] struct {
	cache *genericsCache.Cache[K, V]
}

func (c FifoCache[K, V]) Get(key K) (value V, ok bool) {
	return c.cache.Get(key)
}

func (c FifoCache[K, V]) Set(key K, val V) {
	c.cache.Set(key, val)
}

func (c FifoCache[K, V]) SetWithExp(key K, val V, exp time.Duration) {
	c.cache.Set(key, val, genericsCache.WithExpiration(exp))
}

func (c FifoCache[K, V]) Delete(key K) {
	c.cache.Delete(key)
}

func (c FifoCache[K, V]) Keys() []K {
	return c.cache.Keys()
}
func (c FifoCache[K, V]) Contains(key K) bool {
	return c.cache.Contains(key)
}
func (c FifoCache[K, V]) Empty() {
	keys := c.Keys()
	for _, key := range keys {
		c.Delete(key)
	}
}
func (c FifoCache[K, V]) Len() int {
	return len(c.cache.Keys())
}

func NewFifoCache[K comparable, V any](capacity int) v1.Cache[K, V] {
	opts := genericsCache.AsFIFO[K, V](fifo.WithCapacity(capacity))
	c := genericsCache.New[K, V](opts)
	return &FifoCache[K, V]{
		c,
	}
}
