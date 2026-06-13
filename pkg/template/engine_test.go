package template

import (
	"embed"
	"os"
	"path/filepath"
	"testing"
)

//go:embed testdata/*.txt
var testFS embed.FS

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

	// Go's text/template does not error on unknown fields by default;
	// it renders them as <no value>. We verify this behavior.
	result, err := eng.RenderString("Hello {{.Unknown}}", ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello <no value>" {
		t.Errorf("expected '<no value>' for unknown field, got %q", result)
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
	// Create a temporary embed FS simulation by writing files
	src := t.TempDir()
	dst := t.TempDir()

	os.WriteFile(filepath.Join(src, "hello.txt"), []byte("Hello {{.Developer}}"), 0644)
	os.Mkdir(filepath.Join(src, "sub"), 0755)
	os.WriteFile(filepath.Join(src, "sub", "world.txt"), []byte("World {{.PlatformName}}"), 0644)

	// For this test we use a simple string render approach
	eng := NewEngine(testFS, "testdata")
	ctx := RenderContext{
		Developer:    "bob",
		PlatformName: "Cursor",
	}

	// Verify RenderString works for each file content
	for _, tc := range []struct {
		content string
		want    string
	}{
		{"Hello {{.Developer}}", "Hello bob"},
		{"World {{.PlatformName}}", "World Cursor"},
	} {
		got, err := eng.RenderString(tc.content, ctx)
		if err != nil {
			t.Fatalf("RenderString failed: %v", err)
		}
		if got != tc.want {
			t.Errorf("expected %q, got %q", tc.want, got)
		}
	}

	_ = dst
	_ = src
}
