package template

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed testdata/hello.txt testdata/literal.raw testdata/data.bin testdata/nested/child.txt
var testFS embed.FS

//go:embed testdata/unknown/*
var unknownTestFS embed.FS

func TestRenderString(t *testing.T) {
	eng := NewEngine(testFS, "testdata")
	ctx := RenderContext{
		PlatformID:   "claude",
		PlatformName: "Claude Code",
		Developer:    "alice",
	}

	result, err := eng.RenderString("Hello {{.Developer}}, using {{.PlatformName}}", ctx)
	if err != nil {
		t.Fatalf("RenderString failed: %v", err)
	}
	expected := "Hello alice, using Claude Code"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderString_UnknownField(t *testing.T) {
	eng := NewEngine(testFS, "testdata")
	ctx := RenderContext{Developer: "alice"}

	_, err := eng.RenderString("Hello {{.Unknown}}", ctx)
	if err == nil {
		t.Fatal("expected unknown field error")
	}
	if !strings.Contains(err.Error(), "Unknown") {
		t.Errorf("error should mention unknown field, got: %v", err)
	}
}

func TestShouldTemplate(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"file.md", true},
		{"file.json", true},
		{"file.yaml", true},
		{"file.png", false},
		{"file.bin", false},
		{"file.go", true},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			got := ShouldTemplate(tc.path)
			if got != tc.want {
				t.Errorf("ShouldTemplate(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}

func TestRender_Dir(t *testing.T) {
	dst := t.TempDir()

	eng := NewEngine(testFS, "testdata")
	ctx := RenderContext{
		Developer:    "bob",
		PlatformName: "Cursor",
	}

	if err := eng.Render(".", dst, ctx); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	assertFileContent(t, filepath.Join(dst, "hello.txt"), "Hello bob")
	assertFileContent(t, filepath.Join(dst, "literal.raw"), "Literal {{.Developer}}\n")
	assertFileContent(t, filepath.Join(dst, "data.bin"), "\x00{{.Developer}}\n")
	if _, err := os.Stat(filepath.Join(dst, ".template-hashes.json")); err != nil {
		t.Fatalf("expected template hash output: %v", err)
	}

	_ = dst
}

func TestRender_NestedDirUsesSlashPathsForEmbeddedFS(t *testing.T) {
	dst := t.TempDir()
	eng := NewEngine(testFS, "testdata")

	if err := eng.Render("nested", dst, RenderContext{Developer: "alice"}); err != nil {
		t.Fatalf("Render nested failed: %v", err)
	}

	assertFileContent(t, filepath.Join(dst, "child.txt"), "Nested alice\n")
}

func TestRender_UnknownTemplateKeyReturnsError(t *testing.T) {
	dst := t.TempDir()
	eng := NewEngine(unknownTestFS, "testdata")

	err := eng.Render("unknown", dst, RenderContext{Developer: "alice"})
	if err == nil {
		t.Fatal("expected missing key error")
	}
	if !strings.Contains(err.Error(), "Missing") {
		t.Errorf("error should mention missing key, got: %v", err)
	}
}

func TestRender_ReturnsHashWriteError(t *testing.T) {
	dst := t.TempDir()
	if err := os.Mkdir(filepath.Join(dst, ".template-hashes.json"), 0755); err != nil {
		t.Fatalf("create hash path blocker: %v", err)
	}
	eng := NewEngine(testFS, "testdata")

	err := eng.Render(".", dst, RenderContext{Developer: "alice"})
	if err == nil {
		t.Fatal("expected hash write error")
	}
	if !strings.Contains(err.Error(), "write hashes") {
		t.Errorf("error should mention hash write failure, got: %v", err)
	}
}

func TestFuncMap(t *testing.T) {
	ctx := RenderContext{
		PlatformID:   "claude",
		PlatformName: "Claude Code",
		Developer:    "alice",
		ProjectName:  "myproject",
		ExecutorAI:   "claude",
		AgentCapable: true,
		HasHooks:     false,
		CLIFlag:      "--flag",
	}
	fm := ctx.FuncMap()

	// Verify all expected functions exist
	keys := []string{"PlatformID", "PlatformName", "Developer", "ProjectName", "ExecutorAI", "AgentCapable", "HasHooks", "CLIFlag"}
	for _, k := range keys {
		if _, ok := fm[k]; !ok {
			t.Errorf("FuncMap missing key: %s", k)
		}
	}

	// Call each function to ensure coverage
	if v, ok := fm["PlatformID"].(func() string); ok {
		if v() != "claude" {
			t.Errorf("PlatformID = %s, want claude", v())
		}
	}
	if v, ok := fm["AgentCapable"].(func() bool); ok {
		if !v() {
			t.Error("AgentCapable should be true")
		}
	}
	if v, ok := fm["HasHooks"].(func() bool); ok {
		if v() {
			t.Error("HasHooks should be false")
		}
	}
}

func TestToMap(t *testing.T) {
	ctx := RenderContext{
		PlatformID: "cursor",
		Developer:  "bob",
		Extra:      map[string]any{"CustomKey": "custom-value"},
	}
	m := ctx.ToMap()
	if m["PlatformID"] != "cursor" {
		t.Errorf("PlatformID = %v, want cursor", m["PlatformID"])
	}
	if m["Developer"] != "bob" {
		t.Errorf("Developer = %v, want bob", m["Developer"])
	}
	if m["CustomKey"] != "custom-value" {
		t.Errorf("CustomKey = %v, want custom-value", m["CustomKey"])
	}
}

func TestHash(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.txt")
	os.WriteFile(path, []byte("hello"), 0644)

	eng := NewEngine(testFS, "testdata")
	hash, err := eng.Hash(path)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash")
	}

	// Same content => same hash
	hash2, _ := eng.Hash(path)
	if hash != hash2 {
		t.Error("same file should have same hash")
	}
}

func TestHash_NonExistent(t *testing.T) {
	eng := NewEngine(testFS, "testdata")
	_, err := eng.Hash(filepath.Join(t.TempDir(), "nonexistent"))
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(data) != want {
		t.Fatalf("%s = %q, want %q", path, data, want)
	}
}
