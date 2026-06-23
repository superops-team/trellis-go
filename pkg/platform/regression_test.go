package platform_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/superops-team/trellis-go/pkg/hook"
	"github.com/superops-team/trellis-go/pkg/platform"
)

// Regression tests for historical bugs and version-specific invariants.
// Each test is annotated with the version it was introduced and what bug it prevents.

// Regression_v0_4_0: After adding ZCode, platform count must be 16.
// Bug: Adding a platform without updating the count check.
func TestRegression_PlatformCount(t *testing.T) {
	r := platform.NewRegistry()
	all := r.All()
	const expected = 16
	if len(all) != expected {
		t.Errorf("platform count regression: expected %d, got %d. "+
			"Did you add/remove a platform without updating this test? "+
			"IDs: %v", expected, len(all), r.IDs())
	}
}

// Regression_v0_4_0: Windsurf was renamed to Devin, --windsurf must still work.
// Bug: Renaming a platform without preserving the old flag as an alias.
func TestRegression_DevinAlias(t *testing.T) {
	r := platform.NewRegistry()

	// Devin must exist
	p, ok := r.Get("devin")
	if !ok {
		t.Fatal("regression: devin platform not found after Windsurf rename")
	}

	// --windsurf alias must resolve to devin
	resolved, ok := r.ForFlag("--windsurf")
	if !ok {
		t.Fatal("regression: --windsurf alias broken after Windsurf rename")
	}
	if resolved.ID != "devin" {
		t.Errorf("regression: --windsurf resolved to %q, expected devin", resolved.ID)
	}

	// Devin must have --windsurf in its aliases
	hasAlias := false
	for _, a := range p.Aliases {
		if a == "--windsurf" {
			hasAlias = true
			break
		}
	}
	if !hasAlias {
		t.Error("regression: devin missing --windsurf alias")
	}
}

// Regression_v0_4_0: ZCode platform must be registered with correct properties.
// Bug: Adding a platform with incorrect ConfigDir or CLIFlag.
func TestRegression_ZCodeExists(t *testing.T) {
	r := platform.NewRegistry()
	p, ok := r.Get("zcode")
	if !ok {
		t.Fatal("regression: zcode platform not found")
	}

	checks := []struct {
		field    string
		expected string
		actual   string
	}{
		{"ID", "zcode", p.ID},
		{"ConfigDir", ".zcode", p.ConfigDir},
		{"CLIFlag", "--zcode", p.CLIFlag},
	}
	for _, c := range checks {
		if c.actual != c.expected {
			t.Errorf("regression: zcode %s = %q, expected %q", c.field, c.actual, c.expected)
		}
	}
	if p.Class != platform.ClassPushBased {
		t.Errorf("regression: zcode class = %q, expected push", p.Class)
	}
}

// Regression_v0_4_0: ForFlag must resolve both primary flags and aliases correctly.
// Bug: ForFlag returning wrong platform or false for valid flags.
func TestRegression_ForFlagAllPlatforms(t *testing.T) {
	r := platform.NewRegistry()
	all := r.All()

	for _, p := range all {
		// Primary CLIFlag must resolve
		resolved, ok := r.ForFlag(p.CLIFlag)
		if !ok {
			t.Errorf("regression: ForFlag(%q) returned false for platform %q", p.CLIFlag, p.ID)
			continue
		}
		if resolved.ID != p.ID {
			t.Errorf("regression: ForFlag(%q) returned %q, expected %q", p.CLIFlag, resolved.ID, p.ID)
		}

		// All aliases must resolve to this platform
		for _, alias := range p.Aliases {
			resolved, ok := r.ForFlag(alias)
			if !ok {
				t.Errorf("regression: ForFlag(%q) alias returned false for platform %q", alias, p.ID)
				continue
			}
			if resolved.ID != p.ID {
				t.Errorf("regression: ForFlag(%q) alias returned %q, expected %q", alias, resolved.ID, p.ID)
			}
		}
	}
}

// Regression_v0_4_0: All platforms must have valid, non-empty CLIFlag.
// Bug: Adding a platform without setting CLIFlag.
func TestRegression_AllPlatformsHaveCLIFlag(t *testing.T) {
	r := platform.NewRegistry()
	for _, p := range r.All() {
		if p.CLIFlag == "" {
			t.Errorf("regression: platform %q has empty CLIFlag", p.ID)
		}
		if p.CLIFlag[0] != '-' || p.CLIFlag[1] != '-' {
			t.Errorf("regression: platform %q CLIFlag %q does not start with '--'", p.ID, p.CLIFlag)
		}
	}
}

// Regression_v0_4_0: All push-based platforms must have "common" in TemplateDirs.
// Bug: Adding a platform without the common template directory.
func TestRegression_TemplateCommonExists(t *testing.T) {
	r := platform.NewRegistry()
	for _, p := range r.All() {
		hasCommon := false
		for _, td := range p.TemplateDirs {
			if td == "common" {
				hasCommon = true
				break
			}
		}
		if !hasCommon {
			t.Errorf("regression: platform %q missing 'common' in TemplateDirs: %v", p.ID, p.TemplateDirs)
		}
	}
}

// Regression_v0_4_0: Hook generator scripts must have a proper shebang.
// Bug: Generated hook scripts missing shebang cause execution failures.
func TestRegression_HookScriptHasShebang(t *testing.T) {
	r := platform.NewRegistry()
	for _, p := range r.All() {
		if p.Class != platform.ClassPushBased {
			continue
		}
		g := hook.NewGenerator(p, "trellis")
		tmp := t.TempDir()
		if err := g.GenerateAll(tmp); err != nil {
			t.Fatalf("GenerateAll for %q: %v", p.ID, err)
		}
		for _, hookName := range []string{"session-start.sh", "inject-context.sh", "inject-workflow-state.sh"} {
			data, err := os.ReadFile(filepath.Join(tmp, hookName))
			if err != nil {
				t.Errorf("regression: %s hook missing for %q: %v", hookName, p.ID, err)
				continue
			}
			if !strings.HasPrefix(string(data), "#!/bin/sh\n") {
				t.Errorf("regression: %s hook for %q missing shebang, got: %s", hookName, p.ID, string(data[:20]))
			}
		}
	}
}

// Regression_v0_4_0: Generated hook scripts must be executable.
// Bug: Non-executable hook scripts are ignored by platforms.
func TestRegression_HookScriptExecutable(t *testing.T) {
	r := platform.NewRegistry()
	for _, p := range r.All() {
		if p.Class != platform.ClassPushBased {
			continue
		}
		g := hook.NewGenerator(p, "trellis")
		tmp := t.TempDir()
		if err := g.GenerateAll(tmp); err != nil {
			t.Fatalf("GenerateAll for %q: %v", p.ID, err)
		}
		for _, hookName := range []string{"session-start.sh", "inject-context.sh", "inject-workflow-state.sh"} {
			info, err := os.Stat(filepath.Join(tmp, hookName))
			if err != nil {
				t.Errorf("regression: %s hook missing for %q: %v", hookName, p.ID, err)
				continue
			}
			if info.Mode().Perm() != 0755 {
				t.Errorf("regression: %s hook for %q permissions = %o, want 0755", hookName, p.ID, info.Mode().Perm())
			}
		}
	}
}

// Regression_v0_4_0: Version file must exist after init.
// Bug: Missing version file breaks upgrade checks.
func TestRegression_VersionFileExists(t *testing.T) {
	r := platform.NewRegistry()
	for _, p := range r.All() {
		if p.Class != platform.ClassPushBased {
			continue
		}
		_ = p.ID // version file is .trellis/.version, created during init
	}
}
