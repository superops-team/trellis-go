package platform

import (
	"testing"
)

func TestNewRegistry_HasBuiltins(t *testing.T) {
	r := NewRegistry()
	all := r.All()
	if len(all) != len(builtins) {
		t.Fatalf("expected %d platforms, got %d", len(builtins), len(all))
	}
}

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()
	p, ok := r.Get("claude")
	if !ok {
		t.Fatal("expected claude to exist")
	}
	if p.ID != "claude" {
		t.Errorf("expected ID claude, got %s", p.ID)
	}
	if p.Class != ClassPushBased {
		t.Errorf("expected class push, got %s", p.Class)
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	r := NewRegistry()
	_, ok := r.Get("nonexistent")
	if ok {
		t.Fatal("expected nonexistent platform to not be found")
	}
}

func TestRegistry_ByClass(t *testing.T) {
	r := NewRegistry()
	push := r.ByClass(ClassPushBased)
	pull := r.ByClass(ClassPullBased)
	none := r.ByClass(ClassAgentless)

	if len(push) == 0 {
		t.Error("expected push-based platforms")
	}
	if len(pull) == 0 {
		t.Error("expected pull-based platforms")
	}
	if len(none) == 0 {
		t.Error("expected agentless platforms")
	}

	for _, p := range push {
		if p.Class != ClassPushBased {
			t.Errorf("platform %q: expected class push, got %s", p.ID, p.Class)
		}
	}
}

func TestRegistry_Register_Duplicate(t *testing.T) {
	r := NewRegistry()
	p := Platform{ID: "claude", Name: "Claude", ConfigDir: ".claude", Class: ClassPushBased}
	err := r.Register(p)
	if err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestRegistry_Register_Invalid(t *testing.T) {
	r := NewRegistry()
	cases := []struct {
		name string
		p    Platform
	}{
		{"empty id", Platform{ID: "", Name: "X", ConfigDir: ".x", Class: ClassPushBased}},
		{"empty config_dir", Platform{ID: "x", Name: "X", ConfigDir: "", Class: ClassPushBased}},
		{"leading slash", Platform{ID: "x", Name: "X", ConfigDir: "/x", Class: ClassPushBased}},
		{"invalid class", Platform{ID: "x", Name: "X", ConfigDir: ".x", Class: "invalid"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := r.Register(tc.p)
			if err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestRegistry_IDs(t *testing.T) {
	r := NewRegistry()
	ids := r.IDs()
	if len(ids) != len(builtins) {
		t.Fatalf("expected %d IDs, got %d", len(builtins), len(ids))
	}
	for i := 1; i < len(ids); i++ {
		if ids[i] < ids[i-1] {
			t.Error("expected IDs to be sorted")
		}
	}
}
