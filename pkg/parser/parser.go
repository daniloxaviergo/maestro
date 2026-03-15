package parser

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ParseFile parses a markdown file and extracts the YAML frontmatter
func ParseFile(filePath string) FileData {
	startTime := time.Now()
	result := FileData{
		FilePath:    filePath,
		Frontmatter: Frontmatter{Assignee: []string{}},
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		result.Error = fmt.Errorf("failed to read file: %w", err)
		result.ParseTime = time.Since(startTime)
		return result
	}

	// Try to extract frontmatter
	frontmatter, err := extractFrontmatter(content)
	if err != nil {
		result.Error = fmt.Errorf("failed to extract frontmatter: %w", err)
		result.ParseTime = time.Since(startTime)
		return result
	}

	// Parse YAML frontmatter
	if frontmatter != "" {
		var fm Frontmatter
		if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
			result.Error = fmt.Errorf("failed to parse YAML frontmatter: %w", err)
			result.ParseTime = time.Since(startTime)
			return result
		}

		// Ensure Assignee is always a slice (not nil)
		if fm.Assignee == nil {
			fm.Assignee = []string{}
		}

		result.Frontmatter = fm
	}

	result.ParseTime = time.Since(startTime)
	return result
}

// extractFrontmatter extracts the YAML frontmatter from markdown content
// Returns empty string if no frontmatter is found
func extractFrontmatter(content []byte) (string, error) {
	text := string(content)

	// Frontmatter must start with ---
	if !strings.HasPrefix(text, "---") {
		return "", nil
	}

	// Find the closing ---
	rest := strings.TrimPrefix(text, "---")
	lines := strings.Split(rest, "\n")

	endIndex := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			endIndex = i
			break
		}
	}

	if endIndex == -1 {
		return "", fmt.Errorf("no closing --- found for frontmatter")
	}

	// Extract frontmatter content
	frontmatter := strings.Join(lines[:endIndex], "\n")
	return frontmatter, nil
}
