package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"maestro/pkg/cache"
)

// ErrWatcherStopped is returned when the watcher is stopped
var ErrWatcherStopped = fmt.Errorf("watcher stopped")

// Watcher wraps fsnotify.Watcher for markdown file monitoring
type Watcher struct {
	watcher    *fsnotify.Watcher
	eventChan  chan cache.FileEvent
	done       chan struct{}
	processor  *EventProcessor
	watchPaths []string
}

// WatcherOption configures a Watcher
type WatcherOption func(*Watcher)

// WithWatchPaths sets the watch paths for the watcher
func WithWatchPaths(paths []string) WatcherOption {
	return func(w *Watcher) {
		w.watchPaths = paths
	}
}

// NewWatcher creates a new file watcher
func NewWatcher(opts ...WatcherOption) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		watcher:   fsWatcher,
		eventChan: make(chan cache.FileEvent, 100),
		done:      make(chan struct{}),
		processor: NewEventProcessor(),
	}

	// Apply options
	for _, opt := range opts {
		opt(w)
	}

	return w, nil
}

// AddWatch adds a directory to watch recursively
func (w *Watcher) AddWatch(path string) error {
	// Resolve to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Check if directory exists
	info, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", absPath)
	}

	// Add watch with recursive behavior via fsnotify
	if err := w.watcher.Add(absPath); err != nil {
		return err
	}

	w.watchPaths = append(w.watchPaths, absPath)
	return nil
}

// Watch starts the file watching loop
func (w *Watcher) Watch() error {
	go func() {
		defer close(w.eventChan)
		
		for {
			select {
			case <-w.done:
				return
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				
				// Filter to only markdown files
				if !IsMarkdownFile(event.Name) {
					continue
				}

				// Check for debounce
				if w.processor.cache.ShouldDebounce(event.Name) {
					continue
				}

				var eventType cache.EventType
				switch event.Op {
				case fsnotify.Create:
					eventType = cache.EventCreate
				case fsnotify.Write:
					eventType = cache.EventWrite
				case fsnotify.Remove:
					eventType = cache.EventRemove
				case fsnotify.Rename:
					eventType = cache.EventRename
				default:
					continue
				}

				// Update cache state (except for remove events)
				if eventType != cache.EventRemove {
					state, err := cache.GetFileState(event.Name)
					if err != nil {
						// File may have been deleted between event and read
						if !os.IsNotExist(err) {
							fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", event.Name, err)
						}
						continue
					}
					w.processor.cache.UpdateState(event.Name, state)
				} else {
					w.processor.cache.RemoveState(event.Name)
				}

				// Send event to channel
				select {
				case w.eventChan <- cache.FileEvent{
					Path:      event.Name,
					Type:      eventType,
					Timestamp: time.Now(),
				}:
				case <-w.done:
					return
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				// Log error but continue watching
				fmt.Fprintf(os.Stderr, "watcher error: %v\n", err)
			}
		}
	}()

	return nil
}

// Stop stops the watcher and closes the event channel
func (w *Watcher) Stop() {
	close(w.done)
	w.watcher.Close()
}

// Events returns the channel for receiving file events
func (w *Watcher) Events() <-chan cache.FileEvent {
	return w.eventChan
}

// GetAllMarkdownFiles returns all markdown files in watched paths
func (w *Watcher) GetAllMarkdownFiles() ([]string, error) {
	var files []string
	
	for _, path := range w.watchPaths {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Continue on error
			}
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(filePath), ".md") {
				files = append(files, filePath)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	
	return files, nil
}
