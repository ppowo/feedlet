//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Install development tools
func InstallTools() error {
	fmt.Println("Installing development tools...")
	if err := sh.Run("go", "install", "github.com/air-verse/air@latest"); err != nil {
		return err
	}
	fmt.Println("Tools installed!")
	return nil
}

// Run with auto-reload (development mode)
func Dev() error {
	mg.Deps(InstallTools)
	return sh.RunV("air")
}

// Build the binary for x86_64 Linux (default target)
func Build() error {
	fmt.Println("Building feedlet for x86_64 Linux...")
	// Create target directory if it doesn't exist
	if err := os.MkdirAll("target", 0755); err != nil {
		return err
	}

	env := map[string]string{
		"GOOS":   "linux",
		"GOARCH": "amd64",
	}

	if err := sh.RunWith(env, "go", "build", "-o", "target/feedlet", "."); err != nil {
		return err
	}

	// Make the binary executable
	if err := os.Chmod("target/feedlet", 0755); err != nil {
		return err
	}

	fmt.Println("Build complete! Binary at target/feedlet (executable)")
	return nil
}

// Clean build artifacts
func Clean() error {
	fmt.Println("Cleaning...")
	os.RemoveAll("target")
	os.RemoveAll("tmp")
	fmt.Println("Cleaned!")
	return nil
}

// Tidy go modules and install tools
func Setup() error {
	fmt.Println("Running go mod tidy...")
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return err
	}

	fmt.Println("Installing tools...")
	cmd := exec.Command("go", "install", "github.com/magefile/mage@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return InstallTools()
}
