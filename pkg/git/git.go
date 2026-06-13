package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps Git operations.
type Client struct {
	Dir string
}

// NewClient creates a Git client for the given directory.
func NewClient(dir string) *Client {
	return &Client{Dir: dir}
}

func (c *Client) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = c.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %v: %w\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

// IsRepo checks if the directory is a Git repository.
func (c *Client) IsRepo() bool {
	_, err := c.run("rev-parse", "--git-dir")
	return err == nil
}

// CurrentBranch returns the current branch name.
func (c *Client) CurrentBranch() (string, error) {
	return c.run("branch", "--show-current")
}

// HasChanges checks if the working tree has uncommitted changes.
func (c *Client) HasChanges() (bool, error) {
	out, err := c.run("status", "--porcelain")
	if err != nil {
		return false, err
	}
	return out != "", nil
}

// Status returns the repository status.
func (c *Client) Status() (*Status, error) {
	branch, err := c.CurrentBranch()
	if err != nil {
		return nil, err
	}
	hasChanges, err := c.HasChanges()
	if err != nil {
		return nil, err
	}
	return &Status{
		Branch:  branch,
		IsDirty: hasChanges,
	}, nil
}

// Add stages files.
func (c *Client) Add(paths ...string) error {
	args := append([]string{"add"}, paths...)
	_, err := c.run(args...)
	return err
}

// Commit creates a commit with the given message.
func (c *Client) Commit(message string) error {
	_, err := c.run("commit", "-m", message)
	return err
}

// SafeCommit commits only tracked files matching the given patterns.
func (c *Client) SafeCommit(message string, patterns []string) error {
	if len(patterns) == 0 {
		return c.Commit(message)
	}
	for _, p := range patterns {
		if err := c.Add(p); err != nil {
			return err
		}
	}
	return c.Commit(message)
}
