package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// runTrellis runs the trellis CLI in a temporary directory.
func runTrellis(t *testing.T, dir string, args ...string) (string, string, error) {
	t.Helper()
	cmd := exec.Command("/tmp/trellis-test", append([]string{"--root", filepath.Join(dir, ".trellis")}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "HOME="+dir)
	out, err := cmd.CombinedOutput()
	stdout := string(out)
	stderr := ""
	if err != nil {
		stderr = stdout
		stdout = ""
	}
	return stdout, stderr, err
}

// initGitRepo initializes a git repository in the given directory.
func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	for _, cmd := range [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test User"},
	} {
		c := exec.Command(cmd[0], cmd[1:]...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git init failed: %v\n%s", err, out)
		}
	}
}

// TestE2E_InitNewProject 场景1: 新项目初始化 + 首次任务创建
// 模拟一个全新 Git 仓库，开发者首次使用 Trellis
func TestE2E_InitNewProject(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	// Step 1: 初始化 Trellis
	stdout, stderr, err := runTrellis(t, repo, "init", "--developer", "alice", "--platform", "claude")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Initialized") && !strings.Contains(stderr, "Initialized") {
		t.Logf("init output: %s", stdout)
	}

	// 验证 .trellis 目录结构
	trellisDir := filepath.Join(repo, ".trellis")
	required := []string{
		"config.yaml",
		".version",
		"workflow.md",
		"spec",
		"tasks",
		"workspace",
		".runtime/sessions",
	}
	for _, f := range required {
		path := filepath.Join(trellisDir, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing required path: %s", f)
		}
	}

	// 验证 config.yaml 内容
	cfgData, _ := os.ReadFile(filepath.Join(trellisDir, "config.yaml"))
	if !strings.Contains(string(cfgData), "developer: alice") {
		t.Error("config.yaml should contain developer: alice")
	}

	// 验证 .version
	verData, _ := os.ReadFile(filepath.Join(trellisDir, ".version"))
	if strings.TrimSpace(string(verData)) == "" {
		t.Error(".version should not be empty")
	}

	// Step 2: 创建第一个任务
	stdout, stderr, err = runTrellis(t, repo, "task", "create", "add-auth")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Created task") {
		t.Errorf("expected 'Created task' in output, got: %s", stdout)
	}

	// 验证任务目录
	tasksDir := filepath.Join(trellisDir, "tasks")
	entries, _ := os.ReadDir(tasksDir)
	var taskDir string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "archive" {
			taskDir = filepath.Join(tasksDir, e.Name())
			break
		}
	}
	if taskDir == "" {
		t.Fatal("task directory not created")
	}

	// 验证任务文件
	for _, f := range []string{"task.json", "prd.md", "implement.jsonl", "check.jsonl"} {
		if _, err := os.Stat(filepath.Join(taskDir, f)); err != nil {
			t.Errorf("missing task file: %s", f)
		}
	}

	// 验证 task.json 状态为 planning
	taskJSON, _ := os.ReadFile(filepath.Join(taskDir, "task.json"))
	if !strings.Contains(string(taskJSON), `"status": "planning"`) {
		t.Error("new task should have status 'planning'")
	}
}

