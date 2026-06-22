package agent

import (
	"strings"
	"testing"
)

func TestImplementAgent(t *testing.T) {
	a := ImplementAgent()
	if a.Name != "trellis-implement" {
		t.Errorf("expected name 'trellis-implement', got %q", a.Name)
	}
	if a.Description == "" {
		t.Error("expected non-empty description")
	}
	if len(a.Tools) == 0 {
		t.Error("expected non-empty tools")
	}
	if !strings.Contains(a.Content, "trellis-implement") {
		t.Error("content should contain agent name")
	}
	if !strings.Contains(a.Content, "implement.jsonl") {
		t.Error("content should reference implement.jsonl")
	}
	if !strings.Contains(a.Content, "Do not commit") {
		t.Error("content should mention no-commit rule")
	}
}

func TestCheckAgent(t *testing.T) {
	a := CheckAgent()
	if a.Name != "trellis-check" {
		t.Errorf("expected name 'trellis-check', got %q", a.Name)
	}
	if !strings.Contains(a.Content, "max 3 rounds") {
		t.Error("content should mention 3-round retry limit")
	}
	if !strings.Contains(a.Content, "PASSED or FAILED") {
		t.Error("content should mention PASSED/FAILED report")
	}
}

func TestResearchAgent(t *testing.T) {
	a := ResearchAgent()
	if a.Name != "trellis-research" {
		t.Errorf("expected name 'trellis-research', got %q", a.Name)
	}
	if !strings.Contains(a.Content, "Read-only") {
		t.Error("content should mention read-only constraint")
	}
	if !strings.Contains(a.Content, "No file writes") {
		t.Error("content should mention no file writes")
	}
	// Research agent should not have Write/Edit tools
	for _, tool := range a.Tools {
		if tool == "Write" || tool == "Edit" {
			t.Errorf("research agent should not have %s tool", tool)
		}
	}
}

func TestAllAgents(t *testing.T) {
	agents := AllAgents()
	if len(agents) != 3 {
		t.Errorf("expected 3 agents, got %d", len(agents))
	}
	names := make(map[string]bool)
	for _, a := range agents {
		if names[a.Name] {
			t.Errorf("duplicate agent name: %s", a.Name)
		}
		names[a.Name] = true
	}
	for _, want := range []string{"trellis-implement", "trellis-check", "trellis-research"} {
		if !names[want] {
			t.Errorf("missing agent: %s", want)
		}
	}
}
