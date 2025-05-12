// Package locales provides utilities for working with localization files.
package locales

import (
	"fmt"
	"github.com/up1io/muxo/logger"
	"os"
	"os/exec"
	"path/filepath"
)

// Builder compiles .po files under Root into .mo files.
type Builder struct {
	// Root is the directory containing .po files
	Root string
	// Log is the logger to use for logging messages
	Log logger.Logger
}

// NewBuilder creates a new Builder with the given root directory.
func NewBuilder(root string) *Builder {
	return &Builder{
		Root: root,
		Log:  logger.Default,
	}
}

// WithLogger sets the logger for the Builder.
func (b *Builder) WithLogger(log logger.Logger) *Builder {
	b.Log = log
	return b
}

func (b *Builder) Install() error {
	return b.CheckDependencies()
}

// CheckDependencies checks if the required dependencies are installed.
func (b *Builder) CheckDependencies() error {
	if _, err := exec.LookPath("msgfmt"); err != nil {
		return fmt.Errorf("msgfmt binary not found in PATH: %w", err)
	}
	return nil
}

// Process compiles .po files to .mo files.
// It walks the Root directory and compiles all .po files to .mo files.
func (b *Builder) Process() error {
	if err := b.CheckDependencies(); err != nil {
		return err
	}

	return filepath.Walk(b.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".po" {
			return nil
		}

		rel, err := filepath.Rel(b.Root, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}

		out := filepath.Join(b.Root, rel[:len(rel)-len(filepath.Ext(rel))]+".mo")

		if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", out, err)
		}

		cmd := exec.Command("msgfmt", path, "-o", out)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to compile %s: %w", path, err)
		}

		b.Log.Info("[locale] %s -> %s", path, out)
		return nil
	})
}
