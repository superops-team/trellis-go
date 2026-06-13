package task

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mindfold/trellis/pkg/fsutil"
)

// Status represents the task lifecycle status.
type Status string

const (
	StatusPlanning   Status = "planning"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
)

// Task is the core data structure for a task, compatible with original task.json.
type Task struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Assignee  string    `json:"assignee"`
	Branch    string    `json:"branch"`
	Subtasks  []Subtask `json:"subtasks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Subtask represents a subtask item.
type Subtask struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Validate checks the task data for correctness.
func (t *Task) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("task ID is required")
	}
	if t.Name == "" {
		return fmt.Errorf("task name is required")
	}
	if t.Status != StatusPlanning && t.Status != StatusInProgress && t.Status != StatusCompleted {
		return fmt.Errorf("invalid status: %s", t.Status)
	}
	return nil
}

// DirName generates the task directory name: MM-DD-<name>.
func (t *Task) DirName() string {
	return fmt.Sprintf("%02d-%02d-%s", t.CreatedAt.Month(), t.CreatedAt.Day(), t.Name)
}

// Save writes the task to its task.json file.
func (t *Task) Save(path string) error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal task: %w", err)
	}
	return fsutil.WriteFile(path, data, 0644)
}

// LoadTask reads a Task from a task.json file.
func LoadTask(path string) (*Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read task: %w", err)
	}
	var t Task
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parse task: %w", err)
	}
	return &t, nil
}
