package testutil

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TempRepo creates a temporary directory initialized as a Git repository.
func TempRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	git := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}
	git("init")
	git("config", "user.email", "test@test.com")
	git("config", "user.name", "Test")
	return dir
}

// DiffDir recursively compares two directories and returns a description of differences.
func DiffDir(a, b string) string {
	var diffs []string
	_ = filepath.Walk(a, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(a, path)
		other := filepath.Join(b, rel)
		_, err = os.Stat(other)
		if err != nil {
			diffs = append(diffs, "missing in b: "+rel)
			return nil
		}
		adata, _ := os.ReadFile(path)
		bdata, _ := os.ReadFile(other)
		if string(adata) != string(bdata) {
			diffs = append(diffs, "content differs: "+rel)
		}
		return nil
	})
	_ = filepath.Walk(b, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(b, path)
		other := filepath.Join(a, rel)
		_, err = os.Stat(other)
		if err != nil {
			diffs = append(diffs, "missing in a: "+rel)
		}
		return nil
	})
	return strings.Join(diffs, "\n")
}

// MustRead reads a file and fails the test on error.
func MustRead(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

// MustParseJSON parses a JSON file into v and fails the test on error.
func MustParseJSON(t *testing.T, path string, v any) {
	t.Helper()
	data := MustRead(t, path)
	if err := json.Unmarshal([]byte(data), v); err != nil {
		t.Fatalf("parse JSON %s: %v", path, err)
	}
}
