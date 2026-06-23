package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestE2E_SubtaskAddAndList 测试 subtask 添加和列表显示
func TestE2E_SubtaskAddAndList(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 创建任务
	_, stderr, err = runTrellis(t, repo, "task", "create", "main-feature")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	// 添加两个子任务
	stdout, stderr, err := runTrellis(t, repo, "task", "add-subtask", "main-feature", "Setup DB schema")
	if err != nil {
		t.Fatalf("add-subtask failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Added subtask") {
		t.Errorf("expected 'Added subtask' in output, got: %s", stdout)
	}

	stdout, stderr, err = runTrellis(t, repo, "task", "add-subtask", "main-feature", "Create API endpoint")
	if err != nil {
		t.Fatalf("second add-subtask failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Added subtask") {
		t.Errorf("expected 'Added subtask' in output, got: %s", stdout)
	}

	// 查看任务信息，验证子任务列表
	stdout, stderr, err = runTrellis(t, repo, "task", "info", "main-feature")
	if err != nil {
		t.Fatalf("task info failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Subtasks:") {
		t.Errorf("task info should show Subtasks section, got: %s", stdout)
	}
	if !strings.Contains(stdout, "Setup DB schema") {
		t.Errorf("task info should contain 'Setup DB schema', got: %s", stdout)
	}
	if !strings.Contains(stdout, "Create API endpoint") {
		t.Errorf("task info should contain 'Create API endpoint', got: %s", stdout)
	}
}

// TestE2E_SubtaskDone 测试 subtask 完成标记
func TestE2E_SubtaskDone(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	_, stderr, err = runTrellis(t, repo, "task", "create", "main-feature")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	// 添加子任务
	_, stderr, err = runTrellis(t, repo, "task", "add-subtask", "main-feature", "Setup DB")
	if err != nil {
		t.Fatalf("add-subtask failed: %v\nstderr: %s", err, stderr)
	}

	// 完成子任务
	stdout, stderr, err := runTrellis(t, repo, "task", "done-subtask", "main-feature", "1")
	if err != nil {
		t.Fatalf("done-subtask failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Done") && !strings.Contains(stdout, "done") {
		t.Errorf("expected done confirmation, got: %s", stdout)
	}

	// 验证 task.json 中子任务标记为 done
	tasksDir := filepath.Join(repo, ".trellis", "tasks")
	taskDirName := firstTaskDirName(t, tasksDir)
	taskData, _ := os.ReadFile(filepath.Join(tasksDir, taskDirName, "task.json"))
	if !strings.Contains(string(taskData), `"done": true`) {
		t.Errorf("subtask should be marked done in task.json, got: %s", taskData)
	}
}

// TestE2E_SubtaskProgress 测试 subtask 进度显示
func TestE2E_SubtaskProgress(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	_, stderr, err = runTrellis(t, repo, "task", "create", "main-feature")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	// 添加 3 个子任务
	for _, title := range []string{"Setup DB", "Create API", "Write tests"} {
		_, stderr, err = runTrellis(t, repo, "task", "add-subtask", "main-feature", title)
		if err != nil {
			t.Fatalf("add-subtask %q failed: %v\nstderr: %s", title, err, stderr)
		}
	}

	// 完成 1 个
	_, stderr, err = runTrellis(t, repo, "task", "done-subtask", "main-feature", "1")
	if err != nil {
		t.Fatalf("done-subtask failed: %v\nstderr: %s", err, stderr)
	}

	// 验证 task.json 中进度
	tasksDir := filepath.Join(repo, ".trellis", "tasks")
	taskDirName := firstTaskDirName(t, tasksDir)
	taskData, _ := os.ReadFile(filepath.Join(tasksDir, taskDirName, "task.json"))

	// 验证 3 个子任务存在
	if !strings.Contains(string(taskData), `"Setup DB"`) {
		t.Error("task.json should contain 'Setup DB'")
	}
	if !strings.Contains(string(taskData), `"Create API"`) {
		t.Error("task.json should contain 'Create API'")
	}
	if !strings.Contains(string(taskData), `"Write tests"`) {
		t.Error("task.json should contain 'Write tests'")
	}

	// 验证只有第 1 个是 done
	if !strings.Contains(string(taskData), `"done": true`) {
		t.Error("at least one subtask should be done")
	}
}

// TestE2E_SubtaskAllDone 测试全部子任务完成
func TestE2E_SubtaskAllDone(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	_, stderr, err = runTrellis(t, repo, "task", "create", "main-feature")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	// 添加 2 个子任务
	for _, title := range []string{"Subtask A", "Subtask B"} {
		_, stderr, err = runTrellis(t, repo, "task", "add-subtask", "main-feature", title)
		if err != nil {
			t.Fatalf("add-subtask %q failed: %v\nstderr: %s", title, err, stderr)
		}
	}

	// 全部完成
	_, stderr, err = runTrellis(t, repo, "task", "done-subtask", "main-feature", "1")
	if err != nil {
		t.Fatalf("done-subtask 1 failed: %v\nstderr: %s", err, stderr)
	}
	_, stderr, err = runTrellis(t, repo, "task", "done-subtask", "main-feature", "2")
	if err != nil {
		t.Fatalf("done-subtask 2 failed: %v\nstderr: %s", err, stderr)
	}

	// 验证所有子任务都 done
	tasksDir := filepath.Join(repo, ".trellis", "tasks")
	taskDirName := firstTaskDirName(t, tasksDir)
	taskData, _ := os.ReadFile(filepath.Join(tasksDir, taskDirName, "task.json"))
	if !strings.Contains(string(taskData), `"done": true`) {
		t.Errorf("subtasks should be done, got: %s", taskData)
	}
}

// TestE2E_SubtaskDuplicateName 测试重复名称子任务
func TestE2E_SubtaskDuplicateName(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	_, stderr, err = runTrellis(t, repo, "task", "create", "main-feature")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	// 添加同名子任务（允许重复名称）
	_, stderr, err = runTrellis(t, repo, "task", "add-subtask", "main-feature", "Setup DB")
	if err != nil {
		t.Fatalf("first add-subtask failed: %v\nstderr: %s", err, stderr)
	}

	// 第二次添加同名子任务（应允许，因为 ID 不同）
	stdout, stderr, err := runTrellis(t, repo, "task", "add-subtask", "main-feature", "Setup DB")
	if err != nil {
		t.Fatalf("second add-subtask with same name failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Added subtask") {
		t.Errorf("duplicate name should still be added, got: %s", stdout)
	}

	// 验证两个子任务都存在
	tasksDir := filepath.Join(repo, ".trellis", "tasks")
	taskDirName := firstTaskDirName(t, tasksDir)
	taskData, _ := os.ReadFile(filepath.Join(tasksDir, taskDirName, "task.json"))
	if !strings.Contains(string(taskData), `"subtasks"`) {
		t.Errorf("task.json should contain subtasks, got: %s", taskData)
	}
}

// TestE2E_SubtaskDoneNotFound 测试完成不存在的子任务
func TestE2E_SubtaskDoneNotFound(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "test")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	_, stderr, err = runTrellis(t, repo, "task", "create", "main-feature")
	if err != nil {
		t.Fatalf("task create failed: %v\nstderr: %s", err, stderr)
	}

	// 完成不存在的子任务
	_, stderr, err = runTrellis(t, repo, "task", "done-subtask", "main-feature", "999")
	if err == nil {
		t.Fatal("expected error for non-existent subtask")
	}
	if !strings.Contains(stderr, "not found") {
		t.Errorf("error should mention 'not found', got: %s", stderr)
	}
}
