package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/superops-team/trellis-go/pkg/config"
)

// TestE2E_ConfigInitWithPackages 测试 init 带 packages 配置
func TestE2E_ConfigInitWithPackages(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 手动添加 packages 到 config.yaml
	cfgPath := filepath.Join(repo, ".trellis", "config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	cfg.Packages = []string{"auth", "user", "billing"}
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	// 重新加载验证
	cfg2, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("reload config: %v", err)
	}
	if len(cfg2.Packages) != 3 {
		t.Errorf("expected 3 packages, got %d", len(cfg2.Packages))
	}
	if cfg2.Packages[0] != "auth" {
		t.Errorf("expected first package 'auth', got %q", cfg2.Packages[0])
	}
}

// TestE2E_ConfigInitWithHooks 测试 init 带 hooks 配置
func TestE2E_ConfigInitWithHooks(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 手动添加 hooks 到 config.yaml
	cfgPath := filepath.Join(repo, ".trellis", "config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	cfg.Hooks = map[string]string{
		"pre-commit":  "echo 'pre-commit hook'",
		"post-commit": "echo 'post-commit hook'",
	}
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	// 重新加载验证
	cfg2, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("reload config: %v", err)
	}
	if len(cfg2.Hooks) != 2 {
		t.Errorf("expected 2 hooks, got %d", len(cfg2.Hooks))
	}
	if cfg2.Hooks["pre-commit"] != "echo 'pre-commit hook'" {
		t.Errorf("unexpected pre-commit hook: %q", cfg2.Hooks["pre-commit"])
	}
}

// TestE2E_ConfigInitWithCodex 测试 init 带 codex 配置
func TestE2E_ConfigInitWithCodex(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test", "--platform", "codex")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证默认 codex 配置
	cfgPath := filepath.Join(repo, ".trellis", "config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Codex.DispatchMode != "inline" {
		t.Errorf("expected default codex.dispatch_mode 'inline', got %q", cfg.Codex.DispatchMode)
	}

	// 修改为 sub-agent 模式
	cfg.Codex.DispatchMode = "sub-agent"
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	// 重新加载验证
	cfg2, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("reload config: %v", err)
	}
	if cfg2.Codex.DispatchMode != "sub-agent" {
		t.Errorf("expected codex.dispatch_mode 'sub-agent', got %q", cfg2.Codex.DispatchMode)
	}
}

// TestE2E_ConfigValidateInvalid 测试无效配置验证
func TestE2E_ConfigValidateInvalid(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 设置无效的 codex.dispatch_mode
	cfgPath := filepath.Join(repo, ".trellis", "config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	cfg.Codex.DispatchMode = "invalid-mode"
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	// 验证应该失败
	cfg2, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg2.Validate(); err == nil {
		t.Error("expected validation error for invalid codex.dispatch_mode")
	} else if !strings.Contains(err.Error(), "codex.dispatch_mode") {
		t.Errorf("error should mention codex.dispatch_mode, got: %v", err)
	}
}

// TestE2E_ConfigValidateMissing 测试缺失配置处理
func TestE2E_ConfigValidateMissing(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 测试加载不存在的配置
	nonExistentPath := filepath.Join(repo, ".trellis", "config-not-found.yaml")
	_, err = config.Load(nonExistentPath)
	if err == nil {
		t.Error("expected error for non-existent config file")
	} else if !strings.Contains(err.Error(), "config not found") {
		t.Errorf("error should mention 'config not found', got: %v", err)
	}

	// 测试空 packages 列表（应该通过）
	cfgPath := filepath.Join(repo, ".trellis", "config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Packages = []string{} // 空列表
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	cfg2, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg2.Validate(); err != nil {
		t.Errorf("empty packages list should be valid, got: %v", err)
	}
}
