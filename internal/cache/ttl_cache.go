package cache

import (
	"sync"
	"time"
)

// TTLCache хранит записи с TTL
type TTLCache struct {
	mu    sync.RWMutex
	items map[string]Entry
	ttl   time.Duration
}

func New(ttl time.Duration) *TTLCache {
	return &TTLCache{
		items: make(map[string]Entry),
		ttl:   ttl,
	}
}

func (c *TTLCache) Get(key string) (Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.items[key]
	if !ok {
		return Entry{}, false
	}
	// Проверка TTL
	if time.Now().After(entry.Expiry) {
		return Entry{}, false
	}
	return entry, true
}

func (c *TTLCache) Set(key string, entry Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry.Expiry = time.Now().Add(c.ttl)
	c.items[key] = entry
}

// CleanExpired удаляет устаревшие записи (можно запускать в фоне)
func (c *TTLCache) CleanExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.items {
		if time.Now().After(v.Expiry) {
			delete(c.items, k)
		}
	}
}
