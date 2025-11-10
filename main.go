package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/weipengyu/pod-index/pkg/cache"
	"github.com/weipengyu/pod-index/pkg/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 初始化 Pod 缓存
	podCache, err := cache.NewPodCache()
	if err != nil {
		log.Fatalf("Failed to initialize pod cache: %v", err)
	}

	// 启动 informer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := podCache.Start(ctx); err != nil {
		log.Fatalf("Failed to start pod cache: %v", err)
	}

	// 等待缓存同步
	log.Println("Waiting for cache to sync...")
	if !podCache.WaitForCacheSync(ctx) {
		log.Fatal("Failed to sync cache")
	}
	log.Println("Cache synced successfully")

	// 设置 HTTP 处理器
	h := handler.NewHandler(podCache)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/pod", h.GetPodByUID)
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/ready", h.Ready)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动 HTTP 服务器
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 优雅关闭
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
