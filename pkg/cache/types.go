package cache

import (
	"sync"
	"time"
)

// EventType represents the type of file system event
type EventType int

const (
	EventCreate EventType = iota
	EventWrite
	EventRemove
	EventRename
)

func (e EventType) String() string {
	switch e {
	case EventCreate:
		return "CREATE"
	case EventWrite:
		return "WRITE"
	case EventRemove:
		return "REMOVE"
	case EventRename:
		return "RENAME"
	default:
		return "UNKNOWN"
	}
}

// FileEvent represents a file system event for a markdown file
type FileEvent struct {
	Path    string
	Type    EventType
	Timestamp time.Time
}

// FileState represents the cached state of a file
type FileState struct {
	ContentHash string
	LastModified time.Time
	Size        int64
	Assignee    []string
}

// Cache stores file state with debouncing support
type Cache struct {
	mu          sync.RWMutex
	files       map[string]*FileState
	lastEvent   map[string]time.Time
	debounceMs  int
}
