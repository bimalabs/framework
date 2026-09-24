package utils

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheOwnershipAndUpdate(t *testing.T) {
	cache := NewCache(time.Minute)
	t.Cleanup(cache.Close)
	input := []byte("first")
	cache.Set("key", input)
	input[0] = 'X'

	value, found := cache.Get("key")
	assert.True(t, found)
	assert.Equal(t, []byte("first"), value)
	value[0] = 'Y'
	value, found = cache.Get("key")
	assert.True(t, found)
	assert.Equal(t, []byte("first"), value)

	cache.Set("key", []byte("second"))
	value, found = cache.Get("key")
	assert.True(t, found)
	assert.Equal(t, []byte("second"), value)
	cache.Invalidate("key")
	cache.Set("key", []byte("third"))
	value, found = cache.Get("key")
	assert.True(t, found)
	assert.Equal(t, []byte("third"), value)
}

func TestCacheEmptyAndMissing(t *testing.T) {
	cache := NewCache(time.Minute)
	t.Cleanup(cache.Close)
	value, found := cache.Get("missing")
	assert.False(t, found)
	assert.Nil(t, value)
	cache.Invalidate("missing")
	cache.Set("", []byte{})
	value, found = cache.Get("")
	assert.True(t, found)
	assert.Empty(t, value)
}

func TestCacheConcurrentClose(t *testing.T) {
	cache := NewCache(time.Minute)
	t.Cleanup(cache.Close)
	var workers sync.WaitGroup
	for i := range 16 {
		workers.Go(func() {
			key := fmt.Sprint(i)
			for range 50 {
				cache.Set(key, []byte(key))
				cache.Get(key)
				cache.Invalidate(key)
			}
		})
	}

	workers.Go(cache.Close)
	workers.Wait()
	cache.Close()
	cache.Set("closed", []byte("value"))
	value, found := cache.Get("closed")
	assert.False(t, found)
	assert.Nil(t, value)
}

func TestCacheMinimumLifetime(t *testing.T) {
	for _, lifetime := range []time.Duration{-time.Second, 0, time.Millisecond} {
		t.Run(lifetime.String(), func(t *testing.T) {
			cache := NewCache(lifetime)
			t.Cleanup(cache.Close)
			cache.Set("key", []byte("value"))
			ttl, found := cache.pool.GetTTL("key")
			assert.True(t, found)
			assert.Positive(t, ttl)
			assert.LessOrEqual(t, ttl, time.Second)
		})
	}
}

func Test_Cache_Ttl(t *testing.T) {
	cache := NewCache(time.Millisecond * 100)
	t.Cleanup(cache.Close)
	cache.Set("test", []byte("a"))

	data, found := cache.Get("test")

	assert.True(t, found)
	assert.Equal(t, string(data), "a")

	time.Sleep(3 * time.Second)

	data, found = cache.Get("test")

	assert.False(t, found)
	assert.Nil(t, data)
}

func Test_Cache_Invalidate(t *testing.T) {
	cache := NewCache(time.Millisecond * 100)
	t.Cleanup(cache.Close)
	cache.Set("test", []byte("a"))

	data, found := cache.Get("test")

	assert.True(t, found)
	assert.Equal(t, string(data), "a")

	cache.Invalidate("test")

	data, found = cache.Get("test")

	assert.False(t, found)
	assert.Nil(t, data)
}
