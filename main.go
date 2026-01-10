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
	if err := logging.Setup(); err != nil {
		log.Fatal(err)
	}

	// Load embedded configuration
	cfg := config.GetConfig()

	// Create fetcher with configuration
	f := fetcher.NewFromConfigs(cfg.Sources, cfg.MinFetchInterval, cfg.MaxSubscribers)

	// Start fetcher in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Give fetcher a moment to initialize
	go f.Start(ctx)
	time.Sleep(100 * time.Millisecond)

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
		log.Fatal(err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		shutdownOnce.Do(func() {
			log.Println("Shutting down...")
			cancel()

			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				log.Printf("Server shutdown error: %v", err)
			}

			done := make(chan struct{})
			go func() {
				defer close(done)
				f.Shutdown()
			}()

			select {
			case <-done:
				log.Println("Fetcher shutdown complete")
			case <-time.After(10 * time.Second):
				log.Println("Fetcher shutdown timeout - proceeding anyway")
			}
		})
	}()

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
