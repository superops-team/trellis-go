package prd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrRequired is returned when a PRD is required but missing or empty.
var ErrRequired = errors.New("PRD is required")

// LoadRequired reads and validates a PRD file from a task directory.
func LoadRequired(taskDir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(taskDir, "prd.md"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrRequired
		}
		return "", fmt.Errorf("read PRD: %w", err)
	}
	if strings.TrimSpace(string(data)) == "" {
		return "", ErrRequired
	}
	return string(data), nil
}