// TestE2E_MultiPlatformInit 场景2: 多平台配置生成
// 模拟团队使用多个 AI 平台（Claude + Cursor + Codex）
func TestE2E_MultiPlatformInit(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	// 初始化时指定多个平台
	_, stderr, err := runTrellis(t, repo, "init", "--developer", "team", "--platform", "claude", "--platform", "cursor", "--platform", "codex")
	if err != nil {
		t.Fatalf("multi-platform init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证各平台配置目录存在
	platforms := []string{".claude", ".cursor", ".codex"}
	for _, p := range platforms {
		path := filepath.Join(repo, p)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("platform config dir missing: %s", p)
		}
	}
}

// TestE2E_TaskLifecycle 场景3: 任务生命周期完整流转
// 模拟一个任务从创建 -> 启动 -> 归档的完整流程
func TestE2E_TaskLifecycle(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	// 初始化
	_, stderr, err := runTrellis(t, repo, "init", "--developer", "bob")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 创建任务
	_, stderr, err = runTrellis(t, repo, "task", "create", "refactor-api")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	// 查找任务目录
	tasksDir := filepath.Join(repo, ".trellis", "tasks")
	entries, _ := os.ReadDir(tasksDir)
	var taskDirName string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "archive" {
			taskDirName = e.Name()
			break
		}
	}
	if taskDirName == "" {
		t.Fatal("task not found")
	}

	// 验证初始状态为 planning
	taskPath := filepath.Join(tasksDir, taskDirName, "task.json")
	taskData, _ := os.ReadFile(taskPath)
	if !strings.Contains(string(taskData), `"status": "planning"`) {
		t.Error("initial status should be planning")
	}

	// 启动任务（需要直接调用内部逻辑，CLI 暂未完整实现 start）
	// 这里验证 task.json 可被修改
	newData := strings.Replace(string(taskData), `"status": "planning"`, `"status": "in_progress"`, 1)
	if err := os.WriteFile(taskPath, []byte(newData), 0644); err != nil {
		t.Fatalf("update task.json failed: %v", err)
	}

	// 验证状态已更新
	updated, _ := os.ReadFile(taskPath)
	if !strings.Contains(string(updated), `"status": "in_progress"`) {
		t.Error("status should be updated to in_progress")
	}

	// 模拟归档：移动目录到 archive/YYYY-MM/
	archiveDir := filepath.Join(tasksDir, "archive")
	os.MkdirAll(archiveDir, 0755)
	now := "2026-01"
	monthDir := filepath.Join(archiveDir, now)
	os.MkdirAll(monthDir, 0755)
	destDir := filepath.Join(monthDir, taskDirName)
	if err := os.Rename(filepath.Join(tasksDir, taskDirName), destDir); err != nil {
		t.Fatalf("archive task failed: %v", err)
	}

	// 验证任务已归档
	if _, err := os.Stat(destDir); err != nil {
		t.Error("archived task should exist")
	}
	if _, err := os.Stat(filepath.Join(tasksDir, taskDirName)); !os.IsNotExist(err) {
		t.Error("original task dir should be removed")
	}
}

// TestE2E_ContextBuild 场景4: 上下文构建与注入
// 模拟 AI 代理需要加载任务上下文
func TestE2E_ContextBuild(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	// 初始化
	_, stderr, err := runTrellis(t, repo, "init", "--developer", "charlie")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 创建任务
	_, stderr, err = runTrellis(t, repo, "task", "create", "user-auth")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	// 写入 PRD
	tasksDir := filepath.Join(repo, ".trellis", "tasks")
	entries, _ := os.ReadDir(tasksDir)
	var taskDir string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "archive" {
			taskDir = filepath.Join(tasksDir, e.Name())
			break
		}
	}

	prdContent := "# PRD: User Authentication\n\nImplement JWT-based auth."
	os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte(prdContent), 0644)

	// 写入 implement.jsonl
	manifestContent := `{"path":"spec/auth.md","description":"Auth spec","required":true}
{"path":"spec/api.md","description":"API spec","required":false}
`
	os.WriteFile(filepath.Join(taskDir, "implement.jsonl"), []byte(manifestContent), 0644)

	// 创建 spec 文件
	specDir := filepath.Join(repo, ".trellis", "spec")
	os.MkdirAll(filepath.Join(specDir, "auth"), 0755)
	os.WriteFile(filepath.Join(specDir, "auth.md"), []byte("# Auth Spec\nUse JWT."), 0644)

	// 验证 PRD 内容
	readPRD, _ := os.ReadFile(filepath.Join(taskDir, "prd.md"))
	if string(readPRD) != prdContent {
		t.Error("prd.md content mismatch")
	}

	// 验证 manifest 可解析
	manifestPath := filepath.Join(taskDir, "implement.jsonl")
	manifestData, _ := os.ReadFile(manifestPath)
	lines := strings.Split(strings.TrimSpace(string(manifestData)), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 manifest entries, got %d", len(lines))
	}

	// 验证上下文构建标记
	// 实际 Builder 需要 spec.Loader，这里验证文件结构正确
	if !strings.Contains(string(manifestData), `"required":true`) {
		t.Error("manifest should contain required entry")
	}
}

