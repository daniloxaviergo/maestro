package matcher

import (
	"log"
	"strings"

	"maestro/pkg/agent"
)

// Matcher manages agent matching logic.
// It matches assignee names from task YAML frontmatter to configured agents.
type Matcher struct {
	agentMap map[string]*agent.Agent // lowercase agent name -> Agent
	agents   []*agent.Agent          // original agent list for iteration order
}

// NewMatcher creates a new Matcher with the given list of agents.
// Agent names are stored case-insensitively for matching.
func NewMatcher(agents []*agent.Agent) *Matcher {
	m := &Matcher{
		agentMap: make(map[string]*agent.Agent),
		agents:   agents,
	}

	// Build lookup map with lowercase names
	for _, a := range agents {
		name := strings.ToLower(a.GetName())
		m.agentMap[name] = a
	}

	return m
}

// MatchAssignees matches assignee names to configured agents.
// Returns a list of matching agents (can be empty if no matches found).
// Matching is case-insensitive.
func (m *Matcher) MatchAssignees(assignees []string) []*agent.Agent {
	var matchedAgents []*agent.Agent

	for _, assignee := range assignees {
		lowerAssignee := strings.ToLower(assignee)
		if agent, exists := m.agentMap[lowerAssignee]; exists {
			matchedAgents = append(matchedAgents, agent)
		} else {
			log.Printf("Warning: No agent found for assignee %q (matched %d of %d agents)", assignee, len(matchedAgents), len(assignees))
		}
	}

	return matchedAgents
}
