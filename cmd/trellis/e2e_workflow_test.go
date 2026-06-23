package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestE2E_WorkflowFullCycle 测试 workflow 完整 4 阶段流转
// 验证 workflow.md 生成、各阶段状态注入、完整生命周期
func TestE2E_WorkflowFullCycle(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	// Step 1: init 生成 workflow.md
	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	trellisDir := filepath.Join(repo, ".trellis")
	workflowPath := filepath.Join(trellisDir, "workflow.md")
	workflowData, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("workflow.md not created: %v", err)
	}

	// 验证 workflow.md 包含所有 4 个阶段的状态注入标签
	for _, state := range []string{"PLAN", "IMPLEMENT", "VERIFY", "FINISH"} {
		if !strings.Contains(string(workflowData), "[workflow-state:"+state+"]") {
			t.Errorf("workflow.md should contain [workflow-state:%s]", state)
		}
	}

	// Step 2: 验证每个阶段的 inject-workflow-state 输出
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
		if !strings.Contains(stdout, tc.hint) && !strings.Contains(stdout, strings.ToUpper(tc.state)) {
			t.Errorf("inject-workflow-state %s output should contain %q, got: %s", tc.state, tc.hint, stdout)
		}
	}

	// Step 3: 创建任务并验证完整生命周期
	_, stderr, err = runTrellis(t, repo, "task", "create", "workflow-test")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	tasksDir := filepath.Join(trellisDir, "tasks")
	taskDirName := firstTaskDirName(t, tasksDir)
	taskPath := filepath.Join(tasksDir, taskDirName, "task.json")

	// 验证初始状态为 planning
	taskData, _ := os.ReadFile(taskPath)
	if !strings.Contains(string(taskData), `"status": "planning"`) {
		t.Error("initial status should be planning")
	}

	// 写入 PRD 并启动任务
	if err := os.WriteFile(filepath.Join(tasksDir, taskDirName, "prd.md"), []byte("# PRD\nWorkflow test."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}

	stdout, stderr, err := runTrellis(t, repo, "task", "start", "workflow-test")
	if err != nil {
		t.Fatalf("task start failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Started task: workflow-test") {
		t.Errorf("expected start output, got: %s", stdout)
	}

	// 验证状态为 in_progress
	taskData, _ = os.ReadFile(taskPath)
	if !strings.Contains(string(taskData), `"status": "in_progress"`) {
		t.Errorf("status should be in_progress, got: %s", taskData)
	}

	// 归档任务（完成）
	stdout, stderr, err = runTrellis(t, repo, "task", "archive", "workflow-test")
	if err != nil {
		t.Fatalf("task archive failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Archived task: workflow-test") {
		t.Errorf("expected archive output, got: %s", stdout)
	}

	// 验证归档后状态为 completed
	destDir := archivedTaskDir(t, tasksDir, taskDirName)
	archivedData, _ := os.ReadFile(filepath.Join(destDir, "task.json"))
	if !strings.Contains(string(archivedData), `"status": "completed"`) {
		t.Errorf("archived task should have status completed, got: %s", archivedData)
	}
}

// TestE2E_WorkflowStateInjection 测试 workflow 状态注入输出
// 验证 inject-workflow-state hook 在任务上下文中输出正确的工作流状态
func TestE2E_WorkflowStateInjection(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 创建任务并启动
	_, stderr, err = runTrellis(t, repo, "task", "create", "state-test")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	tasksDir := filepath.Join(repo, ".trellis", "tasks")
	taskDirName := firstTaskDirName(t, tasksDir)
	if err := os.WriteFile(filepath.Join(tasksDir, taskDirName, "prd.md"), []byte("# PRD\nState injection test."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	_, stderr, err = runTrellis(t, repo, "task", "start", "state-test")
	if err != nil {
		t.Fatalf("task start failed: %v\nstderr: %s", err, stderr)
	}

	// 执行 inject-workflow-state hook
	stdout, stderr, err := runTrellis(t, repo, "hook", "inject-workflow-state", "--state", "implement")
	if err != nil {
		t.Fatalf("inject-workflow-state failed: %v\nstderr: %s", err, stderr)
	}

	// 验证输出包含关键信息
	checks := []string{
		"IMPLEMENT",
		"<workflow-state>",
		"</workflow-state>",
	}
	for _, want := range checks {
		if !strings.Contains(stdout, want) {
			t.Errorf("output should contain %q, got: %s", want, stdout)
		}
	}

	// 验证输出有实质内容
	if len(stdout) < 20 {
		t.Errorf("output too short: %s", stdout)
	}
	for _, want := range checks {
		if !strings.Contains(stdout, want) {
			t.Errorf("output should contain %q, got: %s", want, stdout)
		}
	}
}

// TestE2E_WorkflowInvalidTransition 测试无效的阶段跳转被拒绝
// 验证不能从 planning 直接 archive，不能从 completed 重新 start
func TestE2E_WorkflowInvalidTransition(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 创建任务
	_, stderr, err = runTrellis(t, repo, "task", "create", "invalid-flow")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	// 尝试从 planning 直接 archive（应拒绝）
	_, stderr, err = runTrellis(t, repo, "task", "archive", "invalid-flow")
	if err == nil {
		t.Fatal("expected archiving planning task to fail")
	}
	if !strings.Contains(stderr, "invalid task status transition") {
		t.Errorf("archive invalid transition should mention invalid transition, got: %s", stderr)
	}

	// 模拟任务已完成状态，尝试重新 start（应拒绝）
	tasksDir := filepath.Join(repo, ".trellis", "tasks")
	taskDirName := firstTaskDirName(t, tasksDir)
	taskPath := filepath.Join(tasksDir, taskDirName, "task.json")
	taskData, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatalf("read task.json: %v", err)
	}
	completedData := strings.Replace(string(taskData), `"status": "planning"`, `"status": "completed"`, 1)
	if completedData == string(taskData) {
		t.Fatalf("task.json did not contain planning status: %s", taskData)
	}
	if err := os.WriteFile(taskPath, []byte(completedData), 0644); err != nil {
		t.Fatalf("write completed task.json: %v", err)
	}

	_, stderr, err = runTrellis(t, repo, "task", "start", "invalid-flow")
	if err == nil {
		t.Fatal("expected starting completed task to fail")
	}
	if !strings.Contains(stderr, "invalid task status transition") {
		t.Errorf("start invalid transition should mention invalid transition, got: %s", stderr)
	}

	// 验证状态未被修改
	taskData, _ = os.ReadFile(taskPath)
	if !strings.Contains(string(taskData), `"status": "completed"`) {
		t.Errorf("rejected transition should not change status, got: %s", taskData)
	}
}
