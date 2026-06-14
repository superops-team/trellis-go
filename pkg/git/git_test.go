package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestClient_StatusInRepository(t *testing.T) {
	repo := initTestRepo(t)
	client := NewClient(repo)

	if !client.IsRepo() {
		t.Fatal("expected initialized git repository")
	}

	branch, err := client.CurrentBranch()
	if err != nil {
		t.Fatalf("CurrentBranch failed: %v", err)
	}
	if branch == "" {
		t.Fatal("expected non-empty current branch")
	}

	status, err := client.Status()
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if status.Branch != branch {
		t.Fatalf("Status branch = %q, want %q", status.Branch, branch)
	}
	if status.IsDirty {
		t.Fatal("new repository should be clean after initial commit")
	}

	if err := os.WriteFile(filepath.Join(repo, "untracked.txt"), []byte("dirty"), 0644); err != nil {
		t.Fatalf("write untracked file: %v", err)
	}
	hasChanges, err := client.HasChanges()
	if err != nil {
		t.Fatalf("HasChanges failed: %v", err)
	}
	if !hasChanges {
		t.Fatal("expected untracked file to make working tree dirty")
	}
}

func TestClient_CommandFailureIncludesGitContext(t *testing.T) {
	dir := t.TempDir()
	client := NewClient(dir)

	if client.IsRepo() {
		t.Fatal("temporary directory should not be a git repository")
	}
	_, err := client.CurrentBranch()
	if err == nil {
		t.Fatal("expected CurrentBranch to fail outside a git repository")
	}
	if !strings.Contains(err.Error(), "git") || !strings.Contains(err.Error(), "branch") {
		t.Fatalf("error should include git command context, got: %v", err)
	}
}

func initTestRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "test@example.com")
	runGit(t, repo, "config", "user.name", "Test User")
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("# test\n"), 0644); err != nil {
		t.Fatalf("write README: %v", err)
	}
	runGit(t, repo, "add", "README.md")
	runGit(t, repo, "commit", "-m", "initial commit")
	return repo
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}
