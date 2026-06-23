package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/superops-team/trellis-go/pkg/config"
	"github.com/superops-team/trellis-go/pkg/platform"
)

// TestE2E_CodexInitWithConfig 测试 init 带 codex 配置
func TestE2E_CodexInitWithConfig(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test", "--platform", "codex")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证 .codex 目录存在
	codexDir := filepath.Join(repo, ".codex")
	if _, err := os.Stat(codexDir); err != nil {
		t.Errorf(".codex dir missing: %v", err)
	}

	// 验证 config.yaml 中有 codex 配置
	cfgPath := filepath.Join(repo, ".trellis", "config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Codex.DispatchMode != "inline" {
		t.Errorf("expected default codex.dispatch_mode 'inline', got %q", cfg.Codex.DispatchMode)
	}
}

// TestE2E_CodexDispatchInline 测试 inline dispatch 模式
func TestE2E_CodexDispatchInline(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test", "--platform", "codex")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 设置 inline 模式
	cfgPath := filepath.Join(repo, ".trellis", "config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	cfg.Codex.DispatchMode = "inline"
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	// 验证配置正确
	cfg2, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("reload config: %v", err)
	}
	if cfg2.Codex.DispatchMode != "inline" {
		t.Errorf("expected inline mode, got %q", cfg2.Codex.DispatchMode)
	}
}

// TestE2E_CodexDispatchSubAgent 测试 sub-agent dispatch 模式
func TestE2E_CodexDispatchSubAgent(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test", "--platform", "codex")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 设置 sub-agent 模式
	cfgPath := filepath.Join(repo, ".trellis", "config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	cfg.Codex.DispatchMode = "sub-agent"
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	// 验证配置正确
	cfg2, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("reload config: %v", err)
	}
	if cfg2.Codex.DispatchMode != "sub-agent" {
		t.Errorf("expected sub-agent mode, got %q", cfg2.Codex.DispatchMode)
	}
}

// TestE2E_CodexDispatchInvalid 测试无效 dispatch 模式
func TestE2E_CodexDispatchInvalid(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test", "--platform", "codex")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 设置无效的 dispatch 模式
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
		t.Error("expected validation error for invalid dispatch mode")
	} else if !strings.Contains(err.Error(), "dispatch_mode") {
		t.Errorf("error should mention dispatch_mode, got: %v", err)
	}
}

// TestE2E_CodexAgentDef 测试 agent 定义文件生成
func TestE2E_CodexAgentDef(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test", "--platform", "codex")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证 codex 平台目录结构
	codexP, _ := platform.NewRegistry().Get("codex")
	codexDir := filepath.Join(repo, codexP.ConfigDir)
	if _, err := os.Stat(codexDir); err != nil {
		t.Errorf("codex dir missing: %v", err)
	}

	// 验证 agents 目录
	agentsDir := filepath.Join(codexDir, "agents")
	if _, err := os.Stat(agentsDir); err != nil {
		t.Errorf("codex agents dir missing: %v", err)
	}

	// 验证 agent 定义文件存在
	expectedAgents := []string{"trellis-implement.toml", "trellis-check.toml", "trellis-research.toml"}
	for _, agentFile := range expectedAgents {
		agentPath := filepath.Join(agentsDir, agentFile)
		if _, err := os.Stat(agentPath); err != nil {
			t.Errorf("agent file %s missing: %v", agentFile, err)
		} else {
			// 验证文件内容是有效的 TOML
			data, _ := os.ReadFile(agentPath)
			content := string(data)
			if !strings.Contains(content, "name = ") {
				t.Errorf("agent file %s should contain 'name' field", agentFile)
			}
			if !strings.Contains(content, "description = ") {
				t.Errorf("agent file %s should contain 'description' field", agentFile)
			}
		}
	}

	// 验证 AGENTS.md 存在（codex 特有）
	agentsMDPath := filepath.Join(codexDir, "AGENTS.md")
	if _, err := os.Stat(agentsMDPath); err != nil {
		t.Errorf("AGENTS.md missing for codex: %v", err)
	}
}
