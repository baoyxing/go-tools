package v1

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Person struct {
	Name string
}

func TestExpireCache(t *testing.T) {
	cache := NewLruCache[string, Person](10)
	cache.SetWithExp("john", Person{Name: "John Smith"}, time.Second)

	got, ok := cache.Get("john")
	assert.Equal(t, ok, true)
	assert.NotNil(t, got)
	assert.Equal(t, got.Name, "John Smith")

	time.Sleep(2 * time.Second)
	got, ok = cache.Get("john")
	assert.NotEqual(t, ok, true)
	assert.Empty(t, got)
}
