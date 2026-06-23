package upgrade

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func newTestChecker(server *httptest.Server) *Checker {
	return &Checker{
		CurrentVersion: "v0.1.0",
		RepoOwner:      "test-owner",
		RepoName:       "test-repo",
		HTTPClient:     server.Client(),
		BaseURL:        server.URL,
	}
}

func TestCheckLatest_Stable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/repos/test-owner/test-repo/releases/latest") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Release{
			TagName: "v0.2.0",
			Assets: []Asset{
				{Name: "trellis_linux_amd64", DownloadURL: "http://example.com/trellis", Size: 1000},
			},
		})
	}))
	defer server.Close()

	c := newTestChecker(server)
	release, err := c.CheckLatest("latest")
	if err != nil {
		t.Fatalf("CheckLatest() error: %v", err)
	}
	if release.TagName != "v0.2.0" {
		t.Errorf("expected v0.2.0, got %s", release.TagName)
	}
}

func TestCheckLatest_Beta(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Release{
			{TagName: "v0.2.0", Prerelease: false},
			{TagName: "v0.3.0-beta", Prerelease: true},
		})
	}))
	defer server.Close()

	c := newTestChecker(server)
	release, err := c.CheckLatest("beta")
	if err != nil {
		t.Fatalf("CheckLatest(beta) error: %v", err)
	}
	if release.TagName != "v0.3.0-beta" {
		t.Errorf("expected v0.3.0-beta, got %s", release.TagName)
	}
}

func TestNeedsUpdate_True(t *testing.T) {
	c := &Checker{CurrentVersion: "v0.1.0"}
	if !c.NeedsUpdate(&Release{TagName: "v0.2.0"}) {
		t.Error("expected NeedsUpdate to be true")
	}
}

func TestNeedsUpdate_False(t *testing.T) {
	c := &Checker{CurrentVersion: "v0.2.0"}
	if c.NeedsUpdate(&Release{TagName: "v0.2.0"}) {
		t.Error("expected NeedsUpdate to be false")
	}
}

func TestNeedsUpdate_NoVPrefix(t *testing.T) {
	c := &Checker{CurrentVersion: "0.1.0"}
	if !c.NeedsUpdate(&Release{TagName: "v0.2.0"}) {
		t.Error("expected NeedsUpdate to be true (no v prefix)")
	}
}

func TestDownload(t *testing.T) {
	expectedContent := []byte("mock binary content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(expectedContent)))
		w.Write(expectedContent)
	}))
	defer server.Close()

	dir := t.TempDir()
	dest := filepath.Join(dir, "downloaded")

	c := &Checker{HTTPClient: server.Client()}
	release := &Release{
		TagName: "v0.2.0",
		Assets: []Asset{
			{
				Name:        "trellis_" + runtime.GOOS + "_" + runtime.GOARCH,
				DownloadURL: server.URL,
				Size:        int64(len(expectedContent)),
			},
		},
	}

	if err := c.Download(release, dest); err != nil {
		t.Fatalf("Download() error: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("read downloaded: %v", err)
	}
	if string(data) != string(expectedContent) {
		t.Errorf("expected %q, got %q", expectedContent, data)
	}

	info, err := os.Stat(dest)
	if err != nil {
		t.Fatalf("stat downloaded: %v", err)
	}
	if info.Mode().Perm()&0111 == 0 {
		t.Error("downloaded file should be executable")
	}
}

func TestDownload_NoMatchingAsset(t *testing.T) {
	c := &Checker{}
	release := &Release{
		TagName: "v0.2.0",
		Assets:  []Asset{},
	}

	err := c.Download(release, "/tmp/nowhere")
	if err == nil {
		t.Fatal("expected error for no matching asset")
	}
}

func TestFindAsset(t *testing.T) {
	c := &Checker{}
	release := &Release{
		Assets: []Asset{
			{Name: "trellis_darwin_amd64"},
			{Name: "trellis_linux_amd64"},
			{Name: "trellis_windows_amd64.exe"},
		},
	}

	asset := c.findAsset(release)
	if asset == nil {
		t.Fatal("expected to find asset")
	}
	expected := "trellis_" + runtime.GOOS + "_" + runtime.GOARCH
	if !strings.Contains(asset.Name, expected) {
		t.Errorf("expected asset containing %q, got %q", expected, asset.Name)
	}
}

