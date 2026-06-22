package skill

import (
	"strings"
	"testing"
)

func TestBrainstormSkill(t *testing.T) {
	s := BrainstormSkill()
	if s.Name != "trellis-brainstorm" {
		t.Errorf("expected name 'trellis-brainstorm', got %q", s.Name)
	}
	if s.Description == "" {
		t.Error("expected non-empty description")
	}
	if !strings.Contains(s.Content, "prd.md") {
		t.Error("content should reference prd.md")
	}
	if !strings.Contains(s.Content, "one question at a time") {
		t.Error("content should mention one-question-at-a-time rule")
	}
}

func TestBeforeDevSkill(t *testing.T) {
	s := BeforeDevSkill()
	if s.Name != "trellis-before-dev" {
		t.Errorf("expected name 'trellis-before-dev', got %q", s.Name)
	}
	if !strings.Contains(s.Content, "spec index") {
		t.Error("content should mention spec index")
	}
}

func TestCheckSkill(t *testing.T) {
	s := CheckSkill()
	if s.Name != "trellis-check" {
		t.Errorf("expected name 'trellis-check', got %q", s.Name)
	}
	if !strings.Contains(s.Content, "max 3 rounds") {
		t.Error("content should mention 3-round retry limit")
	}
}

func TestUpdateSpecSkill(t *testing.T) {
	s := UpdateSpecSkill()
	if s.Name != "trellis-update-spec" {
		t.Errorf("expected name 'trellis-update-spec', got %q", s.Name)
	}
	if !strings.Contains(s.Content, "target spec layer") {
		t.Error("content should mention target spec layer")
	}
}

func TestBreakLoopSkill(t *testing.T) {
	s := BreakLoopSkill()
	if s.Name != "trellis-break-loop" {
		t.Errorf("expected name 'trellis-break-loop', got %q", s.Name)
	}
	if !strings.Contains(s.Content, "root cause") {
		t.Error("content should mention root cause analysis")
	}
}

func TestAllSkills(t *testing.T) {
	skills := AllSkills()
	if len(skills) != 5 {
		t.Errorf("expected 5 skills, got %d", len(skills))
	}
	names := make(map[string]bool)
	for _, s := range skills {
		if names[s.Name] {
			t.Errorf("duplicate skill name: %s", s.Name)
		}
		names[s.Name] = true
	}
	for _, want := range []string{
		"trellis-brainstorm",
		"trellis-before-dev",
		"trellis-check",
		"trellis-update-spec",
		"trellis-break-loop",
	} {
		if !names[want] {
			t.Errorf("missing skill: %s", want)
		}
	}
}
