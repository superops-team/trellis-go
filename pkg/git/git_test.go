package git

import (
	"fmt"
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

func TestClient_AddAndCommit(t *testing.T) {
	repo := initTestRepo(t)
	client := NewClient(repo)

	// Create a new file
	if err := os.WriteFile(filepath.Join(repo, "new.txt"), []byte("new content"), 0644); err != nil {
		t.Fatalf("write new file: %v", err)
	}

	// Add and commit
	if err := client.Add("new.txt"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if err := client.Commit("add new.txt"); err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	// Verify clean status
	hasChanges, err := client.HasChanges()
	if err != nil {
		t.Fatalf("HasChanges failed: %v", err)
	}
	if hasChanges {
		t.Fatal("expected clean after commit")
	}
}

func TestClient_SafeCommit(t *testing.T) {
	repo := initTestRepo(t)
	client := NewClient(repo)

	// Create files
	if err := os.WriteFile(filepath.Join(repo, "a.txt"), []byte("a"), 0644); err != nil {
		t.Fatalf("write a.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "b.txt"), []byte("b"), 0644); err != nil {
		t.Fatalf("write b.txt: %v", err)
	}

	// SafeCommit with specific patterns
	if err := client.SafeCommit("add a and b", []string{"a.txt", "b.txt"}); err != nil {
		t.Fatalf("SafeCommit failed: %v", err)
	}

	// Verify clean
	hasChanges, err := client.HasChanges()
	if err != nil {
		t.Fatalf("HasChanges failed: %v", err)
	}
	if hasChanges {
		t.Fatal("expected clean after SafeCommit")
	}
}

func TestClient_SafeCommitNoPatterns(t *testing.T) {
	repo := initTestRepo(t)
	client := NewClient(repo)

	if err := os.WriteFile(filepath.Join(repo, "c.txt"), []byte("c"), 0644); err != nil {
		t.Fatalf("write c.txt: %v", err)
	}

	// SafeCommit with empty patterns should commit all staged files
	// First add the file manually
	if err := client.Add("c.txt"); err != nil {
		t.Fatalf("Add c.txt: %v", err)
	}
	if err := client.SafeCommit("add c", nil); err != nil {
		t.Fatalf("SafeCommit with nil patterns failed: %v", err)
	}

	hasChanges, err := client.HasChanges()
	if err != nil {
		t.Fatalf("HasChanges failed: %v", err)
	}
	if hasChanges {
		t.Fatal("expected clean after SafeCommit with nil patterns")
	}
}

func TestClient_RecentCommits(t *testing.T) {
	repo := initTestRepo(t)
	client := NewClient(repo)

	// Make additional commits
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("file%d.txt", i)
		if err := os.WriteFile(filepath.Join(repo, name), []byte(name), 0644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
		if err := client.Add(name); err != nil {
			t.Fatalf("Add %s: %v", name, err)
		}
		if err := client.Commit(fmt.Sprintf("add %s", name)); err != nil {
			t.Fatalf("Commit %s: %v", name, err)
		}
	}

	commits, err := client.RecentCommits(2)
	if err != nil {
		t.Fatalf("RecentCommits failed: %v", err)
	}
	if len(commits) != 2 {
		t.Errorf("expected 2 recent commits, got %d", len(commits))
	}
	if len(commits) > 0 && commits[0].Hash == "" {
		t.Error("commit hash should not be empty")
	}
	if len(commits) > 0 && commits[0].Message == "" {
		t.Error("commit message should not be empty")
	}
}

func TestClient_RecentCommitsEmpty(t *testing.T) {
	repo := initTestRepo(t)
	client := NewClient(repo)

	commits, err := client.RecentCommits(0)
	if err != nil {
		t.Fatalf("RecentCommits(0) failed: %v", err)
	}
	if commits != nil {
		t.Errorf("expected nil for 0 commits, got %d", len(commits))
	}
}

func TestClient_RemoteURL(t *testing.T) {
	repo := initTestRepo(t)
	client := NewClient(repo)

	// No remote configured
	_, err := client.RemoteURL()
	if err == nil {
		t.Fatal("expected error for repo without remote")
	}
}

func TestClient_IsRepoFalse(t *testing.T) {
	dir := t.TempDir()
	client := NewClient(dir)

	if client.IsRepo() {
		t.Fatal("temp dir should not be a git repo")
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
