package platform

import (
	"fmt"
	"strings"
	"testing"
)

// --- Internal consistency invariants ---

func TestInvariant_PlatformCount(t *testing.T) {
	r := NewRegistry()
	all := r.All()
	if len(all) != len(builtins) {
		t.Fatalf("platform count mismatch: All()=%d, builtins=%d", len(all), len(builtins))
	}
}

func TestInvariant_UniqueIDs(t *testing.T) {
	r := NewRegistry()
	seen := make(map[string]bool)
	for _, p := range r.All() {
		if seen[p.ID] {
			t.Errorf("duplicate platform ID: %q", p.ID)
		}
		seen[p.ID] = true
	}
}

func TestInvariant_UniqueConfigDirs(t *testing.T) {
	r := NewRegistry()
	seen := make(map[string]string) // configDir -> platform ID
	for _, p := range r.All() {
		if prev, ok := seen[p.ConfigDir]; ok {
			t.Errorf("duplicate ConfigDir %q: platforms %q and %q", p.ConfigDir, prev, p.ID)
		}
		seen[p.ConfigDir] = p.ID
	}
}

func TestInvariant_UniqueCLIFlags(t *testing.T) {
	r := NewRegistry()
	seen := make(map[string]string) // cliFlag -> platform ID
	for _, p := range r.All() {
		if prev, ok := seen[p.CLIFlag]; ok {
			t.Errorf("duplicate CLIFlag %q: platforms %q and %q", p.CLIFlag, prev, p.ID)
		}
		seen[p.CLIFlag] = p.ID
	}
}

func TestInvariant_NoReservedWordIDs(t *testing.T) {
	r := NewRegistry()
	reserved := map[string]bool{
		"all": true, "list": true, "init": true, "task": true,
		"update": true, "upgrade": true, "help": true, "version": true,
	}
	for _, p := range r.All() {
		if reserved[p.ID] {
			t.Errorf("platform %q uses reserved word as ID", p.ID)
		}
	}
}

func TestInvariant_IDsAreLowercaseAlphanumeric(t *testing.T) {
	r := NewRegistry()
	for _, p := range r.All() {
		for _, ch := range p.ID {
			if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-') {
				t.Errorf("platform %q: ID contains invalid character %q", p.ID, ch)
			}
		}
	}
}

func TestInvariant_ConfigDirsStartWithDot(t *testing.T) {
	r := NewRegistry()
	for _, p := range r.All() {
		if !strings.HasPrefix(p.ConfigDir, ".") {
			t.Errorf("platform %q: ConfigDir %q does not start with '.'", p.ID, p.ConfigDir)
		}
	}
}

func TestInvariant_AllPlatformsPassValidate(t *testing.T) {
	r := NewRegistry()
	for _, p := range r.All() {
		if err := p.Validate(); err != nil {
			t.Errorf("platform %q fails Validate(): %v", p.ID, err)
		}
	}
}

func TestInvariant_IDsSorted(t *testing.T) {
	r := NewRegistry()
	ids := r.IDs()
	for i := 1; i < len(ids); i++ {
		if ids[i] < ids[i-1] {
			t.Errorf("IDs not sorted: %q before %q", ids[i-1], ids[i])
		}
	}
}

func TestInvariant_AllSorted(t *testing.T) {
	r := NewRegistry()
	all := r.All()
	for i := 1; i < len(all); i++ {
		if all[i].ID < all[i-1].ID {
			t.Errorf("All() not sorted by ID: %q before %q", all[i-1].ID, all[i].ID)
		}
	}
}

// --- Alias consistency invariants ---

func TestInvariant_DevinHasWindsurfAlias(t *testing.T) {
	r := NewRegistry()
	p, ok := r.Get("devin")
	if !ok {
		t.Fatal("devin platform not found")
	}
	found := false
	for _, a := range p.Aliases {
		if a == "--windsurf" {
			found = true
			break
		}
	}
	if !found {
		t.Error("devin should have --windsurf alias")
	}
}

func TestInvariant_ForFlagResolvesAliases(t *testing.T) {
	r := NewRegistry()
	// --windsurf alias should resolve to devin
	p, ok := r.ForFlag("--windsurf")
	if !ok {
		t.Fatal("--windsurf alias not resolved")
	}
	if p.ID != "devin" {
		t.Errorf("--windsurf resolved to %q, expected devin", p.ID)
	}
}

func TestInvariant_NoDuplicateAliases(t *testing.T) {
	r := NewRegistry()
	seen := make(map[string]string) // alias -> platform ID
	for _, p := range r.All() {
		for _, a := range p.Aliases {
			if prev, ok := seen[a]; ok {
				t.Errorf("duplicate alias %q: platforms %q and %q", a, prev, p.ID)
			}
			seen[a] = p.ID
		}
	}
}

