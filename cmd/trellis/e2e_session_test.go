package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestE2E_SessionRecordAndList 测试 session 记录 + 列表
func TestE2E_SessionRecordAndList(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 记录第一个 session
	stdout, stderr, err := runTrellis(t, repo, "hook", "record-session",
		"--title", "Implement login",
		"--task", "login-feature",
		"--commits", "abc1234,def5678",
		"--summary", "Implemented login page with JWT auth")
	if err != nil {
		t.Fatalf("record-session failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Session recorded") {
		t.Errorf("expected 'Session recorded' in output, got: %s", stdout)
	}

	// 记录第二个 session
	_, stderr, err = runTrellis(t, repo, "hook", "record-session",
		"--title", "Add tests",
		"--task", "login-feature",
		"--commits", "ghi9012",
		"--summary", "Added unit tests for auth")
	if err != nil {
		t.Fatalf("second record-session failed: %v\nstderr: %s", err, stderr)
	}

	// 列出 sessions
	stdout, stderr, err = runTrellis(t, repo, "hook", "list-sessions")
	if err != nil {
		t.Fatalf("list-sessions failed: %v\nstderr: %s", err, stderr)
	}

	// 验证两个 session 都在列表中
	if !strings.Contains(stdout, "Implement login") {
		t.Errorf("list should contain 'Implement login', got: %s", stdout)
	}
	if !strings.Contains(stdout, "Add tests") {
		t.Errorf("list should contain 'Add tests', got: %s", stdout)
	}

	// 验证 journal 文件存在
	workspaceDir := filepath.Join(repo, ".trellis", "workspace", "test")
	journalPath := filepath.Join(workspaceDir, "journal-1.md")
	journalData, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}
	if !strings.Contains(string(journalData), "login-feature") {
		t.Error("journal should contain task ID")
	}
	if !strings.Contains(string(journalData), "abc1234") {
		t.Error("journal should contain commit hash")
	}
}

// TestE2E_SessionListSearch 测试 session 列表搜索过滤
func TestE2E_SessionListSearch(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 记录多个 session
	for _, s := range []struct {
		title string
		task  string
	}{
		{"Implement login", "login-feature"},
		{"Add tests", "login-feature"},
		{"Fix bug", "bugfix-123"},
	} {
		_, stderr, err := runTrellis(t, repo, "hook", "record-session",
			"--title", s.title,
			"--task", s.task)
		if err != nil {
			t.Fatalf("record-session %s failed: %v\nstderr: %s", s.title, err, stderr)
		}
	}

	// 列出所有 sessions
	stdout, stderr, err := runTrellis(t, repo, "hook", "list-sessions")
	if err != nil {
		t.Fatalf("list-sessions failed: %v\nstderr: %s", err, stderr)
	}

	// 验证所有 session 都在列表中
	for _, want := range []string{"Implement login", "Add tests", "Fix bug"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("list should contain %q, got: %s", want, stdout)
		}
	}
}

// TestE2E_SessionListEmpty 测试空 session 列表
func TestE2E_SessionListEmpty(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 没有记录任何 session 时列出
	stdout, stderr, err := runTrellis(t, repo, "hook", "list-sessions")
	if err != nil {
		t.Fatalf("list-sessions failed: %v\nstderr: %s", err, stderr)
	}

	if !strings.Contains(stdout, "No sessions recorded") {
		t.Errorf("empty list should show 'No sessions recorded', got: %s", stdout)
	}
}

// TestE2E_SessionJournalFileFormat 测试 journal 文件格式
func TestE2E_SessionJournalFileFormat(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 记录一个 session
	_, stderr, err = runTrellis(t, repo, "hook", "record-session",
		"--title", "Format test",
		"--task", "test-task",
		"--commits", "abc1234",
		"--summary", "Testing journal format")
	if err != nil {
		t.Fatalf("record-session failed: %v\nstderr: %s", err, stderr)
	}

	// 验证 journal 文件格式
	workspaceDir := filepath.Join(repo, ".trellis", "workspace", "test")
	journalPath := filepath.Join(workspaceDir, "journal-1.md")
	journalData, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}

	content := string(journalData)

	// 验证格式：以 # Session 开头
	if !strings.HasPrefix(content, "# Session") {
		t.Errorf("journal should start with '# Session', got: %s", content[:20])
	}

	// 验证包含任务 ID
	if !strings.Contains(content, "test-task") {
		t.Errorf("journal should contain task ID, got: %s", content)
	}

	// 验证包含 commit
	if !strings.Contains(content, "abc1234") {
		t.Errorf("journal should contain commit hash, got: %s", content)
	}

	// 验证包含 summary
	if !strings.Contains(content, "Testing journal format") {
		t.Errorf("journal should contain summary, got: %s", content)
	}

	// 验证 index.json 存在且格式正确
	idxPath := filepath.Join(workspaceDir, "index.json")
	idxData, err := os.ReadFile(idxPath)
	if err != nil {
		t.Fatalf("read index.json: %v", err)
	}
	if !strings.Contains(string(idxData), `"title": "Format test"`) {
		t.Errorf("index.json should contain session title, got: %s", idxData)
	}
}
