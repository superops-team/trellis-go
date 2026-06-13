package task

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/mindfold/trellis/pkg/fsutil"
)

var (
	ErrInvalidTransition = errors.New("invalid task status transition")
	ErrTaskNotFound      = errors.New("task not found")
)

// Phase represents the context phase for adding entries.
type Phase string

const (
	PhaseImplement Phase = "implement"
	PhaseCheck     Phase = "check"
	PhaseResearch  Phase = "research"
)

// CreateOptions provides options for task creation.
type CreateOptions struct {
	Assignee string
	Branch   string
	FromRef  string // Optional: copy context from existing task
}

// Manager handles the task lifecycle.
type Manager struct {
	root string // .trellis/tasks/ path
}

// NewManager creates a new task manager.
func NewManager(root string) *Manager {
	return &Manager{root: root}
}

// Create creates a new task and returns the task and its directory path.
func (m *Manager) Create(name string, opts CreateOptions) (*Task, string, error) {
	now := time.Now()
	task := &Task{
		ID:        name,
		Name:      name,
		Status:    StatusPlanning,
		Assignee:  opts.Assignee,
		Branch:    opts.Branch,
		Subtasks:  []Subtask{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	dirName := task.DirName()
	taskDir := filepath.Join(m.root, dirName)
	if _, err := os.Stat(taskDir); err == nil {
		return nil, "", fmt.Errorf("task directory already exists: %s", taskDir)
	}

	if err := fsutil.EnsureDir(taskDir); err != nil {
		return nil, "", err
	}

	// Write task.json
	taskPath := filepath.Join(taskDir, "task.json")
	if err := task.Save(taskPath); err != nil {
		return nil, "", err
	}

	// Create empty files
	for _, f := range []string{"prd.md", "implement.jsonl", "check.jsonl"} {
		if err := os.WriteFile(filepath.Join(taskDir, f), []byte{}, 0644); err != nil {
			return nil, "", err
		}
	}

	// Create research directory
	if err := fsutil.EnsureDir(filepath.Join(taskDir, "research")); err != nil {
		return nil, "", err
	}

	return task, taskDir, nil
}

// Start transitions a task from planning to in_progress.
func (m *Manager) Start(taskID string) error {
	task, path, err := m.findTask(taskID)
	if err != nil {
		return err
	}
	if task.Status != StatusPlanning {
		return fmt.Errorf("%w: cannot start task with status %s", ErrInvalidTransition, task.Status)
	}
	task.Status = StatusInProgress
	task.UpdatedAt = time.Now()
	return task.Save(path)
}

// Archive transitions a task to completed and moves it to archive/YYYY-MM/.
func (m *Manager) Archive(taskID string) error {
	task, path, err := m.findTask(taskID)
	if err != nil {
		return err
	}
	if task.Status != StatusInProgress {
		return fmt.Errorf("%w: cannot archive task with status %s", ErrInvalidTransition, task.Status)
	}

	task.Status = StatusCompleted
	task.UpdatedAt = time.Now()

	// Move to archive
	dir := filepath.Dir(path)
	archiveDir := filepath.Join(m.root, "archive", task.UpdatedAt.Format("2006-01"))
	if err := fsutil.EnsureDir(archiveDir); err != nil {
		return err
	}
	archivePath := filepath.Join(archiveDir, filepath.Base(dir))
	if err := os.Rename(dir, archivePath); err != nil {
		return fmt.Errorf("archive task: %w", err)
	}

	return task.Save(filepath.Join(archivePath, "task.json"))
}

// Current returns the currently active task (from .runtime/sessions/).
func (m *Manager) Current() (*Task, error) {
	// TODO: read from .runtime/sessions/
	return nil, ErrTaskNotFound
}

// List returns all non-archived tasks.
func (m *Manager) List() ([]Task, error) {
	entries, err := os.ReadDir(m.root)
	if err != nil {
		return nil, fmt.Errorf("read tasks: %w", err)
	}

	var tasks []Task
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "archive" {
			continue
		}
		taskPath := filepath.Join(m.root, entry.Name(), "task.json")
		task, err := LoadTask(taskPath)
		if err != nil {
			continue
		}
		tasks = append(tasks, *task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})
	return tasks, nil
}

// Get retrieves a task by ID.
func (m *Manager) Get(taskID string) (*Task, error) {
	task, _, err := m.findTask(taskID)
	return task, err
}

// AddContext adds a context entry to the task's manifest.
func (m *Manager) AddContext(taskID string, phase Phase, entry ContextEntry) error {
	taskDir, err := m.findTaskDir(taskID)
	if err != nil {
		return err
	}

	var manifestFile string
	switch phase {
	case PhaseImplement:
		manifestFile = "implement.jsonl"
	case PhaseCheck:
		manifestFile = "check.jsonl"
	default:
		return fmt.Errorf("unknown phase: %s", phase)
	}

	manifestPath := filepath.Join(taskDir, manifestFile)
	manifest, err := loadManifest(manifestPath)
	if err != nil {
		manifest = &Manifest{Version: "1.0", Entries: []ContextEntry{}}
	}

	manifest.Entries = append(manifest.Entries, entry)
	return manifest.save(manifestPath)
}

// Validate checks the task directory structure for completeness.
func (m *Manager) Validate(taskID string) error {
	taskDir, err := m.findTaskDir(taskID)
	if err != nil {
		return err
	}

	required := []string{"task.json", "prd.md", "implement.jsonl", "check.jsonl"}
	for _, f := range required {
		path := filepath.Join(taskDir, f)
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("missing required file: %s", f)
		}
	}
	return nil
}

func (m *Manager) findTask(taskID string) (*Task, string, error) {
	taskDir, err := m.findTaskDir(taskID)
	if err != nil {
		return nil, "", err
	}
	path := filepath.Join(taskDir, "task.json")
	task, err := LoadTask(path)
	if err != nil {
		return nil, "", err
	}
	return task, path, nil
}

func (m *Manager) findTaskDir(taskID string) (string, error) {
	entries, err := os.ReadDir(m.root)
	if err != nil {
		return "", fmt.Errorf("read tasks: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "archive" {
			continue
		}
		path := filepath.Join(m.root, entry.Name(), "task.json")
		task, err := LoadTask(path)
		if err != nil {
			continue
		}
		if task.ID == taskID {
			return filepath.Join(m.root, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("%w: %s", ErrTaskNotFound, taskID)
}
