package parser

import "time"

// Frontmatter represents the YAML frontmatter section of a markdown file
type Frontmatter struct {
	ID       string   `yaml:"id"`
	Title    string   `yaml:"title"`
	Assignee []string `yaml:"assignee"`
	Status   string   `yaml:"status"`
}

// FileData represents the result of parsing a file
type FileData struct {
	FilePath   string
	Frontmatter Frontmatter
	Error      error
	ParseTime  time.Duration
}
