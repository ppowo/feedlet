package logging

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	currentLogFile *os.File
	logDir         string
	mu             sync.Mutex
	cleanupDone    chan struct{}
)

// Setup configures logging to write to both stdout and a log file in the OS log directory
// Old log files (older than 3 days) are automatically cleaned up
// Starts a background goroutine to rotate logs daily
func Setup(ctx context.Context) error {
	// Get OS-specific log directory
	var err error
	logDir, err = getLogDir()
	if err != nil {
		return fmt.Errorf("failed to get log directory: %w", err)
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open initial log file
	if err := rotateLogFile(); err != nil {
		return err
	}

	// Clean up old log files on startup
	if err := cleanOldLogs(logDir); err != nil {
		log.Printf("Warning: failed to clean old logs: %v", err)
	}

	// Initialize cleanup done channel
	cleanupDone = make(chan struct{})

	// Start background goroutine for daily rotation and cleanup
	go dailyRotation(ctx, cleanupDone)

	return nil
}

// rotateLogFile closes the current log file and opens a new one for today
func rotateLogFile() error {
	mu.Lock()
	defer mu.Unlock()

	// Close existing file if open
	if currentLogFile != nil {
		currentLogFile.Close()
	}

	// Create log file with current date
	logFile := filepath.Join(logDir, fmt.Sprintf("feedlet-%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	currentLogFile = f

	// Log to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Logging to %s", logFile)
	return nil
}

// dailyRotation runs in the background to rotate logs at midnight and clean old logs
func dailyRotation(ctx context.Context, done chan struct{}) {
	for {
		// Calculate time until next midnight
		now := time.Now()
		tomorrow := now.Add(24 * time.Hour)
		nextMidnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, now.Location())
		duration := nextMidnight.Sub(now)

		// Sleep until midnight or context is cancelled
		select {
		case <-time.After(duration):
			// Time to rotate
		case <-ctx.Done():
			log.Println("Daily rotation cancelled via context")
			return
		case <-done:
			log.Println("Daily rotation cancelled")
			return
		}

		// Rotate log file
		if err := rotateLogFile(); err != nil {
			log.Printf("Error rotating log file: %v", err)
		}

		// Clean up old logs
		if err := cleanOldLogs(logDir); err != nil {
			log.Printf("Warning: failed to clean old logs: %v", err)
		}

		// Check if we're shutting down
		select {
		case <-done:
			log.Println("Daily rotation shutting down")
			return
		default:
		}
	}
}

// getLogDir returns the OS-specific log directory for feedlet
func getLogDir() (string, error) {
	switch {
	case os.Getenv("XDG_STATE_HOME") != "":
		// Linux with XDG_STATE_HOME set
		return filepath.Join(os.Getenv("XDG_STATE_HOME"), "feedlet", "logs"), nil
	case os.Getenv("HOME") != "":
		// macOS/Linux fallback
		home := os.Getenv("HOME")
		if _, err := os.Stat(filepath.Join(home, "Library")); err == nil {
			// macOS
			return filepath.Join(home, "Library", "Logs", "feedlet"), nil
		}
		// Linux fallback to ~/.local/state
		return filepath.Join(home, ".local", "state", "feedlet", "logs"), nil
	default:
		// Fallback to current directory
		return filepath.Join(".", "logs"), nil
	}
}

// cleanOldLogs removes log files older than 3 days
func cleanOldLogs(logDir string) error {
	cutoff := time.Now().AddDate(0, 0, -3)

	entries, err := os.ReadDir(logDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process feedlet-*.log files
		matched, err := filepath.Match("feedlet-*.log", entry.Name())
		if err != nil || !matched {
			continue
		}

		fullPath := filepath.Join(logDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			log.Printf("Warning: failed to stat %s: %v", fullPath, err)
			continue
		}

		if info.ModTime().Before(cutoff) {
			if err := os.Remove(fullPath); err != nil {
				log.Printf("Warning: failed to remove old log %s: %v", fullPath, err)
			} else {
				log.Printf("Removed old log file: %s", fullPath)
			}
		}
	}

	return nil
}
