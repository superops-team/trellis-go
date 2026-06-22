package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// SessionConfig holds configuration for session recording.
type SessionConfig struct {
	MaxJournalLines int    // Default: 2000
	AutoCommit      bool   // Default: true
	CommitMessage   string // Default: "chore: record journal"
}

// DefaultSessionConfig returns sensible defaults.
func DefaultSessionConfig() SessionConfig {
	return SessionConfig{
		MaxJournalLines: 2000,
		AutoCommit:      true,
		CommitMessage:   "chore: record journal",
	}
}

// SessionEntry represents one recorded session.
type SessionEntry struct {
	Title   string
	Commits []string
	Summary string
	TaskID  string
}

// SessionIndexEntry is a single entry in the session index.
type SessionIndexEntry struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Commits     []string `json:"commits"`
	JournalFile string   `json:"journal_file"`
	StartedAt   string   `json:"started_at"`
	FinishedAt  string   `json:"finished_at"`
}

// sessionIndex is the structure of index.json.
type sessionIndex struct {
	Sessions []SessionIndexEntry `json:"sessions"`
}

// SessionRecorder manages session journal recording.
type SessionRecorder struct {
	WorkspaceDir string
	Config       SessionConfig
}

// RecordSession appends a session entry to the current journal.
func (r *SessionRecorder) RecordSession(entry SessionEntry) error {
	if err := os.MkdirAll(r.WorkspaceDir, 0755); err != nil {
		return fmt.Errorf("create workspace dir: %w", err)
	}

	now := time.Now()
	journalFile, err := r.currentJournal()
	if err != nil {
		return fmt.Errorf("find current journal: %w", err)
	}

	journalPath := filepath.Join(r.WorkspaceDir, journalFile)

	// Build journal entry
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Session %s\n\n", now.Format("2006-01-02 15:04")))
	if entry.TaskID != "" {
		sb.WriteString(fmt.Sprintf("- Task: %s\n", entry.TaskID))
	}
	if len(entry.Commits) > 0 {
		sb.WriteString(fmt.Sprintf("- Commits: %s\n", strings.Join(entry.Commits, ", ")))
	}
	if entry.Summary != "" {
		sb.WriteString(fmt.Sprintf("- Summary: %s\n", entry.Summary))
	}
	sb.WriteString("\n---\n\n")

	// Append to journal
	f, err := os.OpenFile(journalPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open journal: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(sb.String()); err != nil {
		return fmt.Errorf("write journal: %w", err)
	}

	// Check if rotation needed
	if r.needsRotation(journalPath) {
		if err := r.rotate(journalFile); err != nil {
			return fmt.Errorf("rotate journal: %w", err)
		}
	}

	// Update index (use nanosecond precision to avoid ID collisions)
	sessionID := now.Format("2006-01-02T15:04:05.000000000Z07:00")
	return r.updateIndex(SessionIndexEntry{
		ID:          sessionID,
		Title:       entry.Title,
		Commits:     entry.Commits,
		JournalFile: journalFile,
		StartedAt:   sessionID,
		FinishedAt:  sessionID,
	})
}

// ListSessions returns all recorded sessions from the index.
func (r *SessionRecorder) ListSessions() ([]SessionIndexEntry, error) {
	idx, err := r.readIndex()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return idx.Sessions, nil
}

// currentJournal returns the filename of the current (latest) journal.
func (r *SessionRecorder) currentJournal() (string, error) {
	entries, err := os.ReadDir(r.WorkspaceDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "journal-1.md", nil
		}
		return "", err
	}

	maxNum := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "journal-") && strings.HasSuffix(e.Name(), ".md") {
			numStr := strings.TrimPrefix(e.Name(), "journal-")
			numStr = strings.TrimSuffix(numStr, ".md")
			if n, err := strconv.Atoi(numStr); err == nil && n > maxNum {
				maxNum = n
			}
		}
	}

	if maxNum == 0 {
		return "journal-1.md", nil
	}
	return fmt.Sprintf("journal-%d.md", maxNum), nil
}

// needsRotation checks if the journal file exceeds the line limit.
func (r *SessionRecorder) needsRotation(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	lines := strings.Count(string(data), "\n")
	return lines >= r.Config.MaxJournalLines
}

// rotate creates a new journal file and adds a continuation note.
func (r *SessionRecorder) rotate(oldFile string) error {
	oldPath := filepath.Join(r.WorkspaceDir, oldFile)

	// Extract journal number
	numStr := strings.TrimPrefix(oldFile, "journal-")
	numStr = strings.TrimSuffix(numStr, ".md")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return fmt.Errorf("parse journal number: %w", err)
	}

	newFile := fmt.Sprintf("journal-%d.md", num+1)
	continuation := fmt.Sprintf("\n(continued in %s)\n", newFile)

	f, err := os.OpenFile(oldPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open old journal for continuation: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(continuation); err != nil {
		return fmt.Errorf("write continuation: %w", err)
	}

	// Create the new journal file
	newPath := filepath.Join(r.WorkspaceDir, newFile)
	if err := os.WriteFile(newPath, nil, 0644); err != nil {
		return fmt.Errorf("create new journal: %w", err)
	}

	return nil
}

// readIndex reads the session index from index.json.
func (r *SessionRecorder) readIndex() (*sessionIndex, error) {
	path := filepath.Join(r.WorkspaceDir, "index.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var idx sessionIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parse index: %w", err)
	}
	return &idx, nil
}

// updateIndex adds or updates a session entry in the index.
func (r *SessionRecorder) updateIndex(entry SessionIndexEntry) error {
	idx, err := r.readIndex()
	if err != nil {
		if os.IsNotExist(err) {
			idx = &sessionIndex{}
		} else {
			return err
		}
	}

	// Update existing or append
	found := false
	for i, s := range idx.Sessions {
		if s.ID == entry.ID {
			idx.Sessions[i] = entry
			found = true
			break
		}
	}
	if !found {
		idx.Sessions = append(idx.Sessions, entry)
	}

	// Sort by started_at descending
	sort.Slice(idx.Sessions, func(i, j int) bool {
		return idx.Sessions[i].StartedAt > idx.Sessions[j].StartedAt
	})

	return r.writeIndex(idx)
}

// writeIndex writes the session index to index.json.
func (r *SessionRecorder) writeIndex(idx *sessionIndex) error {
	path := filepath.Join(r.WorkspaceDir, "index.json")
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal index: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write index: %w", err)
	}
	return nil
}
