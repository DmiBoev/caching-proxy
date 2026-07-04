package config

import (
	"flag"
	"time"
)

type Config struct {
	Port              int
	TargetURL         string
	CacheTTL          time.Duration
	MaxCacheSize      int64
	CacheableStatuses []int
	LogLevel          string
}

func Load() (*Config, error) {
	var (
		port     = flag.Int("port", 8080, "Port to serve on")
		target   = flag.String("target", "http://example.com", "Target URL to proxy")
		cacheTTL = flag.Duration("cache-ttl", 5*time.Minute, "Cache TTL")
		maxSize  = flag.Int64("max-cache-size", 100*1024*1024, "Max cache size in bytes (default 100MB)")
		logLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	)
	flag.Parse()

	return &Config{
		Port:              *port,
		TargetURL:         *target,
		CacheTTL:          *cacheTTL,
		MaxCacheSize:      *maxSize,
		CacheableStatuses: []int{200, 301, 302},
		LogLevel:          *logLevel,
	}, nil
}
