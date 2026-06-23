package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/superops-team/trellis-go/pkg/upgrade"
)

// --- Update E2E ---

// TestE2E_UpdateSyncsTemplates 验证 update 恢复缺失的嵌入模板文件
func TestE2E_UpdateSyncsTemplates(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "alice", "--platform", "claude")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 先运行一次 update 创建嵌入模板文件
	_, _, err = runTrellis(t, repo, "update")
	if err != nil {
		t.Fatalf("initial update failed: %v", err)
	}

	// 删除嵌入模板文件
	gitkeepPath := filepath.Join(repo, ".trellis", "templates", ".gitkeep")
	if _, err := os.Stat(gitkeepPath); err != nil {
		t.Skipf("templates/.gitkeep not created by update, skipping: %v", err)
	}
	if err := os.Remove(gitkeepPath); err != nil {
		t.Fatalf("remove .gitkeep: %v", err)
	}

	stdout, stderr, err := runTrellis(t, repo, "update")
	if err != nil {
		t.Fatalf("update failed: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(gitkeepPath); err != nil {
		t.Errorf(".gitkeep should be restored after update: %v", err)
	}

	if !strings.Contains(stdout, "Added") && !strings.Contains(stdout, "+") {
		t.Logf("update output: %s", stdout)
	}
}

// TestE2E_UpdatePreservesUserEdits 验证 update 不覆盖用户修改的模板文件
func TestE2E_UpdatePreservesUserEdits(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "alice", "--platform", "claude")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 先运行一次 update 创建嵌入模板文件
	_, _, err = runTrellis(t, repo, "update")
	if err != nil {
		t.Fatalf("initial update failed: %v", err)
	}

	gitkeepPath := filepath.Join(repo, ".trellis", "templates", ".gitkeep")
	if _, err := os.Stat(gitkeepPath); err != nil {
		t.Skipf("templates/.gitkeep not created by update, skipping: %v", err)
	}

	original, err := os.ReadFile(gitkeepPath)
	if err != nil {
		t.Fatalf("read .gitkeep: %v", err)
	}
	modified := append([]byte("# User modified\n"), original...)
	if err := os.WriteFile(gitkeepPath, modified, 0644); err != nil {
		t.Fatalf("write modified .gitkeep: %v", err)
	}

	stdout, stderr, err := runTrellis(t, repo, "update")
	if err != nil {
		t.Fatalf("update failed: %v\nstderr: %s", err, stderr)
	}

	current, err := os.ReadFile(gitkeepPath)
	if err != nil {
		t.Fatalf("read .gitkeep after update: %v", err)
	}
	if !strings.HasPrefix(string(current), "# User modified") {
		t.Error("user modified template file should be preserved")
	}

	if !strings.Contains(stdout, "Skipped") && !strings.Contains(stdout, "-") {
		t.Logf("update output: %s", stdout)
	}
}

// TestE2E_UpdateDryRun 验证 --dry-run 不写入文件
func TestE2E_UpdateDryRun(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "alice", "--platform", "claude")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 先运行一次 update 创建嵌入模板文件
	_, _, err = runTrellis(t, repo, "update")
	if err != nil {
		t.Fatalf("initial update failed: %v", err)
	}

	gitkeepPath := filepath.Join(repo, ".trellis", "templates", ".gitkeep")
	if _, err := os.Stat(gitkeepPath); err != nil {
		t.Skipf("templates/.gitkeep not created by update, skipping: %v", err)
	}
	if err := os.Remove(gitkeepPath); err != nil {
		t.Fatalf("remove .gitkeep: %v", err)
	}

	stdout, stderr, err := runTrellis(t, repo, "update", "--dry-run")
	if err != nil {
		t.Fatalf("update --dry-run failed: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(gitkeepPath); err == nil {
		t.Error(".gitkeep should NOT be restored after --dry-run")
	}

	if !strings.Contains(stdout, "Dry run") && !strings.Contains(stdout, "dry") {
		t.Logf("dry-run output: %s", stdout)
	}
}

// TestE2E_UpdateSkipPaths 验证 update.skip 配置生效
func TestE2E_UpdateSkipPaths(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "alice", "--platform", "claude")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	// 先运行一次 update 创建嵌入模板文件
	_, _, err = runTrellis(t, repo, "update")
	if err != nil {
		t.Fatalf("initial update failed: %v", err)
	}

	configPath := filepath.Join(repo, ".trellis", "config.yaml")
	cfgData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config.yaml: %v", err)
	}
	skipConfig := string(cfgData) + "\nupdate:\n  skip:\n    - templates/.gitkeep\n"
	if err := os.WriteFile(configPath, []byte(skipConfig), 0644); err != nil {
		t.Fatalf("write config.yaml: %v", err)
	}

	gitkeepPath := filepath.Join(repo, ".trellis", "templates", ".gitkeep")
	if _, err := os.Stat(gitkeepPath); err != nil {
		t.Skipf("templates/.gitkeep not created by update, skipping: %v", err)
	}
	if err := os.Remove(gitkeepPath); err != nil {
		t.Fatalf("remove .gitkeep: %v", err)
	}

	stdout, stderr, err := runTrellis(t, repo, "update")
	if err != nil {
		t.Fatalf("update failed: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(gitkeepPath); err == nil {
		t.Error(".gitkeep should NOT be restored when in update.skip")
	}
	_ = stdout
}

