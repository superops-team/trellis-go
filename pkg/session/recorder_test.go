package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestRecorder(dir string) *SessionRecorder {
	cfg := DefaultSessionConfig()
	cfg.MaxJournalLines = 5 // Small for testing rotation
	return &SessionRecorder{
		WorkspaceDir: filepath.Join(dir, "workspace", "test-dev"),
		Config:       cfg,
	}
}

func TestRecordSession_Basic(t *testing.T) {
	dir := t.TempDir()
	r := newTestRecorder(dir)

	err := r.RecordSession(SessionEntry{
		Title:   "add-login",
		TaskID:  "add-login-feature",
		Commits: []string{"abc1234", "def5678"},
		Summary: "Implemented login page",
	})
	if err != nil {
		t.Fatalf("RecordSession() error: %v", err)
	}

	// Check journal file exists
	journalPath := filepath.Join(r.WorkspaceDir, "journal-1.md")
	data, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "add-login-feature") {
		t.Error("journal should contain task ID")
	}
	if !strings.Contains(content, "abc1234") {
		t.Error("journal should contain commit")
	}
	if !strings.Contains(content, "Implemented login page") {
		t.Error("journal should contain summary")
	}
}

func TestRecordSession_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	r := newTestRecorder(dir)
	r.Config.MaxJournalLines = 100 // Disable rotation for this test

	r.RecordSession(SessionEntry{Title: "session-1", TaskID: "task-1"})
	r.RecordSession(SessionEntry{Title: "session-2", TaskID: "task-2"})

	journalPath := filepath.Join(r.WorkspaceDir, "journal-1.md")
	data, _ := os.ReadFile(journalPath)
	content := string(data)

	if !strings.Contains(content, "task-1") || !strings.Contains(content, "task-2") {
		t.Error("journal should contain both sessions")
	}
}

func TestRecordSession_Index(t *testing.T) {
	dir := t.TempDir()
	r := newTestRecorder(dir)

	r.RecordSession(SessionEntry{Title: "session-1", TaskID: "task-1", Commits: []string{"a"}})
	r.RecordSession(SessionEntry{Title: "session-2", TaskID: "task-2", Commits: []string{"b"}})

	sessions, err := r.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error: %v", err)
	}
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
	if sessions[0].Title != "session-2" {
		t.Error("sessions should be sorted by started_at descending")
	}
}

func TestRecordSession_Rotation(t *testing.T) {
	dir := t.TempDir()
	r := newTestRecorder(dir)
	r.Config.MaxJournalLines = 5 // Rotate after 5 lines

	// Each session entry is ~5 lines, so 2 entries should trigger rotation
	r.RecordSession(SessionEntry{Title: "s1", TaskID: "t1", Summary: "summary"})
	r.RecordSession(SessionEntry{Title: "s2", TaskID: "t2", Summary: "summary"})

	// Check journal-1.md has continuation
	j1, _ := os.ReadFile(filepath.Join(r.WorkspaceDir, "journal-1.md"))
	if !strings.Contains(string(j1), "continued in journal-2.md") {
		t.Error("journal-1 should have continuation note")
	}

	// Check journal-2.md exists
	_, err := os.Stat(filepath.Join(r.WorkspaceDir, "journal-2.md"))
	if err != nil {
		t.Error("journal-2.md should exist after rotation")
	}
}

func TestRecordSession_EmptyEntry(t *testing.T) {
	dir := t.TempDir()
	r := newTestRecorder(dir)

	err := r.RecordSession(SessionEntry{})
	if err != nil {
		t.Fatalf("RecordSession() with empty entry should not error: %v", err)
	}
}

func TestListSessions_Empty(t *testing.T) {
	dir := t.TempDir()
	r := newTestRecorder(dir)

	sessions, err := r.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error: %v", err)
	}
	if sessions != nil {
		t.Errorf("expected nil for no sessions, got %v", sessions)
	}
}

func TestListSessions_NoIndex(t *testing.T) {
	dir := t.TempDir()
	r := &SessionRecorder{
		WorkspaceDir: filepath.Join(dir, "nonexistent"),
		Config:       DefaultSessionConfig(),
	}

	sessions, err := r.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error: %v", err)
	}
	if sessions != nil {
		t.Errorf("expected nil for no index, got %v", sessions)
	}
}

func TestCurrentJournal_New(t *testing.T) {
	dir := t.TempDir()
	r := newTestRecorder(dir)

	journal, err := r.currentJournal()
	if err != nil {
		t.Fatalf("currentJournal() error: %v", err)
	}
	if journal != "journal-1.md" {
		t.Errorf("expected journal-1.md, got %s", journal)
	}
}

func TestCurrentJournal_Existing(t *testing.T) {
	dir := t.TempDir()
	r := newTestRecorder(dir)

	// Create journal-1.md and journal-2.md
	os.MkdirAll(r.WorkspaceDir, 0755)
	os.WriteFile(filepath.Join(r.WorkspaceDir, "journal-1.md"), []byte("old"), 0644)
	os.WriteFile(filepath.Join(r.WorkspaceDir, "journal-2.md"), []byte("current"), 0644)

	journal, err := r.currentJournal()
	if err != nil {
		t.Fatalf("currentJournal() error: %v", err)
	}
	if journal != "journal-2.md" {
		t.Errorf("expected journal-2.md, got %s", journal)
	}
}

func TestNeedsRotation(t *testing.T) {
	dir := t.TempDir()
	r := newTestRecorder(dir)
	r.Config.MaxJournalLines = 3

	path := filepath.Join(dir, "test-journal.md")

	// Under limit
	os.WriteFile(path, []byte("line1\nline2\n"), 0644)
	if r.needsRotation(path) {
		t.Error("should not need rotation for 2 lines with limit 3")
	}

	// At limit
	os.WriteFile(path, []byte("line1\nline2\nline3\n"), 0644)
	if !r.needsRotation(path) {
		t.Error("should need rotation for 3 lines with limit 3")
	}
}

func TestIndexFormat(t *testing.T) {
	dir := t.TempDir()
	r := newTestRecorder(dir)

	r.RecordSession(SessionEntry{
		Title:   "test",
		TaskID:  "task-1",
		Commits: []string{"abc"},
		Summary: "test summary",
	})

	// Read index.json and verify format
	idxPath := filepath.Join(r.WorkspaceDir, "index.json")
	data, err := os.ReadFile(idxPath)
	if err != nil {
		t.Fatalf("read index.json: %v", err)
	}

	var idx sessionIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("parse index.json: %v", err)
	}

	if len(idx.Sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(idx.Sessions))
	}

	s := idx.Sessions[0]
	if s.Title != "test" {
		t.Errorf("expected title 'test', got %q", s.Title)
	}
	if len(s.Commits) != 1 || s.Commits[0] != "abc" {
		t.Errorf("expected commits [abc], got %v", s.Commits)
	}
	if s.JournalFile != "journal-1.md" {
		t.Errorf("expected journal_file 'journal-1.md', got %q", s.JournalFile)
	}
}
