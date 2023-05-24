package v1

import "time"

type Cache[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, val V)
	SetWithExp(key K, val V, exp time.Duration)
	Keys() []K
	Contains(key K) bool
	Delete(key K)
}
