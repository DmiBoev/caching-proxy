package cache

import "time"

type Cache interface {
	Get(key string) (Entry, bool)
	Set(key string, entry Entry)
}

// Entry хранит ответ и метаданные
type Entry struct {
	Data       []byte
	StatusCode int
	Header     map[string][]string
	Expiry     time.Time
}
