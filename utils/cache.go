package utils

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

type Cache struct {
	pool *bigcache.BigCache
}

func NewCache(lifetime time.Duration) *Cache {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(lifetime))

	return &Cache{pool: cache}
}

func (c *Cache) Set(key string, value []byte) {
	_ = c.pool.Set(key, value)
}

func (c *Cache) Invalidate(key string) {
	_ = c.pool.Delete(key)
}

func (c *Cache) Get(key string) ([]byte, bool) {
	data, err := c.pool.Get(key)
	if err != nil {
		return nil, false
	}

	return data, true
}
