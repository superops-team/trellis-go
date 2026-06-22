package task

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildHooksFromConfig_Nil(t *testing.T) {
	r := BuildHooksFromConfig(nil)
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestBuildHooksFromConfig_Empty(t *testing.T) {
	r := BuildHooksFromConfig(map[string]string{})
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestBuildHooksFromConfig_Valid(t *testing.T) {
	r := BuildHooksFromConfig(map[string]string{
		"after_create": "echo 'created'",
	})
	if len(r.Hooks) != 1 {
		t.Fatalf("expected 1 hook, got %d", len(r.Hooks))
	}
	if r.Hooks["after_create"][0] != "echo 'created'" {
		t.Errorf("unexpected hook command")
	}
}

func TestHookRunner_Run_Success(t *testing.T) {
	r := &HookRunner{
		Hooks: map[string][]string{
			"after_create": {"echo 'ok'"},
		},
	}
	// Should not panic or error
	r.Run("after_create", "/tmp/test.json")
}

func TestHookRunner_Run_NoEvent(t *testing.T) {
	r := &HookRunner{
		Hooks: map[string][]string{
			"after_create": {"echo 'ok'"},
		},
	}
	// Should not panic for non-existent event
	r.Run("after_start", "/tmp/test.json")
}

func TestHookRunner_Run_Failure(t *testing.T) {
	r := &HookRunner{
		Hooks: map[string][]string{
			"after_create": {"exit 1"},
		},
	}
	// Should not panic; failure is just a warning
	r.Run("after_create", "/tmp/test.json")
}

func TestHookRunner_Run_EnvVar(t *testing.T) {
	dir := t.TempDir()
	taskPath := filepath.Join(dir, "task.json")
	os.WriteFile(taskPath, []byte(`{"id":"test"}`), 0644)

	r := &HookRunner{
		Hooks: map[string][]string{
			"after_create": {"cat $TASK_JSON_PATH > " + filepath.Join(dir, "env-check.txt")},
		},
	}
	r.Run("after_create", taskPath)

	data, err := os.ReadFile(filepath.Join(dir, "env-check.txt"))
	if err != nil {
		t.Fatalf("read env check: %v", err)
	}
	if string(data) != `{"id":"test"}` {
		t.Errorf("expected task.json content, got: %s", string(data))
	}
}

func TestHookEvents(t *testing.T) {
	events := HookEvents()
	expected := []string{"after_create", "after_start", "after_finish", "after_archive"}
	if len(events) != len(expected) {
		t.Fatalf("expected %d events, got %d", len(expected), len(events))
	}
	for i, e := range expected {
		if events[i] != e {
			t.Errorf("events[%d] = %q, want %q", i, events[i], e)
		}
	}
}
