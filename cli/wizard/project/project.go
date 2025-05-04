package project

import (
	"context"
	"errors"
	"fmt"
	"github.com/up1io/muxo/utils"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const dirPerm os.FileMode = 0750

// Todo(john): Move this to a config file outside of the code
const defaultTemplateRepoUrl string = "https://github.com/up1io/muxo-base-project-template"

// predefinedIgnoreList is a collection of all files/directories that must be ignored.
var predefinedIgnoreList = []string{
	".git",
	".muxoignore",
}

// Execute starts project setup based on the provided config.
func Execute(ctx context.Context) error {
	cfg, ok := ConfigFromContext(ctx)
	if !ok {
		return ErrConfigNotFound
	}

	srcDir := fmt.Sprintf("%s/tmp_%s_%v", os.TempDir(), cfg.ProjectName, time.Now().Unix())
	if err := os.Mkdir(srcDir, dirPerm); err != nil {
		return err
	}

	destDir := cfg.ProjectDir

	if err := cloneRepo(defaultTemplateRepoUrl, srcDir); err != nil {
		return err
	}

	ignoreList, err := loadMuxoIgnore(srcDir)
	if err != nil {
		return err
	}

	// Note(john): Make sure that the following file/dir is not transferred from the src.
	for _, s := range predefinedIgnoreList {
		ignoreList = append(ignoreList, s)
	}

	if err := transferFiles(srcDir, destDir, cfg, ignoreList); err != nil {
		return err
	}

	return nil
}

// loadMuxoIgnore searches for a muxo ignore file in the directory and returns the parsed values.
// If no muxo ignore file is found, an empty slice is returned.
func loadMuxoIgnore(dir string) ([]string, error) {
	var out []string

	b, err := os.ReadFile(fmt.Sprintf("%s/.muxoignore", dir))
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	} else if err != nil {
		return out, err
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			out = append(out, trimmed)
		}
	}

	return out, nil
}

// transferFiles walk to directory and transfer all static and dynamic files and directory to the destination.
func transferFiles(srcDir, destDir string, cfg *Config, shouldIgnore []string) error {
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(srcDir, path)
		if utils.Contains(shouldIgnore, d.Name()) || utils.Contains(shouldIgnore, relPath) {
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		if d.IsDir() {
			if relPath != destDir && relPath != "." {
				if err = os.Mkdir(fmt.Sprintf("%s/%s", destDir, relPath), dirPerm); err != nil {
					return err
				}
			}

			return nil
		}

		if strings.HasSuffix(relPath, ".dyn") {
			if err := renderFile(srcDir, destDir, relPath, cfg); err != nil {
				return err
			}

			return nil
		}

		if err := copyStaticFile(srcDir, destDir, relPath); err != nil {
			return err
		}

		return nil
	})
}

// renderFile renders a .dyn template file with provided data.
func renderFile(srcDir, destDir, relPath string, cfg *Config) error {
	b, err := os.ReadFile(fmt.Sprintf("%s/%s", srcDir, relPath))
	if err != nil {
		return fmt.Errorf("failed to read dyn file. %s", err)
	}

	templ, err := template.New("file").Parse(string(b))
	if err != nil {
		return err
	}

	s, _ := strings.CutSuffix(relPath, ".dyn")
	destFile, err := os.Create(fmt.Sprintf("%s/%s", destDir, s))
	if err := templ.Execute(destFile, *cfg); err != nil {
		return err
	}

	return nil
}

// copyStaticFile copies a regular (non-template) file from source to destination.
func copyStaticFile(srcDir, destDir string, relPath string) error {
	srcFile, err := os.Open(fmt.Sprintf("%s/%s", srcDir, relPath))
	if err != nil {
		return fmt.Errorf("failed to open static file during copy. %s", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(fmt.Sprintf("%s/%s", destDir, relPath))
	if err != nil {
		return fmt.Errorf("failed to create static file on dest dir. %s", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// cloneRepo clone repo inside target directory.
func cloneRepo(repoUrl, targetDir string) error {
	cmd := exec.Command("git", "clone", repoUrl, targetDir)
	return cmd.Run()
}
