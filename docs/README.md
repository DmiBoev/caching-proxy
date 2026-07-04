# Caching Reverse Proxy with LFU

HTTP reverse proxy with in-memory cache using LFU (Least Frequently Used) eviction policy.

## Features
- Reverse proxy to any origin server
- In-memory caching with TTL (Time-To-Live)
- LFU eviction policy using `container/heap`
- Configurable via command line flags
- Graceful shutdown

## How to run
```bash
go run cmd/proxy/main.go --target http://httpbin.org --port 8080