package workflow

import (
	"strings"
	"testing"
)

func TestParse_Phases(t *testing.T) {
	input := `## Phase 1: Plan
#### 1.1 Brainstorm
Think about the problem.

#### 1.2 Write PRD
Document requirements.

## Phase 2: Implement
#### 2.1 Write Code
Implement the solution.
`
	wf, err := (&Parser{}).Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if len(wf.Phases) != 2 {
		t.Fatalf("expected 2 phases, got %d", len(wf.Phases))
	}

	if wf.Phases[0].Title != "Plan" {
		t.Errorf("expected phase 1 title 'Plan', got %q", wf.Phases[0].Title)
	}
	if len(wf.Phases[0].Steps) != 2 {
		t.Errorf("expected 2 steps in phase 1, got %d", len(wf.Phases[0].Steps))
	}
	if wf.Phases[0].Steps[0].Number != "1.1" {
		t.Errorf("expected step 1.1, got %s", wf.Phases[0].Steps[0].Number)
	}
}

func TestParse_Breadcrumbs(t *testing.T) {
	input := `[workflow-state:no_task]No active task. Start by creating one.[/workflow-state:no_task]
[workflow-state:planning]Plan the implementation.[/workflow-state:planning]
[workflow-state:in_progress]Implement the changes.[/workflow-state:in_progress]
[workflow-state:completed]Task is done.[/workflow-state:completed]
`
	wf, err := (&Parser{}).Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	cases := []struct {
		state State
		want  string
	}{
		{State("no_task"), "No active task. Start by creating one."},
		{State("planning"), "Plan the implementation."},
		{State("in_progress"), "Implement the changes."},
		{State("completed"), "Task is done."},
	}
	for _, tc := range cases {
		got, ok := wf.Breadcrumbs[tc.state]
		if !ok {
			t.Errorf("missing breadcrumb for state %s", tc.state)
			continue
		}
		if got != tc.want {
			t.Errorf("breadcrumb[%s] = %q, want %q", tc.state, got, tc.want)
		}
	}
}

func TestParse_SkillRouting(t *testing.T) {
	input := `### Skill Routing
| User intent | Skill |
|-------------|-------|
| New feature | trellis-brainstorm |
| Before coding | trellis-before-dev |
| Verify code | trellis-check |
`
	wf, err := (&Parser{}).Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if len(wf.SkillRouting) != 3 {
		t.Fatalf("expected 3 skill routes, got %d", len(wf.SkillRouting))
	}
	if wf.SkillRouting[0].Intent != "New feature" {
		t.Errorf("expected intent 'New feature', got %q", wf.SkillRouting[0].Intent)
	}
	if wf.SkillRouting[0].Skill != "trellis-brainstorm" {
		t.Errorf("expected skill 'trellis-brainstorm', got %q", wf.SkillRouting[0].Skill)
	}
}

func TestParse_DoNotSkip(t *testing.T) {
	input := `### DO NOT skip skills
trellis-check and trellis-update-spec must always run.
`
	wf, err := (&Parser{}).Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if !strings.Contains(wf.DoNotSkip, "trellis-check") {
		t.Error("DoNotSkip should contain trellis-check")
	}
}

func TestParse_TaskSystem(t *testing.T) {
	input := `### Task System
Tasks are managed via .trellis/tasks/ directory.
`
	wf, err := (&Parser{}).Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if !strings.Contains(wf.TaskSystem, ".trellis/tasks/") {
		t.Error("TaskSystem should contain .trellis/tasks/")
	}
}

