package proxy

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/DmiBoev/caching-proxy/internal/cache"
	"github.com/DmiBoev/caching-proxy/internal/config"
)

// Proxy основной обработчик
type Proxy struct {
	target          *url.URL
	reverseProxy    *httputil.ReverseProxy
	cache           cache.Cache
	cacheTTL        time.Duration
	cacheableStatus map[int]bool
}

// NewReverseProxy создаёт новый прокси
func NewReverseProxy(cfg *config.Config) (*Proxy, error) {
	parsed, err := url.Parse(cfg.TargetURL)
	if err != nil {
		return nil, err
	}

	// Статусы, которые будем кэшировать
	cacheable := map[int]bool{
		http.StatusOK:               true,
		http.StatusMovedPermanently: true,
		http.StatusFound:            true,
		// можно добавить другие
	}

	c := cache.NewLFU(cfg.MaxCacheItems, cfg.CacheTTL)

	return &Proxy{
		target:          parsed,
		reverseProxy:    httputil.NewSingleHostReverseProxy(parsed),
		cache:           c,
		cacheableStatus: cacheable,
	}, nil
}

// ServeHTTP реализует http.Handler
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Кэшируем только GET и HEAD
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		p.reverseProxy.ServeHTTP(w, r)
		return
	}

	// Генерируем ключ
	key := p.generateCacheKey(r)
	if key == "" {
		p.reverseProxy.ServeHTTP(w, r)
		return
	}

	slog.Debug("Generated cache key", "key", key, "method", r.Method, "path", r.URL.Path)

	// Проверяем кэш
	if entry, ok := p.cache.Get(key); ok {
		slog.Info("Cache HIT", "method", r.Method, "path", r.URL.Path)
		p.writeCacheResponse(w, entry)
		return
	}

	slog.Info("Cache MISS", "method", r.Method, "path", r.URL.Path)

	// Перехватчик ответа
	capturer := &responseCapturer{
		ResponseWriter: w,
		buffer:         bytes.NewBuffer(nil),
	}

	// Проксируем запрос
	p.reverseProxy.ServeHTTP(capturer, r)
	slog.Debug("Response captured", "status", capturer.statusCode, "cacheable", p.cacheableStatus[capturer.statusCode])
	// Если статус кэшируемый – сохраняем
	if p.cacheableStatus[capturer.statusCode] {
		entry := cache.Entry{
			Data:       capturer.buffer.Bytes(),
			StatusCode: capturer.statusCode,
			Header:     capturer.header,
		}
		p.cache.Set(key, entry)
	}
}

// writeCacheResponse отправляет ответ из кэша
func (p *Proxy) writeCacheResponse(w http.ResponseWriter, entry cache.Entry) {
	for k, v := range entry.Header {
		for _, val := range v {
			w.Header().Add(k, val)
		}
	}
	w.Header().Set("X-Cache", "HIT")
	w.WriteHeader(entry.StatusCode)
	w.Write(entry.Data)
}
