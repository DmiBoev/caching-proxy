package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DmiBoev/caching-proxy/internal/config"
	"github.com/DmiBoev/caching-proxy/internal/logger"
	"github.com/DmiBoev/caching-proxy/internal/proxy"
)

func main() {
	// 1. Загрузка конфигурации (из флагов или файла)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Инициализация логгера
	level := logger.ParseLevel(cfg.LogLevel)
	logger.Init(level)
	slog.Info("Starting caching reverse proxy",
		"port", cfg.Port,
		"target", cfg.TargetURL,
		"cache_ttl", cfg.CacheTTL,
	)

	// 3. Создание прокси с кэшем (LFU)
	proxy, err := proxy.NewReverseProxy(cfg.TargetURL, cfg.CacheTTL)
	if err != nil {
		slog.Error("Failed to create proxy", "error", err)
		os.Exit(1)
	}

	// 4. Настройка HTTP-сервера
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      proxy,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 5. Запуск сервера в отдельной горутине
	go func() {
		slog.Info("Server is listening", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 6. Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// 7. Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited gracefully")
}
