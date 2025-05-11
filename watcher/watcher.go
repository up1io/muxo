package watcher

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Watcher monitors file changes and triggers callbacks based on extension.
type Watcher struct {
	watcher  *fsnotify.Watcher
	onChange func(path string)
	extMap   map[string]bool
}

// NewWatcher creates a Watcher that watches rootDir for given extensions.
func NewWatcher(rootDir string, extensions []string, onChange func(path string)) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	extMap := make(map[string]bool)
	for _, ext := range extensions {
		extMap[ext] = true
	}
	return &Watcher{watcher: w, onChange: onChange, extMap: extMap}, nil
}

// AddDir recursively adds directories to the watcher.
func (w *Watcher) AddDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return w.watcher.Add(path)
		}
		return nil
	})
}

// Run starts the event loop (blocking).
func (w *Watcher) Run() {
	defer w.watcher.Close()
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				ext := strings.ToLower(filepath.Ext(event.Name))
				if w.extMap[ext] {
					log.Printf("File changed: %s", event.Name)
					w.onChange(event.Name)
				}
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Println("watcher error:", err)
		}
	}
}
