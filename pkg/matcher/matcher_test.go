package matcher

import (
	"fmt"
	"testing"

	"maestro/pkg/agent"
	"maestro/pkg/config"
)

// createTestAgent creates a test agent with the given name
func createTestAgent(name string) *agent.Agent {
	return agent.NewAgent(name, "/fake/path/config.yml")
}

// TestNewMatcher_EmptyAgents tests creating a matcher with no agents
func TestNewMatcher_EmptyAgents(t *testing.T) {
	matcher := NewMatcher(nil)
	if matcher == nil {
		t.Fatal("Expected matcher to be non-nil")
	}

	// Matching with no agents should return empty slice
	matched := matcher.MatchAssignees([]string{"alice"})
	if len(matched) != 0 {
		t.Errorf("Expected no matches, got %d", len(matched))
	}
}

// TestNewMatcher_SingleAgent tests creating a matcher with one agent
func TestNewMatcher_SingleAgent(t *testing.T) {
	agents := []*agent.Agent{createTestAgent("alice")}
	matcher := NewMatcher(agents)

	matched := matcher.MatchAssignees([]string{"alice"})
	if len(matched) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matched))
	}
	if matched[0].GetName() != "alice" {
		t.Errorf("Expected matched agent to be 'alice', got %q", matched[0].GetName())
	}
}

// TestNewMatcher_MultipleAgents tests creating a matcher with multiple agents
func TestNewMatcher_MultipleAgents(t *testing.T) {
	agents := []*agent.Agent{
		createTestAgent("alice"),
		createTestAgent("bob"),
		createTestAgent("charlie"),
	}
	matcher := NewMatcher(agents)

	matched := matcher.MatchAssignees([]string{"alice", "bob"})
	if len(matched) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matched))
	}
}

// TestMatchAssignees_NoMatches tests when no assignees match any agents
func TestMatchAssignees_NoMatches(t *testing.T) {
	agents := []*agent.Agent{
		createTestAgent("alice"),
		createTestAgent("bob"),
	}
	matcher := NewMatcher(agents)

	matched := matcher.MatchAssignees([]string{"charlie", "dave"})
	if len(matched) != 0 {
		t.Errorf("Expected no matches, got %d", len(matched))
	}
}

// TestMatchAssignees_SingleMatch tests when one assignee matches one agent
func TestMatchAssignees_SingleMatch(t *testing.T) {
	agents := []*agent.Agent{createTestAgent("alice")}
	matcher := NewMatcher(agents)

	matched := matcher.MatchAssignees([]string{"alice"})
	if len(matched) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matched))
	}
}

// TestMatchAssignees_MultipleMatches tests when multiple assignees match multiple agents
func TestMatchAssignees_MultipleMatches(t *testing.T) {
	agents := []*agent.Agent{
		createTestAgent("alice"),
		createTestAgent("bob"),
		createTestAgent("charlie"),
	}
	matcher := NewMatcher(agents)

	matched := matcher.MatchAssignees([]string{"alice", "bob", "charlie"})
	if len(matched) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(matched))
	}
}

// TestMatchAssignees_CaseInsensitive tests that matching is case-insensitive
func TestMatchAssignees_CaseInsensitive(t *testing.T) {
	agents := []*agent.Agent{
		createTestAgent("Alice"),
		createTestAgent("Bob"),
	}
	matcher := NewMatcher(agents)

	testCases := []struct {
		assignee string
		expected string
	}{
		{"alice", "Alice"},
		{"ALICE", "Alice"},
		{"Alice", "Alice"},
		{"bob", "Bob"},
		{"BOB", "Bob"},
		{"Bob", "Bob"},
	}

	for _, tc := range testCases {
		matched := matcher.MatchAssignees([]string{tc.assignee})
		if len(matched) != 1 {
			t.Errorf("Expected 1 match for %q, got %d", tc.assignee, len(matched))
		}
		if matched[0].GetName() != tc.expected {
			t.Errorf("Expected matched agent %q, got %q", tc.expected, matched[0].GetName())
		}
	}
}

// TestMatchAssignees_PartialMatch tests when some assignees match and some don't
func TestMatchAssignees_PartialMatch(t *testing.T) {
	agents := []*agent.Agent{
		createTestAgent("alice"),
		createTestAgent("bob"),
	}
	matcher := NewMatcher(agents)

	matched := matcher.MatchAssignees([]string{"alice", "charlie", "bob", "dave"})
	if len(matched) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matched))
	}
	// Verify the matched agents are alice and bob
	matchNames := make(map[string]bool)
	for _, a := range matched {
		matchNames[a.GetName()] = true
	}
	if !matchNames["alice"] || !matchNames["bob"] {
		t.Errorf("Expected to match alice and bob, got %v", matchNames)
	}
}

