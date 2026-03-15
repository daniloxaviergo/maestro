package cache

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"time"
)

const (
	defaultDebounceMs = 50
)

// NewCache creates a new file cache with debouncing
func NewCache() *Cache {
	return &Cache{
		files:      make(map[string]*FileState),
		lastEvent:  make(map[string]time.Time),
		debounceMs: defaultDebounceMs,
	}
}

// GetState retrieves the cached state for a file
func (c *Cache) GetState(path string) (*FileState, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	state, exists := c.files[path]
	return state, exists
}

// UpdateState updates the cached state for a file
func (c *Cache) UpdateState(path string, state *FileState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.files[path] = state
}

// RemoveState removes a file from the cache
func (c *Cache) RemoveState(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.files, path)
}

// ShouldDebounce checks if an event should be debounced
func (c *Cache) ShouldDebounce(path string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	now := time.Now()
	last, exists := c.lastEvent[path]
	if !exists {
		c.lastEvent[path] = now
		return false
	}
	
	elapsed := now.Sub(last)
	if elapsed.Milliseconds() < int64(c.debounceMs) {
		return true
	}
	
	c.lastEvent[path] = now
	return false
}

// GetContentHash computes a hash of file content
func GetContentHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:]), nil
}

// GetFileState retrieves the current state of a file
func GetFileState(path string) (*FileState, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	
	hash, err := GetContentHash(path)
	if err != nil {
		return nil, err
	}
	
	return &FileState{
		ContentHash: hash,
		LastModified: info.ModTime(),
		Size:        info.Size(),
	}, nil
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.files = make(map[string]*FileState)
	c.lastEvent = make(map[string]time.Time)
}

// GetAssignee retrieves the cached assignee for a file
func (c *Cache) GetAssignee(path string) ([]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	state, exists := c.files[path]
	if !exists {
		return nil, false
	}
	// Return a copy to prevent external modification
	assignee := make([]string, len(state.Assignee))
	copy(assignee, state.Assignee)
	return assignee, true
}

// SetAssignee updates the cached assignee for a file
func (c *Cache) SetAssignee(path string, assignee []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if state, exists := c.files[path]; exists {
		state.Assignee = assignee
	}
}

// RemoveAssignee removes the assignee entry for a file
func (c *Cache) RemoveAssignee(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if state, exists := c.files[path]; exists {
		state.Assignee = nil
	}
}
