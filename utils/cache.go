package utils

import (
	"bytes"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

type Cache struct {
	pool     *ristretto.Cache[string, []byte]
	lifetime time.Duration
	mu       sync.RWMutex
}

func NewCache(lifetime time.Duration) *Cache {
	if lifetime < time.Second {
		lifetime = time.Second
	}

	cache, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 1_000_000,
		MaxCost:     64 << 20,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}

	return &Cache{pool: cache, lifetime: lifetime}
}

func (c *Cache) Set(key string, value []byte) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.pool == nil {
		return
	}

	if c.pool.SetWithTTL(key, bytes.Clone(value), int64(len(key))+int64(len(value)), c.lifetime) {
		c.pool.Wait()
	}
}

func (c *Cache) Invalidate(key string) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.pool == nil {
		return
	}

	c.pool.Del(key)
	c.pool.Wait()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.pool == nil {
		return nil, false
	}

	data, found := c.pool.Get(key)
	return bytes.Clone(data), found
}

func (c *Cache) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pool == nil {
		return
	}

	c.pool.Close()
	c.pool = nil
}
