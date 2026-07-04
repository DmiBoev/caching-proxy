package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/DmiBoev/caching-proxy/internal/config"
)

type ReverseProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	//cache    cache.Cache
	cacheTTL time.Duration
	// возможно другие поля еще доп
}

func NewReverseProxy(cfg *config.Config) (*ReverseProxy, error) {
	targetURL, err := url.Parse(cfg.TargetURL)
	if err != nil {
		return nil, err
	}

	// Создаем LFU-кэш с ограничением по размеру
	//cache := cache.NewLFUCache(cfg.MaxCacheSize)

	// Создаем стандартный ReverseProxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return &ReverseProxy{
		target: targetURL,
		proxy:  proxy,
		//	cache:    cache,
		cacheTTL: cfg.CacheTTL,
	}, nil
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ... logic !!!
}
