package change_detect

import (
	"maestro/pkg/cache"
	"maestro/pkg/logs"
	"maestro/pkg/notifier"
	"maestro/pkg/parser"
)

// Detector compares cached assignee values with current values and logs changes
type Detector struct {
	cache      *cache.Cache
	logger     *logs.Logger
	processed  map[string]bool // Track first run for each file
	notifier   *notifier.Notifier
}

// NewDetector creates a new change detector
func NewDetector(logger *logs.Logger) *Detector {
	return &Detector{
		cache:     cache.NewCache(),
		logger:    logger,
		processed: make(map[string]bool),
	}
}

// SetNotifier sets the notifier to use for assignee change notifications
func (d *Detector) SetNotifier(n *notifier.Notifier) {
	d.notifier = n
}

// ProcessFile processes a parsed file and detects assignee changes
// Returns true if a log entry was written, false otherwise
func (d *Detector) ProcessFile(fileData parser.FileData) (bool, error) {
	if fileData.Error != nil {
		// Log error but don't fail - use empty assignee array
		return false, nil
	}

	filePath := fileData.FilePath
	newAssignee := fileData.Frontmatter.Assignee

	// Get cached assignee
	cachedAssignee, exists := d.cache.GetAssignee(filePath)

	// First run: no cached value, just store and return
	if !exists {
		d.cache.SetAssignee(filePath, newAssignee)
		d.processed[filePath] = true
		return false, nil
	}

	// Compare assignees (order-insensitive)
	if assigneesEqual(cachedAssignee, newAssignee) {
		// No change, just update cache
		d.cache.SetAssignee(filePath, newAssignee)
		return false, nil
	}

	// Assignee changed - log the change
	if err := d.logger.LogAssigneeChange(filePath, cachedAssignee, newAssignee); err != nil {
		return false, err
	}

	// Notify tmux if notifier is configured
	if d.notifier != nil {
		event := notifier.AssigneeChangeEvent{
			FilePath:    filePath,
			OldAssignee: cachedAssignee,
			NewAssignee: newAssignee,
		}
		d.notifier.Notify(event)
	}

	// Update cache with new assignee
	d.cache.SetAssignee(filePath, newAssignee)
	return true, nil
}

// assigneesEqual compares two assignee slices order-insensitively
func assigneesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Create maps to count occurrences
	aCount := make(map[string]int)
	for _, val := range a {
		aCount[val]++
	}

	bCount := make(map[string]int)
	for _, val := range b {
		bCount[val]++
	}

	return mapEqual(aCount, bCount)
}

// mapEqual compares two string maps for equality
func mapEqual(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// RemoveFile removes a file from the cache (for delete events)
func (d *Detector) RemoveFile(filePath string) {
	d.cache.RemoveAssignee(filePath)
	delete(d.processed, filePath)
}
