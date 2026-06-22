package task

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// HookRunner executes lifecycle hooks from config.
type HookRunner struct {
	Hooks map[string][]string
}

// Run executes hooks for the given event.
// Hook failures print warnings but do not block the operation.
func (r *HookRunner) Run(event, taskJSONPath string) {
	cmds, ok := r.Hooks[event]
	if !ok || len(cmds) == 0 {
		return
	}
	for _, cmd := range cmds {
		c := exec.Command("sh", "-c", cmd)
		c.Env = append(os.Environ(), "TASK_JSON_PATH="+taskJSONPath)
		if out, err := c.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: hook %q failed: %v\n%s\n", cmd, err, string(out))
		}
	}
}

// HookEvents returns the standard lifecycle event names.
func HookEvents() []string {
	return []string{"after_create", "after_start", "after_finish", "after_archive"}
}

// BuildHooksFromConfig builds a HookRunner from a config's hooks map.
func BuildHooksFromConfig(hooks map[string]string) *HookRunner {
	if hooks == nil {
		return &HookRunner{Hooks: nil}
	}
	result := make(map[string][]string)
	for event, cmd := range hooks {
		event = strings.TrimSpace(event)
		cmd = strings.TrimSpace(cmd)
		if event != "" && cmd != "" {
			result[event] = []string{cmd}
		}
	}
	return &HookRunner{Hooks: result}
}
