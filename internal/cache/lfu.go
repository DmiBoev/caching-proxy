package cache

import (
	"container/heap"
	"log/slog"
	"sync"
	"time"
)

type Item struct {
	key       string
	entry     Entry
	frequency int
	index     int // индекс в куче
}

type ItemHeap []*Item

func (h ItemHeap) Len() int { return len(h) }
func (h ItemHeap) Less(i, j int) bool {
	return h[i].frequency < h[j].frequency
}

func (h ItemHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *ItemHeap) Push(x any) {
	n := len(*h)
	item := x.(*Item)
	item.index = n
	*h = append(*h, item)
}

func (h *ItemHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*h = old[0 : n-1]
	return item
}

type LFUCache struct {
	mu       sync.RWMutex
	items    map[string]*Item // быстрый доступ по ключу
	h        ItemHeap         // куча для выбора жертвы
	capacity int              // макс. количество элементов
	ttl      time.Duration    // время жизни записей
}

func NewLFU(capacity int, ttl time.Duration) *LFUCache {
	c := &LFUCache{
		items:    make(map[string]*Item),
		h:        make(ItemHeap, 0),
		capacity: capacity,
		ttl:      ttl,
	}
	heap.Init(&c.h)
	return c
}

func (c *LFUCache) Get(key string) (Entry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return Entry{}, false
	}

	// Проверка TTL
	if time.Now().After(item.entry.Expiry) {
		c.removeItem(item)
		return Entry{}, false
	}

	// Увеличиваем частоту и обновляем позицию в куче
	item.frequency++
	heap.Fix(&c.h, item.index)

	return item.entry, true
}

func (c *LFUCache) Set(key string, entry Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// элемент уже есть обновление и увеличение частоты
	if item, ok := c.items[key]; ok {
		item.entry = entry
		item.entry.Expiry = time.Now().Add(c.ttl)
		item.frequency++
		heap.Fix(&c.h, item.index)
		return
	}

	// full cache – удаляем наименее частый элемент
	if len(c.items) >= c.capacity {
		victim := heap.Pop(&c.h).(*Item)
		slog.Debug("Evicted LFU item", "key", victim.key, "frequency", victim.frequency)
		delete(c.items, victim.key)
	}

	// Создаём новый элемент
	entry.Expiry = time.Now().Add(c.ttl)
	item := &Item{
		key:       key,
		entry:     entry,
		frequency: 1,
	}
	c.items[key] = item
	heap.Push(&c.h, item)
}

// removeItem удаляет элемент из кэша при истечении TTL
func (c *LFUCache) removeItem(item *Item) {
	delete(c.items, item.key)
	heap.Remove(&c.h, item.index)
}
