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
	"github.com/weipengyu/pod-index/pkg/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize pod cache
	podCache, err := cache.NewPodCache()
	if err != nil {
		log.Fatalf("Failed to initialize pod cache: %v", err)
	}

	// Start informer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := podCache.Start(ctx); err != nil {
		log.Fatalf("Failed to start pod cache: %v", err)
	}

	// Wait for cache sync
	log.Println("Waiting for cache to sync...")
	if !podCache.WaitForCacheSync(ctx) {
		log.Fatal("Failed to sync cache")
	}
	log.Println("Cache synced successfully")

	// Setup HTTP handlers with basic auth
	auth := middleware.NewBasicAuth()
	h := handler.NewHandler(podCache)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/pod", auth.Middleware(h.GetPodByUID))
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/ready", h.Ready)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start HTTP server
	go func() {
		log.Printf("Starting server on port %s", port)
		if auth.IsEnabled() {
			log.Println("Basic authentication is enabled")
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
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
