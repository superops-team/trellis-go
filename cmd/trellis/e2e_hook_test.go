package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/superops-team/trellis-go/pkg/platform"
)

// TestE2E_HookInjectWorkflowState 测试 inject-workflow-state hook 输出
func TestE2E_HookInjectWorkflowState(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 创建任务并启动
	_, stderr, err = runTrellis(t, repo, "task", "create", "workflow-test")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	tasksDir := filepath.Join(repo, ".trellis", "tasks")
	taskDirName := firstTaskDirName(t, tasksDir)
	if err := os.WriteFile(filepath.Join(tasksDir, taskDirName, "prd.md"), []byte("# PRD\nTest workflow."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	_, stderr, err = runTrellis(t, repo, "task", "start", "workflow-test")
	if err != nil {
		t.Fatalf("task start failed: %v\nstderr: %s", err, stderr)
	}

	// 验证各阶段 inject-workflow-state 输出
	for _, tc := range []struct {
		state string
		hint  string
	}{
		{"plan", "PLAN"},
		{"implement", "IMPLEMENT"},
		{"verify", "VERIFY"},
		{"finish", "FINISH"},
	} {
		stdout, stderr, err := runTrellis(t, repo, "hook", "inject-workflow-state", "--state", tc.state)
		if err != nil {
			t.Fatalf("inject-workflow-state %s failed: %v\nstderr: %s", tc.state, err, stderr)
		}
		if !strings.Contains(stdout, tc.hint) {
			t.Errorf("inject-workflow-state %s output should contain %q, got: %s", tc.state, tc.hint, stdout)
		}
		// 验证输出包含 workflow-state 标签
		if !strings.Contains(stdout, "<workflow-state>") {
			t.Errorf("output should contain <workflow-state> tags, got: %s", stdout)
		}
	}
}

// TestE2E_HookSessionStart 测试 session-start hook 输出
func TestE2E_HookSessionStart(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 创建任务
	_, stderr, err = runTrellis(t, repo, "task", "create", "session-test")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	// 执行 session-start hook
	stdout, stderr, err := runTrellis(t, repo, "hook", "session-start")
	if err != nil {
		t.Fatalf("session-start failed: %v\nstderr: %s", err, stderr)
	}

	// 验证输出包含关键信息
	for _, want := range []string{"Trellis session started", "RepoRoot: ", "TrellisDir: "} {
		if !strings.Contains(stdout, want) {
			t.Errorf("session-start output should contain %q, got: %s", want, stdout)
		}
	}
}

// TestE2E_HookInjectContext 测试 inject-context hook 输出
func TestE2E_HookInjectContext(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 创建 spec 文件
	specDir := filepath.Join(repo, ".trellis", "spec")
	os.MkdirAll(specDir, 0755)
	if err := os.WriteFile(filepath.Join(specDir, "auth.md"), []byte("# Auth Spec\nUse JWT authentication."), 0644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	// 创建任务
	_, stderr, err = runTrellis(t, repo, "task", "create", "context-test")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	tasksDir := filepath.Join(repo, ".trellis", "tasks")
	taskDirName := firstTaskDirName(t, tasksDir)

	// 写入 PRD
	if err := os.WriteFile(filepath.Join(tasksDir, taskDirName, "prd.md"), []byte("# PRD\nImplement auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}

	// 添加 spec 到上下文
	_, stderr, err = runTrellis(t, repo, "context", "add", "spec/auth.md", "--task", "context-test", "--phase", "implement", "--required", "--description", "Auth spec")
	if err != nil {
		t.Fatalf("context add failed: %v\nstderr: %s", err, stderr)
	}

	// 执行 inject-context hook
	stdout, stderr, err := runTrellis(t, repo, "hook", "inject-context", "--task", "context-test", "--phase", "implement")
	if err != nil {
		t.Fatalf("inject-context failed: %v\nstderr: %s", err, stderr)
	}

	// 验证输出包含注入标记和 spec 内容
	for _, want := range []string{"<!-- trellis-hook-injected -->", "# PRD\nImplement auth.", "# Auth Spec\nUse JWT authentication."} {
		if !strings.Contains(stdout, want) {
			t.Errorf("inject-context output should contain %q, got: %s", want, stdout)
		}
	}
}

// TestE2E_HookFilesForAllPlatforms 测试多平台 hook 文件生成
func TestE2E_HookFilesForAllPlatforms(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	// 初始化多个平台
	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test",
		"--platform", "claude",
		"--platform", "cursor",
		"--platform", "codex",
		"--platform", "gemini")
	if err != nil {
		t.Fatalf("multi-platform init failed: %v\nstderr: %s", err, stderr)
	}

	// 验证 push-based 平台有 hook 目录和文件
	pushPlatforms := []string{"claude", "cursor"}
	for _, id := range pushPlatforms {
		p, ok := platform.NewRegistry().Get(id)
		if !ok {
			t.Fatalf("platform %s not found", id)
		}
		hookDir := filepath.Join(repo, p.ConfigDir)
		if _, err := os.Stat(hookDir); err != nil {
			t.Errorf("platform dir missing: %s", p.ConfigDir)
			continue
		}

		for _, hook := range []string{"session-start.sh", "inject-context.sh", "inject-workflow-state.sh"} {
			hookPath := filepath.Join(hookDir, hook)
			info, err := os.Stat(hookPath)
			if err != nil {
				t.Errorf("missing hook file %s for platform %s: %v", hook, id, err)
				continue
			}
			// 验证 hook 文件有执行权限
			if info.Mode().Perm() != 0755 {
				t.Errorf("hook file %s/%s permissions = %o, want 0755", p.ConfigDir, hook, info.Mode().Perm())
			}
			// 验证 hook 文件内容调用 trellis hook
			data, _ := os.ReadFile(hookPath)
			if !strings.Contains(string(data), "trellis hook") && !strings.Contains(string(data), "hook") {
				t.Errorf("hook file %s/%s should call trellis hook, got: %s", p.ConfigDir, hook, data)
			}
		}
	}

	// 验证 pull-based 平台（codex）有 agent 定义文件
	codexP, _ := platform.NewRegistry().Get("codex")
	codexAgentDir := filepath.Join(repo, codexP.ConfigDir, "agents")
	if _, err := os.Stat(codexAgentDir); err != nil {
		t.Errorf("codex agents dir missing: %v", err)
	} else {
		agentDefPath := filepath.Join(codexAgentDir, "trellis-implement.toml")
		if _, err := os.Stat(agentDefPath); err != nil {
			t.Errorf("missing agent def file for codex: %v", err)
		}
	}

	// 验证 pull-based 平台（gemini）没有 hook 文件
	geminiP, _ := platform.NewRegistry().Get("gemini")
	geminiDir := filepath.Join(repo, geminiP.ConfigDir)
	if _, err := os.Stat(geminiDir); err != nil {
		t.Errorf("gemini dir missing: %v", err)
	} else {
		// Gemini has HasHooks=false, so no hook scripts should exist
		for _, hook := range []string{"session-start.sh", "inject-context.sh", "inject-workflow-state.sh"} {
			hookPath := filepath.Join(geminiDir, hook)
			if _, err := os.Stat(hookPath); err == nil {
				t.Errorf("gemini should not have hook file %s (HasHooks=false)", hook)
			}
		}
	}
}
