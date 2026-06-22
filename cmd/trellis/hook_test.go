package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func executeRootCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
}

func TestRootCommandExposesHookNamespace(t *testing.T) {
	output, err := executeRootCommand(t, "hook", "--help")
	if err != nil {
		t.Fatalf("hook help failed: %v\n%s", err, output)
	}
	for _, want := range []string{"session-start", "inject-context", "inject-workflow-state"} {
		if !strings.Contains(output, want) {
			t.Fatalf("hook help should list %q, got: %s", want, output)
		}
	}
}

func TestHookInjectWorkflowStatePrintsPrompt(t *testing.T) {
	output, err := executeRootCommand(t, "hook", "inject-workflow-state", "--state", "implement")
	if err != nil {
		t.Fatalf("inject workflow state failed: %v\n%s", err, output)
	}
	if !strings.Contains(output, "IMPLEMENT phase") {
		t.Fatalf("workflow prompt should describe implement phase, got: %s", output)
	}
}

func TestHookInjectWorkflowStateRejectsUnknownState(t *testing.T) {
	output, err := executeRootCommand(t, "hook", "inject-workflow-state", "--state", "unknown")
	if err != nil {
		t.Fatalf("unknown state should not error, got: %v", err)
	}
	if !strings.Contains(output, "Refer to workflow.md") {
		t.Fatalf("unknown state should return generic fallback, got: %s", output)
	}
}

func TestHookSessionStartPrintsDeterministicContext(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".trellis", "spec"), 0755); err != nil {
		t.Fatalf("create trellis dirs: %v", err)
	}

	output, err := executeRootCommand(t, "--root", filepath.Join(root, ".trellis"), "hook", "session-start")
	if err != nil {
		t.Fatalf("session-start failed: %v\n%s", err, output)
	}
	for _, want := range []string{"Trellis session started", "RepoRoot: ", "TrellisDir: "} {
		if !strings.Contains(output, want) {
			t.Fatalf("session-start output should contain %q, got: %s", want, output)
		}
	}
}

func TestHookInjectContextRequiresTaskForImplement(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".trellis", "spec"), 0755); err != nil {
		t.Fatalf("create trellis dirs: %v", err)
	}

	output, err := executeRootCommand(t, "--root", filepath.Join(root, ".trellis"), "hook", "inject-context", "--phase", "implement")
	if err == nil {
		t.Fatalf("expected missing task to fail, got: %s", output)
	}
	if !strings.Contains(err.Error(), "--task is required") {
		t.Fatalf("error should explain task requirement, got: %v", err)
	}
}

func TestHookInjectContextResearchDoesNotRequireTask(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".trellis", "spec"), 0755); err != nil {
		t.Fatalf("create trellis dirs: %v", err)
	}

	output, err := executeRootCommand(t, "--root", filepath.Join(root, ".trellis"), "hook", "inject-context", "--phase", "research")
	if err != nil {
		t.Fatalf("research hook context failed: %v\n%s", err, output)
	}
	if !strings.Contains(output, "<!-- trellis-hook-injected -->") {
		t.Fatalf("research hook output should contain injection marker, got: %s", output)
	}
}

func TestHookInjectContextFailsWhenImplementPRDIsBlank(t *testing.T) {
	root := t.TempDir()
	tasksDir := filepath.Join(root, ".trellis", "tasks")
	taskDir := filepath.Join(tasksDir, "01-01-user-auth")
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "task.json"), []byte(`{
  "id": "user-auth",
  "name": "user-auth",
  "status": "planning",
  "created_at": "2026-06-15T00:00:00Z",
  "updated_at": "2026-06-15T00:00:00Z"
}`), 0644); err != nil {
		t.Fatalf("write task.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte(" \n\t"), 0644); err != nil {
		t.Fatalf("write blank prd: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "implement.jsonl"), nil, 0644); err != nil {
		t.Fatalf("write implement manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "check.jsonl"), nil, 0644); err != nil {
		t.Fatalf("write check manifest: %v", err)
	}

	output, err := executeRootCommand(t, "--root", filepath.Join(root, ".trellis"), "hook", "inject-context", "--task", "user-auth", "--phase", "implement")
	if err == nil {
		t.Fatalf("expected blank PRD to fail, got: %s", output)
	}
	if !strings.Contains(err.Error(), "PRD is required for task user-auth") {
		t.Fatalf("error should mention missing PRD goal, got: %v", err)
	}
}

func TestHookInjectContextUsesCheckBuilderPath(t *testing.T) {
	root := t.TempDir()
	taskDir := filepath.Join(root, ".trellis", "tasks", "01-01-user-auth")
	if err := os.MkdirAll(filepath.Join(root, ".trellis", "spec"), 0755); err != nil {
		t.Fatalf("create spec dir: %v", err)
	}
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "task.json"), []byte(`{
  "id": "user-auth",
  "name": "user-auth",
  "status": "planning",
  "created_at": "2026-06-15T00:00:00Z",
  "updated_at": "2026-06-15T00:00:00Z"
}`), 0644); err != nil {
		t.Fatalf("write task.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nCheck auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".trellis", "spec", "auth-check.md"), []byte("# Auth Check\nVerify JWT."), 0644); err != nil {
		t.Fatalf("write check file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "check.jsonl"), []byte("{\"path\":\"spec/auth-check.md\",\"required\":true}\n"), 0644); err != nil {
		t.Fatalf("write check manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "implement.jsonl"), nil, 0644); err != nil {
		t.Fatalf("write implement manifest: %v", err)
	}

	output, err := executeRootCommand(t, "--root", filepath.Join(root, ".trellis"), "hook", "inject-context", "--task", "user-auth", "--phase", "check")
	if err != nil {
		t.Fatalf("check hook context failed: %v\n%s", err, output)
	}
	for _, want := range []string{"<!-- trellis-hook-injected -->", "# PRD\nCheck auth.", "# Auth Check\nVerify JWT."} {
		if !strings.Contains(output, want) {
			t.Fatalf("check hook output should contain %q, got: %s", want, output)
		}
	}
}
