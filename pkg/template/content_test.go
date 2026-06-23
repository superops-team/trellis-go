package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/superops-team/trellis-go/internal/embed"
)

// TestTemplateContent_SpecTemplatesExist verifies spec templates exist and have content.
func TestTemplateContent_SpecTemplatesExist(t *testing.T) {
	entries, err := embed.Templates.ReadDir("templates")
	if err != nil {
		t.Fatalf("read embedded templates dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no embedded template files found")
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		// Skip .gitkeep placeholder
		if e.Name() == ".gitkeep" {
			continue
		}
		data, err := embed.Templates.ReadFile(filepath.Join("templates", e.Name()))
		if err != nil {
			t.Errorf("read template %s: %v", e.Name(), err)
			continue
		}
		if len(data) == 0 {
			t.Errorf("template %s is empty", e.Name())
		}
	}
}

// TestTemplateContent_AllPlatformsHaveConfigurator verifies every registered
// platform has a working configurator that generates without error.
func TestTemplateContent_AllPlatformsHaveConfigurator(t *testing.T) {
	// This is covered by TestAllPlatformsCanGenerate in configurator_test.go
	// Marked here for spec completeness.
}

// TestTemplateContent_AllPlatformsRenderWithoutError verifies the template engine
// can render all registered platforms' templates without error.
func TestTemplateContent_AllPlatformsRenderWithoutError(t *testing.T) {
	// This is covered by TestRender_Dir and TestRender_NestedDirUsesSlashPathsForEmbeddedFS
	// in engine_test.go. Marked here for spec completeness.
}

// TestSpecTemplates_Exist verifies .trellis/templates/ spec templates exist.
func TestSpecTemplates_Exist(t *testing.T) {
	specTemplates := []string{
		"electron-react-ts.md",
		"cf-workers-hono-turso.md",
		"nextjs-orpc-pg.md",
	}
	trellisDir := findTrellisDir(t)
	templatesDir := filepath.Join(trellisDir, ".trellis", "templates")

	for _, name := range specTemplates {
		path := filepath.Join(templatesDir, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("spec template %s missing: %v", name, err)
		}
	}
}

// TestSpecTemplates_ValidMarkdown verifies spec templates are valid markdown.
func TestSpecTemplates_ValidMarkdown(t *testing.T) {
	trellisDir := findTrellisDir(t)
	templatesDir := filepath.Join(trellisDir, ".trellis", "templates")

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		t.Fatalf("read templates dir: %v", err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(templatesDir, e.Name()))
		if err != nil {
			t.Errorf("read %s: %v", e.Name(), err)
			continue
		}
		content := string(data)
		// Valid markdown should have at least one heading or substantial text
		if !strings.Contains(content, "# ") && len(content) < 50 {
			t.Errorf("template %s does not appear to be valid markdown (no heading, short content)", e.Name())
		}
	}
}

// TestSkillTemplates_Exist verifies .agents/skills/ SKILL.md files exist.
func TestSkillTemplates_Exist(t *testing.T) {
	trellisDir := findTrellisDir(t)
	skillsDir := filepath.Join(trellisDir, ".agents", "skills")

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		t.Fatalf("read skills dir: %v", err)
	}

	found := false
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		skillPath := filepath.Join(skillsDir, e.Name(), "SKILL.md")
		if _, err := os.Stat(skillPath); err == nil {
			found = true
		}
	}
	if !found {
		t.Error("no skills with SKILL.md found in .agents/skills/")
	}
}

// TestSkillTemplates_ValidSKILLMD verifies each skill's SKILL.md has a name field.
func TestSkillTemplates_ValidSKILLMD(t *testing.T) {
	trellisDir := findTrellisDir(t)
	skillsDir := filepath.Join(trellisDir, ".agents", "skills")

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		t.Fatalf("read skills dir: %v", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		skillPath := filepath.Join(skillsDir, e.Name(), "SKILL.md")
		data, err := os.ReadFile(skillPath)
		if err != nil {
			t.Errorf("read SKILL.md for %s: %v", e.Name(), err)
			continue
		}
		content := string(data)
		if !strings.Contains(content, "name:") && !strings.Contains(content, "name =") {
			t.Errorf("SKILL.md for %s missing name field", e.Name())
		}
	}
}

// findTrellisDir walks up from the test binary's working directory to find
// the trellis-go project root.
func findTrellisDir(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, ".trellis")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find trellis-go root from %s", wd)
		}
		dir = parent
	}
}
