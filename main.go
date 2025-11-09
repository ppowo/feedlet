package main

//go:generate go run install_tools.go

import (
	"context"
	_ "embed"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ppowo/feedlet/internal/config"
	"github.com/ppowo/feedlet/internal/fetcher"
	"github.com/ppowo/feedlet/internal/logging"
	"github.com/ppowo/feedlet/internal/server"
	"github.com/ppowo/feedlet/web"
)

var (
	shutdownOnce sync.Once
)

func main() {
	// Setup logging
	if err := logging.Setup(context.Background()); err != nil {
		log.Fatalf("Failed to setup logging: %v", err)
	}

	// Load embedded configuration
	cfg := config.GetConfig()

	// Create fetcher with configuration
	f := fetcher.NewFromConfigs(cfg.Sources, cfg.MinFetchInterval, cfg.MaxSubscribers)

	// Start fetcher in background (no need to wait for initial fetch)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go f.Start(ctx)

	// Create and start server
	port := cfg.Port
	if port == 0 {
		port = 8080
	}

	// Build source limits and days maps
	sourceLimits := make(map[string]int)
	sourceDays := make(map[string]int)
	for _, srcCfg := range cfg.Sources {
		if srcCfg.Limit > 0 {
			sourceLimits[srcCfg.Name] = srcCfg.Limit
		}
		if srcCfg.Days > 0 {
			sourceDays[srcCfg.Name] = srcCfg.Days
		} else {
			sourceDays[srcCfg.Name] = 2 // Default 2 days
		}
	}

	srv, err := server.New(f, web.IndexTemplate, port, sourceLimits, sourceDays)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		shutdownOnce.Do(func() {
			log.Println("Shutting down...")
			cancel()

			// Give goroutines time to shut down gracefully
			shutdownTimeout := time.After(3 * time.Second)

			// Wait for shutdown or timeout
			<-shutdownTimeout

			// Cleanup logging goroutine
			logging.Cleanup()
		})
	}()

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
