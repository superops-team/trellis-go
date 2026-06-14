package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const canonicalModulePath = "github.com/superops-team/trellis-go"
const oldModulePath = "github.com/" + "mindfold/trellis"

func TestModulePathMatchesReadmeInstallCommand(t *testing.T) {
	repo := repoRoot(t)

	goMod, err := os.ReadFile(filepath.Join(repo, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if !strings.Contains(string(goMod), "module "+canonicalModulePath) {
		t.Fatalf("go.mod should declare %s, got:\n%s", canonicalModulePath, goMod)
	}

	for _, name := range []string{"README.md", "README.zh-CN.md"} {
		data, err := os.ReadFile(filepath.Join(repo, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		install := "go install " + canonicalModulePath + "/cmd/trellis@latest"
		if !strings.Contains(string(data), install) {
			t.Fatalf("%s should contain install command %q", name, install)
		}
	}
}

func TestRepositoryOwnedPathsDoNotUseOldModulePath(t *testing.T) {
	repo := repoRoot(t)

	err := filepath.WalkDir(repo, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "vendor":
				return filepath.SkipDir
			}
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(data), oldModulePath) {
			t.Fatalf("%s contains old module path %s; use %s", path, oldModulePath, canonicalModulePath)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scan repository paths: %v", err)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("repo root with go.mod not found")
		}
		dir = parent
	}
}