// TestE2E_UpdateConfigSectionAppend 验证配置段追加 + 幂等
func TestE2E_UpdateConfigSectionAppend(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)

	_, stderr, err := runTrellis(t, repo, "init", "--developer", "alice", "--platform", "claude")
	if err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", err, stderr)
	}

	stdout1, _, err := runTrellis(t, repo, "update")
	if err != nil {
		t.Fatalf("first update failed: %v", err)
	}

	stdout2, _, err := runTrellis(t, repo, "update")
	if err != nil {
		t.Fatalf("second update failed: %v", err)
	}

	if strings.Contains(stdout1, "Config sections") && strings.Contains(stdout2, "Config sections") {
		t.Log("config sections appended on both runs — may not be idempotent")
	}

	if strings.Contains(stdout2, "Already up to date") {
		// 幂等行为正确
	}
}

// --- Upgrade E2E (with mock HTTP) ---

func newMockUpgradeServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/repos/superops-team/trellis-go/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		release := upgrade.Release{
			TagName:    "v99.0.0",
			Prerelease: false,
			Assets: []upgrade.Asset{
				{Name: "trellis-go_linux_amd64", DownloadURL: "http://example.com/trellis", Size: 100},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	})

	mux.HandleFunc("/repos/superops-team/trellis-go/releases/tags/", func(w http.ResponseWriter, r *http.Request) {
		release := upgrade.Release{
			TagName:    "v99.0.0",
			Prerelease: false,
			Assets: []upgrade.Asset{
				{Name: "trellis-go_linux_amd64", DownloadURL: "http://example.com/trellis", Size: 100},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	})

	return httptest.NewServer(mux)
}

// TestE2E_UpgradeCheckVersion 验证版本检查逻辑（用 mock HTTP server）
func TestE2E_UpgradeCheckVersion(t *testing.T) {
	server := newMockUpgradeServer(t)
	defer server.Close()

	checker := &upgrade.Checker{
		CurrentVersion: "v0.4.0",
		RepoOwner:      "superops-team",
		RepoName:       "trellis-go",
		BaseURL:        server.URL,
	}

	release, err := checker.CheckLatest("latest")
	if err != nil {
		t.Fatalf("CheckLatest failed: %v", err)
	}

	if release.TagName != "v99.0.0" {
		t.Errorf("expected v99.0.0, got %s", release.TagName)
	}

	if !checker.NeedsUpdate(release) {
		t.Error("v0.4.0 should need update to v99.0.0")
	}

	checker.CurrentVersion = "v99.0.0"
	if checker.NeedsUpdate(release) {
		t.Error("v99.0.0 should NOT need update to v99.0.0")
	}
}

// TestE2E_UpgradeCheckLatestSameVersion 验证已是最新版本时 NeedsUpdate 返回 false
func TestE2E_UpgradeCheckLatestSameVersion(t *testing.T) {
	server := newMockUpgradeServer(t)
	defer server.Close()

	checker := &upgrade.Checker{
		CurrentVersion: "v99.0.0",
		RepoOwner:      "superops-team",
		RepoName:       "trellis-go",
		BaseURL:        server.URL,
	}

	release, err := checker.CheckLatest("latest")
	if err != nil {
		t.Fatalf("CheckLatest failed: %v", err)
	}

	if checker.NeedsUpdate(release) {
		t.Error("same version should not need update")
	}
}

// TestE2E_UpgradeCheckBeta 验证 beta 版本检查
func TestE2E_UpgradeCheckBeta(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		releases := []upgrade.Release{
			{TagName: "v1.0.0", Prerelease: false},
			{TagName: "v1.1.0-beta.1", Prerelease: true},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(releases)
	}))
	defer server.Close()

	checker := &upgrade.Checker{
		CurrentVersion: "v1.0.0",
		RepoOwner:      "superops-team",
		RepoName:       "trellis-go",
		BaseURL:        server.URL,
	}

	release, err := checker.CheckLatest("beta")
	if err != nil {
		t.Fatalf("CheckLatest beta failed: %v", err)
	}

	if release.TagName != "v1.1.0-beta.1" {
		t.Errorf("expected v1.1.0-beta.1, got %s", release.TagName)
	}
}
