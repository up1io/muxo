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
	Root string // directory containing .po files
}

func (b *Builder) Install() error {
	if _, err := exec.LookPath("msgfmt"); err != nil {
		return fmt.Errorf("msgfmt binary not found in PATH: %w", err)
	}
	return nil
}

// Process runs the locale compilation.
func (b *Builder) Process() error {
	if _, err := exec.LookPath("msgfmt"); err != nil {
		return fmt.Errorf("msgfmt binary not found in PATH: %w", err)
	}

	return filepath.Walk(b.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".po" {
			return nil
		}

		rel, _ := filepath.Rel(b.Root, path)
		out := filepath.Join(b.Root, rel[:len(rel)-len(filepath.Ext(rel))]+".mo")

		if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
			return err
		}

		cmd := exec.Command("msgfmt", path, "-o", out)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to compile %s: %w", path, err)
		}

		logger.Info("[locale] %s -> %s", path, out)
		return nil
	})
}
