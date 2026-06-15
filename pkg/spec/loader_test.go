package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoader_Load(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "guide.md"), []byte("# Guide"), 0644); err != nil {
		t.Fatalf("write spec file: %v", err)
	}

	got, err := NewLoader(root).Load("guide.md")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if got != "# Guide" {
		t.Fatalf("Load() = %q, want %q", got, "# Guide")
	}
}

func TestLoader_Load_MissingFileReturnsError(t *testing.T) {
	root := t.TempDir()

	_, err := NewLoader(root).Load("missing.md")
	if err == nil {
		t.Fatal("expected missing file error")
	}
	if !strings.Contains(err.Error(), "missing.md") {
		t.Fatalf("error should mention missing file, got: %v", err)
	}
}

func TestLoader_LoadPackage(t *testing.T) {
	root := t.TempDir()
	layerDir := filepath.Join(root, "auth", "api")
	if err := os.MkdirAll(layerDir, 0755); err != nil {
		t.Fatalf("create layer dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(layerDir, "index.md"), []byte("# Auth API"), 0644); err != nil {
		t.Fatalf("write layer index: %v", err)
	}

	got, err := NewLoader(root).LoadPackage("auth")
	if err != nil {
		t.Fatalf("LoadPackage failed: %v", err)
	}
	if got["api"] != "# Auth API" {
		t.Fatalf("LoadPackage()[api] = %q", got["api"])
	}
}

func TestLoader_LoadPackage_UnreadableLayerReturnsError(t *testing.T) {
	root := t.TempDir()
	indexPath := filepath.Join(root, "auth", "api", "index.md")
	if err := os.MkdirAll(indexPath, 0755); err != nil {
		t.Fatalf("create unreadable index path: %v", err)
	}

	_, err := NewLoader(root).LoadPackage("auth")
	if err == nil {
		t.Fatal("expected unreadable layer error")
	}
	if !strings.Contains(err.Error(), "auth") || !strings.Contains(err.Error(), "api") {
		t.Fatalf("error should identify package and layer, got: %v", err)
	}
}

func TestLoader_IndexIncludesPackagesAndGuides(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{
		filepath.Join(root, "auth", "api"),
		filepath.Join(root, "auth", "data"),
		filepath.Join(root, "guides"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("create dir %s: %v", dir, err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "auth", "api", "index.md"), []byte("# API"), 0644); err != nil {
		t.Fatalf("write api index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "auth", "data", "index.md"), []byte("# Data"), 0644); err != nil {
		t.Fatalf("write data index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "guides", "getting-started.md"), []byte("# Guide"), 0644); err != nil {
		t.Fatalf("write guide: %v", err)
	}

	idx, err := NewLoader(root).Index()
	if err != nil {
		t.Fatalf("Index failed: %v", err)
	}
	if idx.Packages["auth"].Layers["api"] != filepath.Join("auth", "api", "index.md") {
		t.Fatalf("api layer not indexed: %+v", idx.Packages["auth"].Layers)
	}
	if idx.Packages["auth"].Layers["data"] != filepath.Join("auth", "data", "index.md") {
		t.Fatalf("data layer not indexed: %+v", idx.Packages["auth"].Layers)
	}
	if len(idx.Guides) != 1 || idx.Guides[0] != filepath.Join("guides", "getting-started.md") {
		t.Fatalf("guide not indexed: %+v", idx.Guides)
	}

	markdown := idx.ToMarkdown()
	for _, want := range []string{"# Spec Index", "### auth", "**api**", "**data**", "## Guides", "guides/getting-started.md"} {
		if !strings.Contains(markdown, want) {
			t.Errorf("markdown index should contain %q, got: %s", want, markdown)
		}
	}
}

func TestLoader_LoadPackageSkipsResourceDirectoriesWithoutIndex(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{
		filepath.Join(root, "auth", "api"),
		filepath.Join(root, "auth", "assets"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("create dir %s: %v", dir, err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "auth", "api", "index.md"), []byte("# API"), 0644); err != nil {
		t.Fatalf("write api index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "auth", "assets", "logo.txt"), []byte("logo"), 0644); err != nil {
		t.Fatalf("write asset: %v", err)
	}

	got, err := NewLoader(root).LoadPackage("auth")
	if err != nil {
		t.Fatalf("LoadPackage should skip resource directories without index.md: %v", err)
	}
	if len(got) != 1 || got["api"] != "# API" {
		t.Fatalf("LoadPackage result mismatch: %+v", got)
	}
}

func TestLoader_IndexMissingRootReturnsError(t *testing.T) {
	root := filepath.Join(t.TempDir(), "missing")

	_, err := NewLoader(root).Index()
	if err == nil {
		t.Fatal("expected missing spec root error")
	}
	if !strings.Contains(err.Error(), "read spec root") {
		t.Errorf("error should mention spec root, got: %v", err)
	}
}

func TestLoader_LoadLayer(t *testing.T) {
	root := t.TempDir()
	layerDir := filepath.Join(root, "auth", "api")
	if err := os.MkdirAll(layerDir, 0755); err != nil {
		t.Fatalf("create layer dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(layerDir, "index.md"), []byte("# Auth API"), 0644); err != nil {
		t.Fatalf("write layer index: %v", err)
	}

	got, err := NewLoader(root).LoadLayer("auth", "api")
	if err != nil {
		t.Fatalf("LoadLayer failed: %v", err)
	}
	if got != "# Auth API" {
		t.Fatalf("LoadLayer() = %q, want %q", got, "# Auth API")
	}
}