func TestInvariant_AliasNotEqualToOtherCLIFlag(t *testing.T) {
	r := NewRegistry()
	all := r.All()
	for _, p := range all {
		for _, a := range p.Aliases {
			for _, other := range all {
				if other.ID != p.ID && other.CLIFlag == a {
					t.Errorf("platform %q alias %q conflicts with %q CLIFlag", p.ID, a, other.ID)
				}
			}
		}
	}
}

// --- Platform type consistency invariants ---

func TestInvariant_ValidClasses(t *testing.T) {
	r := NewRegistry()
	valid := map[Class]bool{ClassPushBased: true, ClassPullBased: true, ClassAgentless: true}
	for _, p := range r.All() {
		if !valid[p.Class] {
			t.Errorf("platform %q has invalid class %q", p.ID, p.Class)
		}
	}
}

func TestInvariant_ClassCountsMatchTotal(t *testing.T) {
	r := NewRegistry()
	all := r.All()
	push := r.ByClass(ClassPushBased)
	pull := r.ByClass(ClassPullBased)
	none := r.ByClass(ClassAgentless)

	total := len(push) + len(pull) + len(none)
	if total != len(all) {
		t.Errorf("class counts (%d+%d+%d=%d) != total (%d)",
			len(push), len(pull), len(none), total, len(all))
	}
}

func TestInvariant_EachClassHasAtLeastOnePlatform(t *testing.T) {
	r := NewRegistry()
	for _, c := range []Class{ClassPushBased, ClassPullBased, ClassAgentless} {
		platforms := r.ByClass(c)
		if len(platforms) == 0 {
			t.Errorf("no platforms with class %q", c)
		}
	}
}

// --- Cross-reference invariants ---

func TestInvariant_GetMatchesAll(t *testing.T) {
	r := NewRegistry()
	for _, p := range r.All() {
		got, ok := r.Get(p.ID)
		if !ok {
			t.Errorf("Get(%q) returned not found", p.ID)
			continue
		}
		if got.ID != p.ID {
			t.Errorf("Get(%q) returned platform with ID %q", p.ID, got.ID)
		}
	}
}

func TestInvariant_ForFlagMatchesCLIFlag(t *testing.T) {
	r := NewRegistry()
	for _, p := range r.All() {
		got, ok := r.ForFlag(p.CLIFlag)
		if !ok {
			t.Errorf("ForFlag(%q) returned not found for platform %q", p.CLIFlag, p.ID)
			continue
		}
		if got.ID != p.ID {
			t.Errorf("ForFlag(%q) returned %q, expected %q", p.CLIFlag, got.ID, p.ID)
		}
	}
}

// --- Edge case invariants ---

func TestInvariant_RegisterThenGet(t *testing.T) {
	r := NewRegistry()
	p := Platform{
		ID:        "test-platform",
		Name:      "Test Platform",
		ConfigDir: ".test-platform",
		CLIFlag:   "--test-platform",
		Class:     ClassPushBased,
	}
	if err := r.Register(p); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	got, ok := r.Get("test-platform")
	if !ok {
		t.Fatal("Get after Register returned not found")
	}
	if got.ID != "test-platform" {
		t.Errorf("expected test-platform, got %q", got.ID)
	}
}

func TestInvariant_RegisterRejectsInvalid(t *testing.T) {
	r := NewRegistry()
	invalid := []Platform{
		{ID: "", Name: "NoID", ConfigDir: ".x", Class: ClassPushBased},
		{ID: "x", Name: "NoConfigDir", ConfigDir: "", Class: ClassPushBased},
		{ID: "x", Name: "BadClass", ConfigDir: ".x", Class: "invalid"},
		{ID: "x", Name: "LeadingSlash", ConfigDir: "/x", Class: ClassPushBased},
	}
	for _, p := range invalid {
		t.Run(fmt.Sprintf("invalid_%s", p.Name), func(t *testing.T) {
			if err := r.Register(p); err == nil {
				t.Errorf("expected error for invalid platform %q", p.Name)
			}
		})
	}
}

func TestInvariant_RegisterRejectsDuplicate(t *testing.T) {
	r := NewRegistry()
	p := Platform{
		ID:        "claude",
		Name:      "Claude Duplicate",
		ConfigDir: ".claude-dup",
		CLIFlag:   "--claude-dup",
		Class:     ClassPushBased,
	}
	if err := r.Register(p); err == nil {
		t.Error("expected error for duplicate platform ID")
	}
}
