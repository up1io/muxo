package templater

import (
	"fmt"
	"os"
	"os/exec"
)

// Templater runs the `templ generate` command in a working directory.
type Templater struct {
	Dir string // working directory
}

// Install verifies that the templ binary is available.
func (t *Templater) Install() error {
	if _, err := exec.LookPath("templ"); err != nil {
		return fmt.Errorf("templ binary not found in PATH: %w", err)
	}

	return nil
}

// Process invokes `templ generate DIR`.
func (t *Templater) Process() error {
	cmd := exec.Command("templ", "generate", t.Dir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("templ generate failed: %w", err)
	}

	fmt.Printf("[templater] templates generated in %s\n", t.Dir)
	return nil
}
