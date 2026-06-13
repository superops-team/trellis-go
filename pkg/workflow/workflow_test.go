package workflow

import (
	"strings"
	"testing"
)

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

func TestParser_ExtractState(t *testing.T) {
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
		{"some text\n[workflow-state:PLAN]\nmore", StatePlan, false},
		{"no state here", "", true},
		{"[workflow-state:UNKNOWN]", "", true},
	}
	for _, tc := range cases {
		t.Run(string(tc.want), func(t *testing.T) {
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

func TestParser_InjectPrompt(t *testing.T) {
	p := &Parser{}
	for _, state := range []State{StatePlan, StateImplement, StateVerify, StateFinish} {
		t.Run(string(state), func(t *testing.T) {
			prompt, err := p.InjectPrompt(state)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(prompt, strings.ToUpper(string(state))) {
				t.Errorf("prompt should contain state name: %s", prompt)
			}
		})
	}
}