func TestReplace(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("replace test only works on Linux")
	}

	dir := t.TempDir()
	downloaded := filepath.Join(dir, "downloaded")
	target := filepath.Join(dir, "trellis")

	os.WriteFile(target, []byte("original"), 0644)
	os.WriteFile(downloaded, []byte("new version"), 0644)

	c := &Checker{}
	if err := c.Replace(downloaded, target); err != nil {
		t.Fatalf("Replace() error: %v", err)
	}

	data, _ := os.ReadFile(target)
	if string(data) != "new version" {
		t.Errorf("expected 'new version', got %q", string(data))
	}
}

func TestReplace_DownloadedNotFound(t *testing.T) {
	c := &Checker{}
	err := c.Replace("/nonexistent/downloaded", "/tmp/target")
	if err == nil {
		t.Fatal("expected error when downloaded file not found")
	}
}

func TestCheckLatest_SpecificTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/repos/test-owner/test-repo/releases/tags/v0.5.0") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Release{
			TagName: "v0.5.0",
			Assets:  []Asset{{Name: "trellis_linux_amd64", DownloadURL: "http://example.com/trellis"}},
		})
	}))
	defer server.Close()

	c := newTestChecker(server)
	release, err := c.CheckLatest("v0.5.0")
	if err != nil {
		t.Fatalf("CheckLatest(v0.5.0) error: %v", err)
	}
	if release.TagName != "v0.5.0" {
		t.Errorf("expected v0.5.0, got %s", release.TagName)
	}
}

func TestCheckLatest_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer server.Close()

	c := newTestChecker(server)
	_, err := c.CheckLatest("latest")
	if err == nil {
		t.Fatal("expected error for API 404")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("error should contain 404, got: %v", err)
	}
}

func TestCheckLatest_BetaNoPrerelease(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Release{
			{TagName: "v0.2.0", Prerelease: false},
		})
	}))
	defer server.Close()

	c := newTestChecker(server)
	release, err := c.CheckLatest("beta")
	if err != nil {
		t.Fatalf("CheckLatest(beta) error: %v", err)
	}
	if release.TagName != "v0.2.0" {
		t.Errorf("expected fallback to v0.2.0, got %s", release.TagName)
	}
}

func TestCheckLatest_BetaEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Release{})
	}))
	defer server.Close()

	c := newTestChecker(server)
	_, err := c.CheckLatest("beta")
	if err == nil {
		t.Fatal("expected error for empty releases")
	}
}

func TestDownload_SizeMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("short"))
	}))
	defer server.Close()

	dir := t.TempDir()
	dest := filepath.Join(dir, "downloaded")

	c := &Checker{HTTPClient: server.Client()}
	release := &Release{
		TagName: "v0.2.0",
		Assets: []Asset{
			{
				Name:        "trellis_" + runtime.GOOS + "_" + runtime.GOARCH,
				DownloadURL: server.URL,
				Size:        999999,
			},
		},
	}

	err := c.Download(release, dest)
	if err == nil {
		t.Fatal("expected error for size mismatch")
	}
	if !strings.Contains(err.Error(), "expected") {
		t.Errorf("error should mention expected bytes, got: %v", err)
	}
}

func TestFindAsset_Fallback(t *testing.T) {
	c := &Checker{}
	release := &Release{
		Assets: []Asset{
			{Name: "unrelated_file.zip"},
		},
	}
	asset := c.findAsset(release)
	if asset == nil {
		t.Fatal("expected fallback to first asset")
	}
	if asset.Name != "unrelated_file.zip" {
		t.Errorf("expected fallback asset, got %s", asset.Name)
	}
}

func TestHTTPClient_Default(t *testing.T) {
	c := &Checker{}
	client := c.httpClient()
	if client == nil {
		t.Fatal("expected non-nil default HTTP client")
	}
}
