package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Setup configures logging to write to both stdout and a rotating log file
func Setup() error {
	logDir, err := getLogDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(logDir, "feedlet.log")

	logger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,   // MB
		MaxBackups: 3,    // keep 3 old files
		MaxAge:     3,    // days
		Compress:   false,
	}

	// Log to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, logger)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Logging to %s", logPath)
	return nil
}

// getLogDir returns the OS-specific log directory for feedlet
func getLogDir() (string, error) {
	switch {
	case os.Getenv("XDG_STATE_HOME") != "":
		return filepath.Join(os.Getenv("XDG_STATE_HOME"), "feedlet", "logs"), nil
	case os.Getenv("HOME") != "":
		home := os.Getenv("HOME")
		if _, err := os.Stat(filepath.Join(home, "Library")); err == nil {
			// macOS
			return filepath.Join(home, "Library", "Logs", "feedlet"), nil
		}
		// Linux
		return filepath.Join(home, ".local", "state", "feedlet", "logs"), nil
	default:
		return filepath.Join(".", "logs"), nil
	}
}
