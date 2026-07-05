package cache

import (
	"sync"
	"time"
)

// Entry хранит ответ и метаданные
type Entry struct {
	Data       []byte
	StatusCode int
	Header     map[string][]string
	Expiry     time.Time
}

// Cache хранит записи с TTL
type Cache struct {
	mu    sync.RWMutex
	items map[string]Entry
	ttl   time.Duration
}

func New(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]Entry),
		ttl:   ttl,
	}
}

func (c *Cache) Get(key string) (Entry, bool) {
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

func (c *Cache) Set(key string, entry Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry.Expiry = time.Now().Add(c.ttl)
	c.items[key] = entry
}

// CleanExpired удаляет устаревшие записи (можно запускать в фоне)
func (c *Cache) CleanExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.items {
		if time.Now().After(v.Expiry) {
			delete(c.items, k)
		}
	}
}