func TestExtractState(t *testing.T) {
	p := &Parser{}
	cases := []struct {
		input   string
		want    State
		wantErr bool
	}{
		{"[workflow-state:PLAN]", StatePlan, false},
		{"[workflow-state:IMPLEMENT]", StateImplement, false},
		{"[workflow-state:VERIFY]", StateVerify, false},
		{"[workflow-state:FINISH]", StateFinish, false},
		{"[workflow-state:BLOCKED]", State("blocked"), false},
		{"some text\n[workflow-state:PLAN]\nmore", StatePlan, false},
		{"no state here", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.input[:min(len(tc.input), 30)], func(t *testing.T) {
			got, err := p.ExtractState(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestInjectPrompt_CustomState(t *testing.T) {
	p := &Parser{}
	prompt, err := p.InjectPrompt(State("blocked"))
	if err != nil {
		t.Fatalf("InjectPrompt(blocked) error: %v", err)
	}
	if !strings.Contains(prompt, "Refer to workflow.md") {
		t.Errorf("expected generic fallback for unknown state, got: %s", prompt)
	}
}

func TestInjectPrompt_NoTask(t *testing.T) {
	p := &Parser{}
	prompt, err := p.InjectPrompt(StateNoTask)
	if err != nil {
		t.Fatalf("InjectPrompt(no_task) error: %v", err)
	}
	if !strings.Contains(prompt, "NO TASK") {
		t.Errorf("expected NO TASK prompt, got: %s", prompt)
	}
}

func TestExtractStep(t *testing.T) {
	input := `## Phase 2: Implement
#### 2.1 Write Code
Implement the solution according to the spec.
Make sure to follow conventions.

#### 2.2 Verify
Run tests and lint.
`
	p := &Parser{}
	body, err := p.ExtractStep(input, "2.1")
	if err != nil {
		t.Fatalf("ExtractStep(2.1) error: %v", err)
	}
	if !strings.Contains(body, "Implement the solution") {
		t.Errorf("expected step body, got: %s", body)
	}
	if !strings.Contains(body, "Make sure") {
		t.Errorf("expected multi-line body, got: %s", body)
	}
}

func TestExtractStep_NotFound(t *testing.T) {
	input := `## Phase 1: Plan`
	p := &Parser{}
	_, err := p.ExtractStep(input, "9.9")
	if err == nil {
		t.Fatal("expected error for non-existent step")
	}
}

func TestParseSkillRouting(t *testing.T) {
	input := `### Skill Routing
| User intent | Skill |
|-------------|-------|
| New feature | trellis-brainstorm |
| Bug fix | trellis-break-loop |
`
	p := &Parser{}
	routes, err := p.ParseSkillRouting(input)
	if err != nil {
		t.Fatalf("ParseSkillRouting() error: %v", err)
	}
	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(routes))
	}
	if routes[1].Intent != "Bug fix" {
		t.Errorf("expected 'Bug fix', got %q", routes[1].Intent)
	}
}

func TestStateMachine_CanTransition(t *testing.T) {
	sm := NewStateMachine()
	cases := []struct {
		from, to string
		want     bool
	}{
		{"plan", "implement", true},
		{"plan", "verify", false},
		{"implement", "verify", true},
		{"verify", "implement", true},
		{"verify", "finish", true},
		{"finish", "plan", false},
	}
	for _, tc := range cases {
		t.Run(tc.from+"_"+tc.to, func(t *testing.T) {
			got := sm.CanTransition(State(tc.from), State(tc.to))
			if got != tc.want {
				t.Errorf("CanTransition(%s, %s) = %v, want %v", tc.from, tc.to, got, tc.want)
			}
		})
	}
}

func TestParse_FullWorkflow(t *testing.T) {
	input := `## Phase Index
[workflow-state:no_task]No active task. Start by creating one.[/workflow-state:no_task]
[workflow-state:planning]Plan the implementation.[/workflow-state:planning]
[workflow-state:in_progress]Implement the changes.[/workflow-state:in_progress]
[workflow-state:completed]Task is done.[/workflow-state:completed]

### Phase 1: Plan
#### 1.1 Brainstorm
Think about the problem.

#### 1.2 Write PRD
Document requirements.

### Phase 2: Implement
#### 2.1 Write Code
Implement the solution.

### Phase 3: Finish
#### 3.1 Update Spec
Update spec with learnings.

#### 3.2 Commit
Commit and push.

### Skill Routing
| User intent | Skill |
|-------------|-------|
| New feature | trellis-brainstorm |
| Before coding | trellis-before-dev |

### DO NOT skip skills
trellis-check must always run.

### Task System
Tasks are managed via .trellis/tasks/
`
	wf, err := (&Parser{}).Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if len(wf.Phases) != 3 {
		t.Errorf("expected 3 phases, got %d", len(wf.Phases))
	}
	if len(wf.Breadcrumbs) != 4 {
		t.Errorf("expected 4 breadcrumbs, got %d", len(wf.Breadcrumbs))
	}
	if len(wf.SkillRouting) != 2 {
		t.Errorf("expected 2 skill routes, got %d", len(wf.SkillRouting))
	}
	if !strings.Contains(wf.DoNotSkip, "trellis-check") {
		t.Error("DoNotSkip should contain trellis-check")
	}
	if !strings.Contains(wf.TaskSystem, ".trellis/tasks/") {
		t.Error("TaskSystem should contain .trellis/tasks/")
	}
}
