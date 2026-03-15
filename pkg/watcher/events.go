package watcher

import (
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"maestro/pkg/cache"
)

// EventProcessor processes fsnotify events and converts them to FileEvents
type EventProcessor struct {
	cache *cache.Cache
}

// NewEventProcessor creates a new event processor
func NewEventProcessor() *EventProcessor {
	return &EventProcessor{
		cache: cache.NewCache(),
	}
}

// IsMarkdownFile checks if a file path ends with .md extension
func IsMarkdownFile(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".md")
}

// ProcessEvent converts a fsnotify event to a FileEvent
func (p *EventProcessor) ProcessEvent(event fsnotify.Event) (*cache.FileEvent, error) {
	// Filter to only markdown files
	if !IsMarkdownFile(event.Name) {
		return nil, nil
	}

	// Check for debounce
	if p.cache.ShouldDebounce(event.Name) {
		return nil, nil
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
		return nil, nil
	}

	// Update cache state (except for remove events)
	if eventType != cache.EventRemove {
		state, err := cache.GetFileState(event.Name)
		if err != nil {
			return nil, err
		}
		p.cache.UpdateState(event.Name, state)
	} else {
		p.cache.RemoveState(event.Name)
	}

	return &cache.FileEvent{
		Path:      event.Name,
		Type:      eventType,
		Timestamp: time.Now(),
	}, nil
}