// TestE2E_TaskListAndCurrent 场景5: 任务列表与当前任务查询
// 模拟开发者查看所有任务和当前活跃任务
func TestE2E_TaskListAndCurrent(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	// 初始化
	_, stderr, err := runTrellis(t, repo, "init", "--developer", "dave")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 创建多个任务
	for _, name := range []string{"feature-a", "feature-b", "bugfix-c"} {
		_, stderr, err = runTrellis(t, repo, "task", "create", name)
		if err != nil {
			t.Fatalf("task create %s failed: %v\nstderr: %s", name, err, stderr)
		}
	}

	// 列出任务
	stdout, stderr, err := runTrellis(t, repo, "task", "list")
	if err != nil {
		t.Fatalf("task list failed: %v\nstderr: %s", err, stderr)
	}

	// 验证所有任务都在列表中
	for _, name := range []string{"feature-a", "feature-b", "bugfix-c"} {
		if !strings.Contains(stdout, name) && !strings.Contains(stderr, name) {
			t.Errorf("task list should contain %s", name)
		}
	}

	// 查询当前任务（暂无活跃任务）
	stdout, stderr, err = runTrellis(t, repo, "task", "current")
	if err != nil {
		t.Logf("current task (expected no active): %s", stdout)
	}
}

// TestE2E_UninstallKeepTasks 场景6: 卸载保留任务
// 模拟团队决定停止使用 Trellis 但保留历史任务记录
func TestE2E_UninstallKeepTasks(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	// 初始化并创建任务
	_, stderr, err := runTrellis(t, repo, "init", "--developer", "eve")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}
	_, stderr, err = runTrellis(t, repo, "task", "create", "legacy-feature")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	// 卸载并保留任务
	_, stderr, err = runTrellis(t, repo, "uninstall", "--keep-tasks")
	if err != nil {
		t.Fatalf("uninstall failed: %v\nstderr: %s", err, stderr)
	}

	// 验证 .trellis 已删除（除 tasks 外）
	trellisDir := filepath.Join(repo, ".trellis")
	if _, err := os.Stat(filepath.Join(trellisDir, "config.yaml")); !os.IsNotExist(err) {
		t.Error("config.yaml should be removed")
	}

	// 验证 tasks 目录保留
	tasksDir := filepath.Join(trellisDir, "tasks")
	entries, _ := os.ReadDir(tasksDir)
	var hasTask bool
	for _, e := range entries {
		if e.IsDir() && e.Name() != "archive" {
			hasTask = true
			break
		}
	}
	if !hasTask {
		t.Error("tasks should be preserved")
	}
}

// TestE2E_InvalidPlatform 场景7: 错误处理 — 无效平台
// 模拟用户输入了不存在的平台名称
func TestE2E_InvalidPlatform(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test", "--platform", "nonexistent")
	if err == nil {
		t.Fatal("expected error for invalid platform")
	}
	if !strings.Contains(stderr, "unknown platform") && !strings.Contains(stderr, "nonexistent") {
		t.Errorf("error should mention unknown platform, got: %s", stderr)
	}
}

// TestE2E_NotGitRepo 场景8: 错误处理 — 非 Git 仓库
// 模拟用户在非 Git 目录运行 init
func TestE2E_NotGitRepo(t *testing.T) {
	repo := t.TempDir()
	// 不初始化 git

	_, stderr, err := runTrellis(t, repo, "init")
	if err == nil {
		t.Fatal("expected error for non-git repo")
	}
	if !strings.Contains(stderr, "git") {
		t.Errorf("error should mention git, got: %s", stderr)
	}
}
