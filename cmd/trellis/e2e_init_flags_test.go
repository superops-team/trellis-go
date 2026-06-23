package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/superops-team/trellis-go/pkg/platform"
)

// TestE2E_InitWithPlatformFlag 测试 --platform 单平台
func TestE2E_InitWithPlatformFlag(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test", "--platform", "claude")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证只生成 claude 平台文件
	claudeDir := filepath.Join(repo, ".claude")
	if _, err := os.Stat(claudeDir); err != nil {
		t.Errorf("claude dir missing: %v", err)
	}

	// 验证不生成其他平台文件
	cursorDir := filepath.Join(repo, ".cursor")
	if _, err := os.Stat(cursorDir); err == nil {
		t.Error("cursor dir should not exist for single platform init")
	}
}

// TestE2E_InitWithMultiplePlatformsPlatform 测试多个 --platform
func TestE2E_InitWithMultiplePlatforms(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test",
		"--platform", "claude",
		"--platform", "cursor",
		"--platform", "codex")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证所有指定平台都生成
	for _, id := range []string{"claude", "cursor", "codex"} {
		p, ok := platform.NewRegistry().Get(id)
		if !ok {
			t.Fatalf("platform %s not found", id)
		}
		dir := filepath.Join(repo, p.ConfigDir)
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("platform dir %s missing: %v", id, err)
		}
	}

	// 验证不生成未指定的平台
	geminiDir := filepath.Join(repo, ".gemini")
	if _, err := os.Stat(geminiDir); err == nil {
		t.Error("gemini dir should not exist (not specified)")
	}
}

// TestE2E_InitWithDeveloperFlag 测试 --developer 模式
func TestE2E_InitWithDeveloperFlag(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "alice")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证 .developer 文件内容
	devPath := filepath.Join(repo, ".trellis", ".developer")
	data, err := os.ReadFile(devPath)
	if err != nil {
		t.Fatalf("read .developer: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "alice") {
		t.Errorf(".developer should contain 'alice', got: %s", content)
	}
}

// TestE2E_InitWithAllFlag 测试 --all flag
func TestE2E_InitWithAllFlag(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test", "--all")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证所有已注册平台的文件都生成
	registry := platform.NewRegistry()
	allPlatforms := registry.All()
	if len(allPlatforms) == 0 {
		t.Fatal("registry should have platforms")
	}

	for _, p := range allPlatforms {
		dir := filepath.Join(repo, p.ConfigDir)
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("platform dir %s missing: %v", p.ID, err)
		}
	}
}

// TestE2E_InitWithoutFlags 测试无 flag 默认行为
func TestE2E_InitWithoutFlags(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证默认生成 claude 平台（默认平台）
	claudeDir := filepath.Join(repo, ".claude")
	if _, err := os.Stat(claudeDir); err != nil {
		t.Errorf("claude dir should exist (default platform): %v", err)
	}

	// 验证不生成其他平台
	cursorDir := filepath.Join(repo, ".cursor")
	if _, err := os.Stat(cursorDir); err == nil {
		t.Error("cursor dir should not exist (not default)")
	}

	// 验证 .developer 文件存在（使用环境变量 USER 或默认 "developer"）
	devPath := filepath.Join(repo, ".trellis", ".developer")
	if _, err := os.Stat(devPath); err != nil {
		t.Errorf(".developer file should exist: %v", err)
	}
}

// TestE2E_InitOverwriteProtection 测试二次 init 幂等性
// 当前行为：已存在 .trellis/ 时，带 --platform 的 init 调用 addPlatforms 添加新平台
func TestE2E_InitOverwriteProtection(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	// 第一次 init
	_, stderr, err := runTrellis(t, repo, "init", "--developer", "alice", "--platform", "claude")
	if err != nil {
		t.Fatalf("first init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证 .trellis 存在
	trellisDir := filepath.Join(repo, ".trellis")
	if _, err := os.Stat(trellisDir); err != nil {
		t.Fatalf(".trellis dir missing after init: %v", err)
	}

	// 第二次 init 添加 cursor 平台（不应报错）
	_, stderr, err = runTrellis(t, repo, "init", "--platform", "cursor")
	if err != nil {
		t.Fatalf("second init (add platform) failed: %v\nstderr: %s", err, stderr)
	}

	// 验证 cursor 平台已添加
	cursorDir := filepath.Join(repo, ".cursor")
	if _, err := os.Stat(cursorDir); err != nil {
		t.Errorf("cursor dir should exist after add-platform: %v", err)
	}

	// 验证 claude 文件仍然存在
	claudeDir := filepath.Join(repo, ".claude")
	if _, err := os.Stat(claudeDir); err != nil {
		t.Errorf("claude dir should still exist: %v", err)
	}
}

// TestE2E_InitCreatesCorrectDirectoryStructure 验证完整的目录结构
func TestE2E_InitCreatesCorrectDirectoryStructure(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test", "--platform", "claude")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证 .trellis 目录结构
	trellisDir := filepath.Join(repo, ".trellis")
	expectedDirs := []string{
		"spec",
		"tasks",
		"workspace",
		".runtime/sessions",
	}
	for _, d := range expectedDirs {
		path := filepath.Join(trellisDir, d)
		if info, err := os.Stat(path); err != nil {
			t.Errorf("missing .trellis dir: %s (%v)", d, err)
		} else if !info.IsDir() {
			t.Errorf("expected directory: %s", d)
		}
	}

	expectedFiles := []string{
		"config.yaml",
		".version",
		"workflow.md",
	}
	for _, f := range expectedFiles {
		path := filepath.Join(trellisDir, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing .trellis file: %s (%v)", f, err)
		}
	}

	// 验证 claude 平台目录和 hook 文件
	claudeDir := filepath.Join(repo, ".claude")
	if _, err := os.Stat(claudeDir); err != nil {
		t.Errorf("missing .claude dir: %v", err)
	}

	for _, hook := range []string{"session-start.sh", "inject-context.sh", "inject-workflow-state.sh"} {
		path := filepath.Join(claudeDir, hook)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("missing hook file: %s (%v)", hook, err)
			continue
		}
		if info.Mode().Perm() != 0755 {
			t.Errorf("hook file %s permissions = %o, want 0755", hook, info.Mode().Perm())
		}
	}
}

// TestE2E_InitGeneratesValidConfigYAML 验证 config.yaml 内容正确
func TestE2E_InitGeneratesValidConfigYAML(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "bob", "--platform", "claude")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 解析 config.yaml
	cfgPath := filepath.Join(repo, ".trellis", "config.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read config.yaml: %v", err)
	}

	content := string(data)
	// 验证必填字段
	checks := []struct {
		field   string
		want    string
		missing string
	}{
		{"developer", "bob", "developer field should contain 'bob'"},
		{"packages", "[]", "packages field should exist"},
	}
	for _, c := range checks {
		if !strings.Contains(content, c.field) {
			t.Errorf("config.yaml missing field: %s", c.field)
		}
	}
	if !strings.Contains(content, "bob") {
		t.Errorf("config.yaml should contain developer 'bob', got: %s", content)
	}
}
