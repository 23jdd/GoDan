package cache

import (
	"sync"
	"time"
)

type item struct {
	value  interface{}
	expire int64
}

type MemoryCache struct {
	mu      sync.RWMutex
	items   map[string]*item
	cleanup time.Duration
}

func NewMemoryCache(cleanupInterval time.Duration) *MemoryCache {
	c := &MemoryCache{
		items:   make(map[string]*item),
		cleanup: cleanupInterval,
	}
	go c.cleanExpired()
	return c
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().UnixNano() > it.expire {
		return nil, false
	}
	return it.value, true
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	c.items[key] = &item{
		value:  value,
		expire: time.Now().Add(ttl).UnixNano(),
	}
	c.mu.Unlock()
}

func (c *MemoryCache) Del(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *MemoryCache) cleanExpired() {
	ticker := time.NewTicker(c.cleanup)
	for range ticker.C {
		now := time.Now().UnixNano()
		c.mu.Lock()
		for k, it := range c.items {
			if now > it.expire {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
