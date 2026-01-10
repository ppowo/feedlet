package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

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
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     3,
		Compress:   false,
	}

	multiWriter := io.MultiWriter(os.Stdout, logger)
	log.SetOutput(multiWriter)
	log.SetFlags(0)

	log.Printf("Logging to %s", logPath)
	return nil
}

func getLogDir() (string, error) {
	switch {
	case os.Getenv("XDG_STATE_HOME") != "":
		return filepath.Join(os.Getenv("XDG_STATE_HOME"), "feedlet", "logs"), nil
	case os.Getenv("HOME") != "":
		home := os.Getenv("HOME")
		if _, err := os.Stat(filepath.Join(home, "Library")); err == nil {
			return filepath.Join(home, "Library", "Logs", "feedlet"), nil
		}
		return filepath.Join(home, ".local", "state", "feedlet", "logs"), nil
	default:
		return filepath.Join(".", "logs"), nil
	}
}