// TestMatchAssignees_EmptyInput tests when assignee list is empty
func TestMatchAssignees_EmptyInput(t *testing.T) {
	agents := []*agent.Agent{createTestAgent("alice")}
	matcher := NewMatcher(agents)

	matched := matcher.MatchAssignees([]string{})
	if len(matched) != 0 {
		t.Errorf("Expected no matches for empty input, got %d", len(matched))
	}
}

// TestMatchAssignees_DuplicateAssignees tests behavior with duplicate assignees
func TestMatchAssignees_DuplicateAssignees(t *testing.T) {
	agents := []*agent.Agent{createTestAgent("alice")}
	matcher := NewMatcher(agents)

	matched := matcher.MatchAssignees([]string{"alice", "alice"})
	if len(matched) != 2 {
		t.Errorf("Expected 2 matches (duplicates allowed), got %d", len(matched))
	}
}

// TestMatchAssignees_ReturnsOriginalAgentOrder tests that matched agents preserve original order
func TestMatchAssignees_ReturnsOriginalAgentOrder(t *testing.T) {
	agents := []*agent.Agent{
		createTestAgent("charlie"),
		createTestAgent("alice"),
		createTestAgent("bob"),
	}
	matcher := NewMatcher(agents)

	matched := matcher.MatchAssignees([]string{"bob", "alice", "charlie"})
	if len(matched) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(matched))
	}

	// Verify order is preserved based on assignee order
	expectedOrder := []string{"bob", "alice", "charlie"}
	for i, name := range expectedOrder {
		if matched[i].GetName() != name {
			t.Errorf("Expected matched agent %d to be %q, got %q", i, name, matched[i].GetName())
		}
	}
}

// TestNewMatcher_NilAgents tests creating a matcher with nil agents
func TestNewMatcher_NilAgents(t *testing.T) {
	matcher := NewMatcher(nil)
	if matcher == nil {
		t.Fatal("Expected matcher to be non-nil")
	}

	// Matching with nil agents should return empty slice
	matched := matcher.MatchAssignees([]string{"alice"})
	if len(matched) != 0 {
		t.Errorf("Expected no matches, got %d", len(matched))
	}
}

// TestNewMatcher_SingleAgentEmptyName tests creating a matcher with an agent with empty name
func TestNewMatcher_SingleAgentEmptyName(t *testing.T) {
	// Use the NewTestAgent helper to create agent with empty name
	agentInstance := agent.NewTestAgent("", config.AgentConfig{})

	agents := []*agent.Agent{agentInstance}
	matcher := NewMatcher(agents)

	// Matching with empty name should work if assignee is also empty
	matched := matcher.MatchAssignees([]string{""})
	if len(matched) != 1 {
		t.Errorf("Expected 1 match for empty assignee, got %d", len(matched))
	}
}

// TestMatchAssignees_SpecialCharacters tests matching with special characters in agent names
func TestMatchAssignees_SpecialCharacters(t *testing.T) {
	// Create agents with special characters
	specialChars := []string{"alice Smith", "bob_jones", "charlie.brown", "dave123"}
	var agents []*agent.Agent
	for _, name := range specialChars {
		agents = append(agents, createTestAgent(name))
	}

	matcher := NewMatcher(agents)

	for _, assignee := range specialChars {
		matched := matcher.MatchAssignees([]string{assignee})
		if len(matched) != 1 {
			t.Errorf("Expected 1 match for %q, got %d", assignee, len(matched))
		}
	}
}

// TestMatchAssignees_WhitespaceHandling tests handling of leading/trailing whitespace
func TestMatchAssignees_WhitespaceHandling(t *testing.T) {
	agents := []*agent.Agent{
		createTestAgent("alice"),
	}
	matcher := NewMatcher(agents)

	// Exact match should work
	matched := matcher.MatchAssignees([]string{"alice"})
	if len(matched) != 1 {
		t.Errorf("Expected 1 match for 'alice', got %d", len(matched))
	}
}

// TestMatchAssignees_LargeAssigneeList tests with a large number of assignees
func TestMatchAssignees_LargeAssigneeList(t *testing.T) {
	var agents []*agent.Agent
	var assignees []string

	// Create 100 agents
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("agent%d", i)
		agents = append(agents, createTestAgent(name))
		assignees = append(assignees, name)
	}

	matcher := NewMatcher(agents)

	matched := matcher.MatchAssignees(assignees)
	if len(matched) != 100 {
		t.Errorf("Expected 100 matches, got %d", len(matched))
	}
}

// TestNewMatcher_DuplicateAgents tests handling of duplicate agents in input
func TestNewMatcher_DuplicateAgents(t *testing.T) {
	agent1 := createTestAgent("alice")
	agent2 := createTestAgent("alice") // Same name

	agents := []*agent.Agent{agent1, agent2}
	matcher := NewMatcher(agents)

	// The last one should overwrite the first in the map
	matched := matcher.MatchAssignees([]string{"alice"})
	if len(matched) != 1 {
		t.Errorf("Expected 1 match (duplicate agent names), got %d", len(matched))
	}
}
